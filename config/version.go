package config

import (
	_ "embed"
	"strings"

	"github.com/Masterminds/semver"
)

var (
	Version       = strings.TrimSpace(version)
	VersionSemver = semver.MustParse(Version)

	//go:embed version.txt
	version string
)
