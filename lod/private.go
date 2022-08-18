// Copyright 2022 Vitalii Noha vitalii.noga@gmail.com. All rights reserved.

package lod

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

const (
	moduleExt string = ".sb"
)

type module struct {
	name  string
	kind  string
	items Items
}

type moduleAsync struct {
	mod *module
	err error
}

type Item = map[string]string
type Items = map[string]Item
type modules []module

func getModuleName(fileName string) string {
	if strings.HasSuffix(fileName, moduleExt) {
		return fileName[0 : len(fileName)-len(moduleExt)]
	} else {
		return fileName
	}
}

func split(line string) []string {
	var res []string
	its := strings.Split(line, " ")
	add := true
	ind := -1
	for _, it := range its {
		if add {
			res = append(res, it)
			ind++
			if len(it) > 0 && it[0] == '"' {
				add = false
			}
		} else {
			res[ind] = res[ind] + " " + it
			if len(it) > 0 && it[len(it)-1] == '"' {
				add = true
			}
		}
	}
	return res
}

func loadModule(name string, kind string) (*module, error) {
	mod := module{}
	mod.name = name
	mod.items = Items{}

	fileName := GetModuleFileName(name)
	file, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := bufio.NewReader(file)
	item := ""
	line := ""
	token1 := ""
	token2 := ""
	pos := 0
	index := 1
	cindex := 0
	trimChars := " \t\n\r"
	for {
		line, err = reader.ReadString('\n')
		if err != nil && err != io.EOF {
			return nil, err
		}
		// remove comment
		cindex = strings.Index(line, "//")
		if cindex > 0 {
			line = line[:cindex]
		}
		// process the line
		line = strings.Trim(line, trimChars)
		if line != "" {
			pos = strings.Index(line, " ")
			if pos > 0 {
				token1 = strings.Trim(line[0:pos], trimChars)
				token2 = strings.Trim(line[pos:], trimChars)
			} else {
				token1 = line
				token2 = ""
			}
			if index == 1 {
				// check and initialize a kind of module
				if token1 != kind || token2 != "" {
					return nil, fmt.Errorf("the first token should be \"%s\"", kind)
				}
				mod.kind = token1
			} else {
				// process items
				if token1[len(token1)-1:] == ":" {
					if token2 != "" {
						return nil, fmt.Errorf("invalid syntax in \"%s\" line", line)
					}
					// parse the next item
					item = token1[:len(token1)-1]
				} else {
					if item != "apps" && token2 == "" {
						return nil, fmt.Errorf("invalid syntax in \"%s\" line", line)
					}
					// add new dependency item
					if mod.items[item] == nil {
						mod.items[item] = Item{}
					}
					mod.items[item][token1] = token2
				}
			}
			index++
		}
		// check the EOF
		if err != nil {
			break
		}
	}
	return &mod, nil
}

func loadModuleAsync(name string, kind string, res chan<- moduleAsync) {
	m, e := loadModule(name, kind)
	res <- moduleAsync{m, e}
}

func loadModules(kind string) (modules, error) {
	if kind == "" {
		return nil, fmt.Errorf("kind of modules to load is not specified")
	}
	// read and check all modules in the working directory
	files, err := ioutil.ReadDir(".")
	if err != nil {
		return nil, err
	}
	mods := modules{}
	modExt := fmt.Sprintf(".%s", kind)
	var item chan moduleAsync
	items := []chan moduleAsync{}
	for _, f := range files {
		fname := f.Name()
		if filepath.Ext(fname) != modExt {
			continue
		}
		item = make(chan moduleAsync)
		items = append(items, item)
		go loadModuleAsync(fname, kind, item)
	}
	// wait and process all loaded modules
	for _, it := range items {
		item := <-it
		if err != nil {
			continue
		}
		if item.err != nil {
			err = item.err
			continue
		}
		// validate the loaded module
		if kind != item.mod.kind {
			continue
		}
		// add module
		mods = append(mods, module{name: getModuleName(item.mod.name), kind: item.mod.kind, items: item.mod.items})
	}
	if err != nil {
		return nil, err
	} else if len(mods) > 0 {
		return mods, nil
	} else {
		wd, _ := os.Getwd()
		return nil, fmt.Errorf(ModuleFilesMissingF, wd)
	}
}

