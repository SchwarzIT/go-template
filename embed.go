package gotemplate

import "embed"

//go:embed _template
var FS embed.FS

const Key = "_template"
