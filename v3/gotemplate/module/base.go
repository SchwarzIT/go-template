package module

import (
	"errors"
	"fmt"

	"github.com/schwarzit/go-template/v3/gotemplate/option"
)

var (
	ErrMissingName         = errors.New("missing module name")
	ErrMissingOptions      = errors.New("missing module options")
	ErrUnimplementedMethod = errors.New("method is not implemented")
)

// TemplateFile represents a file that will be generated from a template.
type BaseTemplateFile struct {
	// SourcePath is the path to the template file.
	SourcePath string

	// TargetPath is the path where the generated file will be created.
	TargetPath string
}

// Source returns the path to the template file.
func (f *BaseTemplateFile) Source() string {
	return f.SourcePath
}

// Target returns the path where the generated file will be created.
func (f *BaseTemplateFile) Target() string {
	return f.TargetPath
}

// BaseModule is a struct that implements the Module interface.
type BaseModule struct {
	name    string
	options []option.Option
}

// NewBaseModule returns a new instance of the BaseModule struct.
func NewBaseModule(name string, options []option.Option) (*BaseModule, error) {
	if name == "" {
		return nil, ErrMissingName
	}

	if options == nil || len(options) == 0 {
		return nil, ErrMissingOptions
	}

	return &BaseModule{
		name:    name,
		options: options,
	}, nil
}

// GetName returns the name of the module.
func (m *BaseModule) GetName() (option.ModuleName, error) {
	if m.name == "" {
		return "", ErrMissingName
	}

	return option.ModuleName(m.name), nil
}

// GetOptions returns the options that can be configured for the module.
func (m *BaseModule) GetOptions() ([]option.Option, error) {
	if m.options == nil || len(m.options) == 0 {
		return nil, ErrMissingOptions
	}

	return m.options, nil
}

// Generate generates the project files for the module using the current
// option values and the provided templatePath.
func (m *BaseModule) Generate(files []TemplateFile, options []option.Option) error {
	for _, file := range files {
		fmt.Println(fmt.Sprintf("Generating file %s to %s", file.Source(), file.Target()))
	}

	// print the options
	for _, opt := range options {
		name, err := opt.GetTitle()
		if err != nil {
			return err
		}
		value, err := opt.GetCurrentValue()
		if err != nil {
			return err
		}
		fmt.Println("Option Name:", name)
		fmt.Println("Option Value:", value)
	}

	return nil
}

// GetTemplateFiles returns the template files for the module.
func (m *BaseModule) GetTemplateFiles() ([]TemplateFile, error) {
	return nil, ErrUnimplementedMethod
}
