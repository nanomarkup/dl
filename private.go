// Copyright 2022 Vitalii Noha vitalii.noga@gmail.com. All rights reserved.

package dl

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

type reader struct {
	name string
	buf  *bufio.Reader
}

type module struct {
	name  string
	items Items
}

type moduleAsync struct {
	mod *module
	err error
}

type Item = [][]string
type Items = map[string]Item
type modules []module

var groupId = 0
var trimChars = " \t\n\r"

func getModuleExt(kind string) string {
	return "." + kind
}

func getModuleName(fileName, kind string) string {
	ext := getModuleExt(kind)
	if strings.HasSuffix(fileName, ext) {
		return fileName[0 : len(fileName)-len(ext)]
	} else {
		return fileName
	}
}

func getModuleFileName(name, kind string) string {
	ext := getModuleExt(kind)
	if strings.HasSuffix(name, ext) {
		return name
	} else {
		return name + ext
	}
}

func getNextGroupName() string {
	groupId++
	return "Gen" + strconv.Itoa(groupId)
}

func loadModules(kind string) (modules, error) {
	if kind == "" {
		return nil, fmt.Errorf(ModuleKindIsMissing)
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

		file, err := os.Open(getModuleFileName(fname, kind))
		if err != nil {
			return nil, err
		}
		defer file.Close()
		go newReader(file).readAsync(item)
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
		// add module
		mods = append(mods, module{name: getModuleName(item.mod.name, kind), items: item.mod.items})
	}
	if err != nil {
		return nil, err
	} else if len(mods) > 0 {
		return mods, nil
	} else {
		wd, _ := os.Getwd()
		return nil, fmt.Errorf(ModuleFilesMissingF, kind, wd)
	}
}

func loadItems(mods modules) (*module, error) {
	all := Items{}
	for _, m := range mods {
		// read all items and validate them
		for name, data := range m.items {
			if _, found := all[name]; found {
				return nil, fmt.Errorf(ItemExistsInModuleF, name, m.name)
			}
			all[name] = data
		}
	}
	// process defines
	newItem := ""
	var l int
	var err error
	if defines, found := all[DefinesOptCode]; found && len(defines) > 0 {
		for item, deps := range all {
			if item == DefinesOptCode {
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
			for i, row := range deps {
				l = len(row)
				if l > 0 {
					// update dependency name
					newItem, err = applyDefines(row[0], defines)
					if err != nil {
						return nil, err
					}
					if newItem != row[0] {
						deps[i][0] = newItem
					}
				}
				if l > 1 {
					// update resolver
					newItem, err = applyDefines(row[1], defines)
					if err != nil {
						return nil, err
					}
					if newItem != row[1] {
						deps[i][1] = newItem
					}
				}
			}
		}
		delete(all, DefinesOptCode)
	}
	return &module{name: "", items: all}, nil
}

func saveModule(module *module, kind string) error {
	fileName := getModuleFileName(module.name, kind)
	exists := isModuleExists(fileName)
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
	if found, modName := isItemExists(kind, item); found {
		return fmt.Errorf(ItemExistsInModuleF, item, modName)
	}
	// load the existing module or create a new one
	var mod *module
	var err error
	fileName := getModuleFileName(moduleName, kind)
	if isModuleExists(fileName) {
		file, err := os.Open(fileName)
		if err != nil {
			return err
		}
		defer file.Close()
		if mod, err = newReader(file).read(); err != nil {
			return err
		}
	} else {
		mod = &module{name: moduleName, items: Items{}}
	}
	// add the item to the selected module
	if err = mod.AddItem(item); err != nil {
		return err
	} else {
		return saveModule(mod, kind)
	}
}

func findItem(kind, item string) (*module, error) {
	wd, _ := os.Getwd()
	mods, err := loadModules(kind)
	if (err != nil) && (err.Error() != fmt.Sprintf(ModuleFilesMissingF, kind, wd)) {
		return nil, err
	}
	for _, m := range mods {
		if _, found := m.items[item]; found {
			return &m, nil
		}
	}
	return nil, nil
}

func applyDefines(item string, defines [][]string) (string, error) {
	beg := strings.Index(item, DefineBegOptCode)
	end := -1
	defineOrg := ""
	defineName := ""
	trimDefineChars := fmt.Sprintf(" %s%s", DefineBegOptCode, DefineEndOptCode)
	value := ""
	found := false
	for beg > -1 {
		end = strings.Index(item, DefineEndOptCode)
		if end < beg {
			return "", fmt.Errorf(ItemNameInvalidF, item)
		}
		defineOrg = item[beg : end+1]
		defineName = strings.Trim(defineOrg, trimDefineChars)
		value = ""
		found = false
		for _, def := range defines {
			if def[0] == defineName {
				found = true
				if len(def) > 1 {
					value = def[1]
				}
				break
			}
		}
		if found {
			item = strings.Replace(item, defineOrg, value, 1)
		} else {
			return "", fmt.Errorf(DefineIsMissingF, defineName)
		}
		beg = strings.Index(item, DefineBegOptCode)
	}
	return item, nil
}

func isItemExists(kind, item string) (bool, string) {
	wd, _ := os.Getwd()
	mods, err := loadModules(kind)
	if (err != nil) && (err.Error() != fmt.Sprintf(ModuleFilesMissingF, kind, wd)) {
		return false, ""
	}
	for _, m := range mods {
		if _, found := m.items[item]; found {
			return true, m.name
		}
	}
	return false, ""
}

func isModuleExists(fileName string) bool {
	_, err := os.Stat(fileName)
	return err == nil
}
