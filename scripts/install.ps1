<#
.SYNOPSIS
    Builds and installs PatchPrompt locally for the current user. No admin needed.
#>
[CmdletBinding()]
param(
    [switch]$UpdateProfile,
    [string]$Theme = "neon"
)

$ErrorActionPreference = "Stop"
$RepoRoot   = Split-Path -Parent $PSScriptRoot
$InstallDir = Join-Path $env:LOCALAPPDATA "Programs\PatchPrompt"
$ConfigDir  = Join-Path $env:USERPROFILE ".patchprompt"
$ThemesDir  = Join-Path $ConfigDir "themes"
$Exe        = Join-Path $InstallDir "patchprompt.exe"
$ProfileLine = "patchprompt init pwsh --theme $Theme | Invoke-Expression"

Write-Host "PatchPrompt local installer" -ForegroundColor Cyan
Write-Host ("-" * 40)

# 1. Build the binary (requires Go).
$built = $false
if (Get-Command go -ErrorAction SilentlyContinue) {
    Write-Host "Building patchprompt.exe with Go..."
    Push-Location $RepoRoot
    try {
        $env:GOOS = "windows"; $env:GOARCH = "amd64"
        & go build -trimpath -ldflags "-s -w" -o "$env:TEMP\patchprompt.exe" ./cmd/patchprompt
        if ($LASTEXITCODE -ne 0) { throw "go build failed" }
        $built = $true
    } finally { Pop-Location }
} else {
    Write-Warning "Go not found. Will use a prebuilt patchprompt.exe if present."
}

# 2. Create install dir and copy the binary.
New-Item -ItemType Directory -Force -Path $InstallDir | Out-Null
if ($built) {
    Copy-Item "$env:TEMP\patchprompt.exe" $Exe -Force
} elseif (Test-Path (Join-Path $RepoRoot "dist\patchprompt.exe")) {
    Copy-Item (Join-Path $RepoRoot "dist\patchprompt.exe") $Exe -Force
} elseif (Test-Path (Join-Path $RepoRoot "patchprompt.exe")) {
    Copy-Item (Join-Path $RepoRoot "patchprompt.exe") $Exe -Force
} else {
    throw "No patchprompt.exe available. Install Go and re-run, or build-release.ps1 first."
}
Write-Host "Installed binary -> $Exe" -ForegroundColor Green

# 3. Create config + copy themes.
New-Item -ItemType Directory -Force -Path $ThemesDir | Out-Null
Copy-Item (Join-Path $RepoRoot "themes\*.json") $ThemesDir -Force
& $Exe config generate | Out-Null
& $Exe theme set $Theme | Out-Null
Write-Host "Config + themes -> $ConfigDir" -ForegroundColor Green

# 4. Add install dir to the USER PATH if missing.
$userPath = [Environment]::GetEnvironmentVariable("Path", "User")
if (-not $userPath) { $userPath = "" }
if (($userPath -split ';') -notcontains $InstallDir) {
    $newPath = if ($userPath.TrimEnd(';')) { "$($userPath.TrimEnd(';'));$InstallDir" } else { $InstallDir }
    [Environment]::SetEnvironmentVariable("Path", $newPath, "User")
    $env:Path = "$env:Path;$InstallDir"
    Write-Host "Added $InstallDir to your user PATH (restart terminals to pick it up)." -ForegroundColor Green
} else {
    Write-Host "Install dir already on user PATH." -ForegroundColor DarkGray
}

# 5. Print the activation line.
Write-Host ("-" * 40)
Write-Host "Add this line to your PowerShell profile (notepad `$PROFILE):" -ForegroundColor Yellow
Write-Host "  $ProfileLine" -ForegroundColor White

# 6. Optional automatic profile update.
$doUpdate = $UpdateProfile
if (-not $doUpdate) {
    $ans = Read-Host "Add it automatically to your `$PROFILE now? (y/N)"
    $doUpdate = ($ans -match '^(y|yes)$')
}
if ($doUpdate) {
    $profilePath = $PROFILE
    $profileParent = Split-Path -Parent $profilePath
    if (-not (Test-Path $profileParent)) { New-Item -ItemType Directory -Force -Path $profileParent | Out-Null }
    if (-not (Test-Path $profilePath)) { New-Item -ItemType File -Force -Path $profilePath | Out-Null }
    $existing = Get-Content $profilePath -Raw -ErrorAction SilentlyContinue
    if ($existing -notmatch 'patchprompt init pwsh') {
        Add-Content $profilePath "`n# >>> PatchPrompt >>>`n$ProfileLine`n# <<< PatchPrompt <<<"
        Write-Host "Profile updated: $profilePath" -ForegroundColor Green
    } else {
        Write-Host "Profile already contains a PatchPrompt line; left unchanged." -ForegroundColor DarkGray
    }
}

Write-Host ("-" * 40)
Write-Host "Done. Reopen Windows Terminal. BAM COLOR." -ForegroundColor Magenta
