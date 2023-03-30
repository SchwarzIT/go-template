package module

import (
	"fmt"

	"github.com/schwarzit/go-template/v3/gotemplate/option"
)

type BaseModule struct {
	name    string
	options []option.Option
}

func NewBaseModule(name string, options []option.Option) *BaseModule {
	return &BaseModule{
		name:    name,
		options: options,
	}
}

func (m *BaseModule) GetName() string {
	return m.name
}

func (m *BaseModule) GetOptions() []option.Option {
	return m.options
}

func (m *BaseModule) GetOptionValue(title string) (string, error) {
	for _, opt := range m.options {
		if opt.GetTitle() == title {
			return opt.GetValue(), nil
		}
	}
	return "", fmt.Errorf("option not found: %s", title)
}

func (m *BaseModule) Generate() error {
	// Implementation for generating the module's content goes here.
	return nil
}

func (m *BaseModule) GetState() State {
	state := make(State)
	for _, opt := range m.options {
		state[opt.GetTitle()] = opt.GetValue()
	}
	return state
}
