package segments

// Custom returns the user-configured custom text, if any.
func Custom(ctx Context) (string, bool) {
	if ctx.Cfg.Custom.Text == "" {
		return "", false
	}
	return ctx.Cfg.Custom.Text, true
}
