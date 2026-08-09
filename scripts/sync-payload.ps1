# Copy the current app binary into setup\payload\ for the next Setup embed.
# Preference order:
#   1) ..\slate-windows\Slate.exe  (Slate-win sibling, root only)
#   2) app\Slate.exe
$ErrorActionPreference = "Stop"
$root = Split-Path -Parent $PSScriptRoot
$payloadDir = Join-Path $root "setup\payload"
$payloadExe = Join-Path $payloadDir "Slate.exe"
$siblingExe = Join-Path (Split-Path -Parent $root) "slate-windows\Slate.exe"
$appRootExe = Join-Path $root "app\Slate.exe"

$src = $null
foreach ($cand in @($siblingExe, $appRootExe)) {
    if (Test-Path $cand) {
        $src = (Resolve-Path $cand).Path
        break
    }
}
if (-not $src) {
    throw "No app binary found. Build slate-windows (.\scripts\build.ps1) or app\Slate.exe first."
}

New-Item -ItemType Directory -Force -Path $payloadDir | Out-Null
Copy-Item -Force $src $payloadExe
# Keep monorepo app root in sync when source was the sibling tree.
if ($src -ne $appRootExe) {
    New-Item -ItemType Directory -Force -Path (Split-Path $appRootExe) | Out-Null
    Copy-Item -Force $src $appRootExe
}
Write-Host "Synced payload from:`n  $src`n→ $payloadExe" -ForegroundColor Green
