package module

import "github.com/schwarzit/go-template/v3/gotemplate/option"

// TemplateFile represents a file to be generated by a module.
type TemplateFile interface {
	// Source path sets the source path of the template file.
	Source() string

	// Target path sets the target path of the generated file.
	Target() string

	// FilePermission returns the file permissions for the generated file.
	FilePermission() uint
}

// Module is an interface that describes a module that can be used to generate
// project files.
type Module interface {
	// GetName returns the name of the module.
	GetName() (option.ModuleName, error)

	// GetOptions returns the options that can be configured for the module.
	GetOptions() ([]option.Option, error)

	// Generate generates the project files for the module using the current
	// option values.
	Generate(files []TemplateFile, options []option.Option, outPath string) error

	// GetTemplateFiles returns the template files for the module.
	GetTemplateFiles() ([]TemplateFile, error)
}
