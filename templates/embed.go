package templates

import "embed"

//go:embed */*.tmpl */*/*.tmpl */*/*/*.tmpl */*/*/*/*.tmpl
var FS embed.FS
