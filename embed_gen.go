package gotemplate

import "embed"

//go:embed _template _template/.dockerignore _template/.editorconfig _template/.githooks _template/.gitignore _template/.golangci.yml _template/assets/.gitkeep _template/configs/.gitkeep _template/deployments/.gitkeep _template/internal/.gitkeep _template/pkg/.gitkeep
var FS embed.FS
