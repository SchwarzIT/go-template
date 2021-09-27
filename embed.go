package gotemplate

import "embed"

var (
	//go:embed _template
	FS embed.FS

	//go:embed options.yml
	Options []byte
)

const Key = "_template"
