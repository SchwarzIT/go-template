package gotemplate

// QuestionType is an enum that defines the types of questions that can be presented to the user.
type QuestionType int

const (
	// Choices is a question with pre-defined choices.
	Choices QuestionType = iota
	// SingleChoice is a question with a single choice.
	SingleChoice
	// String is a question with a string input.
	String
	// MultiLineText is a question with a multi-line text input.
	MultiLineText
)

// ResponseValue is the interface that describes the possible values that can be set as a response to a TemplateQuestion.
type ResponseValue interface{}

// TemplateOptionName represents the name of a template option.
type TemplateOptionName string

// TemplateQuestion represents a single question in a module.
type TemplateQuestion interface {
	// Name is the name of the question.
	Name() TemplateOptionName

	// Description is a short description of the question.
	Description() string

	// DefaultValue is the default value for the question. This is optional.
	DefaultValue() *string

	// Choices is a list of pre-defined choices for the question. This is optional.
	Choices() []string

	// IsValid is a function that validates the input value for the question. This is optional.
	// If the input is valid, the function should return (true, "").
	// If the input is invalid, the function should return (false, "reason why input is invalid") or (false, "") if no reason is provided.
	IsValid(answer []string) (isValid bool, reason *string, err error)

	// IsEnabled returns whether or not the question should be enabled based on the current module data.
	IsEnabled(moduleData map[ModuleName]ModuleData) (bool, error)

	// ResponseValue is the value that the user provided in response to the question.
	// This is optional and will be set after the question has been answered.
	ResponseValue(answer []string) (interface{}, error)

	// Type is the type of question.
	Type() QuestionType
}

// Template represents a template that can be generated.
type Template struct {
	// InModules is a map of modules that can be included in the template.
	// The keys in the map are the names of the modules.
	// The values are the module data, including the template files and template data.
	InModules map[ModuleName]Module

	// OutModuleData is a map of module data that is output by the template generation process.
	// The keys in the map are the names of the modules.
	// The values are the module data that was generated during the template generation process.
	OutModuleData map[ModuleName]ModuleData

	// View is the user interface that will be used to present questions to the user.
	View View
}

// NewTemplate creates a new Template with the given View.
func NewTemplate(view View) *Template {
	return &Template{
		InModules:     make(map[ModuleName]Module),
		OutModuleData: make(map[ModuleName]ModuleData),
		View:          view,
	}
}

// AddModules adds modules to the template.
func (t *Template) AddModules(modules []Module) {
	for _, module := range modules {
		t.InModules[module.GetModuleName()] = module
	}
}

// ProcessModules processes all modules in the template and generates the module data.
func (t *Template) ProcessModules() error {
	// Initialize the module data for this module.
	t.OutModuleData = map[ModuleName]ModuleData{}

	// Iterate over each module and process the questions.
	for moduleName, module := range t.InModules {
		t.OutModuleData[moduleName] = ModuleData{
			Name:         moduleName,
			TemplateData: map[TemplateOptionName]interface{}{},
		}

		// // Get the questions for this module and iterate over them.
		// question := module.GetNextQuestion(t.OutModuleData)
		// for question != nil {
		// 	// Present the question to the user and get the response.
		// 	response, err := t.View.PresentQuestion(*question)
		// 	if err != nil {
		// 		return err
		// 	}

		// 	t.OutModuleData[moduleName].TemplateData[question.Name] = response.ResponseValue

		// 	// Get the next question.
		// 	question = module.GetNextQuestion(t.OutModuleData)
		// }

		// Store the module data in the template.
		t.OutModuleData[moduleName] = *module.GetModuleData()
	}

	return nil
}

// ExecuteWizard executes the template wizard using the current View.
func (t *Template) ExecuteWizard() error {
	// // Display a welcome message.
	// err := t.View.ShowMessage("Welcome to the template wizard!")
	// if err != nil {
	// 	return err
	// }

	// Process modules in the template and generate the module data.
	err := t.ProcessModules()
	if err != nil {
		return err
	}

	// // Display a message indicating the end of the wizard.
	// err = t.View.ShowMessage("The template wizard is complete!")
	// if err != nil {
	// 	return err
	// }

	return nil
}
