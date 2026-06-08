package segments

import (
	"os/user"
	"strings"
)

// User returns the current username without any domain prefix.
func User(_ Context) (string, bool) {
	u, err := user.Current()
	if err != nil || u.Username == "" {
		return "", false
	}
	name := u.Username
	if i := strings.LastIndexAny(name, "\\/"); i >= 0 {
		name = name[i+1:]
	}
	return name, true
}
