# Self-QA — Slate Windows port

## Automated

| Check | Command | Result (2026-08-06/07 session) |
|-------|---------|--------------------------------|
| Go unit tests | `go test ./...` | **PASS** (audio, brain, projects) |
| Frontend typecheck/build | `cd frontend && npm run build` | **PASS** (Vite production bundle) |
| Wails production build | `.\scripts\build.ps1` | **PASS** → root `Slate.exe` (~12.7 MB) |
| Control descriptor | Launch app | **PASS** `%APPDATA%\slate\control.json` written |
| MCP tools catalog | GET `/tools` | **PASS** — 13 tools |
| MCP create/list/get/add_scene/set_brain | live invoke | **PASS** |
| MCP get/set_shot_prompt + history | live invoke | **PASS** |
| MCP brain_status / list_local_models | live invoke | **PASS** |
| Project JSON on disk | `Documents\Slate\<id>\project.json` | **PASS** |
| Reference tree `slate-0.3.2` | not modified | **PASS** (read-only) |

## Manual checklist (parity)

Run `.\Slate.exe` and verify:

- [x] App launches (WebView2 window titled Slate)
- [ ] Home: list/create/open/delete project (partially via MCP)
- [ ] Project bible fields + brain picker
- [ ] Add scene + shot; editor syntax highlighting
- [ ] Picture-lock / mute lines
- [ ] Coverage plan / setups insert
- [ ] Deliver compile + copy
- [ ] Brain test pill (needs local or CLI brain)
- [ ] First AD panel
- [ ] Studios / Sound / Refs / Stills
- [x] About shows Sam Wasserman credit + version `1.0.0`
- [x] Help / About menu items wired
- [x] MCP bridge tools (all 13)

## Notes

- Claude/Codex not installed on this machine during QA → `available: false` expected; local brain path verified empty when no server.
- ffmpeg present on PATH (`D:\Studio\ffmpeg\bin\ffmpeg.exe` and another path).
- Drag-drop path resolution remains a known WebView2 limitation (U01).
