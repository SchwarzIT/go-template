package base

import "github.com/schwarzit/go-template/v2/gotemplate"

type Module struct {
	gotemplate.ModuleData
	questions map[string]gotemplate.TemplateQuestion
}

func New() Module {
	return Module{
		ModuleData: gotemplate.ModuleData{
			Name:         "base",
			TemplatePath: "base",
			TemplateFiles: []string{
				"README.md",
				"LICENSE",
				"Makefile",
			},
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
	}
}
