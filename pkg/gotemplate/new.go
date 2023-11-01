package gotemplate

import (
	"bytes"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path"
	"reflect"
	"strconv"
	"strings"
	"text/template"

	"github.com/Masterminds/semver/v3"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"

	gotemplate "github.com/schwarzit/go-template"
	ownexec "github.com/schwarzit/go-template/pkg/exec"
	"github.com/schwarzit/go-template/pkg/gocli"
)

const (
	minGoVersion  = "1.21"
	permissionRWX = 0755
	permissionRW  = 0644
)

var (
	ErrAlreadyExists         = errors.New("already exists")
	ErrParameterNotSet       = errors.New("parameter not set")
	ErrMalformedInput        = errors.New("malformed input")
	ErrParameterSet          = errors.New("parameter set but has no effect in this context")
	ErrGoVersionNotSupported = fmt.Errorf("go version is not supported, gt requires at least %s", minGoVersion)

	minGoVersionSemver = semver.MustParse(minGoVersion) //nolint:gochecknoglobals // parsed semver from const minGoVersion
)

type ErrTypeMismatch struct {
	Expected string
	Actual   string
}

func (e *ErrTypeMismatch) Error() string {
	return fmt.Sprintf("type mismatch, got %s, expected %s", e.Actual, e.Expected)
}

type NewRepositoryOptions struct {
	OutputDir    string
	OptionValues *OptionValues
}

// Validate validates all properties of NewRepositoryOptions except the ConfigValues, since those are validated by the Load functions.
func (opts NewRepositoryOptions) Validate() error {
	if opts.OutputDir == "" {
		return nil
	}

	if _, err := os.Stat(opts.OutputDir); err != nil {
		return err
	}

	return nil
}

// LoadConfigValuesFromFile loads value for the options from a file and validates the inputs
func (gt *GT) LoadConfigValuesFromFile(file string) (*OptionValues, error) { //nolint:cyclop // todo refactor
	fileBytes, err := os.ReadFile(file)
	if err != nil {
		return nil, err
	}

	var optionValues OptionValues

	if err := yaml.Unmarshal(fileBytes, &optionValues); err != nil {
		return nil, err
	}

	for _, option := range gt.Options.Base {
		val, ok := optionValues.Base[option.Name()]
		if !ok || reflect.ValueOf(val).IsZero() {
			return nil, errors.Wrap(ErrParameterNotSet, option.Name())
		}

		if err := validateFileOption(option, val, optionValues); err != nil {
			return nil, err
		}
	}

	for _, category := range gt.Options.Extensions {
		if optionValues.Extensions == nil {
			optionValues.Extensions = map[string]OptionNameToValue{}
		}
		for _, option := range category.Options {
			if optionValues.Extensions[category.Name] == nil {
				optionValues.Extensions[category.Name] = OptionNameToValue{}
			}
			val, ok := optionValues.Extensions[category.Name][option.Name()]
			if !ok {
				// set defaults for all unset optionValues, no need to validate
				optionValues.Extensions[category.Name][option.Name()] = option.Default(&optionValues)
				continue
			}

			if err := validateFileOption(option, val, optionValues); err != nil {
				return nil, err
			}
		}
	}

	return &optionValues, nil
}

func validateFileOption(option Option, value interface{}, optionValues OptionValues) error {
	valType := reflect.TypeOf(value)
	defaultVal := option.Default(&optionValues)
	defaultType := reflect.TypeOf(defaultVal)
	if valType != defaultType {
		return &ErrTypeMismatch{
			Expected: defaultType.Name(),
			Actual:   valType.Name(),
		}
	}

	if err := option.Validate(value); err != nil {
		return errors.Wrap(ErrMalformedInput, fmt.Sprintf("%s: %s", option.Name(), err.Error()))
	}

	// if it is set to sth else than default with shouldDisplay returning false it means the parameters does not have any effect
	if value != defaultVal && !option.ShouldDisplay(&optionValues) {
		return errors.Wrap(ErrParameterSet, option.Name())
	}

	return nil
}

func (gt *GT) LoadConfigValuesInteractively() (*OptionValues, error) {
	gt.printBanner()
	optionValues := NewOptionValues()

	for i := range gt.Options.Base {
		val := gt.loadOptionValueInteractively(&gt.Options.Base[i], optionValues)

		if val == nil {
			continue
		}

		optionValues.Base[gt.Options.Base[i].Name()] = val
	}

	gt.printProgressf("\nYou now have the option to enable additional extensions (organized in different categories)...\n\n")
	for _, category := range gt.Options.Extensions {
		gt.printCategory(category.Name)
		optionValues.Extensions[category.Name] = OptionNameToValue{}

		for i := range category.Options {
			val := gt.loadOptionValueInteractively(&category.Options[i], optionValues)

			if val == nil {
				continue
			}

			optionValues.Extensions[category.Name][category.Options[i].Name()] = val
		}
	}

	return optionValues, nil
}

func (gt *GT) loadOptionValueInteractively(option *Option, optionValues *OptionValues) interface{} {
	if !option.ShouldDisplay(optionValues) {
		return option.Default(optionValues)
	}

	val, err := gt.readOptionValue(option, optionValues)
	for err != nil {
		gt.printWarningf(err.Error())
		val, err = gt.readOptionValue(option, optionValues)
	}

	return val
}

