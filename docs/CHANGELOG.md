# Changelog — Win-Slate

## 1.0.0 (released)

First stable public release. See [Releases](https://github.com/ChrisToast89/Win-Slate/releases/tag/v1.0.0).

## Post-1.0 maintenance (main, not yet retagged)

- **Studios sub-tabs:** labels no longer clip; wrap onto a second row instead of hidden horizontal pan-scroll (`app/frontend` CSS).
- **Build:** `app/scripts/build.ps1` writes `app/Win-Slate.exe` and alias `app/Slate.exe` at app folder root.
- **CI release:** `scripts/build-release.ps1` no longer loses the app path when capturing `wails` stdout; builds only from `app/` (no sibling port tree).
- **Setup bindings:** regenerate includes `Uninstall()` for the Setup UI.
- **Product identity:** `AGENTS.md`, `docs/PRODUCT-IDENTITY.md` — this repo is the **Wails port**, not the Electron installer (slate-windows-setup).
- **Scripts:** `sync-payload.ps1` / release packaging use `setup/payload/Win-Slate.exe` only.

When cutting a patch release (e.g. 1.0.1), bump `VERSION`, tag `v*`, and run the Release workflow.
