<#
.SYNOPSIS
    Builds Windows release artifacts for PatchPrompt into dist/.
    Produces: patchprompt-windows-amd64.zip, checksums.txt, and (if Inno Setup
    is installed) PatchPromptSetup-1.0.0.exe.
#>
[CmdletBinding()]
param([string]$Version = "1.0.0")

$ErrorActionPreference = "Stop"
$RepoRoot = Split-Path -Parent $PSScriptRoot
$Dist     = Join-Path $RepoRoot "dist"
$Staging  = Join-Path $Dist "patchprompt-windows-amd64"

Write-Host "PatchPrompt release builder v$Version" -ForegroundColor Cyan
Write-Host ("-" * 40)

if (-not (Get-Command go -ErrorAction SilentlyContinue)) {
    throw "Go is required to build a release. Install Go and re-run."
}

Push-Location $RepoRoot
try {
    # 1. gofmt.
    Write-Host "Running gofmt..."
    $unformatted = & gofmt -l .
    if ($unformatted) {
        Write-Warning "These files are not gofmt-clean:`n$unformatted"
        & gofmt -w .
        Write-Host "Reformatted." -ForegroundColor Green
    }

    # 2. Tests (only if any *_test.go exist).
    $hasTests = Get-ChildItem -Recurse -Filter *_test.go -ErrorAction SilentlyContinue
    if ($hasTests) {
        Write-Host "Running go test ./..."
        & go test ./...
        if ($LASTEXITCODE -ne 0) { throw "tests failed" }
    } else {
        Write-Host "No tests found; skipping." -ForegroundColor DarkGray
    }

    # 3. Build Windows amd64 binary.
    Write-Host "Building windows/amd64 binary..."
    if (Test-Path $Dist) { Remove-Item $Dist -Recurse -Force }
    New-Item -ItemType Directory -Force -Path $Staging | Out-Null
    $env:GOOS = "windows"; $env:GOARCH = "amd64"
    & go build -trimpath -ldflags "-s -w" -o (Join-Path $Staging "patchprompt.exe") ./cmd/patchprompt
    if ($LASTEXITCODE -ne 0) { throw "go build failed" }
    Copy-Item (Join-Path $Staging "patchprompt.exe") (Join-Path $Dist "patchprompt.exe") -Force

    # 4. Stage themes, README, scripts.
    Copy-Item (Join-Path $RepoRoot "themes") (Join-Path $Staging "themes") -Recurse -Force
    Copy-Item (Join-Path $RepoRoot "README.md") $Staging -Force
    New-Item -ItemType Directory -Force -Path (Join-Path $Staging "scripts") | Out-Null
    Copy-Item (Join-Path $RepoRoot "scripts\*.ps1") (Join-Path $Staging "scripts") -Force
    Copy-Item (Join-Path $RepoRoot "scripts\patchprompt.cmd") (Join-Path $Staging "scripts") -Force

    # 5. Zip.
    $Zip = Join-Path $Dist "patchprompt-windows-amd64.zip"
    Write-Host "Creating $Zip..."
    Compress-Archive -Path (Join-Path $Staging "*") -DestinationPath $Zip -Force

    # 6. Checksums.
    $checks = Join-Path $Dist "checksums.txt"
    Get-FileHash $Zip -Algorithm SHA256 |
        ForEach-Object { "{0}  {1}" -f $_.Hash.ToLower(), (Split-Path $_.Path -Leaf) } |
        Set-Content $checks
    Write-Host "Wrote checksums.txt" -ForegroundColor Green

    # 7. Optional Inno Setup installer.
    $iscc = $null
    foreach ($p in @(
        "${env:ProgramFiles(x86)}\Inno Setup 6\ISCC.exe",
        "$env:ProgramFiles\Inno Setup 6\ISCC.exe")) {
        if (Test-Path $p) { $iscc = $p; break }
    }
    if (-not $iscc) { $iscc = (Get-Command ISCC.exe -ErrorAction SilentlyContinue).Source }

    if ($iscc) {
        Write-Host "Inno Setup found; building installer..."
        & $iscc "/DAppVersion=$Version" (Join-Path $RepoRoot "packaging\inno\PatchPrompt.iss")
        if ($LASTEXITCODE -ne 0) { Write-Warning "Inno Setup build failed; continuing." }
        else {
            Get-FileHash (Join-Path $Dist "PatchPromptSetup-$Version.exe") -Algorithm SHA256 -ErrorAction SilentlyContinue |
                ForEach-Object { Add-Content $checks ("{0}  {1}" -f $_.Hash.ToLower(), (Split-Path $_.Path -Leaf)) }
            Write-Host "Built PatchPromptSetup-$Version.exe" -ForegroundColor Green
        }
    } else {
        Write-Warning "Inno Setup (ISCC.exe) not found - skipping installer."
        Write-Host "  To build the installer: install Inno Setup 6 from https://jrsoftware.org/isdl.php" -ForegroundColor DarkGray
        Write-Host "  then re-run this script, or compile packaging\inno\PatchPrompt.iss manually." -ForegroundColor DarkGray
    }

    Write-Host ("-" * 40)
    Write-Host "Release artifacts in: $Dist" -ForegroundColor Green
    Get-ChildItem $Dist | Select-Object Name, Length | Format-Table -AutoSize
}
finally { Pop-Location }
