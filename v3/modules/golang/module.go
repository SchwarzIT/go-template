package golang

import (
	"fmt"
	// "io/ioutil"
	// "os"
	// "path/filepath"

	"github.com/schwarzit/go-template/v3/gotemplate/module"
	"github.com/schwarzit/go-template/v3/gotemplate/option"
)

type ReadmeModule struct {
	*module.BaseModule
}

func NewReadmeModule() *ReadmeModule {
	nameOption := option.NewBaseOption(
		"Project Name",
		"The name of the project",
		"",
		nil,
		func(value string) error {
			if value == "" {
				return fmt.Errorf("project name cannot be empty")
			}
			return nil
		},
	)

	descOption := option.NewBaseOption(
		"Project Description",
		"A short description of the project",
		"",
		nil,
		func(value string) error {
			if value == "" {
				return fmt.Errorf("project description cannot be empty")
			}
			return nil
		},
	)

	return &ReadmeModule{
		module.NewBaseModule("README", []option.Option{nameOption, descOption}),
	}
}

func (m *ReadmeModule) Generate() error {
	// name, err := m.GetOptionValue("Project Name")
	// if err != nil {
	// 	return err
	// }

	// desc, err := m.GetOptionValue("Project Description")
	// if err != nil {
	// 	return err
	// }

	// content := fmt.Sprintf("# %s\n\n%s\n", name, desc)

	// dir, err := os.Getwd()
	// if err != nil {
	// 	return err
	// }

	// filename := filepath.Join(dir, "README.md")

	// if err := ioutil.WriteFile(filename, []byte(content), 0644); err != nil {
	// 	return err
	// }

	return nil
}
