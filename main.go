package main

import (
	"bufio"
	"bytes"
	"embed"
	"fmt"
	"github.com/fatih/color"
	"io/fs"
	"log"
	"os"
	"path"
	"runtime"
	"strconv"
	"strings"
	"text/template"

	"github.com/Masterminds/sprig"
	"github.com/schwarzit/go-template/pkg/repos"
	"sigs.k8s.io/yaml"
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
	Files       Files       `json:"files"`
}

type Files struct {
	Add    []string `json:"add"`
	Remove []string `json:"remove"`
}

func run() error {
	funcMap := sprig.TxtFuncMap()
	funcMap["latestReleaseTag"] = repos.LatestReleaseTag

	var options []Option
	if err := yaml.Unmarshal(optionBytes, &options); err != nil {
		return err
	}

	printBanner()
	optionNameToValue := make(map[string]interface{}, len(options))
	for _, currentOption := range options {
		// TODO: implement dependsOn
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

		optionNameToValue[currentOption.Name] = val
	}

	printGenerating()

	err := fs.WalkDir(dirTemplate, templateFolder, func(path string, d fs.DirEntry, err error) error {
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
	if err != nil {
		return err
	}

	postHook(options, optionNameToValue)
	// TODO: run certain setup commands
	//  - git init
	//  - go mod init
	//initRepo(optionNameToValue["projectSlug"])

	return nil
}

func postHook(options []Option, optionNameToValue map[string]interface{}) error {
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

func printGenerating() {
	_, _ = color.New(color.FgCyan, color.Bold).Println("Generating repo folder...")
}

// readValue reads a value from the cli.
func readValue(option Option) (interface{}, error) {
	printOption(option)
	defer fmt.Println()

	s, err := readStdin()
	if err != nil {
		return nil, err
	}

	if s == "" {
		return option.Default, nil
	}

	switch option.Default.(type) {
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

func readStdin() (string, error) {
	reader := bufio.NewReader(os.Stdin)
	s, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}

	s = strings.TrimSuffix(s, "\n")
	if runtime.GOOS == "windows" {
		s = strings.TrimSuffix(s, "\r")
	}

	return s, nil
}

func printOption(option Option) {
	highlight := color.New(color.FgCyan).SprintFunc()
	underline := color.New(color.FgHiYellow, color.Underline).SprintFunc()
	fmt.Printf("%s\n", underline(option.Description))
	fmt.Printf("%s: (%v) ", highlight(option.Name), option.Default)
}

func printBanner() {
	highlight := color.New(color.FgCyan).SprintFunc()
	fmt.Printf("Hi! Welcome to the %s cli.\n", highlight("go/template"))
	fmt.Printf("This command will walk you through creating a new project.\n\n")
	fmt.Printf("Enter a value or leave blank to accept the (default), and press %s.\n", highlight("<ENTER>"))
	fmt.Printf("Press %s at any time to quit.\n\n", highlight("^C"))
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
