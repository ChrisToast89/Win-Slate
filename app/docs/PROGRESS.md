# Progress log — Slate Windows port

**Session goal:** Full-fidelity Windows binary (Go + Wails) of Slate 0.3.2 without modifying `slate-0.3.2/`.

## Timeline

| Milestone | Status |
|-----------|--------|
| Leave `slate-0.3.2` untouched | Done |
| Scaffold `slate-windows/` + copy UI/data/LICENSE | Done |
| Go host: projects, brain, media, stills, audio, control | Done |
| Frontend host bridge + LocalImg + About attribution | Done |
| MCP bridge with full tool catalog (13) | Done |
| Docs: NOTICE, PARITY, UNRESOLVED, QA | Done |
| `go test ./...` | Pass |
| `wails build -clean` → `Slate.exe` | Pass |
| Live MCP smoke (create project, scenes, shots, brain) | Pass |

## Architecture decisions

1. **Do not touch** `../slate-0.3.2` — reference only.
2. **Reuse React UI + TS domain logic** for behavior parity.
3. **Rewrite host only** in Go (Electron main/preload surface).
4. **Apache-2.0 NOTICE** retained + Windows modification notes.
5. **Version:** `0.3.2-win.1`

## Deliverable

```
slate-windows/build/bin/Slate.exe
```

## 2026-08-07 — Brain menu (Claude Code)

- Native **Brain** menu: Connect Claude account…, Test Claude brain…, Connect Codex…, Refresh status…
- Host launches `claude auth login` / `codex login` only — **no API keys** stored
- UI listens for `brain:refresh` after menu actions
- Home warning points at **Brain → Connect Claude account…**

## MCP usage

```powershell
# App must be running
claude mcp add slate -- node "M:\Users\Chris\Documents\_code-projects\slate\slate-windows\mcp\slate-mcp.mjs"
# or
powershell -File slate-windows\scripts\smoke-mcp.ps1
```

## Rebuild

```powershell
cd slate-windows
wails build
```
