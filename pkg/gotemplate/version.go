package gotemplate

import (
	"fmt"

	"github.com/schwarzit/go-template/config"
	"github.com/schwarzit/go-template/pkg/repos"
)

const goTemplateGithubRepo = "https://github.com/schwarzit/go-template"

func (gt *GT) CheckVersion() {
	tag, err := repos.LatestReleaseTag(goTemplateGithubRepo)
	if err != nil {
		gt.printWarning("unable to fetch version information. There could be newer release for go/template.")
		return
	}

	if tag != config.Version {
		gt.printWarning(fmt.Sprintf("newer version available: %s. Pls make sure to stay up to date to enjoy the latest features.", tag))
	}
}
