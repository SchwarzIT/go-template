package module

import (
	"errors"
	"html/template"
	"io/fs"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/schwarzit/go-template/v3/gotemplate/option"
)

const (
	// DefaultFilePermissionRWX is the default file permission for a generated
	FilePermissionRWX = 0755
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

	// File Permission is the file permission for the generated file.
	Permission uint
}

// Source returns the path to the template file.
func (f *BaseTemplateFile) Source() string {
	return f.SourcePath
}

// Target returns the path where the generated file will be created.
func (f *BaseTemplateFile) Target() string {
	return f.TargetPath
}

// FilePermission returns the file permission for the generated file.
func (f *BaseTemplateFile) FilePermission() uint {
	if f.Permission == 0 {
		return FilePermissionRWX
	}

	return f.Permission
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
func (m *BaseModule) Generate(files []TemplateFile, options []option.Option, outPath string) error {
	tmpl := template.New("ModuleTemplates")

	for _, file := range files {
		content, err := ioutil.ReadFile(file.Source())
		if err != nil {
			return err
		}

		tmpl, err = tmpl.New(file.Source()).Parse(string(content))
		if err != nil {
			return err
		}
	}

	optionMap := make(map[string]interface{})
	err := updateOptionMap(optionMap, options)
	if err != nil {
		return err
	}

	for _, file := range files {
		targetPath := filepath.Join(outPath, file.Target())
		if err := os.MkdirAll(filepath.Dir(targetPath), fs.FileMode(file.FilePermission())); err != nil {
			return err
		}

		outFile, err := os.Create(targetPath)
		if err != nil {
			return err
		}

		if err := tmpl.ExecuteTemplate(outFile, file.Source(), optionMap); err != nil {
			outFile.Close()
			return err
		}

		if err := outFile.Close(); err != nil {
			return err
		}
	}

	return nil
}

// GetTemplateFiles returns the template files for the module.
func (m *BaseModule) GetTemplateFiles() ([]TemplateFile, error) {
	return nil, ErrUnimplementedMethod
}

func updateOptionMap(optionMap map[string]interface{}, options []option.Option) error {
	for _, opt := range options {
		templateKey, err := opt.GetTemplateKey()
		if err != nil {
			return err
		}

		value, err := opt.GetCurrentValue()
		if err != nil {
			return err
		}

		optionMap[templateKey] = value
	}

	return nil
}
