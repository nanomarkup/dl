// Copyright 2022 Vitalii Noha vitalii.noga@gmail.com. All rights reserved.

// Package smodule manages modules.
package lod

type Manager struct {
	Logger Logger
}

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

type Module interface {
	Kind() string
	Items() map[string]map[string]string
	Dependency(string, string) string
}

type Formatter struct {
}

const (
	// application
	ItemSeparator    string = " "
	ItemOptCode      string = ":"
	DefinesOptCode   string = "defines"
	DefineBegOptCode string = "{"
	DefineEndOptCode string = "}"
	CommentOptCode   string = "//"
	// notifications
	ModuleIsCreatedF string = "%s file has been created\n"
	// errors
	ItemExistsF           string = "the \"%s\" item already exists"
	ItemExistsInModuleF   string = "the \"%s\" item already exists in \"%s\" module"
	ItemIsMissingF        string = "the \"%s\" item does not exist"
	ItemNameInvalidF      string = "\"%s\" incorrect item name"
	DepItemExistsF        string = "\"%s\" already exists for \"%s\" item"
	DefineIsMissingF      string = "\"%s\" define is not declared"
	ModuleFilesMissingF   string = "no sb files in \"%s\""
	ModuleKindIsMissing   string = "kind of modules to load is not specified"
	ModuleKindMismatchF   string = "the \"%s\" kind of \"%s\" module is mismatch the \"%s\" selected kind"
	ModuleErrorOnLoadingF string = "cannot load \"%s\" module/s"
	FirstTokenInvalidF    string = "the first token should be \"%s\""
	LineSyntaxInvalidF    string = "invalid syntax in \"%s\" line"
)
