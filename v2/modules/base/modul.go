package base

import (
	"github.com/schwarzit/go-template/v2/gotemplate"
)

type Module struct {
	gotemplate.ModuleData
	gotemplate.ModuleBase
	// questions map[string]gotemplate.TemplateQuestion
}

func New() (*Module, error) {
	files, err := gotemplate.FindFiles("v2/modules/base/template")
	if err != nil {
		return nil, err
	}

	return &Module{
		ModuleData: gotemplate.ModuleData{
			Name:          "base",
			TemplatePath:  "v2/modules/base/template",
			TemplateFiles: files,
		},
		ModuleBase: gotemplate.ModuleBase{
			Questions: []gotemplate.TemplateQuestion{
				// {
				// 	Name:        "project-name",
				// 	Description: "Project name",
				// 	IsValid: func(value interface{}) (bool, string) {
				// 		return true, ""
				// 	},
				// 	DefaultValue: gotemplate.StringPtr("my-project"),
				// },
				// {
				// 	Name:        "project-description",
				// 	Description: "Project description",
				// 	IsValid: func(value interface{}) (bool, string) {
				// 		return true, ""
				// 	},
				// },
			},
		},
	}, nil
}
