package segments

import "time"

// Time returns the current local time formatted per config.
func Time(ctx Context) (string, bool) {
	format := ctx.Cfg.Time.Format
	if format == "" {
		format = "15:04:05"
	}
	return time.Now().Format(format), true
}
