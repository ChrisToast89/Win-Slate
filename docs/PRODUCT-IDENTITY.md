# Product identity — Win-Slate

## One-sentence definition

**Win-Slate** is a ready-to-run **Windows port** of Sam Wasserman’s Slate (prompt studio), implemented as a **Go + Wails** desktop app with its own Setup — not a wrapper around the Electron build.

## Sibling products (do not mix)

1. **Slate Setup for Windows** (`slate-installer` / GitHub **slate-windows-setup**)  
   - Installs **upstream Electron Slate** via Node/npm tooling.  
   - Output app: Programs **Slate**.  

2. **Archived early port** (`_archive/early-wails-port-NOT-ACTIVE/`)  
   - Former path: `Slate-win/slate-windows/`.  
   - Early Wails experiment only — **not** the published product. All new port work is **this** repo (`Win-Slate/`).

## User-facing names

- App: **Win-Slate**  
- Installer package: **Win-Slate-Setup.zip**  
- Version: see root `VERSION` and Git tags `v*`

## Developer binary locations

| Artifact | Path after build |
|----------|------------------|
| App (primary) | `app/Win-Slate.exe` |
| App (alias) | `app/Slate.exe` (copy of primary) |
| Setup | `Win-Slate-Setup.exe` / zip at repo root (from release script) |

## Further reading

- [AGENTS.md](../AGENTS.md)  
- [app/docs/PARITY.md](../app/docs/PARITY.md)  
- Workspace [PRODUCT-MAP.md](../../../PRODUCT-MAP.md)
