package segments

import (
	"os"
	"path/filepath"
	"strings"
)

// Runtime detects an active Python venv, a Node project, or a .NET project in
// the working directory. Fails silently when none are present.
func Runtime(_ Context) (string, bool) {
	var parts []string
	if venv := os.Getenv("VIRTUAL_ENV"); venv != "" {
		parts = append(parts, "py:"+filepath.Base(venv))
	}
	wd, err := os.Getwd()
	if err == nil {
		if fileExists(filepath.Join(wd, "package.json")) {
			parts = append(parts, "node")
		}
		if globExists(wd, "*.csproj") || globExists(wd, "*.sln") || globExists(wd, "*.fsproj") {
			parts = append(parts, "dotnet")
		}
	}
	if len(parts) == 0 {
		return "", false
	}
	return strings.Join(parts, " "), true
}

func fileExists(p string) bool {
	info, err := os.Stat(p)
	return err == nil && !info.IsDir()
}

func globExists(dir, pattern string) bool {
	matches, err := filepath.Glob(filepath.Join(dir, pattern))
	return err == nil && len(matches) > 0
}
