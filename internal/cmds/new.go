package cmds

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

	"github.com/Masterminds/sprig"
	"github.com/fatih/color"
	"github.com/schwarzit/go-template/config"
	"github.com/schwarzit/go-template/pkg/repos"
	"github.com/spf13/cobra"
	"sigs.k8s.io/yaml"
)

func init() {
	rootCmd.AddCommand(newCmd)

	funcMap = sprig.TxtFuncMap()
	funcMap["latestReleaseTag"] = latestReleaseTagWithDefault
	newCmd.Flags().StringVarP(
		&configFile,
		"config", "c", "",
		`Config file that defines all parameters.
This is helpful if you don't want to run the CLI interactively.
It should either be a json or a yaml file.`,
	)
}

var (
	newCmd = &cobra.Command{
		Use:   "new",
		Short: "Create a new project repository.",
		Long:  "Fill out all given parameters to configure and jump start your next project repository.",
		RunE: func(cmd *cobra.Command, args []string) error {
			return run()
		},
	}
	funcMap    template.FuncMap
	configFile string
)

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

func latestReleaseTagWithDefault(repo, defaultTag string) string {
	tag, err := repos.LatestReleaseTag(repo)
	if err != nil {
		return defaultTag
	}

	return tag
}

// TODO: only use globals in toplevel run. Pass as values to other functions to make testing easier.
func run() error {
	var options []Option
	if err := yaml.Unmarshal(config.Options, &options); err != nil {
		return err
	}

	optionNameToValue, err := getValues(options)
	if err != nil {
		return err
	}

	// TODO: add validation for value (probably regex pattern)

	printProgress("Generating repo folder...")

	targetDir := optionNameToValue["projectSlug"].(string)
	if _, err := os.Stat(targetDir); !os.IsNotExist(err) {
		return fmt.Errorf("directory %s already exists", targetDir)
	}

	err = fs.WalkDir(config.TemplateDir, config.TemplateKey, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		pathToWrite, err := executeTemplateString(path, funcMap, optionNameToValue)
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

		data, err := executeTemplateString(string(fileBytes), funcMap, optionNameToValue)
		if err != nil {
			return err
		}

		return os.WriteFile(pathToWrite, []byte(data), os.ModePerm)
	})
	if err != nil {
		return err
	}

	// TODO: delete created directory if any error occurs

	printProgress("Removing obsolete files...")
	if err := postHook(options, optionNameToValue); err != nil {
		return err
	}

	printProgress("Initializing git and Go modules...")
	if err := initRepo(optionNameToValue); err != nil {
		return err
	}

	return nil
}

func getValues(options []Option) (map[string]interface{}, error) {
	if configFile != "" {
		return getValuesFromFile(options, configFile)
	}
	return getValuesInteractively(options)
}

func getValuesFromFile(options []Option, file string) (map[string]interface{}, error) {
	fileBytes, err := os.ReadFile(file)
	if err != nil {
		return nil, err
	}

	optionNameToValue := make(map[string]interface{}, len(options))
	if err := yaml.Unmarshal(fileBytes, &optionNameToValue); err != nil {
		return nil, err
	}

	return optionNameToValue, nil
}

func getValuesInteractively(options []Option) (map[string]interface{}, error) {
	printBanner()
	optionNameToValue := make(map[string]interface{}, len(options))
	for _, currentOption := range options {
		if !dependenciesMet(currentOption, optionNameToValue) {
			continue
		}

		// default value could contain templating functions
		var err error
		currentOption.Default, err = applyTemplate(currentOption.Default, funcMap, optionNameToValue)
		if err != nil {
			return nil, err
		}

		val, err := readValue(currentOption)
		if err != nil {
			return nil, err
		}

		optionNameToValue[currentOption.Name] = val
	}

	return optionNameToValue, nil
}

func dependenciesMet(option Option, optionNameToValue map[string]interface{}) bool {
	if len(option.DependsOn) == 0 {
		return true
	}

	for _, dep := range option.DependsOn {
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

	return strings.TrimSpace(s), nil
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
