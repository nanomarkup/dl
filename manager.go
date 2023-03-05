// Copyright 2022 Vitalii Noha vitalii.noga@gmail.com. All rights reserved.

package dl

import (
	"fmt"
	"os"
	"path/filepath"
)

func (m *Manager) AddItem(module, item string) error {
	return addItem(module, m.Kind, item)
}

func (m *Manager) AddDependency(item, dependency, resolver string, update bool) error {
	mod, err := findItem(m.Kind, item)
	if err != nil {
		return err
	} else if mod == nil {
		return fmt.Errorf(ItemIsMissingF, item)
	} else {
		if err = mod.AddDependency(item, dependency, resolver, update); err != nil {
			return err
		} else {
			return saveModule(mod, m.Kind)
		}
	}
}

func (m *Manager) DeleteItem(item string) error {
	mod, err := findItem(m.Kind, item)
	if err != nil {
		return err
	} else if mod != nil {
		if err = mod.DeleteItem(item); err != nil {
			return err
		} else {
			return saveModule(mod, m.Kind)
		}
	} else {
		return nil
	}
}

func (m *Manager) DeleteDependency(item, dependency string) error {
	mod, err := findItem(m.Kind, item)
	if err != nil {
		return err
	} else if mod != nil {
		if err = mod.DeleteDependency(item, dependency); err != nil {
			return err
		} else {
			return saveModule(mod, m.Kind)
		}
	} else {
		return nil
	}
}

func (m *Manager) Read(filePath string) (Module, error) {
	m.logTrace(fmt.Sprintf("loading \"%s\" file", filePath))
	ext := filepath.Ext(filePath)
	name := filePath[0 : len(filePath)-len(ext)]
	file, err := os.Open(getModuleFileName(name, ext[1:]))
	if err != nil {
		return &module{}, err
	}
	defer file.Close()

	mod, err := newReader(file).read()
	if err == nil {
		if mod == nil {
			return &module{}, fmt.Errorf(ModuleErrorOnLoadingF, filePath)
		}
		m.logTrace("reading items")
		mods := modules{}
		mods = append(mods, *mod)
		return loadItems(mods)
	} else {
		return &module{}, err
	}
}

func (m *Manager) ReadAll() (Module, error) {
	m.logTrace(fmt.Sprintf("loading \"%s\" modules", m.Kind))
	mods, err := loadModules(m.Kind)
	if err == nil {
		if mods == nil {
			return &module{}, fmt.Errorf(ModuleErrorOnLoadingF, m.Kind)
		}
		m.logTrace("reading items")
		return loadItems(mods)
	} else {
		return &module{}, err
	}
}

func (m *Manager) SetLogger(logger Logger) {
	m.Logger = logger
}

//lint:ignore U1000 Ignore unused function
func (m *Manager) logTrace(message string) {
	if m.Logger != nil {
		m.Logger.Trace(message)
	}
}

//lint:ignore U1000 Ignore unused function
func (m *Manager) logDebug(message string) {
	if m.Logger != nil {
		m.Logger.Debug(message)
	}
}

//lint:ignore U1000 Ignore unused function
func (m *Manager) logInfo(message string) {
	if m.Logger != nil {
		m.Logger.Info(message)
	}
}

//lint:ignore U1000 Ignore unused function
func (m *Manager) logWarn(message string) {
	if m.Logger != nil {
		m.Logger.Warn(message)
	}
}

//lint:ignore U1000 Ignore unused function
func (m *Manager) logError(message string) {
	if m.Logger != nil {
		m.Logger.Error(message)
	}
}
