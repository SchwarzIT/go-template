package gotemplate

import (
	"bufio"
	"io"
	"text/template"

	"github.com/Masterminds/sprig"
	"github.com/schwarzit/go-template/config"
	"github.com/schwarzit/go-template/pkg/option"
	"github.com/schwarzit/go-template/pkg/repos"
	"sigs.k8s.io/yaml"
)

type GT struct {
	Streams
	FuncMap template.FuncMap
	Options []option.Option
}

type Streams struct {
	Out       io.Writer
	Err       io.Writer
	InScanner *bufio.Scanner
}

func New() *GT {
	var options []option.Option
	// panic error since the embedded file should always be valid
	if err := yaml.Unmarshal(config.Options, &options); err != nil {
		panic("embedded options are invalid")
	}

	funcMap := sprig.TxtFuncMap()
	funcMap["latestReleaseTag"] = latestReleaseTagWithDefault

	return &GT{
		FuncMap: funcMap,
		Options: options,
	}
}

func latestReleaseTagWithDefault(repo, defaultTag string) string {
	tag, err := repos.LatestReleaseTag(repo)
	if err != nil {
		return defaultTag
	}

	return tag
}
