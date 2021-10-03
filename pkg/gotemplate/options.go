package gotemplate

// TODO: rebase!

import (
	"errors"
	"os"
	"regexp"
	"strings"
)

// TODO: remove interface?
type Option interface {
	Name() string
	Description(currentValues OptionValues) string
	Default(currentValues OptionValues) interface{}
	ShouldDisplay(currentValues OptionValues) bool
	// validate receives a value that should be set for the option and validates it
	Validate(value interface{}) error
	// TODO: should the value be kept in the struct and reused here instead of passing it?
	PostHook(value interface{}) error
}

var _ Option = &SomeOption{}

type Validator interface {
	Validate(value interface{}) error
}

type ValidatorFunc func(value interface{}) error

func (f ValidatorFunc) Validate(value interface{}) error {
	return f(value)
}

type SomeOption struct {
	name          string
	description   StringValuer
	hook          func(interface{}) error
	defaultValue  Valuer
	validator     Validator
	shouldDisplay BoolValuer
}

func (s *SomeOption) Name() string {
	return s.name
}

func (s *SomeOption) Description(currentValues OptionValues) string {
	return s.description.Value(currentValues)
}

// Value either returns the default value (possibly calculated with currentValues)
// or the actually set value
func (s *SomeOption) Default(currentValues OptionValues) interface{} {
	return s.defaultValue.Value(currentValues)
}

func (s *SomeOption) ShouldDisplay(currentValues OptionValues) bool {
	if s.shouldDisplay != nil {
		return s.shouldDisplay.Value(currentValues)
	}

	return true
}

func (s *SomeOption) Validate(value interface{}) error {
	if s.validator != nil {
		return s.validator.Validate(value)
	}

	return nil
}

func (s *SomeOption) PostHook(v interface{}) error {
	if s.hook != nil {
		return s.hook(v)
	}

	return nil
}

type Regex struct {
	Pattern     regexp.Regexp
	Description string
}

type Category struct {
	Name    string
	Options []Option
}

// would be easier to make the Options hardcoded or maps to access them later one typesafe
// but there would be no option to read in the things dynamically without reflection
// -> either dynamic during reading the values
// -> or while evaluating them and dependencies...
// can't be maps since the order of looping through maps is random
type Options struct {
	Base       []Option
	Extensions []Category
}

type OptionValues struct {
	Base       OptionNameToValue
	Extensions map[string]OptionNameToValue
}

type OptionNameToValue map[string]interface{}

var ErrTypeMismatch = errors.New("type mismatch")
var ErrInvalidPattern = errors.New("invalid pattern")

var options = Options{
	Base: []Option{
		&SomeOption{
			name:         "projectName",
			defaultValue: StaticValue("Awesome Project"),
			description:  StringValue("Name of the project"),
		},
		&SomeOption{
			name: "projectSlug",
			defaultValue: DynamicValue(func(ov OptionValues) interface{} {
				projectName := ov.Base["projectName"].(string)
				return strings.ReplaceAll(strings.ToLower(projectName), " ", "-")
			}),
			validator: RegexValidator(`^[a-z1-9]+(-[a-z1-9]+)*$`),
		},
	},
	Extensions: []Category{
		{
			Name: "cicd",
			Options: []Option{
				&SomeOption{
					name:         "pipeline",
					defaultValue: StaticValue("1"),
					validator:    RegexValidator(`^[1-2]$`),
					description: StringValue(`Set a pipelining system.
Options:
    1. Github
	2. Azure Devops`),
					hook: func(v interface{}) error {
						val := v.(string)
						dirs := []string{".github", ".azuredevops"}
						switch val {
						case "1":
							return removeAllBut(dirs, ".github")
						case "2":
							return removeAllBut(dirs, ".azuredevops")
						}
						return nil
					},
				},
			},
		},
		{
			Name: "grpc",
			Options: []Option{
				&SomeOption{
					name:         "base",
					defaultValue: StaticValue(false),
					hook: func(v interface{}) error {
						set := v.(bool)
						files := []string{"api/proto", "tools.go", "buf.gen.yaml", "buf.yaml", "api/openapi.v1.yaml"}

						if set {
							return os.RemoveAll("api/openapi.v1.yaml")
						}
						return removeAllBut(files, "api/openapi.v1.yaml")
					},
				},
				&SomeOption{
					name:         "grpcGateway",
					defaultValue: StaticValue(false),
					shouldDisplay: DynamicBoolValue(func(vals OptionValues) bool {
						return vals.Extensions["grpc"]["base"].(bool)
					}),
				},
			},
		},
	},
}

// removeAllBut removes all files in the toRemove slice except for the exception
func removeAllBut(toRemove []string, exception string) error {
	for _, item := range toRemove {
		if item == exception {
			continue
		}

		if err := os.RemoveAll(item); err != nil {
			return err
		}
	}

	return nil
}

type RegexValidator string

func (v RegexValidator) Validate(i interface{}) error {
	str, ok := i.(string)
	if !ok {
		return ErrTypeMismatch
	}

	matched, err := regexp.MatchString(string(v), str)
	if err != nil {
		return err
	}

	if !matched {
		return ErrInvalidPattern
	}

	return nil
}
