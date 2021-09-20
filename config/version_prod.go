//go:build prod
// +build prod

package config

import (
	_ "embed"
)

var (
	Version string = strings.TrimSpace(version)

	//go:embed version.txt
	version string
)
