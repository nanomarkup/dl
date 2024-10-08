package dl // import "github.com/nanomarkup/dl"
Package smodule manages modules.
CONSTANTS
const (
	// application
	AppCode          string = "dl"
	ItemSeparator    string = " "
	ItemOptCode      string = ":"
	InitBegOptCode   string = "{"
	InitEndOptCode   string = "}"
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
	ModuleFilesMissingF   string = "no .%s files in \"%s\""
	ModuleKindIsMissing   string = "kind of modules to load is not specified"
	ModuleErrorOnLoadingF string = "cannot load \"%s\" module/s"
	FirstTokenInvalid     string = "incorrect type of file"
	FirstTokenIsMissing   string = "type of file is missing"
	LineSyntaxInvalidF    string = "invalid syntax in \"%s\" line"
)
TYPES
type Formatter struct {
}
func (f *Formatter) Item(name string, deps [][]string) string
func (f *Formatter) String(module Module) string
type Item = [][]string
type Items = map[string]Item
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
	Kind   string
	Logger Logger
}
func (m *Manager) AddDependency(item, dependency, resolver string, update bool) error
func (m *Manager) AddItem(module, item string) error
func (m *Manager) DeleteDependency(item, dependency string) error
func (m *Manager) DeleteItem(item string) error
func (m *Manager) Read(filePath string) (Module, error)
func (m *Manager) ReadAll() (Module, error)
func (m *Manager) SetLogger(logger Logger)
type Module interface {
	Items() map[string][][]string
	Dependency(string, string) string
}
