# Slate for Windows

**Windows desktop build** of [**Slate**](https://github.com/wassermanproductions/slate) — the prompt studio for AI filmmaking — plus a Setup installer.

> **Slate is by [Sam Wasserman](https://wassermanproductions.com)**  
> Upstream (authoritative): [github.com/wassermanproductions/slate](https://github.com/wassermanproductions/slate) · Apache-2.0  
> This repository is a **Windows port and packaging** of that work. It is **not** an official Wasserman release.

---

## Who made what

| | |
|--|--|
| **Slate** (product, design, UI concepts, prompt/compile logic, data profiles) | **[Sam Wasserman](https://wassermanproductions.com)** — [wassermanproductions/slate](https://github.com/wassermanproductions/slate) |
| **This repo** | Windows **Go + Wails** host, installer, and release packaging so Slate can run as a local Windows app |

Please keep Sam’s credit in About screens, NOTICE, and documentation when redistributing.

---

## Download

**[Slate for Windows v0.3.2-win.1 (zip)](https://github.com/ChrisToast89/slate-for-windows/releases/download/v0.3.2-win.1/SlateForWindows-v0.3.2-win.1.zip)**  

All releases: [github.com/ChrisToast89/slate-for-windows/releases](https://github.com/ChrisToast89/slate-for-windows/releases)

1. Download and unzip the package.
2. Run **`SlateForWindows-Setup.exe`** (or portable `Slate.exe`).
3. Follow **Check this PC** → choose an install folder → **Install**.
4. Optional: install/connect **Claude Code** when Setup offers it (AI brain; no API keys stored in Slate).

- Default install: `%LOCALAPPDATA%\Programs\Slate for Windows`
- Your projects (never deleted by Setup): `%USERPROFILE%\Documents\Slate`

---

## What the installer does

1. **Audits** the PC (Windows 10/11, disk space, WebView2, ffmpeg, Claude Code, existing install)
2. **Checks this GitHub repo** for a newer release than the installed version
3. **Installs or updates** program files only (you pick the folder)
4. **Dependencies:** offers to install **ffmpeg** (winget) when missing; checks **WebView2**
5. **Claude Code:** detects CLI; can run `npm install -g @anthropic-ai/claude-code` when Node is available, and open `claude auth login` guidance
6. Creates Start Menu / optional Desktop shortcuts
7. **Smoke-tests** that `Slate.exe` starts
8. **Never touches** `Documents\Slate` project files

---

## Requirements

- Windows 10 or 11 (x64)
- [WebView2 Runtime](https://developer.microsoft.com/microsoft-edge/webview2/) (usually already installed)
- Optional: [ffmpeg](https://ffmpeg.org/) on PATH (video/audio features)
- Optional AI brain: [Claude Code](https://claude.com/claude-code), Codex CLI, or a local OpenAI-compatible server

Slate does **not** store API keys. Claude/Codex use your own CLI sign-in.

---

## Repository layout

```text
slate-for-windows/
  app/          # Windows port of Slate (Wails + React UI)
  setup/        # Setup / installer (Wails wizard)
  scripts/      # build-release.ps1
  NOTICE        # Required attribution
  LICENSE       # Apache-2.0
```

### Build the app

```powershell
cd app
wails build
Copy-Item -Force build\bin\Slate.exe .\Slate.exe
```

### Build Setup (embeds the app binary)

```powershell
.\scripts\build-release.ps1
```

Produces under `dist/`:

- `SlateForWindows-Setup.exe` (and a zip package)
- `Slate.exe` (portable app binary)

---

## License

Apache License 2.0. See [LICENSE](./LICENSE) and [NOTICE](./NOTICE).

Original Slate © 2026 Sam Wasserman. Windows port and installer are derivative packaging under the same license terms with retained attribution.
