package base

import (
	"github.com/schwarzit/go-template/v2/gotemplate"
)

type Module struct {
	gotemplate.ModuleData
	questions map[string]gotemplate.TemplateQuestion
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
		questions: map[string]gotemplate.TemplateQuestion{
			"project-name": {
				Name:        "project-name",
				Description: "Project name",
				IsValid: func(value interface{}) (bool, string) {
					return true, ""
				},
			},
			"project-description": {
				Name:        "project-description",
				Description: "Project description",
				IsValid: func(value interface{}) (bool, string) {
					return true, ""
				},
			},
		},
	}, nil
}
