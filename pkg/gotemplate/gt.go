package gotemplate

import (
	"bufio"
	"context"
	"io"
	"net/http"
	"text/template"
	"time"

	"github.com/Masterminds/sprig/v3"
	"github.com/google/go-github/v39/github"
	"github.com/schwarzit/go-template/pkg/repos"
)

type GT struct {
	Streams
	Options         *Options
	FuncMap         template.FuncMap
	GithubTagLister repos.GithubTagLister
}

type Streams struct {
	Out       io.Writer
	Err       io.Writer
	InScanner *bufio.Scanner
}

func New() *GT {
	githubClient := github.NewClient(&http.Client{Timeout: time.Second})
	githubTagLister := repos.GithubTagListerFunc(func(ctx context.Context, owner, repo string) ([]string, error) {
		tags, _, err := githubClient.Repositories.ListTags(ctx, owner, repo, nil)
		if err != nil {
			return nil, err
		}

		var tagStrings []string
		for _, tag := range tags {
			tagStrings = append(tagStrings, tag.GetName())
		}

		return tagStrings, nil
	})

	return &GT{
		Options:         NewOptions(githubTagLister),
		GithubTagLister: githubTagLister,
		FuncMap:         sprig.TxtFuncMap(),
	}
}
