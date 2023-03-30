package question

import "github.com/schwarzit/go-template/v2/gotemplate"

type Question struct {
	// Name is the name of the question.
	name gotemplate.TemplateOptionName

	// Description is a short description of the question.
	description string

	// DefaultValue is the default value for the question. This is optional.
	defaultValue *string

	// Choices is a list of pre-defined choices for the question. This is optional.
	choices []string

	// Type is the type of question.
	questionType gotemplate.QuestionType
}

func (q *Question) Name() gotemplate.TemplateOptionName {
	return q.name
}

func (q *Question) Description() string {
	return q.description
}

func (q *Question) DefaultValue() *string {
	return q.defaultValue
}

func (q *Question) Choices() []string {
	return q.choices
}

func (q *Question) IsValid(answer []string) (isValid bool, reason *string, err error) {
	if q.IsValid != nil {
		return q.IsValid(answer)
	}
	return true, nil, nil
}

func (q *Question) IsEnabled(moduleData map[gotemplate.ModuleName]gotemplate.ModuleData) (bool, error) {
	if q.IsEnabled != nil {
		return q.IsEnabled(moduleData)
	}
	return true, nil
}

func (q *Question) ResponseValue(answer []string) (interface{}, error) {
	if q.ResponseValue != nil {
		return q.ResponseValue(answer)
	}
	return answer[0], nil
}

func (q *Question) Type() gotemplate.QuestionType {
	return q.questionType
}
