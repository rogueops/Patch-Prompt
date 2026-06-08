// Package render turns a theme and runtime context into an ANSI prompt string.
package render

import (
	"strings"

	"patchprompt/internal/ansi"
	"patchprompt/internal/config"
	"patchprompt/internal/segments"
)

// piece is one resolved, non-empty segment ready to be drawn.
type piece struct {
	fg, bg string
	icon   string
	text   string
}

// icon picks the glyph for the active font mode, falling back across modes.
func iconFor(spec config.SegmentSpec, fontMode string) string {
	if spec.Icon == nil {
		return ""
	}
	if v, ok := spec.Icon[fontMode]; ok {
		return v
	}
	if v, ok := spec.Icon["plain"]; ok {
		return v
	}
	return ""
}

// label joins an icon and text with a single space when both are present.
func label(icon, text string) string {
	icon = strings.TrimSpace(icon)
	text = strings.TrimSpace(text)
	switch {
	case icon != "" && text != "":
		return icon + " " + text
	case icon != "":
		return icon
	default:
		return text
	}
}

// Prompt renders the full prompt for a theme. status is the previous command's
// exit code; trailingSpace adds a final space so input is readable.
func Prompt(cfg *config.Config, theme *config.Theme, status int) string {
	fontMode := config.ResolveFontMode(cfg, theme)
	ctx := segments.Context{Cfg: cfg, Status: status}

	var pieces []piece
	for _, spec := range theme.Segments {
		text, ok := segments.Content(spec.Type, ctx)
		if !ok {
			continue // segment failed or opted out: skip, don't break the prompt
		}
		fg, bg := spec.Foreground, spec.Background
		// Status segment communicates failure through color.
		if spec.Type == "status" {
			if status != 0 {
				fg, bg = "white", "red"
				if fontMode == "plain" || bg == "" {
					fg, bg = "brightRed", ""
				}
			}
		}
		lbl := label(iconFor(spec, fontMode), text)
		if lbl == "" && spec.Type != "status" {
			continue
		}
		pieces = append(pieces, piece{fg: fg, bg: bg, icon: iconFor(spec, fontMode), text: text})
	}
	if len(pieces) == 0 {
		return "$ "
	}

	switch fontMode {
	case "nerd":
		return renderPowerline(theme, pieces) + " "
	default:
		return renderPlain(theme, pieces, fontMode) + " "
	}
}

// renderPowerline draws filled segments with powerline separators between them.
func renderPowerline(theme *config.Theme, pieces []piece) string {
	sep := theme.Separator["nerd"]
	if sep == "" {
		sep = "\ue0b0"
	}
	var b strings.Builder
	for i, p := range pieces {
		hasBG := p.bg != ""
		if hasBG {
			b.WriteString(ansi.BG(p.bg))
		}
		b.WriteString(ansi.FG(p.fg))
		b.WriteString(" " + label(p.icon, p.text) + " ")
		// Separator into the next segment.
		next := i + 1
		if next < len(pieces) {
			nb := pieces[next].bg
			b.WriteString(ansi.Reset)
			if hasBG {
				b.WriteString(ansi.FG(p.bg)) // arrow colored as current bg
			}
			if nb != "" {
				b.WriteString(ansi.BG(nb)) // on next bg
			}
			b.WriteString(sep)
		} else {
			// Closing separator on the default terminal background.
			b.WriteString(ansi.Reset)
			if hasBG {
				b.WriteString(ansi.FG(p.bg) + sep)
			}
		}
	}
	b.WriteString(ansi.Reset)
	return b.String()
}

// renderPlain draws colored segments joined by the plain/emoji separator.
func renderPlain(theme *config.Theme, pieces []piece, fontMode string) string {
	sep := theme.Separator[fontMode]
	if sep == "" {
		sep = theme.Separator["plain"]
	}
	if sep == "" {
		sep = ">"
	}
	var parts []string
	for _, p := range pieces {
		seg := ansi.FG(p.fg)
		if p.bg != "" {
			seg += ansi.BG(p.bg)
		}
		seg += label(p.icon, p.text) + ansi.Reset
		parts = append(parts, seg)
	}
	return strings.Join(parts, " "+sep+" ")
}
