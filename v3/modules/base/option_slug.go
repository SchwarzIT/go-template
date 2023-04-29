package base

import (
	"regexp"
	"strings"

	"github.com/schwarzit/go-template/v3/gotemplate/option"
)

type SlugOption struct {
	*option.BaseOption
}

func (o SlugOption) GetDefaultValue(state map[option.ModuleName]option.State) ([]string, error) {
	// Get the value of "Project Name" from the "README" module
	projectName := state[ModuleName][OptionProjectName]
	if len(projectName) == 1 {
		return []string{slugify(projectName[0])}, nil
	}

	return []string{}, nil
}

func slugify(input string) string {
	// Remove any non-alphanumeric or non-hyphen characters and replace with hyphens
	reg, err := regexp.Compile("[^a-zA-Z0-9-]+")
	if err != nil {
		panic(err)
	}

	// Replace any multiple consecutive hyphens with a single hyphen
	input = reg.ReplaceAllString(strings.TrimSpace(input), "-")

	// Convert the string to lowercase
	return strings.ToLower(input)
}
