package segments

import (
	"os"
	"path/filepath"
	"strings"
)

// Path returns the working directory with the home dir collapsed to ~ and
// truncated to the configured maxDepth.
func Path(ctx Context) (string, bool) {
	wd, err := os.Getwd()
	if err != nil {
		return "", false
	}
	if home, err := os.UserHomeDir(); err == nil && home != "" {
		if wd == home {
			return "~", true
		}
		if strings.HasPrefix(wd, home+string(os.PathSeparator)) {
			wd = "~" + wd[len(home):]
		}
	}
	wd = filepath.ToSlash(wd)
	parts := strings.Split(wd, "/")
	max := ctx.Cfg.Path.MaxDepth
	if max > 0 && len(parts) > max {
		parts = append([]string{"\u2026"}, parts[len(parts)-max:]...)
	}
	return strings.Join(parts, "/"), true
}
