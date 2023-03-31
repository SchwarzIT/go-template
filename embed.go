package gotemplate

import "embed"

const Key = "_template"

//go:embed all:_template
var FS embed.FS
