package render

import (
	"strings"
	"testing"

	"patchprompt/internal/ansi"
	"patchprompt/internal/config"
)

func TestANSINamedColor(t *testing.T) {
	if got := ansi.FG("red"); got != "\x1b[31m" {
		t.Errorf("FG(red) = %q", got)
	}
	if got := ansi.BG("red"); got != "\x1b[41m" {
		t.Errorf("BG(red) = %q", got)
	}
	if got := ansi.FG("brightCyan"); got != "\x1b[96m" {
		t.Errorf("FG(brightCyan) = %q", got)
	}
	if got := ansi.FG("nope"); got != "" {
		t.Errorf("unknown color should be empty, got %q", got)
	}
}

func TestANSIHex(t *testing.T) {
	if got := ansi.FG("#ff8800"); got != "\x1b[38;2;255;136;0m" {
		t.Errorf("FG hex = %q", got)
	}
}

func TestRenderPlainContainsContent(t *testing.T) {
	cfg := config.Default()
	cfg.FontMode = "plain"
	theme := &config.Theme{
		Name:      "t",
		FontMode:  "plain",
		Separator: map[string]string{"plain": ">"},
		Segments: []config.SegmentSpec{
			{Type: "host", Foreground: "blue", Icon: map[string]string{"plain": ""}, Style: "plain"},
		},
	}
	out := Prompt(cfg, theme, 0)
	if !strings.Contains(out, ansi.Reset) {
		t.Errorf("expected a reset sequence in %q", out)
	}
	if strings.TrimSpace(out) == "" {
		t.Error("expected non-empty prompt")
	}
}

func TestRenderStatusFailureColors(t *testing.T) {
	cfg := config.Default()
	cfg.FontMode = "plain"
	theme := &config.Theme{
		Name:      "t",
		FontMode:  "plain",
		Separator: map[string]string{"plain": ">"},
		Segments: []config.SegmentSpec{
			{Type: "status", Foreground: "white", Icon: map[string]string{"plain": ">"}, Style: "plain"},
		},
	}
	out := Prompt(cfg, theme, 1)
	if !strings.Contains(out, ansi.FG("brightRed")) {
		t.Errorf("expected failure color in %q", out)
	}
}
