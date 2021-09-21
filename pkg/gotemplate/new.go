package gotemplate

import (
	"bufio"
	"bytes"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path"
	"runtime"
	"strconv"
	"strings"
	"text/template"

	"github.com/pkg/errors"
	"github.com/schwarzit/go-template/config"
	"github.com/schwarzit/go-template/pkg/option"
	"sigs.k8s.io/yaml"
)

var ErrAlreadyExists = errors.New("already exists")

type NewRepositoryOptions struct {
	OptionNameToValue map[string]interface{}
}

// LoadOptionToValueFromFile loads value for the options from a file.
func (gt *GT) LoadOptionToValueFromFile(file string) (map[string]interface{}, error) {
	fileBytes, err := os.ReadFile(file)
	if err != nil {
		return nil, err
	}

	optionNameToValue := make(map[string]interface{}, len(gt.Options))
	if err := yaml.Unmarshal(fileBytes, &optionNameToValue); err != nil {
		return nil, err
	}

	return optionNameToValue, nil
}

func (gt *GT) GetOptionToValueInteractively() (map[string]interface{}, error) {
	// TODO: add validation for value (probably regex pattern)

	gt.printBanner()
	optionNameToValue := make(map[string]interface{}, len(gt.Options))
	for _, currentOption := range gt.Options {
		// Fix implicit memory aliasing (gosec G601)
		currentOption := currentOption
		if !dependenciesMet(&currentOption, optionNameToValue) {
			continue
		}

		// default value could contain templating functions
		var err error
		currentOption.Default, err = gt.applyTemplate(currentOption.Default, optionNameToValue)
		if err != nil {
			return nil, err
		}

		val, err := gt.readOptionValue(&currentOption)
		if err != nil {
			return nil, err
		}

		optionNameToValue[currentOption.Name] = val
	}

	return optionNameToValue, nil
}

func dependenciesMet(opt *option.Option, optionNameToValue map[string]interface{}) bool {
	if len(opt.DependsOn) == 0 {
		return true
	}

	for _, dep := range opt.DependsOn {
		depVal, ok := optionNameToValue[dep]
		if !ok {
			// if not found it means it not set
			return false
		}

		depBoolVal, ok := depVal.(bool)
		if !ok {
			// value will only be checked for bool values
			continue
		}

		if !depBoolVal {
			return false
		}
	}

	return true
}

func (gt *GT) InitNewProject(opts *NewRepositoryOptions) (err error) {
	gt.printProgress("Generating repo folder...")

	targetDir := opts.OptionNameToValue["projectSlug"].(string)
	if _, err := os.Stat(targetDir); !os.IsNotExist(err) {
		return errors.Wrapf(ErrAlreadyExists, "directory %s", targetDir)
	}

	defer func() {
		if err != nil {
			// ignore error to not overwrite original error
			_ = os.RemoveAll(targetDir)
		}
	}()
	err = fs.WalkDir(config.TemplateDir, config.TemplateKey, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		pathToWrite, err := gt.executeTemplateString(path, opts.OptionNameToValue)
		if err != nil {
			return err
		}

		pathToWrite = strings.ReplaceAll(pathToWrite, config.TemplateKey, targetDir)
		if d.IsDir() {
			return os.MkdirAll(pathToWrite, os.ModePerm)
		}

		fileBytes, err := fs.ReadFile(config.TemplateDir, path)
		if err != nil {
			return err
		}

		data, err := gt.executeTemplateString(string(fileBytes), opts.OptionNameToValue)
		if err != nil {
			return err
		}

		return os.WriteFile(pathToWrite, []byte(data), os.ModePerm)
	})
	if err != nil {
		return err
	}

	gt.printProgress("Removing obsolete files...")
	if err := postHook(gt.Options, opts.OptionNameToValue); err != nil {
		return err
	}

	gt.printProgress("Initializing git and Go modules...")
	if err := initRepo(opts.OptionNameToValue); err != nil {
		return err
	}

	return nil
}

func initRepo(optionToNameValue map[string]interface{}) error {
	targetDir := optionToNameValue["projectSlug"].(string)

	gitInit := exec.Command("git", "init")
	gitInit.Dir = targetDir

	if err := gitInit.Run(); err != nil {
		return err
	}

	// nolint: gosec // no security issue possible with go mod init
	goModInit := exec.Command("go", "mod", "init", optionToNameValue["moduleName"].(string))
	goModInit.Dir = targetDir

	return goModInit.Run()
}

func postHook(options []option.Option, optionNameToValue map[string]interface{}) error {
	var toDelete []string

	for _, opt := range options {
		optEnabled, ok := optionNameToValue[opt.Name].(bool)
		if !ok {
			// if not bool value, files will be ignored
			continue
		}

		if optEnabled {
			toDelete = append(toDelete, opt.Files.Remove...)
			continue
		}
		// the files are added in the loop anyways, but if the option is disabled they should be removed again
		toDelete = append(toDelete, opt.Files.Add...)
	}

	for _, item := range toDelete {
		if err := os.RemoveAll(path.Join(optionNameToValue["projectSlug"].(string), item)); err != nil {
			return err
		}
	}

	return nil
}

// readOptionValue reads a value for an option from the cli.
func (gt *GT) readOptionValue(opts *option.Option) (interface{}, error) {
	gt.printOption(opts)
	defer fmt.Fprintln(gt.Out)

	s, err := gt.readStdin()
	if err != nil {
		return nil, err
	}

	if s == "" {
		return opts.Default, nil
	}

	switch opts.Default.(type) {
	case string:
		return s, nil
	case bool:
		return strconv.ParseBool(s)
	case int:
		return strconv.Atoi(s)
	default:
		panic("unsupported type")
	}
}

func (gt *GT) readStdin() (string, error) {
	reader := bufio.NewReader(gt.In)
	s, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}

	s = strings.TrimSuffix(s, "\n")
	if runtime.GOOS == "windows" {
		s = strings.TrimSuffix(s, "\r")
	}

	return strings.TrimSpace(s), nil
}

// applyTemplate executes a the template in the defaultValue with the valueMap as data.
// If the defaultValue is not a string, the input defaultValue will be returned.
func (gt *GT) applyTemplate(defaultValue interface{}, valueMap map[string]interface{}) (interface{}, error) {
	defaultStr, ok := defaultValue.(string)
	if !ok {
		return defaultValue, nil
	}

	return gt.executeTemplateString(defaultStr, valueMap)
}

// executeTemplateString executes the template in input str with the default p.FuncMap and valueMap as data.
func (gt *GT) executeTemplateString(str string, valueMap map[string]interface{}) (string, error) {
	tmpl, err := template.New("").Funcs(gt.FuncMap).Parse(str)
	if err != nil {
		return "", err
	}

	var buffer bytes.Buffer
	if err := tmpl.Execute(&buffer, valueMap); err != nil {
		return "", err
	}

	return buffer.String(), nil
}
