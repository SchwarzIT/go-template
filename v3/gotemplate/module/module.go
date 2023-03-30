package module

import "github.com/schwarzit/go-template/v3/gotemplate/option"

// State represents the state of a module, as a map of string keys to string values.
type State map[string]string

// TemplateModule represents a module in a template engine.
type Module interface {
	// GetName returns the name of the module.
	GetName() string

	// GetOptions returns a slice of options for the module.
	GetOptions() []option.Option

	// GetOptionValue returns the value of an option by title.
	GetOptionValue(title string) (string, error)

	// Generate generates the module's content.
	Generate() error

	// GetState returns the state of the module.
	GetState() State
}
