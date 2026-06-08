// Package config handles PatchPrompt configuration and theme loading.
package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"patchprompt/themes"
)

// BuiltinThemes is the canonical list of theme names shipped in the binary.
var BuiltinThemes = []string{"default", "neon", "minimal", "rainbow", "stealth"}

// SegmentSpec describes one prompt segment in a theme.
type SegmentSpec struct {
	Type       string            `json:"type"`
	Foreground string            `json:"foreground"`
	Background string            `json:"background"`
	Icon       map[string]string `json:"icon"`
	Style      string            `json:"style"`
}

// Theme is a full prompt theme definition.
type Theme struct {
	Name      string            `json:"name"`
	FontMode  string            `json:"fontMode"`
	Separator map[string]string `json:"separator"`
	Segments  []SegmentSpec     `json:"segments"`
}

// Config is the user configuration.
type Config struct {
	ActiveTheme string `json:"activeTheme"`
	FontMode    string `json:"fontMode"`
	Path        struct {
		MaxDepth int `json:"maxDepth"`
	} `json:"path"`
	Time struct {
		Format string `json:"format"`
	} `json:"time"`
	Custom struct {
		Text string `json:"text"`
		Icon string `json:"icon"`
	} `json:"custom"`
}

// Default returns a Config populated with sensible defaults.
func Default() *Config {
	c := &Config{ActiveTheme: "neon", FontMode: ""}
	c.Path.MaxDepth = 3
	c.Time.Format = "15:04:05"
	return c
}

// Dir returns the PatchPrompt config directory, honoring PATCHPROMPT_CONFIG.
// If PATCHPROMPT_CONFIG points at a file, its directory is used.
func Dir() string {
	if v := os.Getenv("PATCHPROMPT_CONFIG"); v != "" {
		if filepath.Ext(v) == ".json" {
			return filepath.Dir(v)
		}
		return v
	}
	home, err := os.UserHomeDir()
	if err != nil {
		home = "."
	}
	return filepath.Join(home, ".patchprompt")
}

// ConfigPath returns the full path to config.json.
func ConfigPath() string {
	if v := os.Getenv("PATCHPROMPT_CONFIG"); v != "" && filepath.Ext(v) == ".json" {
		return v
	}
	return filepath.Join(Dir(), "config.json")
}

// themesDir returns the on-disk themes directory inside the config dir.
func themesDir() string { return filepath.Join(Dir(), "themes") }

// Load reads config.json, falling back to defaults if it is missing.
func Load() *Config {
	c := Default()
	data, err := os.ReadFile(ConfigPath())
	if err != nil {
		return c
	}
	_ = json.Unmarshal(data, c) // partial/invalid config falls back to defaults
	if c.ActiveTheme == "" {
		c.ActiveTheme = "neon"
	}
	if c.Path.MaxDepth <= 0 {
		c.Path.MaxDepth = 3
	}
	if c.Time.Format == "" {
		c.Time.Format = "15:04:05"
	}
	return c
}

// Save writes the config to disk, creating the directory if needed.
func (c *Config) Save() error {
	if err := os.MkdirAll(Dir(), 0o755); err != nil {
		return err
	}
	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(ConfigPath(), data, 0o644)
}

// readThemeBytes returns raw JSON for a theme, preferring the on-disk copy and
// falling back to the embedded built-in.
func readThemeBytes(name string) ([]byte, error) {
	if data, err := os.ReadFile(filepath.Join(themesDir(), name+".json")); err == nil {
		return data, nil
	}
	return themes.FS.ReadFile(name + ".json")
}

// LoadTheme loads a theme by name.
func LoadTheme(name string) (*Theme, error) {
	data, err := readThemeBytes(name)
	if err != nil {
		return nil, fmt.Errorf("theme %q not found", name)
	}
	var t Theme
	if err := json.Unmarshal(data, &t); err != nil {
		return nil, fmt.Errorf("theme %q is invalid: %w", name, err)
	}
	if t.Name == "" {
		t.Name = name
	}
	return &t, nil
}

// ListThemes returns the names of all available themes (on-disk + built-in).
func ListThemes() []string {
	seen := map[string]bool{}
	var out []string
	for _, n := range BuiltinThemes {
		if !seen[n] {
			seen[n] = true
			out = append(out, n)
		}
	}
	if entries, err := os.ReadDir(themesDir()); err == nil {
		for _, e := range entries {
			if filepath.Ext(e.Name()) == ".json" {
				n := e.Name()[:len(e.Name())-len(".json")]
				if !seen[n] {
					seen[n] = true
					out = append(out, n)
				}
			}
		}
	}
	return out
}

// ThemeExists reports whether a theme can be loaded.
func ThemeExists(name string) bool {
	_, err := readThemeBytes(name)
	return err == nil
}

// WriteDefaults writes the default config and all built-in themes to disk.
func WriteDefaults(force bool) error {
	if err := os.MkdirAll(themesDir(), 0o755); err != nil {
		return err
	}
	for _, name := range BuiltinThemes {
		data, err := themes.FS.ReadFile(name + ".json")
		if err != nil {
			return err
		}
		dst := filepath.Join(themesDir(), name+".json")
		if _, err := os.Stat(dst); err == nil && !force {
			continue
		}
		if err := os.WriteFile(dst, data, 0o644); err != nil {
			return err
		}
	}
	if _, err := os.Stat(ConfigPath()); err != nil || force {
		if err := Default().Save(); err != nil {
			return err
		}
	}
	return nil
}

// ResolveFontMode picks the effective font mode: config override wins, then the
// theme's declared mode, then "nerd".
func ResolveFontMode(c *Config, t *Theme) string {
	if c != nil && c.FontMode != "" {
		return c.FontMode
	}
	if t != nil && t.FontMode != "" {
		return t.FontMode
	}
	return "nerd"
}
