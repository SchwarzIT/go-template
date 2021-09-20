package config

import (
	_ "embed"
	"strings"
)

var (
	Version string = strings.TrimSpace(version)

	//go:embed version.txt
	version string
)
