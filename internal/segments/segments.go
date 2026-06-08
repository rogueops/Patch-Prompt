// Package segments computes the textual content for each prompt segment.
// Each segment returns its content and whether it should be shown. Segments
// fail silently (return ok=false) rather than breaking the prompt.
package segments

import "patchprompt/internal/config"

// Context carries the data segments need to render.
type Context struct {
	Cfg    *config.Config
	Status int
}

// Content returns the rendered text for a segment type and whether to show it.
func Content(segType string, ctx Context) (string, bool) {
	switch segType {
	case "path":
		return Path(ctx)
	case "git":
		return Git(ctx)
	case "user":
		return User(ctx)
	case "host":
		return Host(ctx)
	case "time":
		return Time(ctx)
	case "status":
		return Status(ctx)
	case "admin":
		return Admin(ctx)
	case "custom":
		return Custom(ctx)
	case "runtime":
		return Runtime(ctx)
	default:
		return "", false
	}
}
