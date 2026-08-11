# Win-Slate — application (`app/`)

This folder is the **Wails application** for **Win-Slate**: the self-contained Windows **port** of [Slate](https://github.com/wassermanproductions/slate) by **Sam Wasserman** (Apache-2.0).

- Host: Go + Wails  
- UI/product design remain Sam’s  
- Repo root: [../README.md](../README.md) · [../AGENTS.md](../AGENTS.md) · [../NOTICE](../NOTICE)

## Not the Electron installer

Do not confuse this tree with:

- **`slate-installer`** — Setup that installs **original Electron/npm Slate**  
- **`Slate-win/slate-windows`** — legacy early port (superseded by this product)

## Build (binary must land in this folder’s root)

```powershell
cd app
.\scripts\build.ps1
```

Produces:

| File | Role |
|------|------|
| `.\Win-Slate.exe` | Primary product binary |
| `.\Slate.exe` | Same binary (alias for older shortcuts / muscle memory) |

Launch from **this directory**:

```powershell
.\Win-Slate.exe
# or
.\Slate.exe
```

Do **not** use plain `go build` — Wails tags are required.

Repo-level release packaging (Setup + zips): `..\scripts\build-release.ps1`.
