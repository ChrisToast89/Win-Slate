# Feature parity — Windows port vs Slate 0.3.2

Legend: ✅ parity · 🟡 adapted for Windows · ⚠️ deferred / limited · ❌ not applicable

| Feature | Status | Notes |
|---------|--------|-------|
| Project CRUD under Documents/Slate | ✅ | Same JSON schema; atomic save |
| SLATE_DATA_DIR override | ✅ | |
| Scenes / shots / specs / history / takes | ✅ | Frontend store + model unchanged |
| Prompt editor (CodeMirror, lock/mute) | ✅ | Same component |
| Coverage plans | ✅ | Same `data/coverage-plans.json` |
| Setups library | ✅ | Same `data/setups.json` |
| Model profiles + Deliver compile | ✅ | Same profiles + `compile.ts` |
| Brain: Claude Code CLI | 🟡 | PATH / npm global resolution (no Homebrew paths) |
| Brain: Codex CLI | 🟡 | PATH only (no ChatGPT.app bundle) |
| Brain: local OpenAI-compatible | ✅ | Same ports + custom endpoint |
| Brain cancel | ✅ | |
| Brain connectivity test pill | ✅ | |
| First AD | ✅ | Same `firstAD.ts` + panel |
| Director notes / transforms | ✅ | Same `brainTasks.ts` |
| Studios (cast/art/loc/look) | ✅ | LocalImg for stills |
| Sound department + audio fingerprint | ✅ | Go DSP port of `audio.ts` |
| References + ffmpeg frames | ✅ | |
| Stills extract (ffmpeg) | ✅ | |
| Circle Take discover | 🟡 | Also probes `%APPDATA%\circle-take`; mac Library path retained |
| MCP control server | ✅ | All tools incl. brain_* |
| MCP stdio bridge | ✅ | `mcp/slate-mcp.mjs` |
| About / Help / Ko-fi / credits | ✅ | Attribution + Windows note |
| App menu | 🟡 | File/Help via Wails menus |
| macOS hiddenInset titlebar | ❌ | Standard Windows frame |
| macOS package-macos / install.sh | ❌ | `wails build` instead |
| Screenshot / snap capture harness | ⚠️ | Not ported (dev tooling) |
| Headless QA scripts (Playwright-style) | ⚠️ | Replaced with Go unit tests + manual checklist |
| Electron browser-preview mock | ✅ | Still available if Wails bindings missing |

## Project file compatibility

Project JSON written by either host should open in the other (same schema). Media paths are absolute and OS-specific.
