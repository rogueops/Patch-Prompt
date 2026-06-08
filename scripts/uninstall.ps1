<#
.SYNOPSIS
    Removes PatchPrompt: profile block, PATH entry, and optionally config/themes.
    Only touches what the installer added. No admin required.
#>
[CmdletBinding()]
param([switch]$RemoveConfig)

$ErrorActionPreference = "Stop"
$InstallDir = Join-Path $env:LOCALAPPDATA "Programs\PatchPrompt"
$ConfigDir  = Join-Path $env:USERPROFILE ".patchprompt"

Write-Host "PatchPrompt uninstaller" -ForegroundColor Cyan
Write-Host ("-" * 40)

# 1. Remove the PatchPrompt block from the PowerShell profile.
if (Test-Path $PROFILE) {
    $content = Get-Content $PROFILE -Raw
    if ($content -match '# >>> PatchPrompt >>>') {
        $clean = [regex]::Replace($content,
            '(?s)\r?\n?# >>> PatchPrompt >>>.*?# <<< PatchPrompt <<<\r?\n?', "`n")
        Set-Content $PROFILE $clean.TrimEnd() -NoNewline
        Write-Host "Removed PatchPrompt block from $PROFILE" -ForegroundColor Green
    } else {
        Write-Host "No PatchPrompt block found in profile (you may have added the line manually)." -ForegroundColor DarkGray
    }
}

# 2. Remove install dir from user PATH.
$userPath = [Environment]::GetEnvironmentVariable("Path", "User")
if ($userPath -and (($userPath -split ';') -contains $InstallDir)) {
    $newPath = (($userPath -split ';') | Where-Object { $_ -and $_ -ne $InstallDir }) -join ';'
    [Environment]::SetEnvironmentVariable("Path", $newPath, "User")
    Write-Host "Removed $InstallDir from user PATH." -ForegroundColor Green
}

# 3. Remove the installed binary/dir.
if (Test-Path $InstallDir) {
    Remove-Item $InstallDir -Recurse -Force
    Write-Host "Removed $InstallDir" -ForegroundColor Green
}

# 4. Optionally remove config/themes (ask first).
if (Test-Path $ConfigDir) {
    $remove = $RemoveConfig
    if (-not $remove) {
        $ans = Read-Host "Also delete your config and themes at $ConfigDir? (y/N)"
        $remove = ($ans -match '^(y|yes)$')
    }
    if ($remove) {
        Remove-Item $ConfigDir -Recurse -Force
        Write-Host "Removed $ConfigDir" -ForegroundColor Green
    } else {
        Write-Host "Kept your config at $ConfigDir" -ForegroundColor DarkGray
    }
}

Write-Host ("-" * 40)
Write-Host "PatchPrompt uninstalled. Reopen your terminals." -ForegroundColor Magenta
