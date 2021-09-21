package config

import (
	_ "embed"
	"strings"
)

var (
	Version = strings.TrimSpace(version)

	//go:embed version.txt
	version string
)
