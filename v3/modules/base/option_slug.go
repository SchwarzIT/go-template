package base

import (
	"fmt"

	"github.com/schwarzit/go-template/v3/gotemplate/option"
	"github.com/schwarzit/go-template/v3/gotemplate/slugify"
)

type SlugOption struct {
	*option.BaseOption
}

func (o SlugOption) GetDefaultValue(state map[option.ModuleName]option.State) ([]string, error) {
	// Get the value of "Project Name" from the "README" module
	fmt.Println(state[ModuleName])
	projectName := state[ModuleName][TemplateTagProjectName]
	if len(projectName) == 1 {
		return []string{slugify.Slugify(projectName[0])}, nil
	}

	return []string{}, nil
}
