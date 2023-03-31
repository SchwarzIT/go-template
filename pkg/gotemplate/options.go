package gotemplate

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path"
	"regexp"
	"strings"

	"github.com/schwarzit/go-template/pkg/repos"
)

type ErrOutOfRange struct {
	Value int
	Min   int
	Max   int
}

func (e *ErrOutOfRange) Error() string {
	return fmt.Sprintf("%d: value out of range (min: %d, max: %d)", e.Value, e.Min, e.Max)
}

// ErrInvalidPattern indicates that an error occurred while matching
// a value with a pattern.
// The pattern as well as a description for the pattern is included in the error message.
type ErrInvalidPattern struct {
	Value       string
	Pattern     string
	Description string
}

func (e *ErrInvalidPattern) Error() string {
	return fmt.Sprintf("%s: invalid pattern (expected %s (pattern: %s))", e.Value, e.Description, e.Pattern)
}

// Validator is a single method interface that validates that a given value is valid.
// If any error happens during validation or if the value is not valid an error will be returned.
type Validator interface {
	Validate(value interface{}) error
}

// ValidatorFunc is a function implementing the Validator interface.
type ValidatorFunc func(value interface{}) error

func (f ValidatorFunc) Validate(value interface{}) error {
	return f(value)
}

// Option is a struct containing all needed configuration for options to customize the template.
type Option struct {
	// name is the name of the option that will be used to reference it and also that will be shown on the cli.
	name string
	// description is the description of the option that should be shown.
	// It's a StringValuer since it could depend on some earlier input.
	description string
	// defaultValue is the default value of the option.
	// It's a Valuer since it could be depend on earlier inputs or some http call.
	defaultValue Valuer
	// validator is used to validate an input value if it can be used as the value for this option.
	// If it is not set it will by default by valid.
	validator Validator
	// shouldDisplay decides whether the option is shown when the values are loaded interactively.
	// In most cases this is used to ensure options are only shown if needed values have been supplied earlier.
	// If it is not set it will by default be shown.
	shouldDisplay BoolValuer
	// postHook is some function that will be executed after all options are loaded.
	// This can for example be used to remove files from the created project folder or initialize tools based on inputs.
	// The passed interface contains the value of the option for convenience (technically also contained in optionValues)
	// targetDir indicates the working directory of the postHook
	postHook PostHookFunc
}

type PostHookFunc func(value interface{}, optionValues *OptionValues, targetDir string) error

func NewOption(name, description string, defaultValue Valuer, opts ...NewOptionOption) Option {
	option := Option{
		name:         name,
		description:  description,
		defaultValue: defaultValue,
	}

	for _, opt := range opts {
		opt(&option)
	}

	return option
}

type NewOptionOption func(*Option)

func WithValidator(validator Validator) NewOptionOption {
	return func(o *Option) {
		o.validator = validator
	}
}

func WithShouldDisplay(shouldDisplay BoolValuer) NewOptionOption {
	return func(o *Option) {
		o.shouldDisplay = shouldDisplay
	}
}

func WithPosthook(postHook PostHookFunc) NewOptionOption {
	return func(o *Option) {
		o.postHook = postHook
	}
}

func (s *Option) Name() string {
	return s.name
}

func (s *Option) Description() string {
	return s.description
}

// Default either returns the default value (possibly calculated with currentValues).
func (s *Option) Default(currentValues *OptionValues) interface{} {
	return s.defaultValue.Value(currentValues)
}

// ShouldDisplay returns a bool value indicating whether the option should be shown or not.
// If shouldDisplay variable is not set on the option true is returned.
func (s *Option) ShouldDisplay(currentValues *OptionValues) bool {
	if s.shouldDisplay != nil {
		return s.shouldDisplay.Value(currentValues)
	}

	return true
}

// Validate validates the value if a validator is specified.
func (s *Option) Validate(value interface{}) error {
	if s.validator != nil {
		return s.validator.Validate(value)
	}

	return nil
}

// PostHook executes the registered postHook if there is any.
func (s *Option) PostHook(v interface{}, optionValues *OptionValues, targetDir string) error {
	if s.postHook != nil {
		return s.postHook(v, optionValues, targetDir)
	}

	return nil
}

