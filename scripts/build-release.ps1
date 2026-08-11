# Build Win-Slate app + Setup installer, stage distributable package.
# App binary: folder root only (app\Win-Slate.exe or sibling slate-windows\Win-Slate.exe).
# Requires: Go, Wails, Node/npm
$ErrorActionPreference = "Stop"
$root = Split-Path -Parent $PSScriptRoot
Set-Location $root

$versionFile = Join-Path $root "VERSION"
$version = if (Test-Path $versionFile) { (Get-Content $versionFile -Raw).Trim() } else { "1.0.0" }
$dist = Join-Path $root "distributable"
$stage = Join-Path $dist "Win-Slate-v$version"
$payloadDir = Join-Path $root "setup\payload"
$payloadExe = Join-Path $payloadDir "Win-Slate.exe"
$appRootExe = Join-Path $root "app\Win-Slate.exe"
$siblingExe = Join-Path (Split-Path -Parent $root) "slate-windows\Win-Slate.exe"
$siblingLegacy = Join-Path (Split-Path -Parent $root) "slate-windows\Slate.exe"

New-Item -ItemType Directory -Force -Path $dist, $stage, $payloadDir | Out-Null

function Resolve-AppBinary {
    # Prefer monorepo app build; sibling trees are fallback only.
    foreach ($cand in @($appRootExe, $siblingExe, (Join-Path $root "app\Slate.exe"), $siblingLegacy)) {
        if (Test-Path -LiteralPath $cand) {
            Write-Host "Using app binary: $cand" -ForegroundColor Cyan
            return (Resolve-Path -LiteralPath $cand).Path
        }
    }
    Write-Host "==> Building app (root Win-Slate.exe only)…" -ForegroundColor Cyan
    $appDir = Join-Path $root "app"
    $buildPs1 = Join-Path $appDir "scripts\build.ps1"
    Push-Location $appDir
    try {
        if (Test-Path -LiteralPath $buildPs1) {
            # External tools (wails) write to the success stream. If that output
            # is not consumed, it becomes part of this function's return value and
            # $appExe can end up empty/array — breaking Copy-Item on CI.
            & $buildPs1 *>&1 | ForEach-Object { Write-Host $_ }
            if ($LASTEXITCODE -ne 0 -and $null -ne $LASTEXITCODE) { throw "app build.ps1 failed" }
        } else {
            wails build *>&1 | ForEach-Object { Write-Host $_ }
            if ($LASTEXITCODE -ne 0) { throw "app wails build failed" }
            $wailsOut = Join-Path $appDir "build\bin\Win-Slate.exe"
            if (-not (Test-Path -LiteralPath $wailsOut)) {
                $wailsOut = Join-Path $appDir "build\bin\Slate.exe"
            }
            if (-not (Test-Path -LiteralPath $wailsOut)) { throw "missing wails output binary" }
            Move-Item -Force $wailsOut $appRootExe
        }
    } finally {
        Pop-Location
    }
    if (-not (Test-Path -LiteralPath $appRootExe)) {
        throw "App binary missing after build: $appRootExe"
    }
    return (Resolve-Path -LiteralPath $appRootExe).Path
}

$appExe = Resolve-AppBinary
# Defensive: if multiple pipeline objects leaked, keep the last path-like string.
if ($appExe -is [array]) {
    $appExe = @($appExe | Where-Object { $_ -is [string] -and $_ -and (Test-Path -LiteralPath $_) })[-1]
}
if (-not $appExe -or -not (Test-Path -LiteralPath $appExe)) {
    throw "Resolve-AppBinary returned empty or missing path: '$appExe' (expected $appRootExe)"
}
Write-Host "App binary resolved: $appExe" -ForegroundColor Cyan
Copy-Item -Force -LiteralPath $appExe -Destination $payloadExe
Copy-Item -Force -LiteralPath $appExe -Destination $appRootExe -ErrorAction SilentlyContinue
Copy-Item -Force -LiteralPath $appExe -Destination (Join-Path $stage "Win-Slate.exe")
Write-Host "Payload ready: $payloadExe" -ForegroundColor Green

Write-Host "==> Building Setup (embeds setup\payload\Win-Slate.exe)…" -ForegroundColor Cyan
Set-Location (Join-Path $root "setup")
wails build *>&1 | ForEach-Object { Write-Host $_ }
if ($LASTEXITCODE -ne 0) { throw "setup wails build failed" }

$setupWailsOut = Join-Path $root "setup\build\bin\Win-Slate-Setup.exe"
$setupRepoRoot = Join-Path $root "Win-Slate-Setup.exe"
$setupZipRoot = Join-Path $root "Win-Slate-Setup.zip"
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

# Repo-root download zip (Setup only) — preferred end-user download vs raw .exe
$setupOnlyDir = Join-Path $root "download\Win-Slate-Setup"
New-Item -ItemType Directory -Force -Path $setupOnlyDir | Out-Null
Copy-Item -Force $setupRepoRoot (Join-Path $setupOnlyDir "Win-Slate-Setup.exe")
Copy-Item -Force (Join-Path $root "LICENSE") (Join-Path $setupOnlyDir "LICENSE.txt")
Copy-Item -Force (Join-Path $root "NOTICE") (Join-Path $setupOnlyDir "NOTICE.txt")
Copy-Item -Force (Join-Path $stage "INSTALL.txt") (Join-Path $setupOnlyDir "INSTALL.txt")
if (Test-Path $setupZipRoot) { Remove-Item -Force $setupZipRoot }
Compress-Archive -Path (Join-Path $setupOnlyDir "*") -DestinationPath $setupZipRoot -Force

$zip = Join-Path $dist "Win-Slate-v$version.zip"
if (Test-Path $zip) { Remove-Item -Force $zip }
Compress-Archive -Path (Join-Path $stage "*") -DestinationPath $zip -Force

Write-Host "OK: $zip" -ForegroundColor Green
Write-Host "    App source:  $appExe" -ForegroundColor Green
Write-Host "    Payload:     $payloadExe" -ForegroundColor Green
Write-Host "    Setup exe:   $setupRepoRoot (local build only)" -ForegroundColor Green
Write-Host "    Setup zip:   $setupZipRoot (repo download)" -ForegroundColor Green
Write-Host "    Full package:$stage" -ForegroundColor Green
