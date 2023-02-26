package gotemplate

// Module represents a module that can be included in a template.
type Module interface {
	// GetModuleName returns the name of the module.
	GetModuleName() ModuleName

	// GetNextQuestion returns the next question in the module.
	// The `moduleData` argument is a map of the module data that has been collected so far.
	// If there are no more questions, it returns nil.
	GetNextQuestion(moduleData map[ModuleName]ModuleData) *TemplateQuestion

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
	TemplateData map[TemplateOptionName]interface{}
}

// ModuleBase implements the Module interface.
type ModuleBase struct {
	ModuleData
	Questions []TemplateQuestion
}

// GetNextQuestion returns the next question in the module.
// If there are no more questions, it returns nil.
func (m *ModuleBase) GetNextQuestion(moduleData map[ModuleName]ModuleData) *TemplateQuestion {
	if len(m.Questions) == 0 {
		return nil
	}
	question := m.Questions[0]
	m.Questions = m.Questions[1:]
	return &question
}

// GetModuleData returns the complete module data when all questions have been answered.
func (m *ModuleBase) GetModuleData() *ModuleData {
	if len(m.Questions) == 0 {
		return &m.ModuleData
	}
	return nil
}

// GetModuleName returns the name of the module.
func (m *ModuleBase) GetModuleName() ModuleName {
	return m.ModuleData.Name
}
