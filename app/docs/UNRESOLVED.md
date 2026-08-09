# Unresolved issues — executive review

Items that may need follow-up after the initial Windows port.

| ID | Severity | Issue | Mitigation / next step |
|----|----------|-------|------------------------|
| U01 | Medium | Drag-and-drop of files may not yield a filesystem path under WebView2 (`pathForFile` returns empty). | Use Pick Media / Pick Audio dialogs; investigate WebView2 drag path APIs later. |
| U02 | Low | Circle Take suite is mac-first; Windows discovery only works if recents exist under `%APPDATA%\circle-take`. | Document; optional manual clip import already works. |
| U03 | Medium | Claude/Codex must be on PATH (or npm global). No ChatGPT desktop bundle on Windows. | Document install steps; local brain is the reliable default. |
| U04 | Low | Reveal Project uses `explorer /select,` — behavior differs slightly from macOS Finder. | Acceptable Windows equivalent. |
| U05 | Low | Title bar is standard Windows (not traffic-light / hiddenInset chrome). | Cosmetic only. |
| U06 | Medium | Full GUI regression suite (original `qa-smoke.mjs` / snap) not automated on Windows. | Manual QA checklist in QA.md; unit tests cover host logic. |
| U07 | Info | Publishing requires original author permission; this build is for **personal use**. | Keep NOTICE/LICENSE; do not redistribute without compliance + permission if required by policy. |
| U08 | Low | Optional `*float64` for stills in/out via Wails JSON binding — verify edge cases with live UI. | Covered by API design; add integration test if bugs appear. |
| U09 | Low | `wails generate module` JS stubs must be regenerated after binding changes. | Run as part of build; host.ts uses loose typing as fallback. |
| U10 | Info | Upstream may ship post-0.3.2 changes; this port targets **0.3.2** snapshot only. | Re-sync from new tags if needed. |

## Non-goals (this session)

- Upstream PR / official multi-platform productization
- Signing / SmartScreen reputation
- Bundling ffmpeg inside the .exe
- Linux packaging
