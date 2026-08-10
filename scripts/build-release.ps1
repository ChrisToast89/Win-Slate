# Build Win-Slate app + Setup installer, stage distributable package.
# App binary: folder root only (app\Win-Slate.exe or sibling slate-windows\Win-Slate.exe).
# Requires: Go, Wails, Node/npm
$ErrorActionPreference = "Stop"
$root = Split-Path -Parent $PSScriptRoot
Set-Location $root

$versionFile = Join-Path $root "VERSION"
$version = if (Test-Path $versionFile) { (Get-Content $versionFile -Raw).Trim() } else { "0.3.2-win.1" }
$dist = Join-Path $root "distributable"
$stage = Join-Path $dist "Win-Slate-v$version"
$payloadDir = Join-Path $root "setup\payload"
$payloadExe = Join-Path $payloadDir "Win-Slate.exe"
$appRootExe = Join-Path $root "app\Win-Slate.exe"
$siblingExe = Join-Path (Split-Path -Parent $root) "slate-windows\Win-Slate.exe"
$siblingLegacy = Join-Path (Split-Path -Parent $root) "slate-windows\Slate.exe"

New-Item -ItemType Directory -Force -Path $dist, $stage, $payloadDir | Out-Null

function Resolve-AppBinary {
    foreach ($cand in @($siblingExe, $appRootExe, $siblingLegacy, (Join-Path $root "app\Slate.exe"))) {
        if (Test-Path $cand) {
            Write-Host "Using app binary: $cand" -ForegroundColor Cyan
            return (Resolve-Path $cand).Path
        }
    }
    Write-Host "==> Building app (root Win-Slate.exe only)…" -ForegroundColor Cyan
    $appDir = Join-Path $root "app"
    $buildPs1 = Join-Path $appDir "scripts\build.ps1"
    Push-Location $appDir
    try {
        if (Test-Path $buildPs1) {
            & $buildPs1
            if ($LASTEXITCODE -ne 0 -and $null -ne $LASTEXITCODE) { throw "app build.ps1 failed" }
        } else {
            wails build
            if ($LASTEXITCODE -ne 0) { throw "app wails build failed" }
            $wailsOut = Join-Path $appDir "build\bin\Win-Slate.exe"
            if (-not (Test-Path $wailsOut)) {
                $wailsOut = Join-Path $appDir "build\bin\Slate.exe"
            }
            if (-not (Test-Path $wailsOut)) { throw "missing wails output binary" }
            Move-Item -Force $wailsOut $appRootExe
        }
    } finally {
        Pop-Location
    }
    if (-not (Test-Path $appRootExe)) {
        throw "App binary missing after build: $appRootExe"
    }
    return (Resolve-Path $appRootExe).Path
}

$appExe = Resolve-AppBinary
Copy-Item -Force $appExe $payloadExe
Copy-Item -Force $appExe $appRootExe -ErrorAction SilentlyContinue
Copy-Item -Force $appExe (Join-Path $stage "Win-Slate.exe")
Write-Host "Payload ready: $payloadExe" -ForegroundColor Green

Write-Host "==> Building Setup (embeds setup\payload\Win-Slate.exe)…" -ForegroundColor Cyan
Set-Location (Join-Path $root "setup")
wails build
if ($LASTEXITCODE -ne 0) { throw "setup wails build failed" }

$setupWailsOut = Join-Path $root "setup\build\bin\Win-Slate-Setup.exe"
$setupRepoRoot = Join-Path $root "Win-Slate-Setup.exe"
if (-not (Test-Path $setupWailsOut)) { throw "missing $setupWailsOut" }
Copy-Item -Force $setupWailsOut $setupRepoRoot
Copy-Item -Force $setupWailsOut (Join-Path $stage "Win-Slate-Setup.exe")

Copy-Item -Force (Join-Path $root "README.md") (Join-Path $stage "README.md")
Copy-Item -Force (Join-Path $root "LICENSE") (Join-Path $stage "LICENSE.txt")
Copy-Item -Force (Join-Path $root "NOTICE") (Join-Path $stage "NOTICE.txt")
if (Test-Path (Join-Path $root "ATTRIBUTION.md")) {
    Copy-Item -Force (Join-Path $root "ATTRIBUTION.md") (Join-Path $stage "ATTRIBUTION.md")
}
@"
Win-Slate v$version
===========================

Unofficial Windows build — derivative of Slate by Sam Wasserman (Apache-2.0).
Not an official Wasserman release. No warranty. See LICENSE.txt, NOTICE.txt, ATTRIBUTION.md.

Win-Slate is a standalone executable. It is separate from the npm/Electron
install of Sam's Slate (usually under Programs\Slate). Both may share
Documents\Slate projects.

1. Run Win-Slate-Setup.exe
2. Check this PC → choose a folder (default Programs\Win-Slate) → Install
3. Optional: Claude Code for the AI brain

Portable option: run Win-Slate.exe directly (WebView2 required).

Upstream: https://github.com/wassermanproductions/slate
This package: https://github.com/ChrisToast89/Win-Slate
"@ | Set-Content -Encoding utf8 (Join-Path $stage "INSTALL.txt")

$zip = Join-Path $dist "Win-Slate-v$version.zip"
if (Test-Path $zip) { Remove-Item -Force $zip }
Compress-Archive -Path (Join-Path $stage "*") -DestinationPath $zip -Force

Write-Host "OK: $zip" -ForegroundColor Green
Write-Host "    App source: $appExe" -ForegroundColor Green
Write-Host "    Payload:    $payloadExe" -ForegroundColor Green
Write-Host "    Setup root: $setupRepoRoot" -ForegroundColor Green
Write-Host "    Package:    $stage" -ForegroundColor Green
