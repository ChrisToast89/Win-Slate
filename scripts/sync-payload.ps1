# Copy the current app binary into setup\payload\Win-Slate.exe for Setup embed.
# Preference: sibling slate-windows\Win-Slate.exe → app\Win-Slate.exe → legacy Slate.exe names
$ErrorActionPreference = "Stop"
$root = Split-Path -Parent $PSScriptRoot
$payloadDir = Join-Path $root "setup\payload"
$payloadExe = Join-Path $payloadDir "Win-Slate.exe"
$cands = @(
    (Join-Path (Split-Path -Parent $root) "slate-windows\Win-Slate.exe"),
    (Join-Path $root "app\Win-Slate.exe"),
    (Join-Path (Split-Path -Parent $root) "slate-windows\Slate.exe"),
    (Join-Path $root "app\Slate.exe")
)

$src = $null
foreach ($cand in $cands) {
    if (Test-Path $cand) {
        $src = (Resolve-Path $cand).Path
        break
    }
}
if (-not $src) {
    throw "No app binary found. Build app with .\scripts\build.ps1 first."
}

New-Item -ItemType Directory -Force -Path $payloadDir | Out-Null
Copy-Item -Force $src $payloadExe
$appRoot = Join-Path $root "app\Win-Slate.exe"
if ($src -ne $appRoot) {
    New-Item -ItemType Directory -Force -Path (Split-Path $appRoot) | Out-Null
    Copy-Item -Force $src $appRoot
}
Write-Host "Synced payload from:`n  $src`n→ $payloadExe" -ForegroundColor Green
