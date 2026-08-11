# Win-Slate

**Version 1.0.0**

> **Product identity (maintainers / agents)**  
> This repo is the **ported Windows app** (Go + Wails), **not** the installer for Sam’s Electron/npm Slate.  
> Electron installer → workspace `slate-installer` / GitHub `slate-windows-setup`.  
> Early Wails experiment → workspace `_archive/early-wails-port-NOT-ACTIVE/` (not active).  
> See [AGENTS.md](./AGENTS.md) and [docs/PRODUCT-IDENTITY.md](./docs/PRODUCT-IDENTITY.md).

Windows version of **[Slate](https://github.com/wassermanproductions/slate)** — the prompt studio for AI filmmaking — created by **[Sam Wasserman](https://wassermanproductions.com)**.

<img width="320" alt="image" src="https://github.com/user-attachments/assets/1e88f42c-9285-4b79-9836-9406424b508a" />

This is a ready-to-run **Windows binary** of Slate (Go + Wails), with a Setup installer. It is an **unofficial** derivative packaging — not an official Wasserman release.

Download the **Setup zip**, extract it, and run the Setup program. You do **not** need to compile anything.

---

## Download

### Installer (recommended)
<img width="320" alt="image" src="https://github.com/user-attachments/assets/571986e6-f136-4cef-8595-bb5c8a6e2fb0" />

**→ [Win-Slate-Setup.zip](./Win-Slate-Setup.zip)** — repository root  

Contains `Win-Slate-Setup.exe` plus LICENSE / NOTICE / INSTALL notes.  
A **zip** is usually smoother than downloading a bare `.exe` (Windows SmartScreen).

- On GitHub: open **`Win-Slate-Setup.zip`** → **Download**  
- Direct: [Win-Slate-Setup.zip](https://github.com/ChrisToast89/Win-Slate/raw/main/Win-Slate-Setup.zip)

### Release v1.0.0 (full package)

**[Win-Slate v1.0.0 on Releases](https://github.com/ChrisToast89/Win-Slate/releases/tag/v1.0.0)**  

Assets include **`Win-Slate-Setup.zip`** and the full package zip (Setup + portable `Win-Slate.exe` + licenses).

---

## Install (simple steps)

1. Download **`Win-Slate-Setup.zip`** (link above).
2. Right‑click the zip → **Extract All…** (or use any unzip tool).
3. Open the extracted folder and double‑click **`Win-Slate-Setup.exe`**.
4. If Windows SmartScreen appears, choose **More info** → **Run anyway** (common for apps that are not Microsoft-signed yet).
5. Follow the on‑screen steps: check your PC → pick a folder (or accept the default) → **Install**.
6. When it finishes, open **Win-Slate** from the Start Menu, or use **Launch** in Setup.

**That’s it.** You do not need to build or compile the program yourself.

### Portable option

If a full release zip is available, it may also include **`Win-Slate.exe`** to run without Setup. The root **`Win-Slate-Setup.zip`** is still the recommended path for most people.

### After install (optional — AI features)

Slate’s AI “brain” uses tools on **your** computer (for example Claude Code). Setup never stores API keys. If you want those features:

1. Install [Claude Code](https://claude.com/claude-code) if you don’t have it (Setup may offer help).
2. In a terminal, run: `claude auth login` and finish sign‑in in the browser.
3. In Slate, set the brain to Claude Code and use the Brain control to test.

You can use much of Slate without a brain; agent-style tools need one of the supported options (Claude, Codex, or a local model server).

---

## Where things go

| | |
|--|--|
| **Win-Slate** (default) | `%LOCALAPPDATA%\Programs\Win-Slate` |
| **Your projects** | `%USERPROFILE%\Documents\Slate` |

Setup **never deletes or overwrites** your project files under Documents\Slate.

---

## What this is (and isn’t)

- **Is:** A ready‑to‑run Windows program based on Sam Wasserman’s Slate, plus a Setup helper that installs it for you.
- **Is not:** An official release from Sam or Wasserman Productions.
- **Is not:** Something you need to compile. The download is already a finished program.

Original Slate (authoritative source): [github.com/wassermanproductions/slate](https://github.com/wassermanproductions/slate)

---

## Requirements

- A 64‑bit PC running **Windows 10** or **Windows 11**
- Internet only if you use cloud AI tools or Setup’s optional downloads (for example ffmpeg or Claude Code)
- [WebView2](https://developer.microsoft.com/microsoft-edge/webview2/) — usually already installed with Windows / Edge

Optional:

- [ffmpeg](https://ffmpeg.org/) for some video/audio features (Setup may offer to install it)
- Claude Code, Codex, or a local AI server if you want agent features

---

## What Setup does for you

You can ignore the details and just click through; this is only for the curious:

- Checks that your PC can run the app  
- Installs or updates Win-Slate into the folder you choose  
- Can uninstall an existing Win-Slate install (projects kept)  
- Can help with optional tools (ffmpeg, Claude Code)  
- Adds Start Menu / optional Desktop shortcuts  
- Does **not** touch your Documents\Slate projects  

---

## What’s new in 1.0.0

- First stable public release of **Win-Slate**  
- Setup: system check with live progress, install/update, uninstall, folder browse  
- Standalone **`Win-Slate.exe`** (does not replace Sam’s npm/Electron Slate install)  
- Installer distributed as **`Win-Slate-Setup.zip`** for fewer SmartScreen issues  
- Full credit to Sam Wasserman / upstream Slate (Apache-2.0)  

---

## Credit

| | |
|--|--|
| **Slate** (the creative tool) | **[Sam Wasserman](https://wassermanproductions.com)** · [wassermanproductions/slate](https://github.com/wassermanproductions/slate) |
| **This Windows package** | Unofficial Windows packaging of his work |

Please keep Sam’s name and the NOTICE file when you share copies. That respects both him and the license.

**Generative coding tools were used to help create this Windows port and installer.** The original product and design are Sam’s.

---

## No warranty

**This software is provided as‑is, with no warranty of any kind.**

It may not work on every computer, may differ from Sam’s original app in places, and may break when Windows or AI tools change. You use it at your own risk. See the full legal warranty disclaimer in [LICENSE](./LICENSE) (Apache License 2.0, Sections 7 and 8).

---

## Licenses

| | |
|--|--|
| Original **Slate** (Sam Wasserman) | [Apache License 2.0](https://github.com/wassermanproductions/slate/blob/main/LICENSE) |
| This Windows package | Apache License 2.0 — see [LICENSE](./LICENSE) and [NOTICE](./NOTICE) |

By using this package you acknowledge Sam Wasserman as the author of Slate, that this is an unofficial Windows package of that work, that the NOTICE and license terms apply, that there is no warranty, and that generative tools assisted the Windows packaging.

---

## For developers only

*Skip this section unless you are intentionally building from source. Everyday users only need the Download and Install sections above.*

Source for the app and Setup lives in this repository. Prebuilt install files for users are published on **[Releases](https://github.com/ChrisToast89/Win-Slate/releases)** — that is the supported way to get the product.

If you maintain this project: release packages are assembled into a local `distributable/` folder by `scripts/build-release.ps1` and uploaded to GitHub Releases. That folder is not part of what end users need to understand.
