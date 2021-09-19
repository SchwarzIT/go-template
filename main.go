package main

import (
	"bufio"
	"bytes"
	"embed"
	"fmt"
	"github.com/Masterminds/sprig"
	"io/fs"
	"log"
	"os"
	"sigs.k8s.io/yaml"
	"strconv"
	"strings"
	"text/template"
)

//go:embed _template
var dirTemplate embed.FS

//go:embed options.yml
var optionBytes []byte

const templateFolder = "_template"

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

type Option struct {
	Name        string      `json:"name"`
	Default     interface{} `json:"default"`
	Description string      `json:"description"`
	DependsOn   []string    `json:"dependsOn"`
}

func run() error {
	funcMap := sprig.TxtFuncMap()
	funcMap["latestReleaseTag"] = repos.LatestReleaseTag

	var options []Option
	if err := yaml.Unmarshal(optionBytes, &options); err != nil {
		return err
	}

	optionNameToValue := make(map[string]interface{}, len(options))
	for _, currentOption := range options {
		log.Printf("key: %s", currentOption.Name)

		// default value could contain templating functions
		var err error
		currentOption.Default, err = applyTemplate(currentOption.Default, funcMap, optionNameToValue)
		if err != nil {
			return err
		}

		// TODO: only read value if dependencies.len == 0 or dependencies fulfilled
		val, err := readValue(currentOption)
		if err != nil {
			return err
		}
		log.Printf("new value: %v", val)

		optionNameToValue[currentOption.Name] = val
		log.Printf("--------------")
	}

	return fs.WalkDir(dirTemplate, templateFolder, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		pathToWrite, err := executeTemplateString(path, funcMap, optionNameToValue)
		if err != nil {
			return err
		}

		pathToWrite = strings.ReplaceAll(pathToWrite, templateFolder, optionNameToValue["projectSlug"].(string))
		if d.IsDir() {
			return os.MkdirAll(pathToWrite, os.ModePerm)
		}

		fileBytes, err := fs.ReadFile(dirTemplate, path)
		if err != nil {
			return err
		}

		data, err := executeTemplateString(string(fileBytes), funcMap, optionNameToValue)
		if err != nil {
			return err
		}

		return os.WriteFile(pathToWrite, []byte(data), os.ModePerm)
	})

	// TODO: implement post hook to remove files that are not needed
	//  - either manually, or assing define files that are specific to an option in the options
	//  - defined
	//      - should only be allowed for bool values (working as switches?)
	//      - validate options in pipeline to make sure they are valid
	//      - write file only if the file is they are not included in any disabled switch

	// TODO: run certain setup commands
	//  - git init
	//  - go mod init
}

// readValue reads a value from the cli.
func readValue(option Option) (interface{}, error) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Printf("Enter text for %s (%s, default: %v): ", option.Name, option.Description, option.Default)
	text, err := reader.ReadString('\n')
	if err != nil {
		return nil, err
	}

	text = strings.TrimRight(text, "\n")
	if text == "" {
		return option.Default, nil
	}

	switch option.Default.(type) {
	case string:
		return text, nil
	case bool:
		return strconv.ParseBool(text)
	case int:
		return strconv.Atoi(text)
	default:
		panic("unsupported type")
	}
}

// applyTemplate executes a the template in the defaultValue with the valueMap as data.
// If the defaultValue is not a string, the input defaultValue will be returned.
func applyTemplate(defaultValue interface{}, funcMap template.FuncMap, valueMap map[string]interface{}) (interface{}, error) {
	defaultStr, ok := defaultValue.(string)
	if !ok {
		return defaultValue, nil
	}

	return executeTemplateString(defaultStr, funcMap, valueMap)
}

// executeTemplateString executes the template in input str with the default funcMap and valueMap as data.
func executeTemplateString(str string, funcMap template.FuncMap, valueMap map[string]interface{}) (string, error) {
	tmpl, err := template.New("").Funcs(funcMap).Parse(str)
	if err != nil {
		return "", err
	}

	var buffer bytes.Buffer
	if err := tmpl.Execute(&buffer, valueMap); err != nil {
		return "", err
	}

	return buffer.String(), nil
}
