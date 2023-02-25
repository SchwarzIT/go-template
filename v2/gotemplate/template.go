package gotemplate

// TemplateOptionName represents the name of a template option.
type TemplateOptionName string

// TemplateOptionResponse represents the user's response to a template option.
type TemplateOptionResponse struct {
	// Name is the name of the option.
	Name TemplateOptionName
	// Value is the user's selected value for the option.
	Value interface{}
}

// TemplateQuestion represents a single question in a module.
type TemplateQuestion struct {
	// Name is the name of the question.
	Name TemplateOptionName
	// Description is a short description of the question.
	Description string
	// DefaultValue is the default value for the question. This is optional.
	DefaultValue interface{}
	// IsValid is a function that validates the input value for the question. This is optional.
	// If the input is valid, the function should return (true, "").
	// If the input is invalid, the function should return (false, "reason why input is invalid") or (false, "") if no reason is provided.
	IsValid func(value interface{}) (isValid bool, reason string)
	// PredefinedValues is a list of predefined values for the question. This is optional.
	// If this is set, the user will be presented with a list of predefined values to choose from instead of an input field.
	PredefinedValues []interface{}
}

// Template represents a template that can be generated.
type Template struct {
	// Modules is a map of modules that can be included in the template.
	// The keys in the map are the names of the modules.
	// The values are the module data, including the template files and template data.
	Modules map[ModuleName]Module
}
