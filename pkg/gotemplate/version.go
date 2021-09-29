package gotemplate

import (
	"github.com/schwarzit/go-template/config"
	"github.com/schwarzit/go-template/pkg/repos"
)

const (
	goTemplateGithubOwner = "schwarzit"
	goTemplateGithubRepo  = "go-template"
)

func (gt *GT) PrintVersion() {
	gt.printf(config.Version)
}

func (gt *GT) CheckVersion() {
	tag, err := repos.LatestGithubReleaseTag(gt.GithubTagLister, goTemplateGithubOwner, goTemplateGithubRepo)
	if err != nil {
		gt.printWarningf("unable to fetch version information. There could be newer release for go/template.")
		return
	}

	if tag.GreaterThan(config.VersionSemver) {
		gt.printWarningf("newer version available: %s. Pls make sure to stay up to date to enjoy the latest features.", tag)
	}
}