func (gt *GT) InitNewProject(opts *NewRepositoryOptions) (err error) { //nolint:cyclop // todo refactor
	gt.printProgressf("Generating repo folder...")

	targetDir := path.Join(opts.OutputDir, opts.OptionValues.Base["projectSlug"].(string))
	gt.printProgressf("Writing to %s...", targetDir)

	if _, err := os.Stat(targetDir); !os.IsNotExist(err) {
		return errors.Wrapf(ErrAlreadyExists, "directory %s", targetDir)
	}

	defer func() {
		if err != nil {
			// ignore error to not overwrite original error
			_ = os.RemoveAll(targetDir)
		}
	}()
	err = fs.WalkDir(gotemplate.FS, gotemplate.Key, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		pathToWrite, err := gt.executeTemplateString(path, opts.OptionValues)
		if err != nil {
			return err
		}

		pathToWrite = strings.ReplaceAll(pathToWrite, gotemplate.Key, targetDir)
		if d.IsDir() {
			return os.MkdirAll(pathToWrite, permissionRWX)
		}

		fileBytes, err := fs.ReadFile(gotemplate.FS, path)
		if err != nil {
			return err
		}

		data, err := gt.executeTemplateString(string(fileBytes), opts.OptionValues)
		if err != nil {
			return err
		}

		filePermissions := fs.FileMode(permissionRW)
		// files that contain a shebang should be executable
		if strings.HasPrefix(strings.TrimSpace(data), "#!") {
			filePermissions = permissionRWX
		}

		return os.WriteFile(pathToWrite, []byte(data), filePermissions)
	})
	if err != nil {
		return err
	}

	gt.printProgressf("Removing obsolete files of unused integrations...")
	if err := postHook(gt.Options, opts.OptionValues, targetDir); err != nil {
		return err
	}

	gt.printProgressf("Initializing git and Go modules...")
	gt.initRepo(targetDir, opts.OptionValues.Base["moduleName"].(string))

	return nil
}

func (gt *GT) initRepo(targetDir, moduleName string) {
	commandGroups := []ownexec.CommandGroup{
		{
			Commands: []*exec.Cmd{
				exec.Command("git", "init"),
			},
			TargetDir: targetDir,
		},
		{
			PreRun: checkGoVersion,
			Commands: []*exec.Cmd{
				exec.Command("go", "mod", "init", moduleName),
				exec.Command("go", "mod", "tidy"),
			},
			TargetDir: targetDir,
		},
	}

	failedCGs := 0
	for _, cg := range commandGroups {
		if err := cg.Run(); err != nil {
			gt.printWarningf(err.Error())
			failedCGs++
		}
	}

	if failedCGs > 0 {
		gt.printWarningf("one or more initialization steps failed, pls see warnings for more info.")
	}
}

func checkGoVersion() error {
	goSemver, err := gocli.Semver()
	if err != nil {
		return err
	}

	if goSemver.LessThan(minGoVersionSemver) {
		return errors.Wrap(ErrGoVersionNotSupported, goSemver.String())
	}

	return nil
}

func postHook(options *Options, optionValues *OptionValues, targetDir string) error {
	for _, option := range options.Base {
		optionValue, ok := optionValues.Base[option.Name()]
		if !ok {
			return nil
		}

		if err := option.PostHook(optionValue, optionValues, targetDir); err != nil {
			return err
		}
	}

	for _, category := range options.Extensions {
		for _, option := range category.Options {
			optionValue, ok := optionValues.Extensions[category.Name][option.Name()]
			if !ok {
				return nil
			}

			if err := option.PostHook(optionValue, optionValues, targetDir); err != nil {
				return err
			}
		}
	}

	return nil
}

// readOptionValue reads a value for an option from the cli.
func (gt *GT) readOptionValue(opt *Option, optionValues *OptionValues) (interface{}, error) {
	gt.printOption(opt, optionValues)
	defer fmt.Fprintln(gt.Out)

	s, err := gt.readStdin()
	if err != nil {
		return nil, err
	}

	defaultVal := opt.Default(optionValues)

	var returnVal interface{}

	// TODO: cleanup somehow
	if s == "" {
		returnVal = defaultVal
	} else {
		switch defaultVal.(type) {
		case string:
			returnVal = s
		case bool:
			boolVal, err := strconv.ParseBool(s)
			if err != nil {
				return nil, err
			}
			returnVal = boolVal
		case int:
			intVal, err := strconv.Atoi(s)
			if err != nil {
				return nil, err
			}
			returnVal = intVal
		default:
			panic("unsupported type")
		}
	}

	if err := opt.Validate(returnVal); err != nil {
		gt.printf("\n")
		gt.printWarningf("Validation failed: %s", err.Error())
		return gt.readOptionValue(opt, optionValues)
	}

	return returnVal, nil
}

func (gt *GT) readStdin() (string, error) {
	if ok := gt.InScanner.Scan(); !ok {
		return "", gt.InScanner.Err()
	}

	return strings.TrimSpace(gt.InScanner.Text()), nil
}

// executeTemplateString executes the template in input str with the default p.FuncMap and valueMap as data.
func (gt *GT) executeTemplateString(str string, optionValues *OptionValues) (string, error) {
	tmpl, err := template.New("").Funcs(gt.FuncMap).Parse(str)
	if err != nil {
		return "", err
	}

	var buffer bytes.Buffer
	if err := tmpl.Execute(&buffer, optionValues); err != nil {
		return "", err
	}

	return buffer.String(), nil
}
