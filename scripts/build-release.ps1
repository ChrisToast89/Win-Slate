# Build Slate for Windows app + Setup installer, stage dist package.
# Requires: Go, Wails, Node/npm
$ErrorActionPreference = "Stop"
$root = Split-Path -Parent $PSScriptRoot
Set-Location $root

$version = "0.3.2-win.1"
$dist = Join-Path $root "dist"
$stage = Join-Path $dist "SlateForWindows-v$version"
New-Item -ItemType Directory -Force -Path $dist, $stage, (Join-Path $root "setup\payload") | Out-Null

Write-Host "==> Building app (wails)…" -ForegroundColor Cyan
Set-Location (Join-Path $root "app")
wails build
if ($LASTEXITCODE -ne 0) { throw "app wails build failed" }
$appExe = Join-Path $root "app\build\bin\Slate.exe"
if (-not (Test-Path $appExe)) { throw "missing $appExe" }
Copy-Item -Force $appExe (Join-Path $root "app\Slate.exe")
Copy-Item -Force $appExe (Join-Path $root "setup\payload\Slate.exe")
Copy-Item -Force $appExe (Join-Path $stage "Slate.exe")

Write-Host "==> Building Setup (wails)…" -ForegroundColor Cyan
Set-Location (Join-Path $root "setup")
wails build
if ($LASTEXITCODE -ne 0) { throw "setup wails build failed" }
$setupExe = Join-Path $root "setup\build\bin\SlateForWindows-Setup.exe"
if (-not (Test-Path $setupExe)) { throw "missing $setupExe" }
Copy-Item -Force $setupExe (Join-Path $stage "SlateForWindows-Setup.exe")
Copy-Item -Force $setupExe (Join-Path $root "SlateForWindows-Setup.exe")

# Package docs
Copy-Item -Force (Join-Path $root "README.md") (Join-Path $stage "README.md")
Copy-Item -Force (Join-Path $root "LICENSE") (Join-Path $stage "LICENSE.txt")
Copy-Item -Force (Join-Path $root "NOTICE") (Join-Path $stage "NOTICE.txt")
@"
Slate for Windows v$version
===========================

1. Run SlateForWindows-Setup.exe
2. Check this PC → choose install folder → Install
3. Optional: Install / sign in to Claude Code for the AI brain

Portable option: run Slate.exe directly (WebView2 required).

Slate by Sam Wasserman — https://github.com/wassermanproductions/slate
This package is a Windows port + installer (not an official Wasserman release).
Projects: %USERPROFILE%\Documents\Slate (never deleted by Setup)
"@ | Set-Content -Encoding utf8 (Join-Path $stage "INSTALL.txt")

$zip = Join-Path $dist "SlateForWindows-v$version.zip"
if (Test-Path $zip) { Remove-Item -Force $zip }
Compress-Archive -Path (Join-Path $stage "*") -DestinationPath $zip -Force

Write-Host "OK: $zip" -ForegroundColor Green
Write-Host "    Setup: $setupExe" -ForegroundColor Green
Write-Host "    App:   $appExe" -ForegroundColor Green
