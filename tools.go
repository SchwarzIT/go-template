//go:build tools

package main

import (
	// Go vulnerability scanner
	// https://go.dev/blog/vuln
	_ "golang.org/x/vuln/cmd/govulncheck"
)