// Category is used to wrap multiple extensions into one organizational unit.
// This is to reduce the amount of required user input if certain categories if extensions
// can be skipped as a category instead of needing to skip all one by one.
type Category struct {
	Name    string
	Options []Option
}

// Options is the main struct wrapping the configuration
// for all allowed parameters and extensions.
// Slices are used instead of maps since the iteration order of maps is undefined/random.
// which could lead to confusion.
type Options struct {
	Base       []Option
	Extensions []Category
}

// OptionValues is a struct mirroring the structure of Options but using maps.
// Instead of the whole option only the set value of the option is kept.
// This makes looking up already supplied option values easier than it would
// be in the Options struct.
type OptionValues struct {
	Base       OptionNameToValue            `yaml:"base"`
	Extensions map[string]OptionNameToValue `yaml:"extensions"`
}

func NewOptionValues() *OptionValues {
	return &OptionValues{
		Base:       OptionNameToValue{},
		Extensions: map[string]OptionNameToValue{},
	}
}

type OptionNameToValue map[string]interface{}

// NewOptions returns all of go/template's options.
// Keeping repos.GithubTagLister in case it's needed in the future
func NewOptions(_ repos.GithubTagLister) *Options { //nolint:funlen,cyclop // Static initialization
	return &Options{
		Base: []Option{
			{
				name:         "projectName",
				defaultValue: StaticValue("Awesome Project"),
				description:  "Name of the project",
			},
			{
				name: "projectSlug",
				defaultValue: DynamicValue(func(ov *OptionValues) interface{} {
					projectName := ov.Base["projectName"].(string)
					return strings.ReplaceAll(strings.ToLower(projectName), " ", "-")
				}),
				description: "Technical name of the project for folders and names. This will also be used as output directory.",
				validator:   RegexValidator(`^[a-z1-9]+(-[a-z1-9]+)*$`, "only lowercase letters, numbers and dashes"),
			},
			{
				name:         "projectDescription",
				defaultValue: StaticValue("The awesome project provides awesome features to awesome people."),
				description:  "Description of the project used in the README.",
			},
			{
				name: "appName",
				defaultValue: DynamicValue(func(ov *OptionValues) interface{} {
					return ov.Base["projectSlug"].(string)
				}),
				description: `The name of the binary that you want to create.
Could be the same as your "projectSlug" but since Go supports multiple apps in one repo it could also be sth. else.
For example if your project is for some API there could be one app for the server and one CLI client.`,
				validator: RegexValidator(`^[a-z1-9]+(-[a-z1-9]+)*$`, "only lowercase letters, numbers and dashes"),
			},
			{
				name: "moduleName",
				defaultValue: DynamicValue(func(vals *OptionValues) interface{} {
					projectSlug := vals.Base["projectSlug"].(string)
					return fmt.Sprintf("github.com/user/%s", projectSlug)
				}),
				description: `The name of the Go module defined in the "go.mod" file.
This is used if you want to "go get" the module.
Please be aware that this depends on your version control system.
The default points to "github.com" but for devops for example it would look sth. like this "dev.azure.com/org/project/repo.git"`,
				validator: RegexValidator(`^[\S]+$`, "no whitespaces"),
			},
		},
		Extensions: []Category{
			{
				Name: "openSource",
				Options: []Option{
					{
						name:         "license",
						defaultValue: StaticValue(1),
						description: `Set an OpenSource license.
Unsure which to pick? Checkout Github's https://choosealicense.com/
Options:
	0: Add no license
	1: MIT License
	2: Apache License 2.0
	3: GNU AGPLv3
	4: GNU GPLv3
	5: GNU LGPLv3
	6: Mozilla Public License 2.0
	7: Boost Software License 1.0
	8: The Unlicense`,
						postHook: func(v interface{}, _ *OptionValues, targetDir string) error {
							if v.(int) == 0 {
								files := []string{"LICENSE", "CODEOWNERS"}
								for _, file := range files {
									if err := os.RemoveAll(path.Join(targetDir, file)); err != nil {
										return err
									}
								}
							}
							return nil
						},
					},
					{
						name: "author",
						defaultValue: DynamicValue(func(vals *OptionValues) interface{} {
							buffer := &bytes.Buffer{}
							gitName := exec.Command("git", "config", "--get", "user.name")
							gitName.Stdout = buffer
							if err := gitName.Run(); err != nil || len(buffer.Bytes()) == 0 {
								return "Marty Mc Fly"
							}
							return strings.TrimSpace(buffer.String())
						}),
						description: `License author`,
						shouldDisplay: DynamicBoolValue(func(vals *OptionValues) bool {
							switch vals.Extensions["openSource"]["license"].(int) {
							case 1, 2: //nolint:gomnd // 1: MIT License; 2: Apache License 2.0
								return true
							}
							return false
						}),
					},
					{
						name: "codeowner",
						defaultValue: DynamicValue(func(vals *OptionValues) interface{} {
							buffer := &bytes.Buffer{}
							gitMail := exec.Command("git", "config", "--get", "user.email")
							gitMail.Stdout = buffer
							if err := gitMail.Run(); err != nil || len(buffer.Bytes()) == 0 {
								return "Marty.Mc.Fly@future.back"
							}
							return strings.TrimSpace(buffer.String())
						}),
						description: "Set the codeowner of the project",
						shouldDisplay: DynamicBoolValue(func(vals *OptionValues) bool {
							return vals.Extensions["openSource"]["license"].(int) != 0 // 0 == no license
						}),
					},
				},
			},
			{
				Name: "ci",
				Options: []Option{
					{
						name:         "provider",
						defaultValue: StaticValue(1),
						description: `Set an CI pipeline provider integration
			Options:
			0: No CI
			1: Github
			2: Gitlab
			3: Azure DevOps`,
						postHook: func(v interface{}, _ *OptionValues, targetDir string) error {
							ciFiles := map[int][]string{
								0: {},
								1: {".github"},
								2: {".gitlab-ci.yml"},
								3: {".azure-pipelines.yml"},
							}

							for i, files := range ciFiles {
								if i == v.(int) {
									continue
								}
								for _, file := range files {
									if err := os.RemoveAll(path.Join(targetDir, file)); err != nil {
										return err
									}
								}
							}
							return nil
						},
					},
				},
			},
			{
				Name: "grpc",
				Options: []Option{
					{
						name:         "base",
						defaultValue: StaticValue(false),
						description:  "Base configuration for gRPC",
						postHook: func(v interface{}, _ *OptionValues, targetDir string) error {
							set := v.(bool)
							files := []string{"api/proto", "buf.gen.yaml", "buf.work.yaml", "api/openapi.v1.yml"}

							if set {
								return os.RemoveAll(path.Join(targetDir, "api/openapi.v1.yml"))
							}
							return removeAllBut(targetDir, files, "api/openapi.v1.yml")
						},
					},
					{
						name:         "grpcGateway",
						defaultValue: StaticValue(false),
						description:  "Extend gRPC configuration with grpc-gateway",
						shouldDisplay: DynamicBoolValue(func(vals *OptionValues) bool {
							return vals.Extensions["grpc"]["base"].(bool)
						}),
					},
				},
			},
		},
	}
}

// removeAllBut removes all files in the toRemove slice except for the exception.
func removeAllBut(targetDir string, toRemove []string, exception string) error {
	for _, item := range toRemove {
		if item == exception {
			continue
		}

		if err := os.RemoveAll(path.Join(targetDir, item)); err != nil {
			return err
		}
	}

	return nil
}

// RangeValidator validates that value is in between or equal to min and max.
func RangeValidator(min, max int) ValidatorFunc {
	return func(value interface{}) error {
		val := value.(int)

		if val < min || val > max {
			return &ErrOutOfRange{
				Value: val,
				Min:   min,
				Max:   max,
			}
		}

		return nil
	}
}

// RegexValidator returns a ValidatorFunc to validate a given value against a regex pattern.
// If the pattern doesn't match a ErrInvalidPattern is returned with a description on what the pattern means.
func RegexValidator(pattern, description string) ValidatorFunc {
	return func(value interface{}) error {
		str := value.(string)

		matched, err := regexp.MatchString(pattern, str)
		if err != nil {
			return err
		}

		if !matched {
			return &ErrInvalidPattern{Value: str, Pattern: pattern, Description: description}
		}

		return nil
	}
}
