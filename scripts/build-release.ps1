# Build Slate for Windows app + Setup installer, stage dist package.
# App binary location: folder root only (app\Slate.exe or sibling slate-windows\Slate.exe).
# Never ships a build\bin duplicate of Slate.exe.
# Requires: Go, Wails, Node/npm
$ErrorActionPreference = "Stop"
$root = Split-Path -Parent $PSScriptRoot
Set-Location $root

$version = "0.3.2-win.1"
# End-user package only (not the full dev tree) lives under distributable/
$dist = Join-Path $root "distributable"
$stage = Join-Path $dist "SlateForWindows-v$version"
$payloadDir = Join-Path $root "setup\payload"
$payloadExe = Join-Path $payloadDir "Slate.exe"
$appRootExe = Join-Path $root "app\Slate.exe"
# Sibling checkout: Slate-win\slate-windows\Slate.exe
$siblingExe = Join-Path (Split-Path -Parent $root) "slate-windows\Slate.exe"

New-Item -ItemType Directory -Force -Path $dist, $stage, $payloadDir | Out-Null

function Resolve-AppBinary {
    # Prefer already-built root binaries (no build\bin).
    foreach ($cand in @($siblingExe, $appRootExe)) {
        if (Test-Path $cand) {
            Write-Host "Using app binary: $cand" -ForegroundColor Cyan
            return (Resolve-Path $cand).Path
        }
    }
    Write-Host "==> Building app (root Slate.exe only)…" -ForegroundColor Cyan
    $appDir = Join-Path $root "app"
    $buildPs1 = Join-Path $appDir "scripts\build.ps1"
    if (Test-Path $buildPs1) {
        Push-Location $appDir
        try {
            & $buildPs1
            if ($LASTEXITCODE -ne 0 -and $null -ne $LASTEXITCODE) { throw "app build.ps1 failed" }
        } finally {
            Pop-Location
        }
    } else {
        Push-Location $appDir
        try {
            wails build
            if ($LASTEXITCODE -ne 0) { throw "app wails build failed" }
            $wailsOut = Join-Path $appDir "build\bin\Slate.exe"
            if (-not (Test-Path $wailsOut)) { throw "missing $wailsOut" }
            Move-Item -Force $wailsOut $appRootExe
        } finally {
            Pop-Location
        }
    }
    if (-not (Test-Path $appRootExe)) {
        throw "App binary missing after build: $appRootExe"
    }
    return (Resolve-Path $appRootExe).Path
}

$appExe = Resolve-AppBinary
# Stage for embed + portable package (single source of truth path for Setup).
Copy-Item -Force $appExe $payloadExe
Copy-Item -Force $appExe $appRootExe -ErrorAction SilentlyContinue
Copy-Item -Force $appExe (Join-Path $stage "Slate.exe")
Write-Host "Payload ready: $payloadExe" -ForegroundColor Green

Write-Host "==> Building Setup (embeds setup\payload\Slate.exe)…" -ForegroundColor Cyan
Set-Location (Join-Path $root "setup")
wails build
if ($LASTEXITCODE -ne 0) { throw "setup wails build failed" }

$setupWailsOut = Join-Path $root "setup\build\bin\SlateForWindows-Setup.exe"
$setupRepoRoot = Join-Path $root "SlateForWindows-Setup.exe"
if (-not (Test-Path $setupWailsOut)) { throw "missing $setupWailsOut" }
# Tracked at repo root for end users + copy into distributable package.
Copy-Item -Force $setupWailsOut $setupRepoRoot
Copy-Item -Force $setupWailsOut (Join-Path $stage "SlateForWindows-Setup.exe")

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
Write-Host "    App source: $appExe" -ForegroundColor Green
Write-Host "    Payload:    $payloadExe" -ForegroundColor Green
Write-Host "    Package:    $stage" -ForegroundColor Green
