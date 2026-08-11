# Production build for Slate (Windows).
# Always use this (or `wails build`) — plain `go build` lacks Wails tags and will fail at launch.
#
# Moves the binary to the app folder root: .\Win-Slate.exe

$ErrorActionPreference = "Stop"
$root = Split-Path -Parent $PSScriptRoot
Set-Location $root

Write-Host "Building Win-Slate (wails)…" -ForegroundColor Cyan
# Keep wails on the host stream so callers that assign script output don't get pollution.
wails build *>&1 | ForEach-Object { Write-Host $_ }
if ($LASTEXITCODE -ne 0) {
    throw "wails build failed with exit code $LASTEXITCODE"
}

$src = Join-Path $root "build\bin\Win-Slate.exe"
if (-not (Test-Path $src)) {
    # transitional name if wails.json not yet applied
    $alt = Join-Path $root "build\bin\Slate.exe"
    if (Test-Path $alt) { $src = $alt }
}
$dst = Join-Path $root "Win-Slate.exe"
if (-not (Test-Path $src)) {
    throw "Expected binary not found under build\bin"
}

Move-Item -Force $src $dst
Write-Host "OK: $dst" -ForegroundColor Green
Write-Host "Launch from folder root:  .\Win-Slate.exe" -ForegroundColor Green
