package assets

import "embed"

//go:embed static
var StaticFiles embed.FS

//go:embed templates
var TemplateFS embed.FS
