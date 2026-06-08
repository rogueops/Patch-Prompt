// Package themes embeds the built-in PatchPrompt theme definitions so the
// binary works before any files are installed to disk. install.ps1 also copies
// these same *.json files to the user config directory.
package themes

import "embed"

//go:embed *.json
var FS embed.FS
