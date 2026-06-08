# PatchPrompt

A Windows-first terminal prompt customizer for **PowerShell** and **Command Prompt**.
PatchPrompt makes your prompt look good with colors, icons, Nerd Font powerline
separators, emojis, and useful segments (path, git, time, status, and more).

PatchPrompt is an independent project. It is **not** affiliated with, and does
not reuse the code, themes, names, or assets of, any other prompt tool.

## Version 1.0 status

Working MVP. The Go CLI builds and runs, all listed commands are implemented,
PowerShell integration is fully dynamic, cmd integration is functional (with
documented limits), themes/config work, and the release build produces a zip
plus a winget-ready Inno Setup installer when Inno Setup is present.

## Recommended Nerd Font

For full icons and powerline separators, install a Nerd Font. Recommended:
**CaskaydiaCove Nerd Font**. Then set it as your terminal font (Windows
Terminal: Settings -> your profile -> Appearance -> Font face).

## Quick Start

1. Install a Nerd Font:
   Recommended: CaskaydiaCove Nerd Font

2. Install PatchPrompt locally:
   .\scripts\install.ps1

3. Open your PowerShell profile:
   notepad $PROFILE

4. Add this line:
   patchprompt init pwsh --theme neon | Invoke-Expression

5. Save, close, reopen Windows Terminal.

BAM COLOR.

## Local install

```powershell
.\scripts\install.ps1
```

The installer (no admin required):
- builds `patchprompt.exe` with Go (if Go is installed),
- copies it to `%LOCALAPPDATA%\Programs\PatchPrompt`,
- writes default config + themes to `%USERPROFILE%\.patchprompt`,
- adds the install dir to your user `PATH`,
- prints the profile line and optionally adds it for you.

Add `-UpdateProfile` to update `$PROFILE` without prompting, or `-Theme minimal`
to install with a different default theme.

## PowerShell setup

After installing, add to your profile (`notepad $PROFILE`):

```powershell
patchprompt init pwsh --theme neon | Invoke-Expression
```

`init pwsh` emits a `prompt` function that calls `patchprompt render` on every
prompt. It preserves `$LASTEXITCODE`, reflects command success/failure in the
status segment, and falls back to a plain `PS path>` prompt if the binary is
missing — so it never breaks your shell.

## Command Prompt setup

cmd.exe **cannot run an external program on every prompt**, so PatchPrompt
applies a static, ANSI-colored `PROMPT` instead of dynamic git/runtime segments.
This is a real cmd limitation, not a missing feature.

Run in a session:

```bat
scripts\patchprompt.cmd
```

Or wire it into cmd AutoRun (see `examples\cmd-autorun-example.cmd`). For the
full dynamic experience, use PowerShell.

## Theme usage

```powershell
patchprompt theme list           # list themes (* = active)
patchprompt theme set minimal    # change the active theme
```

Or pin a theme at activation: `patchprompt init pwsh --theme neon | Invoke-Expression`.
A `--theme` passed to `init`/`render` overrides the config; otherwise the
config's `activeTheme` (default `neon`) is used.

Built-in themes: `default`, `neon`, `minimal`, `rainbow`, `stealth`.

## Config usage

```powershell
patchprompt config path       # show config.json location
patchprompt config generate   # write default config + themes to disk
```

Config lives at `%USERPROFILE%\.patchprompt\config.json`. Override the location
with the `PATCHPROMPT_CONFIG` environment variable.

Config controls the active theme, font mode (`nerd` / `emoji` / `plain`), path
`maxDepth`, time format, and custom segment text. Themes control colors, icons,
separators, and which segments are visible.

### Font modes

- `nerd` — Nerd Font icons + powerline separators (needs a Nerd Font).
- `emoji` — emoji icons + simple separators.
- `plain` — safe ASCII only (works in any terminal).

If a theme declares `nerd` but you set `fontMode` to `plain` in config, the
plain fallback is used.

## Release build

```powershell
.\scripts\build-release.ps1
```

Produces in `dist\`:
- `patchprompt-windows-amd64.zip`
- `checksums.txt`
- `PatchPromptSetup-1.0.0.exe` (only if Inno Setup 6 is installed; skipped
  gracefully otherwise with instructions).

## Winget packaging path

See [`packaging/winget/README.md`](packaging/winget/README.md) for the full flow:
build release -> upload installer to GitHub Releases -> generate SHA256 ->
`wingetcreate` -> submit to `microsoft/winget-pkgs`.

## Future Winget Install

After PatchPrompt is approved for winget, the install command will be:

winget install PatchPrompt.PatchPrompt

## Troubleshooting

- **Icons show as boxes/question marks** — install a Nerd Font and set it as the
  terminal font, or use `fontMode: "plain"` in config.
- **`patchprompt` not found** — reopen the terminal so the updated `PATH` loads,
  or run `patchprompt doctor`.
- **No colors** — use Windows Terminal (full ANSI support); legacy consoles may
  need `fontMode: "plain"`.
- **Prompt looks wrong / want to debug** — run `patchprompt doctor`.

## Uninstall

```powershell
.\scripts\uninstall.ps1
```

Removes the profile block, the `PATH` entry, and the installed binary. It asks
before deleting your config/themes and never touches unrelated files. Use
`-RemoveConfig` to also remove `%USERPROFILE%\.patchprompt` without prompting.

## Repo layout

```
cmd/patchprompt   CLI entrypoint
internal/ansi     color/escape sequences
internal/config   config + theme loading
internal/render   prompt rendering engine
internal/segments path, git, user, host, time, status, admin, custom, runtime
internal/shell    PowerShell + cmd integration code
themes            built-in themes (embedded + copied on install)
scripts           install / uninstall / build-release / cmd helper
packaging         Inno Setup (.iss) + winget docs
examples          sample profile and cmd AutoRun
```

## License

MIT.
