package segments

import "os"

// Host returns the computer name.
func Host(_ Context) (string, bool) {
	h, err := os.Hostname()
	if err != nil || h == "" {
		return "", false
	}
	return h, true
}
