package gotemplate

import (
	"bufio"
	"context"
	"io"
	"net/http"
	"text/template"
	"time"

	"github.com/Masterminds/sprig"
	"github.com/google/go-github/v39/github"
	gotemplate "github.com/schwarzit/go-template"
	"github.com/schwarzit/go-template/pkg/option"
	"github.com/schwarzit/go-template/pkg/repos"
	"sigs.k8s.io/yaml"
)

type GT struct {
	Streams
	FuncMap         template.FuncMap
	Configs         option.Configuration
	GithubTagLister repos.GithubTagLister
}

type Streams struct {
	Out       io.Writer
	Err       io.Writer
	InScanner *bufio.Scanner
}

func New() *GT {
	var configs option.Configuration
	// panic error since the embedded file should always be valid
	if err := yaml.Unmarshal(gotemplate.Options, &configs); err != nil {
		panic("embedded options are invalid")
	}

	githubClient := github.NewClient(&http.Client{Timeout: time.Second})

	gt := &GT{
		Configs: configs,
		GithubTagLister: repos.GithubTagListerFunc(func(ctx context.Context, owner, repo string) ([]string, error) {
			tags, _, err := githubClient.Repositories.ListTags(ctx, owner, repo, nil)
			if err != nil {
				return nil, err
			}

			var tagStrings []string
			for _, tag := range tags {
				tagStrings = append(tagStrings, tag.GetName())
			}

			return tagStrings, nil
		}),
	}

	gt.FuncMap = sprig.TxtFuncMap()
	gt.FuncMap["latestReleaseTag"] = gt.latestReleaseTagWithDefault

	return gt
}

func (gt *GT) latestReleaseTagWithDefault(owner, repo, defaultTag string) string {
	tag, err := repos.LatestGithubReleaseTag(gt.GithubTagLister, owner, repo)
	if err != nil {
		return defaultTag
	}

	return tag.String()
}
