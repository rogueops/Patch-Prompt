// Package shell generates shell integration code for PatchPrompt.
package shell

import (
	"fmt"
	"strings"
)

// PowerShell returns PowerShell code that overrides the prompt function to call
// `patchprompt render`. The generated code preserves $LASTEXITCODE, falls back
// to a default prompt if the binary is missing, and pins the theme if given.
func PowerShell(theme string) string {
	themeArg := ""
	if strings.TrimSpace(theme) != "" {
		themeArg = fmt.Sprintf(" --theme %s", theme)
	}
	return fmt.Sprintf(`# PatchPrompt PowerShell integration
function global:prompt {
    $__pp_code = $LASTEXITCODE
    if ($null -eq $__pp_code) { $__pp_code = 0 }
    $__pp_loc = (Get-Location).Path
    try {
        $__pp_out = & patchprompt render --shell pwsh --status $__pp_code%s 2>$null
        if ($LASTEXITCODE -ne 0 -or [string]::IsNullOrEmpty($__pp_out)) {
            $__pp_out = "PS $__pp_loc> "
        }
    } catch {
        $__pp_out = "PS $__pp_loc> "
    }
    $global:LASTEXITCODE = $__pp_code
    return $__pp_out
}
`, themeArg)
}
