# Smoke-test the live control server / MCP tool surface.
# Prerequisite: Slate.exe already running.
$ErrorActionPreference = "Stop"
$desc = Join-Path $env:APPDATA "slate\control.json"
if (-not (Test-Path $desc)) { throw "control.json missing — start Slate.exe first" }
$meta = Get-Content $desc -Raw | ConvertFrom-Json
$headers = @{ Authorization = "Bearer $($meta.token)" }

function Invoke-Tool([string]$Tool, $ArgsObj) {
  $body = @{ tool = $Tool; args = $ArgsObj } | ConvertTo-Json -Depth 8 -Compress
  Invoke-RestMethod -Uri "http://127.0.0.1:$($meta.port)/invoke" -Method POST -Headers $headers -ContentType "application/json" -Body $body
}

$tools = Invoke-RestMethod -Uri "http://127.0.0.1:$($meta.port)/tools" -Headers $headers
Write-Host "tools: $($tools.tools.Count)"
$created = Invoke-Tool "create_project" @{ name = "MCP Smoke $(Get-Date -Format o)" }
$projId = $created.result.id
Invoke-Tool "add_scene" @{ projectId = $projId; name = "S1"; synopsis = "test" } | Out-Null
Invoke-Tool "brain_status" @{} | Out-Null
Invoke-Tool "list_local_models" @{} | Out-Null
Invoke-Tool "set_brain" @{ projectId = $projId; brain = "local" } | Out-Null
Write-Host "OK project=$projId"
