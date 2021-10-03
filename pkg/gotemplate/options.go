package gotemplate

import (
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/schwarzit/go-template/pkg/repos"
)

// ErrInvalidaPattern indicates that an error occured while matching
// a value with a pattern.
// The pattern is included in the error message.

type ErrInvalidPattern struct {
	Value   string
	Pattern string
}

func (e *ErrInvalidPattern) Error() string {
	return fmt.Sprint("%s: invalid pattern (expected %s)", e.Value, e.Pattern)
}

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

func NewOptions(githubTagLister repos.GithubTagLister) *Options {
	return &Options{
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
				description: StringValue("Technical name of the project for folders and names. This will also be used as output directory."),
				validator:   RegexValidator(`^[a-z1-9]+(-[a-z1-9]+)*$`, "only lowercase letters and dashes"),
			},
			&SomeOption{
				name:         "projectDescription",
				defaultValue: StaticValue("The awesome project provides awesome features to awesome people."),
				description:  StringValue("Description of the project used in the README."),
			},
			&SomeOption{
				name:         "appName",
				defaultValue: StaticValue("awesomecli"),
				description:  StringValue("The name of the binary that you want to create. Could be the same your `project_slug` but since Go supports multiple apps in one repo it could also be sth. else. For example if your project is for some API there could be one app for the server and one CLI client."),
				validator:    RegexValidator(`^[a-z]+$`, "only lowercase letters"),
			},
			&SomeOption{
				name: "moduleName",
				defaultValue: DynamicValue(func(vals OptionValues) interface{} {
					projectSlug := vals.Base["projectSlug"].(string)
					return fmt.Sprintf("github.com/user/%s", projectSlug)
				}),
				description: StringValue("The name of the Go module defined in the `go.mod` file. This is used if you want to `go get` the module. Please be aware that this depends on your version control system. The default points to `github.com` but for devops for example it would look sth. like this `dev.azure.com/org/project/repo.git`"),
				validator:   RegexValidator(`^[\S]+$`, "no whitespaces"),
			},
			&SomeOption{
				name: "golangciVersion",
				defaultValue: DynamicValue(func(_ OptionValues) interface{} {
					latestTag, err := repos.LatestGithubReleaseTag(githubTagLister, "golangci", "golangci-lint")
					if err != nil {
						return "1.42.1"
					}

					return latestTag.String()
				}),
				description: StringValue("Golangci-lint version to use."),
				validator: RegexValidator(
					`^(?P<major>0|[1-9]\d*)\.(?P<minor>0|[1-9]\d*)\.(?P<patch>0|[1-9]\d*)(?:-(?P<prerelease>(?:0|[1-9]\d*|\d*[a-zA-Z-][0-9a-zA-Z-]*)(?:\.(?:0|[1-9]\d*|\d*[a-zA-Z-][0-9a-zA-Z-]*))*))?(?:\+(?P<buildmetadata>[0-9a-zA-Z-]+(?:\.[0-9a-zA-Z-]+)*))?$`,
					"valid semver version string",
				),
			},
		},
		Extensions: []Category{
			{
				Name: "cicd",
				Options: []Option{
					&SomeOption{
						name:         "pipeline",
						defaultValue: StaticValue(1),
						validator:    RegexValidator(`^[1-2]$`, "number between 1-2"),
						description: StringValue(`Set a pipelining system.
	Options:
		1. Github
		2. Azure Devops`),
						hook: func(v interface{}) error {
							val := v.(int)
							dirs := []string{".github", ".azuredevops"}
							switch val {
							case 1:
								return removeAllBut(dirs, ".github")
							case 2:
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

func RegexValidator(pattern, description string) ValidatorFunc {
	return func(value interface{}) error {
		str := value.(string)

		matched, err := regexp.MatchString(string(pattern), str)
		if err != nil {
			return err
		}

		if !matched {
			return &ErrInvalidPattern{Value: str, Pattern: pattern}
		}

		return nil
	}
}
