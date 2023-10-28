package gotemplate

import (
	"bufio"
	"context"
	"io"
	"net/http"
	"sync"
	"text/template"
	"time"

	"github.com/Masterminds/sprig/v3"
	"github.com/google/go-github/v56/github"
	"github.com/muesli/termenv"
	"github.com/schwarzit/go-template/pkg/repos"
)

type GT struct {
	Streams
	Options         *Options
	FuncMap         template.FuncMap
	GithubTagLister repos.GithubTagLister
	once            sync.Once
	output          *termenv.Output
}

func (gt *GT) styler() *termenv.Output {
	if gt.output != nil {
		return gt.output
	}

	if gt.Out == nil {
		// panic here since it's a package user error
		// that it is not set
		panic("gt out stream not set")
	}

	gt.once.Do(func() {
		gt.output = termenv.NewOutput(gt.Out, termenv.WithProfile(termenv.EnvColorProfile()))
	})

	return gt.output
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
