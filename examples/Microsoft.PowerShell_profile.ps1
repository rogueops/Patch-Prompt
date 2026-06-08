# Example PowerShell profile for PatchPrompt.
# Copy the line below into your real profile: notepad $PROFILE

# >>> PatchPrompt >>>
patchprompt init pwsh --theme neon | Invoke-Expression
# <<< PatchPrompt <<<

# Tip: change the theme by editing --theme above, or run:
#   patchprompt theme set minimal
# Available themes: default, neon, minimal, rainbow, stealth
