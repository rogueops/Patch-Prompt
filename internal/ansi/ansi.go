// Package ansi builds ANSI escape sequences for terminal colors and styling.
package ansi

import (
	"fmt"
	"strconv"
	"strings"
)

const (
	esc   = "\x1b["
	Reset = "\x1b[0m"
)

// fgCodes maps named colors to their foreground SGR codes.
var fgCodes = map[string]int{
	"black": 30, "red": 31, "green": 32, "yellow": 33,
	"blue": 34, "magenta": 35, "cyan": 36, "white": 37,
	"gray": 90, "grey": 90,
	"brightBlack": 90, "brightRed": 91, "brightGreen": 92, "brightYellow": 93,
	"brightBlue": 94, "brightMagenta": 95, "brightCyan": 96, "brightWhite": 97,
}

// hexToRGB parses a "#rrggbb" string. Returns false if it is not valid hex.
func hexToRGB(s string) (int, int, int, bool) {
	s = strings.TrimPrefix(s, "#")
	if len(s) != 6 {
		return 0, 0, 0, false
	}
	v, err := strconv.ParseInt(s, 16, 32)
	if err != nil {
		return 0, 0, 0, false
	}
	return int(v>>16) & 0xff, int(v>>8) & 0xff, int(v) & 0xff, true
}

// FG returns a foreground color escape for a named color or "#rrggbb" hex.
// An empty or unknown name returns an empty string (no color change).
func FG(name string) string {
	if name == "" {
		return ""
	}
	if strings.HasPrefix(name, "#") {
		if r, g, b, ok := hexToRGB(name); ok {
			return fmt.Sprintf("%s38;2;%d;%d;%dm", esc, r, g, b)
		}
		return ""
	}
	if code, ok := fgCodes[name]; ok {
		return fmt.Sprintf("%s%dm", esc, code)
	}
	return ""
}

// BG returns a background color escape for a named color or "#rrggbb" hex.
func BG(name string) string {
	if name == "" {
		return ""
	}
	if strings.HasPrefix(name, "#") {
		if r, g, b, ok := hexToRGB(name); ok {
			return fmt.Sprintf("%s48;2;%d;%d;%dm", esc, r, g, b)
		}
		return ""
	}
	if code, ok := fgCodes[name]; ok {
		return fmt.Sprintf("%s%dm", esc, code+10)
	}
	return ""
}

// FGAsBG converts a foreground color name into a background escape. This is
// used to color powerline separators with the previous segment's background.
func FGAsBG(name string) string {
	return BG(name)
}
