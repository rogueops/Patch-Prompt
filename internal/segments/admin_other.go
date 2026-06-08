//go:build !windows

package segments

import "os"

// Admin reports elevation on non-Windows platforms (uid 0 == root).
func Admin(_ Context) (string, bool) {
	if os.Geteuid() == 0 {
		return "ROOT", true
	}
	return "", false
}
