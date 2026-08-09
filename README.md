# Slate for Windows

A **self-contained Windows executable** of [**Slate**](https://github.com/wassermanproductions/slate) — the prompt studio for AI filmmaking — based on the public work of **[Sam Wasserman](https://wassermanproductions.com)**.

This is **not** Sam’s original Electron / macOS-first distribution, and it is **not** an official Wasserman Productions release. It is a **recompile / host port**: the same product idea, UI, and knowledge bases, running as a native Windows binary (Go + Wails) with an optional Setup installer. You download a package, run Setup or `Slate.exe`, and go — no Node build step required to *use* the app.

> **Upstream (authoritative):** [github.com/wassermanproductions/slate](https://github.com/wassermanproductions/slate)  
> **This repo:** Windows packaging and host conversion only

---

## What you’re getting

Slate for Windows ships as a **standalone `.exe`** (plus a Setup wizard that can install it for you). Under the hood it is still Slate: scenes, shots, prompt craft, First AD, brains, and the rest of the studio workflow. The conversion swaps the Electron shell for a Wails host so Windows users get a single local binary instead of an npm-driven build of the upstream tree.

If you want a helper that downloads and builds **Sam’s official npm/Electron sources** on Windows, that is a **different** project (`slate-windows-setup` / `slate-installer`). **This** project is the **binary port** and its installer.

---

## Credit and authorship

| | |
|--|--|
| **Slate** — product design, UI concepts, prompt and compile logic, coverage plans, model profiles, and the overall creative tool | **[Sam Wasserman](https://wassermanproductions.com)** · [wassermanproductions/slate](https://github.com/wassermanproductions/slate) |
| **Windows host, installer, and packaging in this repository** | Derivative work for Windows · not affiliated with or endorsed by Sam as an official release |

Please keep Sam’s credit in About screens, the `NOTICE` file, and any documentation when you redistribute. That is required by the Apache License 2.0 notice rules and is the right thing to do.

---

## Use of generative tools in this port

**Generative coding tools were used to help port, package, and document** this Windows build (including host code, installer flow, and repository tooling). The creative product and original design remain Sam Wasserman’s. Treat AI-assisted portions of the Windows packaging as assistance on a derivative engineering effort, not as a claim of original authorship of Slate itself.

---

## Disclaimer — no warranty

**This software is offered without warranty of any kind.**

It is provided **“AS IS”**, with no guarantees that it will work on every machine, match every behavior of the upstream app, or remain compatible with future changes to Claude Code, Codex, WebView2, or Windows. You use it at your own risk. The authors and redistributors of this Windows port are not liable for data loss, failed installs, or any other damages arising from use of this package. See the full warranty disclaimer and limitation of liability in the [Apache License 2.0](./LICENSE) (Sections 7 and 8).

Upstream Slate is likewise licensed under Apache-2.0 with its own terms; this port does not replace or weaken those terms.

---

## Licenses and acknowledgments

Please read and respect **all** of the following when you use or redistribute this package:

| Work | License / notes | Where |
|------|-----------------|--------|
| **Slate** (original application by Sam Wasserman) | **Apache License 2.0** | Upstream [LICENSE](https://github.com/wassermanproductions/slate/blob/main/LICENSE) · retained as [LICENSE](./LICENSE) |
| **Attribution / NOTICE** | Required under Apache-2.0 §4(d) | [NOTICE](./NOTICE) · credit Sam Wasserman (wassermanproductions.com) |
| **This Windows port and Setup installer** | Same **Apache-2.0** terms as a derivative / packaging work | This repository |
| **Wails / Go / third-party libraries** | Their respective open-source licenses (as used by the build toolchains and dependencies) | See module and package managers when building from source (`go.mod`, `package.json`, Wails stack) |

By running or redistributing Slate for Windows you acknowledge:

1. Sam Wasserman is the author of Slate.  
2. This Windows binary is a recompile-based derivative packaging of that work.  
3. Licenses and the NOTICE file must travel with redistributed copies.  
4. There is **no warranty**.  
5. Generative tooling assisted the Windows port and packaging.

---

## Download

**[Slate for Windows v0.3.2-win.1 (zip)](https://github.com/ChrisToast89/slate-for-windows/releases/download/v0.3.2-win.1/SlateForWindows-v0.3.2-win.1.zip)**  

All releases: [github.com/ChrisToast89/slate-for-windows/releases](https://github.com/ChrisToast89/slate-for-windows/releases)

### Install (end users)

1. Download and unzip the package (from GitHub Releases, or the local `distributable/` folder if you built it yourself).  
2. Run **`SlateForWindows-Setup.exe`**, or the portable **`Slate.exe`**.  
3. In Setup: **Check this PC** → choose an install folder → **Install**.  
4. Optionally install or sign in to **Claude Code** when Setup offers it (AI “brain”). Slate does **not** store API keys; it uses your own CLI sign-in or a local model server.

- Default install location: `%LOCALAPPDATA%\Programs\Slate for Windows`  
- Your projects (never deleted by Setup): `%USERPROFILE%\Documents\Slate`

---

## What the installer does

1. Audits the PC (Windows 10/11, disk space, WebView2, ffmpeg, Claude Code, existing install)  
2. Compares the installed version to this GitHub repo’s latest release  
3. Installs or updates **program files only** into a folder you choose  
4. Can offer **ffmpeg** via winget when missing; checks **WebView2**  
5. Can help install **Claude Code** (`npm install -g @anthropic-ai/claude-code`) and point you at `claude auth login`  
6. Creates Start Menu and optional Desktop shortcuts  
7. Smoke-tests that `Slate.exe` starts  
8. **Never touches** `Documents\Slate` project files  

---

## Requirements

- Windows 10 or 11 (x64)  
- [WebView2 Runtime](https://developer.microsoft.com/microsoft-edge/webview2/) (usually already present)  
- Optional: [ffmpeg](https://ffmpeg.org/) on PATH (video / audio features)  
- Optional AI brain: [Claude Code](https://claude.com/claude-code), Codex CLI, or a local OpenAI-compatible server  

---

## Distributable package (local)

End-user install bits live under **`distributable/`** — not the full development tree (`app/`, `setup/` sources, etc.):

```text
distributable/
  SlateForWindows-v0.3.2-win.1/     # unzipped package
    SlateForWindows-Setup.exe
    Slate.exe
    INSTALL.txt
    LICENSE.txt
    NOTICE.txt
    README.md
  SlateForWindows-v0.3.2-win.1.zip  # same package, zipped
```

That folder is what you hand to someone who only wants to install the app.

---

## Repository layout (developers)

```text
slate-for-windows/
  app/              # Windows port of Slate (Wails + React UI)
  setup/            # Setup / installer wizard (source)
  scripts/          # build-release.ps1, sync-payload.ps1
  distributable/    # end-user install package only
  NOTICE            # required attribution
  LICENSE           # Apache-2.0
```

### App binary location

The app binary is **`Slate.exe` at the app folder root** (not a long-lived `build\bin` copy):

| Source | Path |
|--------|------|
| Dev tree (sibling) | `../slate-windows/Slate.exe` |
| Monorepo app | `app/Slate.exe` |
| Embedded in Setup | `setup/payload/Slate.exe` (copied at build time) |

```powershell
.\scripts\sync-payload.ps1    # refresh Setup embed from root binary
.\scripts\build-release.ps1   # build app + Setup → distributable/
```

---

## License

Licensed under the **Apache License, Version 2.0**. See [LICENSE](./LICENSE) and [NOTICE](./NOTICE).

**Original Slate** © 2026 **Sam Wasserman**.  
**Windows port and installer** are a derivative packaging of that work under the same license, with attribution retained, **no warranty**, and with generative assistance used in the port as described above.
