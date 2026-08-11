# AGENTS — Win-Slate

## What this repository is

**Win-Slate** is a **self-contained Windows port** of [Slate](https://github.com/wassermanproductions/slate) by Sam Wasserman (Apache-2.0).

- Host: **Go + Wails** (not Electron)
- Product binary: **`Win-Slate.exe`**
- Product Setup: **`Win-Slate-Setup.exe`** / **`Win-Slate-Setup.zip`**
- Public repo: https://github.com/ChrisToast89/Win-Slate
- Version file: `VERSION` (current line of truth for releases)

This is an **unofficial derivative**. Credit Sam Wasserman; do not present as official Wasserman software.

## What this is NOT

| Not this | That lives at |
|----------|----------------|
| Installer for Sam’s **Electron/npm** Slate | Workspace `slate-installer/` · GitHub `ChrisToast89/slate-windows-setup` |
| Early Wails port (archived) | Workspace `_archive/early-wails-port-NOT-ACTIVE/` (was `Slate-win/slate-windows/`) |
| Upstream Electron source | Workspace `slate/` or `slate-0.3.2/` |

See workspace root `PRODUCT-MAP.md` and `../README.md` (parent `Slate-win/`).  
Pointer: `../MOVED-slate-windows.md`.

## Layout

```
Win-Slate/                 ← git root (this repo)
  app/                     ← Wails application (the port)
    Win-Slate.exe          ← build output (primary; also Slate.exe alias)
    scripts/build.ps1      ← production build → root of app/
    frontend/              ← React UI (parity with upstream)
    internal/              ← Go brain, projects, audio, …
  setup/                   ← Win-Slate Setup (embeds payload)
    payload/Win-Slate.exe  ← filled by release scripts
  scripts/build-release.ps1
  VERSION
  ATTRIBUTION.md, NOTICE, LICENSE
```

## Build & binary placement (required)

Always use Wails — not plain `go build`.

```powershell
cd app
.\scripts\build.ps1
```

- **Must** leave a launchable binary in **`app\` root** (next to sources):  
  - `app\Win-Slate.exe` (primary)  
  - `app\Slate.exe` (same file; legacy/alias for shortcuts)
- Do **not** leave the only copy under `app\build\bin\` for local testing.
- For Setup packaging / GitHub release: `.\scripts\build-release.ps1` from repo root.

## Install identity (deconflict)

| | Win-Slate | Electron Slate (via slate-installer) |
|--|-----------|--------------------------------------|
| Binary name | Win-Slate.exe | Slate (Electron) |
| Typical folder | `Programs\Win-Slate` | `Programs\Slate` |
| Projects | `Documents\Slate` (shared JSON schema) | same |

Setup must not uninstall or overwrite the other product’s install tree.

## When to edit here

- Port UI, brain runner, studios/tabs, Windows-specific host behavior  
- Win-Slate Setup UX  
- Releases / CI for **this** repo  

## When to leave this tree alone

- “Install original Slate with npm/Electron” → `slate-installer`  
- “Fix Sam’s Electron app only” → `slate/` (upstream), not this port  

## Attribution

Keep `ATTRIBUTION.md`, `NOTICE`, About credits, and release notes accurate. Apache-2.0.
