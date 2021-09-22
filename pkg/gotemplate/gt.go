package gotemplate

import (
	"bufio"
	"github.com/google/go-github/v39/github"
	"io"
	"text/template"

	"github.com/Masterminds/sprig"
	"github.com/schwarzit/go-template/config"
	"github.com/schwarzit/go-template/pkg/option"
	"sigs.k8s.io/yaml"
)

type GT struct {
	Streams
	FuncMap      template.FuncMap
	Options      []option.Option
	GitHubClient *github.Client
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

	gt := &GT{
		Options:      options,
		GitHubClient: github.NewClient(nil),
	}

	gt.FuncMap = sprig.TxtFuncMap()
	gt.FuncMap["latestReleaseTag"] = gt.latestReleaseTagWithDefault

	return gt
}

func (gt *GT) latestReleaseTagWithDefault(owner, repo, defaultTag string) string {
	tag, err := gt.getLatestGithubVersion(owner, repo)
	if err != nil {
		return defaultTag
	}

	return tag
}
