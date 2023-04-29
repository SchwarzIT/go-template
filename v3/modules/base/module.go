package base

import (
	"errors"
	"fmt"
	"path/filepath"
	"runtime"

	"github.com/schwarzit/go-template/v3/gotemplate/module"
	"github.com/schwarzit/go-template/v3/gotemplate/option"
)

var (
	ErrUnableFindCurrentPath = errors.New("unable to get current file path")
)

const (
	ModuleName             = "README"
	TemplateTagProjectName = "NAME"
	TemplateTagProjectSlug = "SLUG"
)

type ReadmeModule struct {
	*module.BaseModule
}

func NewReadmeModule() (*ReadmeModule, error) {
	nameOption := option.NewBaseOption(option.NewBaseOptionArgs{
		Title:       "Project Name",
		TemplateTag: TemplateTagProjectName,
		Description: "The name of the project",
	})

	slugOption := SlugOption{
		BaseOption: option.NewBaseOption(option.NewBaseOptionArgs{
			Title:       "Project Slug",
			TemplateTag: TemplateTagProjectSlug,
			Description: "The name of the project slug",
		}),
	}

	baseModule, err := module.NewBaseModule(
		ModuleName,
		[]option.Option{nameOption, slugOption},
	)
	if err != nil {
		return nil, err
	}

	return &ReadmeModule{baseModule}, nil
}

func (m *ReadmeModule) GetTemplateFiles() ([]module.TemplateFile, error) {
	_, currentFilePath, _, ok := runtime.Caller(0)
	if !ok {
		return nil, ErrUnableFindCurrentPath
	}
	currentPackagePath := filepath.Dir(currentFilePath)

	fmt.Println("Current package path:", currentPackagePath)

	return []module.TemplateFile{
		&module.BaseTemplateFile{
			SourcePath: filepath.Join(currentPackagePath, "template", "README.md"),
			TargetPath: "README.md",
		},
	}, nil
}
