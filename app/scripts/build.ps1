# Production build for Slate (Windows).
# Always use this (or `wails build`) — plain `go build` lacks Wails tags and will fail at launch.
#
# Copies the binary to the repo root so you can run:  .\Slate.exe

$ErrorActionPreference = "Stop"
$root = Split-Path -Parent $PSScriptRoot
Set-Location $root

Write-Host "Building Slate (wails)…" -ForegroundColor Cyan
wails build
if ($LASTEXITCODE -ne 0) {
    throw "wails build failed with exit code $LASTEXITCODE"
}

$src = Join-Path $root "build\bin\Slate.exe"
$dst = Join-Path $root "Slate.exe"
if (-not (Test-Path $src)) {
    throw "Expected binary not found: $src"
}

# Keep a single launchable binary at the folder root (no build\bin duplicate).
Move-Item -Force $src $dst
Write-Host "OK: $dst" -ForegroundColor Green
Write-Host "Launch from folder root:  .\Slate.exe" -ForegroundColor Green
