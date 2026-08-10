# Build (optional), package, and publish a GitHub Release for this derivative work.
#
# Always attaches Apache-2.0 credit notes (ATTRIBUTION / NOTICE / template).
#
# Usage:
#   .\scripts\publish-release.ps1
#   .\scripts\publish-release.ps1 -Version 0.3.2-win.2 -Changelog "Fix install path."
#   .\scripts\publish-release.ps1 -SkipBuild          # reuse existing binaries
#   .\scripts\publish-release.ps1 -DryRun
#
# Prerequisites: Go, Wails, Node (unless -SkipBuild), gh CLI authenticated, git.
param(
    [string]$Version = "",
    [string]$Changelog = "Maintenance release of the Windows port and Setup installer.",
    [switch]$SkipBuild,
    [switch]$SkipRootCommit,
    [switch]$DryRun
)

$ErrorActionPreference = "Stop"
$root = Split-Path -Parent $PSScriptRoot
Set-Location $root

if (-not $Version) {
    $Version = (Get-Content (Join-Path $root "VERSION") -Raw).Trim()
}
if (-not $Version) {
    throw "VERSION is empty — set VERSION file or pass -Version"
}

$tag = if ($Version.StartsWith("v")) { $Version } else { "v$Version" }
$versionBare = $tag.TrimStart("v")
# Keep VERSION file in sync with what we publish
[System.IO.File]::WriteAllText((Join-Path $root "VERSION"), ($versionBare + "`n"))

Write-Host "=== Publish Win-Slate $tag (derivative of Sam Wasserman's Slate) ===" -ForegroundColor Cyan

# --- Build / package ---
if (-not $SkipBuild) {
    Write-Host "==> build-release.ps1" -ForegroundColor Cyan
    & (Join-Path $root "scripts\build-release.ps1")
    if ($LASTEXITCODE -ne 0 -and $null -ne $LASTEXITCODE) {
        throw "build-release.ps1 failed"
    }
} else {
    Write-Host "==> SkipBuild: expecting existing distributable + root Setup" -ForegroundColor Yellow
}

$dist = Join-Path $root "distributable"
$stage = Join-Path $dist "Win-Slate-v$versionBare"
$zip = Join-Path $dist "Win-Slate-v$versionBare.zip"
$setupRoot = Join-Path $root "Win-Slate-Setup.exe"
$setupStage = Join-Path $stage "Win-Slate-Setup.exe"

if (-not (Test-Path $setupRoot) -and (Test-Path $setupStage)) {
    Copy-Item -Force $setupStage $setupRoot
}
if (-not (Test-Path $setupRoot)) {
    throw "Missing $setupRoot — build first or place Setup at repo root"
}
if (-not (Test-Path $zip)) {
    if (-not (Test-Path $stage)) {
        throw "Missing package folder $stage and zip $zip"
    }
    Compress-Archive -Path (Join-Path $stage "*") -DestinationPath $zip -Force
}

# Ensure license artifacts in stage
foreach ($pair in @(
        @{ Src = "LICENSE"; Dst = "LICENSE.txt" },
        @{ Src = "NOTICE"; Dst = "NOTICE.txt" },
        @{ Src = "ATTRIBUTION.md"; Dst = "ATTRIBUTION.md" },
        @{ Src = "README.md"; Dst = "README.md" }
    )) {
    $s = Join-Path $root $pair.Src
    if (Test-Path $s) {
        Copy-Item -Force $s (Join-Path $stage $pair.Dst)
    }
}

# Refresh zip after doc copy
if (Test-Path $zip) { Remove-Item -Force $zip }
Compress-Archive -Path (Join-Path $stage "*") -DestinationPath $zip -Force

$notesFile = & (Join-Path $root "scripts\generate-release-notes.ps1") -Version $versionBare -Changelog $Changelog
if ($notesFile -is [array]) { $notesFile = $notesFile[-1] }

Write-Host "Assets:" -ForegroundColor Green
Write-Host "  Setup: $setupRoot"
Write-Host "  Zip:   $zip"
Write-Host "  Notes: $notesFile"

if ($DryRun) {
    Write-Host "DryRun: would create tag $tag and gh release with Setup + zip + attribution notes." -ForegroundColor Yellow
    exit 0
}

# --- Git: optional commit of root Setup.exe ---
if (-not $SkipRootCommit) {
    git add VERSION Win-Slate-Setup.exe ATTRIBUTION.md docs/RELEASE_NOTES_TEMPLATE.md 2>$null
    git add scripts/*.ps1 .github 2>$null
    $status = git status --porcelain
    if ($status) {
        git add -A -- VERSION Win-Slate-Setup.exe
        # Only stage publish tooling if present
        foreach ($p in @(
                "ATTRIBUTION.md", "VERSION", "NOTICE", "README.md",
                "scripts/build-release.ps1", "scripts/publish-release.ps1",
                "scripts/generate-release-notes.ps1", "docs/RELEASE_NOTES_TEMPLATE.md",
                ".github/workflows/release.yml", ".github/workflows/publish.yml"
            )) {
            if (Test-Path (Join-Path $root $p)) { git add $p }
        }
        git commit -m "Release $tag — Win-Slate (Apache-2.0 derivative of Sam Wasserman's Slate)" 2>&1
        if ($LASTEXITCODE -eq 0) {
            git push origin HEAD 2>&1
        }
    }
}

# --- Tag ---
$existing = git tag -l $tag
if ($existing) {
    Write-Host "Tag $tag already exists locally." -ForegroundColor Yellow
} else {
    git tag -a $tag -m "Win-Slate $tag — unofficial Windows derivative of Slate by Sam Wasserman (Apache-2.0)"
}
git push origin $tag 2>&1

# --- GitHub Release ---
$relExists = gh release view $tag --repo ChrisToast89/Win-Slate 2>$null
if ($LASTEXITCODE -eq 0) {
    Write-Host "Release $tag exists — uploading/clobbering assets…" -ForegroundColor Yellow
    gh release upload $tag $setupRoot $zip --clobber --repo ChrisToast89/Win-Slate
    gh release edit $tag --notes-file $notesFile --repo ChrisToast89/Win-Slate
} else {
    gh release create $tag $setupRoot $zip `
        --title "Win-Slate $tag" `
        --notes-file $notesFile `
        --repo ChrisToast89/Win-Slate
}

Write-Host ""
Write-Host "Published: https://github.com/ChrisToast89/Win-Slate/releases/tag/$tag" -ForegroundColor Green
Write-Host "Credit block included in release notes (Sam Wasserman / Apache-2.0 / unofficial)." -ForegroundColor Green
