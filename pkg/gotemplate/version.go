package gotemplate

import (
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
	tag, err := repos.LatestGithubReleaseTag(gt.GithubTagLister, goTemplateGithubOwner, goTemplateGithubRepo)
	if err != nil {
		gt.printWarning("unable to fetch version information. There could be newer release for go/template.")
		return
	}

	if tag.GreaterThan(config.VersionSemver) {
		gt.printWarning(fmt.Sprintf("newer version available: %s. Pls make sure to stay up to date to enjoy the latest features.", tag))
	}
}
