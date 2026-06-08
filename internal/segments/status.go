package segments

import "strconv"

// Status returns the previous command's exit code when non-zero (empty on
// success). The status segment is always shown so the prompt marker is stable;
// the render layer colors it based on the code.
func Status(ctx Context) (string, bool) {
	if ctx.Status != 0 {
		return strconv.Itoa(ctx.Status), true
	}
	return "", true
}
