package gotemplate

// Module represents a module that can be included in a template.
type Module interface {
	// GetNextQuestion returns the next question in the module.
	// If there are no more questions, it returns nil.
	GetNextQuestion() *TemplateQuestion
	// GetModuleData returns the complete module data when all questions have been answered.
	GetModuleData() *ModuleData
}

// ModuleName represents the name of a module that can be included in a template.
type ModuleName string

// ModuleData represents a module that can be included in a template.
type ModuleData struct {
	// Name is the name of the module.
	Name ModuleName
	// TemplatePath is the path to the template files for this module.
	TemplatePath string
	// TemplateFiles is a list of file paths relative to TemplatePath.
	// These files will be generated when the template is generated.
	TemplateFiles []string
	// TemplateData is the data that will be used to generate the template for this module.
	// The data can be of any type, depending on the needs of the module.
	TemplateData interface{}
}

// ModuleBase implements the Module interface.
type ModuleBase struct {
	ModuleData
	questions []TemplateQuestion
}

// GetNextQuestion returns the next question in the module.
// If there are no more questions, it returns nil.
func (m *ModuleBase) GetNextQuestion() *TemplateQuestion {
	if len(m.questions) == 0 {
		return nil
	}
	question := m.questions[0]
	m.questions = m.questions[1:]
	return &question
}

// GetModule returns the complete module data when all questions have been answered.
func (m *ModuleBase) GetModule() *ModuleData {
	if len(m.questions) == 0 {
		return &m.ModuleData
	}
	return nil
}
