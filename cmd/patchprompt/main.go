// Command patchprompt is a Windows-first terminal prompt customizer.
package main

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"patchprompt/internal/config"
	"patchprompt/internal/render"
	"patchprompt/internal/shell"
)

const version = "1.0.0"

func main() {
	if len(os.Args) < 2 {
		usage()
		os.Exit(1)
	}
	switch os.Args[1] {
	case "init":
		cmdInit(os.Args[2:])
	case "render":
		cmdRender(os.Args[2:])
	case "config":
		cmdConfig(os.Args[2:])
	case "theme":
		cmdTheme(os.Args[2:])
	case "doctor":
		cmdDoctor()
	case "version", "--version", "-v":
		fmt.Println("PatchPrompt " + version)
	case "help", "--help", "-h":
		usage()
	default:
		fmt.Fprintf(os.Stderr, "unknown command: %s\n", os.Args[1])
		usage()
		os.Exit(1)
	}
}

func usage() {
	fmt.Println(`PatchPrompt ` + version + ` - Windows terminal prompt customizer

Usage:
  patchprompt init pwsh|powershell|cmd [--theme <name>]
  patchprompt render --shell pwsh|cmd [--status <code>] [--theme <name>]
  patchprompt config path
  patchprompt config generate
  patchprompt theme list
  patchprompt theme set <theme-name>
  patchprompt doctor
  patchprompt version`)
}

// flagValue extracts "--name value" from args; returns value and remaining args.
func flagValue(args []string, name string) (string, []string) {
	var rest []string
	val := ""
	for i := 0; i < len(args); i++ {
		if args[i] == name && i+1 < len(args) {
			val = args[i+1]
			i++
			continue
		}
		rest = append(rest, args[i])
	}
	return val, rest
}

func cmdInit(args []string) {
	if len(args) < 1 {
		fmt.Fprintln(os.Stderr, "usage: patchprompt init pwsh|powershell|cmd [--theme <name>]")
		os.Exit(1)
	}
	target := args[0]
	theme, _ := flagValue(args[1:], "--theme")
	if theme == "" {
		theme = config.Load().ActiveTheme
	}
	switch strings.ToLower(target) {
	case "pwsh", "powershell":
		fmt.Print(shell.PowerShell(theme))
	case "cmd":
		fmt.Println(shell.Cmd(theme))
	default:
		fmt.Fprintf(os.Stderr, "unknown shell: %s (use pwsh, powershell, or cmd)\n", target)
		os.Exit(1)
	}
}

func cmdRender(args []string) {
	shellName, args := flagValue(args, "--shell")
	statusStr, args := flagValue(args, "--status")
	themeName, _ := flagValue(args, "--theme")

	status := 0
	if statusStr != "" {
		if n, err := strconv.Atoi(strings.TrimSpace(statusStr)); err == nil {
			status = n
		}
	}

	cfg := config.Load()
	if themeName == "" {
		themeName = cfg.ActiveTheme
	}
	theme, err := config.LoadTheme(themeName)
	if err != nil {
		// Fall back to a minimal usable prompt rather than failing the shell.
		if shellName == "cmd" {
			fmt.Print("$ ")
		} else {
			fmt.Print("PS> ")
		}
		return
	}
	fmt.Print(render.Prompt(cfg, theme, status))
}

func cmdConfig(args []string) {
	if len(args) < 1 {
		fmt.Fprintln(os.Stderr, "usage: patchprompt config path|generate")
		os.Exit(1)
	}
	switch args[0] {
	case "path":
		fmt.Println(config.ConfigPath())
	case "generate":
		if err := config.WriteDefaults(false); err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Wrote config and themes to %s\n", config.Dir())
	default:
		fmt.Fprintf(os.Stderr, "unknown config subcommand: %s\n", args[0])
		os.Exit(1)
	}
}

func cmdTheme(args []string) {
	if len(args) < 1 {
		fmt.Fprintln(os.Stderr, "usage: patchprompt theme list|set <name>")
		os.Exit(1)
	}
	switch args[0] {
	case "list":
		active := config.Load().ActiveTheme
		for _, name := range config.ListThemes() {
			marker := "  "
			if name == active {
				marker = "* "
			}
			fmt.Printf("%s%s\n", marker, name)
		}
	case "set":
		if len(args) < 2 {
			fmt.Fprintln(os.Stderr, "usage: patchprompt theme set <name>")
			os.Exit(1)
		}
		name := args[1]
		if !config.ThemeExists(name) {
			fmt.Fprintf(os.Stderr, "theme %q not found. Available: %s\n",
				name, strings.Join(config.ListThemes(), ", "))
			os.Exit(1)
		}
		cfg := config.Load()
		cfg.ActiveTheme = name
		if err := cfg.Save(); err != nil {
			fmt.Fprintf(os.Stderr, "error saving config: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Active theme set to %q\n", name)
	default:
		fmt.Fprintf(os.Stderr, "unknown theme subcommand: %s\n", args[0])
		os.Exit(1)
	}
}

func cmdDoctor() {
	fmt.Println("PatchPrompt doctor")
	fmt.Println(strings.Repeat("-", 30))

	// Binary on PATH.
	if p, err := exec.LookPath("patchprompt"); err == nil {
		ok("patchprompt on PATH: " + p)
	} else {
		warn("patchprompt not found on PATH (run scripts/install.ps1)")
	}

	// Config.
	if _, err := os.Stat(config.ConfigPath()); err == nil {
		ok("config found: " + config.ConfigPath())
	} else {
		warn("config not found (run: patchprompt config generate)")
	}

	// Active theme loadable.
	cfg := config.Load()
	if _, err := config.LoadTheme(cfg.ActiveTheme); err == nil {
		ok("active theme loads: " + cfg.ActiveTheme)
	} else {
		warn("active theme failed to load: " + cfg.ActiveTheme)
	}

	// Git.
	if _, err := exec.LookPath("git"); err == nil {
		ok("git available (git segment enabled)")
	} else {
		warn("git not found (git segment will be hidden)")
	}

	// Font mode hint.
	theme, _ := config.LoadTheme(cfg.ActiveTheme)
	fm := config.ResolveFontMode(cfg, theme)
	ok("font mode: " + fm)
	if fm == "nerd" {
		fmt.Println("  -> install a Nerd Font (e.g. CaskaydiaCove Nerd Font) for full icons")
	}

	fmt.Println(strings.Repeat("-", 30))
	fmt.Println("Add to $PROFILE:")
	fmt.Println("  patchprompt init pwsh --theme " + cfg.ActiveTheme + " | Invoke-Expression")
}

func ok(msg string)   { fmt.Println("[ ok ] " + msg) }
func warn(msg string) { fmt.Println("[warn] " + msg) }
