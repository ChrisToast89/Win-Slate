# Slate for Windows

Windows version of **[Slate](https://github.com/wassermanproductions/slate)** — the prompt studio for AI filmmaking — created by **[Sam Wasserman](https://wassermanproductions.com)**.

This is a binary executable for locally running Sam's Slate application.

Some additional configuration may be required.

Download the zip, unzip it, and run the Setup program.

---

## Download

### Installer (recommended)

**→ [SlateForWindows-Setup.exe](./SlateForWindows-Setup.exe)**

On GitHub: open that file in the file list, then **Download** (or use [this direct link](https://github.com/ChrisToast89/slate-for-windows/raw/main/SlateForWindows-Setup.exe)).

### Full package (zip)

Also available as a zip (Setup + portable app + license files):  
**[SlateForWindows-v0.3.2-win.1.zip](https://github.com/ChrisToast89/slate-for-windows/releases/download/v0.3.2-win.1/SlateForWindows-v0.3.2-win.1.zip)** · [all releases](https://github.com/ChrisToast89/slate-for-windows/releases)

---

## Install (simple steps)

1. Download **`SlateForWindows-Setup.exe`** (link above).
2. Double‑click **`SlateForWindows-Setup.exe`**.
3. If Windows SmartScreen appears, choose **More info** → **Run anyway** (common for new apps that are not Microsoft-signed).
4. Follow the on‑screen steps: check your PC → pick a folder (or accept the default) → **Install**.
5. When it finishes, open **Slate for Windows** from the Start Menu, or use **Launch** in Setup.

**That’s it.** You do not need to build or compile the program yourself.

### Portable option

The [release zip](https://github.com/ChrisToast89/slate-for-windows/releases/download/v0.3.2-win.1/SlateForWindows-v0.3.2-win.1.zip) also includes **`Slate.exe`**, which can run without Setup. Setup is still recommended for Start Menu shortcuts and a simple install path.

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
| **The app** (default) | `%LOCALAPPDATA%\Programs\Slate for Windows` |
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
- Can look for a newer version on this GitHub project  
- Installs or updates the program into the folder you choose  
- Can help with optional tools (ffmpeg, Claude Code)  
- Adds Start Menu / optional Desktop shortcuts  
- Does **not** touch your Documents\Slate projects  

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

Source for the app and Setup lives in this repository. Prebuilt install files for users are published on **[Releases](https://github.com/ChrisToast89/slate-for-windows/releases)** — that is the supported way to get the product.

If you maintain this project: release packages are assembled into a local `distributable/` folder by `scripts/build-release.ps1` and uploaded to GitHub Releases. That folder is not part of what end users need to understand.
