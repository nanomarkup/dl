// Copyright 2022 Vitalii Noha vitalii.noga@gmail.com. All rights reserved.

// Package smodule manages modules.
package lod

import (
	"fmt"
	"os"
	"strings"
)

const (
	AppsItemName string = "apps"
	// notifications
	ModuleIsCreatedF string = "%s file has been created\n"
	// errors
	ItemExistsF         string = "the %s item already exists in %s module"
	ItemIsMissingF      string = "the %s item does not exist"
	ModuleFilesMissingF string = "no sb files in %s"
	ModuleKindMismatchF string = "the %s kind of %s module is mismatch the %s selected kind"
)

type Logger interface {
	Trace(msg string, args ...interface{})
	Debug(msg string, args ...interface{})
	Info(msg string, args ...interface{})
	Warn(msg string, args ...interface{})
	Error(msg string, args ...interface{})
	IsTrace() bool
	IsDebug() bool
	IsInfo() bool
	IsWarn() bool
	IsError() bool
}

type Manager struct {
	Logger Logger
}

type Reader interface {
	Kind() string
	Items() map[string]map[string]string
	Dependency(string, string) string
	App(string) (map[string]string, error)
	Apps() (map[string]string, error)
}

type Formatter struct {
}

func GetModuleFileName(name string) string {
	if strings.HasSuffix(name, moduleExt) {
		return name
	} else {
		return name + moduleExt
	}
}

func IsModuleExists(name string) bool {
	_, err := os.Stat(GetModuleFileName(name))
	return err == nil
}

func IsItemExists(kind, item string) (bool, string) {
	wd, _ := os.Getwd()
	mods, err := loadModules(kind)
	if (err != nil) && (err.Error() != fmt.Sprintf(ModuleFilesMissingF, wd)) {
		return false, ""
	}
	for _, m := range mods {
		if _, found := m.items[item]; found {
			return true, m.name
		}
	}
	return false, ""
}
