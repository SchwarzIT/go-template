package option

import (
	"errors"
	"fmt"
)

var (
	ErrInvalidAnswer = errors.New("invalid answer")
	ErrInvalidValue  = errors.New("invalid value")
	ErrMissingTitle  = errors.New("missing title")
	ErrMissingDesc   = errors.New("missing description")
)

type BaseOption struct {
	Title             string
	Description       string
	DefaultValue      []string
	AvailableAnswers  []string
	CurrentValue      []string
	ValidationFunc    func([]string) error
	ShouldDisplayFunc func(map[ModuleName]State) (bool, error)
}

type NewBaseOptionArgs struct {
	Title             string
	Description       string
	DefaultValue      []string
	AvailableAnswers  []string
	ValidationFunc    func([]string) error
	ShouldDisplayFunc func(map[ModuleName]State) (bool, error)
}

func NewBaseOption(args NewBaseOptionArgs) *BaseOption {
	return &BaseOption{
		Title:             args.Title,
		Description:       args.Description,
		DefaultValue:      args.DefaultValue,
		AvailableAnswers:  args.AvailableAnswers,
		CurrentValue:      args.DefaultValue,
		ValidationFunc:    args.ValidationFunc,
		ShouldDisplayFunc: args.ShouldDisplayFunc,
	}
}

func (o *BaseOption) GetTitle() (string, error) {
	if o.Title == "" {
		return "", ErrMissingTitle
	}
	return o.Title, nil
}

func (o *BaseOption) GetDescription() (string, error) {
	if o.Description == "" {
		return "", ErrMissingDesc
	}
	return o.Description, nil
}

func (o *BaseOption) GetDefaultValue(state map[ModuleName]State) ([]string, error) {
	return o.DefaultValue, nil
}

func (o *BaseOption) GetAvailableAnswers() ([]string, error) {
	if o.AvailableAnswers == nil {
		return nil, nil
	}
	return o.AvailableAnswers, nil
}

func (o *BaseOption) GetCurrentValue() ([]string, error) {
	if o.CurrentValue == nil {
		return nil, ErrInvalidValue
	}
	return o.CurrentValue, nil
}

func (o *BaseOption) SetCurrentValue(value []string) error {
	if err := o.Validate(value); err != nil {
		return err
	}
	o.CurrentValue = value
	return nil
}

func (o *BaseOption) Validate(value []string) error {
	if o.ValidationFunc == nil {
		return nil
	}

	if err := o.ValidationFunc(value); err != nil {
		return fmt.Errorf("invalid value: %w", err)
	}
	return nil
}

func (o *BaseOption) ShouldDisplay(state map[ModuleName]State) (bool, error) {
	if o.ShouldDisplayFunc == nil {
		return true, nil
	}
	return o.ShouldDisplayFunc(state)
}
