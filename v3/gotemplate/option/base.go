package option

import "fmt"

type BaseOption struct {
	title            string
	description      string
	defaultValue     string
	availableAnswers []string
	value            string
	validationFunc   func(string) error
}

func NewBaseOption(title, description, defaultValue string, availableAnswers []string, validationFunc func(string) error) *BaseOption {
	return &BaseOption{
		title:            title,
		description:      description,
		defaultValue:     defaultValue,
		availableAnswers: availableAnswers,
		value:            defaultValue,
		validationFunc:   validationFunc,
	}
}

func (o *BaseOption) GetTitle() string {
	return o.title
}

func (o *BaseOption) GetDescription() string {
	return o.description
}

func (o *BaseOption) GetDefaultValue() string {
	return o.defaultValue
}

func (o *BaseOption) GetAvailableAnswers() []string {
	return o.availableAnswers
}

func (o *BaseOption) GetValue() string {
	return o.value
}

func (o *BaseOption) SetValue(value string) error {
	if err := o.Validate(value); err != nil {
		return err
	}
	o.value = value
	return nil
}

func (o *BaseOption) Validate(value string) error {
	if o.availableAnswers != nil {
		found := false
		for _, answer := range o.availableAnswers {
			if answer == value {
				found = true
				break
			}
		}
		if !found {
			return fmt.Errorf("invalid answer: %s", value)
		}
	}
	if o.validationFunc != nil {
		if err := o.validationFunc(value); err != nil {
			return err
		}
	}
	return nil
}
