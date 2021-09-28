package config

import "embed"

//go:embed _template
var TemplateDir embed.FS

const TemplateKey = "_template"

//go:embed options.yml
var Options []byte
