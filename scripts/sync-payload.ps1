# Copy the current app binary into setup\payload\Win-Slate.exe for Setup embed.
# Sources: app\Win-Slate.exe, then app\Slate.exe alias only (no archived port trees).
$ErrorActionPreference = "Stop"
$root = Split-Path -Parent $PSScriptRoot
$payloadDir = Join-Path $root "setup\payload"
$payloadExe = Join-Path $payloadDir "Win-Slate.exe"
$cands = @(
    (Join-Path $root "app\Win-Slate.exe"),
    (Join-Path $root "app\Slate.exe")
)

$src = $null
foreach ($cand in $cands) {
    if (Test-Path -LiteralPath $cand) {
        $src = (Resolve-Path -LiteralPath $cand).Path
        break
    }
}
if (-not $src) {
    throw "No app binary found. Build app with app\scripts\build.ps1 first."
}

New-Item -ItemType Directory -Force -Path $payloadDir | Out-Null
Copy-Item -Force $src $payloadExe
$appRoot = Join-Path $root "app\Win-Slate.exe"
if ($src -ne $appRoot) {
    New-Item -ItemType Directory -Force -Path (Split-Path $appRoot) | Out-Null
    Copy-Item -Force $src $appRoot
}
Write-Host "Synced payload from:`n  $src`n→ $payloadExe" -ForegroundColor Green
