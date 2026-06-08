package segments

import (
	"context"
	"os/exec"
	"strings"
	"time"
)

// runGit runs a git subcommand with a short timeout, returning trimmed stdout.
func runGit(args ...string) (string, bool) {
	ctx, cancel := context.WithTimeout(context.Background(), 800*time.Millisecond)
	defer cancel()
	out, err := exec.CommandContext(ctx, "git", args...).Output()
	if err != nil {
		return "", false
	}
	return strings.TrimSpace(string(out)), true
}

// Git returns the current branch with a dirty indicator. Fails silently if git
// is missing or the directory is not a repository.
func Git(_ Context) (string, bool) {
	if _, err := exec.LookPath("git"); err != nil {
		return "", false
	}
	if inside, ok := runGit("rev-parse", "--is-inside-work-tree"); !ok || inside != "true" {
		return "", false
	}
	branch, ok := runGit("rev-parse", "--abbrev-ref", "HEAD")
	if !ok || branch == "" {
		return "", false
	}
	if branch == "HEAD" {
		if short, ok := runGit("rev-parse", "--short", "HEAD"); ok {
			branch = short
		}
	}
	if status, ok := runGit("status", "--porcelain"); ok && status != "" {
		branch += "*"
	}
	return branch, true
}
