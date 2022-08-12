package config

import (
	_ "embed"
	"strings"

	"github.com/Masterminds/semver/v3"
)

var (
	Version       = strings.TrimSpace(version) //nolint:gochecknoglobals // version
	VersionSemver = semver.MustParse(Version)  //nolint:gochecknoglobals // version

	//go:embed version.txt
	version string
)
