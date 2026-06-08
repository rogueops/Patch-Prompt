package shell

import (
	"fmt"
	"strings"
)

// Cmd returns a cmd.exe PROMPT command. cmd cannot invoke an external program
// on every prompt the way PowerShell can, so PatchPrompt emits a static,
// ANSI-colored PROMPT string (path + ">"). This is intentionally limited; see
// the README "Command Prompt" section.
//
// $E is the cmd escape token (Windows 10+). $P is the path, $G is ">".
func Cmd(theme string) string {
	// Color the path cyan and the ">" magenta. These map to the general feel of
	// the neon/default themes without copying any third-party theme data.
	pathColor := "36"
	arrowColor := "35"
	switch strings.ToLower(strings.TrimSpace(theme)) {
	case "stealth", "minimal":
		pathColor, arrowColor = "90", "37"
	case "rainbow":
		pathColor, arrowColor = "32", "33"
	}
	return fmt.Sprintf("PROMPT $E[%sm$P$E[0m $E[%sm$G$E[0m $S", pathColor, arrowColor)
}
