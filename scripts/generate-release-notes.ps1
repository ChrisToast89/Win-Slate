# Expand docs/RELEASE_NOTES_TEMPLATE.md with VERSION and optional changelog.
param(
    [string]$Version = "",
    [string]$Changelog = "See commit history for this tag.",
    [string]$OutFile = ""
)

$ErrorActionPreference = "Stop"
$root = Split-Path -Parent $PSScriptRoot

if (-not $Version) {
    $Version = (Get-Content (Join-Path $root "VERSION") -Raw).Trim()
}
$templatePath = Join-Path $root "docs\RELEASE_NOTES_TEMPLATE.md"
if (-not (Test-Path $templatePath)) {
    throw "Missing template: $templatePath"
}

$body = Get-Content $templatePath -Raw
$body = $body.Replace("{{VERSION}}", $Version).Replace("{{CHANGELOG}}", $Changelog)

if (-not $OutFile) {
    $OutFile = Join-Path $root "distributable\RELEASE_NOTES-v$Version.md"
}
$dir = Split-Path -Parent $OutFile
if ($dir) {
    New-Item -ItemType Directory -Force -Path $dir | Out-Null
}
Set-Content -Path $OutFile -Value $body -Encoding utf8
Write-Host "Wrote $OutFile"
return $OutFile