func loadItems(mods modules) (*module, error) {
	all := Items{}
	kind := ""
	if len(mods) > 0 {
		kind = mods[0].kind
	}
	for _, m := range mods {
		// read all items and validate them
		for name, data := range m.items {
			if _, found := all[name]; found {
				return nil, fmt.Errorf("\"%s\" item of \"%s\" module already exists", name, m.name)
			}
			all[name] = data
		}
	}
	// process defines
	newItem := ""
	var err error
	if defines, found := all["defines"]; found && len(defines) > 0 {
		for item, deps := range all {
			if item == "defines" {
				continue
			}
			// update item name
			newItem, err = applyDefines(item, defines)
			if err != nil {
				return nil, err
			}
			if newItem != item {
				all[newItem] = deps
				delete(all, item)
			}
			// process all dependencies
			for dk, dv := range deps {
				// update dependency name
				newItem, err = applyDefines(dk, defines)
				if err != nil {
					return nil, err
				}
				if newItem != dk {
					deps[newItem] = dv
					delete(deps, dk)
					dk = newItem
				}
				// update resolver
				newItem, err = applyDefines(dv, defines)
				if err != nil {
					return nil, err
				}
				if newItem != dv {
					deps[dk] = newItem
				}
			}
		}
		delete(all, "defines")
	}
	return &module{name: "", kind: kind, items: all}, nil
}

func saveModule(module *module) error {
	fileName := GetModuleFileName(module.name)
	exists := IsModuleExists(fileName)
	file, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer file.Close()
	// notify about a new module has been created
	defer func() {
		if !exists {
			fmt.Printf(ModuleIsCreatedF, fileName)
		}
	}()
	// save the module
	writer := bufio.NewWriter(file)
	defer writer.Flush()
	f := Formatter{}
	_, err = writer.WriteString(f.String(module))
	return err
}

func addItem(moduleName, kind, item string) error {
	// check the item is exist
	if found, modName := IsItemExists(kind, item); found {
		return fmt.Errorf(ItemExistsF, item, modName)
	}
	// load the existing module or create a new one
	var mod *module
	var err error
	if IsModuleExists(moduleName) {
		if mod, err = loadModule(moduleName, kind); err != nil {
			return err
		}
		// check kind of the selected module
		if mod.kind != kind {
			return fmt.Errorf(ModuleKindMismatchF, mod.kind, mod.name, kind)
		}
	} else {
		mod = &module{name: moduleName, kind: kind, items: Items{}}
	}
	// add the item to the selected module
	if err = mod.AddItem(item); err != nil {
		return err
	} else {
		return saveModule(mod)
	}
}

func findItem(kind, item string) (*module, error) {
	wd, _ := os.Getwd()
	mods, err := loadModules(kind)
	if (err != nil) && (err.Error() != fmt.Sprintf(ModuleFilesMissingF, wd)) {
		return nil, err
	}
	for _, m := range mods {
		if _, found := m.items[item]; found {
			return &m, nil
		}
	}
	return nil, nil
}

func applyDefines(item string, defines map[string]string) (string, error) {
	beg := strings.Index(item, "{")
	end := -1
	defineOrg := ""
	defineName := ""
	for beg > -1 {
		end = strings.Index(item, "}")
		if end < beg {
			return "", fmt.Errorf("\"%s\" incorrect item name", item)
		}
		defineOrg = item[beg : end+1]
		defineName = strings.Trim(defineOrg, " {}")
		if value, found := defines[defineName]; found {
			item = strings.Replace(item, defineOrg, value, 1)
		} else {
			return "", fmt.Errorf("\"%s\" define is not declared", defineName)
		}
		beg = strings.Index(item, "{")
	}
	return item, nil
}
