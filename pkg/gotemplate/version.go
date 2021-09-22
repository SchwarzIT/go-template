package gotemplate

import (
	"context"
	"fmt"

	"github.com/schwarzit/go-template/config"
	"github.com/schwarzit/go-template/pkg/repos"
)

const (
	goTemplateGithubOwner = "schwarzit"
	goTemplateGithubRepo  = "go-template"
)

func (gt *GT) PrintVersion() {
	_, _ = fmt.Fprintln(gt.Out, config.Version)
}

func (gt *GT) CheckVersion() {
	tag, err := gt.getLatestGithubVersion(goTemplateGithubOwner, goTemplateGithubRepo)
	if err != nil {
		gt.printWarning("unable to fetch version information. There could be newer release for go/template.")
		return
	}

	if tag != config.Version {
		gt.printWarning(fmt.Sprintf("newer version available: %s. Pls make sure to stay up to date to enjoy the latest features.", tag))
	}
}

func (gt *GT) getLatestGithubVersion(owner string, repo string) (string, error) {
	return repos.LatestGithubReleaseTag(repos.GithubTagListerFunc(func(ctx context.Context, owner, repo string) ([]string, error) {
		tags, _, err := gt.GitHubClient.Repositories.ListTags(ctx, owner, repo, nil)
		if err != nil {
			return nil, err
		}

		var tagStrings []string
		for _, tag := range tags {
			tagStrings = append(tagStrings, tag.GetName())
		}

		return tagStrings, nil
	}), owner, repo)
}
