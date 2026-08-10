package audit

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/ChrisToast89/Win-Slate/setup/internal/logx"
	"github.com/ChrisToast89/Win-Slate/setup/internal/manifest"
	"github.com/ChrisToast89/Win-Slate/setup/internal/paths"
	"golang.org/x/sys/windows"
)

// ProgressFn reports live activity during Check PC (never blocks the audit itself).
type ProgressFn func(step, detail string, percent int)

type Check struct {
	ID       string `json:"id"`
	Label    string `json:"label"`
	OK       bool   `json:"ok"`
	Required bool   `json:"required"`
	Detail   string `json:"detail"`
	Action   string `json:"action"`
}

type Report struct {
	Checks           []Check `json:"checks"`
	CanProceed       bool    `json:"canProceed"`
	Summary          string  `json:"summary"`
	WindowsOK        bool    `json:"windowsOk"`
	DiskOK           bool    `json:"diskOk"`
	WebView2OK       bool    `json:"webView2Ok"`
	FFmpegOK         bool    `json:"ffmpegOk"`
	WingetOK         bool    `json:"wingetOk"`
	NodeOK           bool    `json:"nodeOk"`
	ClaudeOK         bool    `json:"claudeOk"`
	CodexOK          bool    `json:"codexOk"`
	AlreadyInstalled bool   `json:"alreadyInstalled"`
	InstallPath      string `json:"installPath"`
	InstalledVersion string `json:"installedVersion"`
	NpmSlatePresent  bool   `json:"npmSlatePresent"`
	NpmSlatePath     string `json:"npmSlatePath"`
	ProjectsDir      string `json:"projectsDir"`
}

// Run performs a finite system audit. Every external process has a hard timeout.
// There are no loops over unbounded user data — only fixed candidate lists.
func Run(progress ProgressFn) Report {
	if progress == nil {
		progress = func(string, string, int) {}
	}
	logx.Log("Starting system audit")
	var r Report
	r.ProjectsDir = paths.ProjectsDir()
	r.Checks = []Check{}

	progress("Windows", "Reading Windows version…", 5)
	ver, winOK := windowsVersion()
	r.WindowsOK = winOK
	r.Checks = append(r.Checks, Check{
		ID: "windows", Label: "Windows 10/11", OK: winOK, Required: true,
		Detail: ver,
		Action: ternary(winOK, "Ready", "Windows 10 or 11 (64-bit) is required"),
	})

	progress("Disk", "Checking free disk space…", 12)
	freeGB, diskOK := freeSpaceGB()
	r.DiskOK = diskOK
	r.Checks = append(r.Checks, Check{
		ID: "disk", Label: "Free disk space", OK: diskOK, Required: true,
		Detail: fmt.Sprintf("%.1f GB free (need ~0.5 GB)", freeGB),
		Action: ternary(diskOK, "Ready", "Free up disk space, then try again"),
	})

	progress("WebView2", "Looking for WebView2 / Edge…", 22)
	wvOK, wvDetail := webView2OK()
	r.WebView2OK = wvOK
	r.Checks = append(r.Checks, Check{
		ID: "webview2", Label: "WebView2 Runtime", OK: wvOK, Required: true,
		Detail: wvDetail,
		Action: ternary(wvOK, "Ready", "Install Microsoft Edge WebView2 Runtime, then retry"),
	})

	progress("winget", "Checking Windows Package Manager…", 32)
	r.WingetOK = which("winget") != ""
	r.Checks = append(r.Checks, Check{
		ID: "winget", Label: "Windows Package Manager (winget)", OK: r.WingetOK, Required: false,
		Detail: ternary(r.WingetOK, which("winget"), "Not found — optional auto-install of ffmpeg/tools"),
		Action: ternary(r.WingetOK, "Can auto-install ffmpeg when missing", "Install tools manually if needed"),
	})

	progress("ffmpeg", "Looking for ffmpeg on PATH…", 42)
	ff := which("ffmpeg")
	r.FFmpegOK = ff != ""
	r.Checks = append(r.Checks, Check{
		ID: "ffmpeg", Label: "ffmpeg (video / audio features)", OK: r.FFmpegOK, Required: false,
		Detail: ternary(r.FFmpegOK, ff, "Not on PATH"),
		Action: ternary(r.FFmpegOK, "Ready", "Setup can install via winget, or install later"),
	})

	progress("Node.js", "Checking Node.js (optional, for Claude Code)…", 55)
	nodeV, nodeOK := nodeVersion()
	r.NodeOK = nodeOK
	r.Checks = append(r.Checks, Check{
		ID: "node", Label: "Node.js (for Claude Code install only)", OK: nodeOK, Required: false,
		Detail: ternary(nodeOK, "Found: "+nodeV, "Not required to run Win-Slate; needed only to npm-install Claude Code"),
		Action: ternary(nodeOK, "Can install Claude Code via npm", "Install Node LTS if you want one-click Claude Code setup"),
	})

	progress("Claude Code", "Looking for Claude Code CLI…", 68)
	claude := resolveClaude()
	r.ClaudeOK = claude != ""
	r.Checks = append(r.Checks, Check{
		ID: "claude", Label: "Claude Code (AI brain)", OK: r.ClaudeOK, Required: false,
		Detail: ternary(r.ClaudeOK, claude, "Not found — optional; recommended for agent tools"),
		Action: ternary(r.ClaudeOK, "Ready — sign in with: claude auth login", "Setup can try npm install -g @anthropic-ai/claude-code"),
	})

	progress("Codex", "Looking for Codex CLI…", 78)
	codex := which("codex")
	r.CodexOK = codex != ""
	r.Checks = append(r.Checks, Check{
		ID: "codex", Label: "Codex CLI (optional brain)", OK: r.CodexOK, Required: false,
		Detail: ternary(r.CodexOK, codex, "Not installed"),
		Action: "Optional alternative to Claude Code",
	})

	progress("Win-Slate", "Checking for an existing Win-Slate install…", 88)
	if dir, m, ok := manifest.Discover(); ok {
		r.AlreadyInstalled = true
		r.InstallPath = dir
		if m != nil {
			r.InstalledVersion = m.AppVersion
			if m.ReleaseTag != "" {
				r.InstalledVersion = m.ReleaseTag
			}
		} else {
			r.InstalledVersion = "(unknown)"
		}
		r.Checks = append(r.Checks, Check{
			ID: "existing-winslate", Label: "Win-Slate (this app)", OK: true, Required: false,
			Detail: r.InstallPath + " · " + r.InstalledVersion,
			Action: "Update in place or choose a different folder for Win-Slate",
		})
	} else {
		r.Checks = append(r.Checks, Check{
			ID: "existing-winslate", Label: "Win-Slate (this app)", OK: true, Required: false,
			Detail: "Not installed yet — will use the folder you choose (default: Programs\\Win-Slate)",
			Action: "Fresh Win-Slate install",
		})
	}

	progress("npm Slate", "Checking for Sam's npm/Electron Slate (separate product)…", 94)
	npm := manifest.DetectNpmSlate()
	r.NpmSlatePresent = npm.Present
	r.NpmSlatePath = npm.Path
	r.Checks = append(r.Checks, Check{
		ID: "npm-slate", Label: "Sam's Slate (npm / Electron package)", OK: true, Required: false,
		Detail: npm.Detail,
		Action: ternary(npm.Present,
			"Leave it alone — Win-Slate is a different binary and installs elsewhere",
			"Optional; install via the separate Slate Setup helper if you want the Electron build"),
	})

	r.Checks = append(r.Checks, Check{
		ID: "projects", Label: "Projects folder (shared, protected)", OK: true, Required: false,
		Detail: r.ProjectsDir + " — used by both Win-Slate and Sam's Slate; never modified by this Setup",
		Action: "Never modified by Win-Slate Setup",
	})

	r.CanProceed = r.WindowsOK && r.DiskOK && r.WebView2OK
	fail := 0
	for _, c := range r.Checks {
		if c.Required && !c.OK {
			fail++
		}
	}
	if r.CanProceed {
		r.Summary = "This PC can install Win-Slate."
	} else {
		r.Summary = fmt.Sprintf("%d required check(s) failed — fix those, then try again.", fail)
	}
	progress("Done", r.Summary, 100)
	logx.Log("Audit done: canProceed=%v", r.CanProceed)
	return r
}

func resolveClaude() string {
	if p := which("claude"); p != "" {
		return p
	}
	home, _ := os.UserHomeDir()
	cands := []string{
		filepath.Join(os.Getenv("APPDATA"), "npm", "claude.cmd"),
		filepath.Join(os.Getenv("APPDATA"), "npm", "node_modules", "@anthropic-ai", "claude-code", "bin", "claude.exe"),
		filepath.Join(home, "AppData", "Roaming", "npm", "claude.cmd"),
	}
	for _, c := range cands {
		if st, err := os.Stat(c); err == nil && !st.IsDir() {
			return c
		}
	}
	return ""
}

func which(name string) string {
	p, err := exec.LookPath(name)
	if err != nil {
		return ""
	}
	return p
}

// runTimed runs a short CLI probe with a hard deadline (avoids hung node/winget aliases).
func runTimed(timeout time.Duration, name string, args ...string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	cmd := exec.CommandContext(ctx, name, args...)
	hideConsole(cmd)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	if ctx.Err() == context.DeadlineExceeded {
		return "", fmt.Errorf("timed out after %s", timeout)
	}
	return strings.TrimSpace(out.String()), err
}

func nodeVersion() (string, bool) {
	p := which("node")
	if p == "" {
		return "", false
	}
	// Windows App Execution Alias stubs under WindowsApps can hang forever if executed.
	if strings.Contains(strings.ToLower(p), `\windowsapps\`) {
		logx.Log("node on PATH is WindowsApps alias (%s) — skipping live -v probe", p)
		return "(WindowsApps alias — not verified)", false
	}
	out, err := runTimed(3*time.Second, p, "-v")
	if err != nil {
		logx.Log("node -v failed: %v out=%s", err, out)
		return "", false
	}
	v := strings.TrimSpace(out)
	v = strings.TrimPrefix(v, "v")
	// First line only
	if i := strings.IndexAny(v, "\r\n"); i >= 0 {
		v = v[:i]
	}
	parts := strings.Split(v, ".")
	if len(parts) == 0 {
		return "v" + v, true
	}
	maj, _ := strconv.Atoi(parts[0])
	return "v" + v, maj >= 18
}

func webView2OK() (bool, string) {
	keys := []string{
		`SOFTWARE\WOW6432Node\Microsoft\EdgeUpdate\Clients\{F3017226-FE2A-4295-8BDF-00C3A9A7E4C5}`,
		`SOFTWARE\Microsoft\EdgeUpdate\Clients\{F3017226-FE2A-4295-8BDF-00C3A9A7E4C5}`,
	}
	for _, k := range keys {
		opened, err := registryOpen(k)
		if err == nil {
			return true, "WebView2 runtime detected (" + opened + ")"
		}
	}
	if fileExists(`C:\Program Files (x86)\Microsoft\Edge\Application\msedge.exe`) ||
		fileExists(`C:\Program Files\Microsoft\Edge\Application\msedge.exe`) {
		return true, "Microsoft Edge present (WebView2 typically available)"
	}
	// LookPath("msedge") can be slow/noisy — only use fixed paths above
	return false, "WebView2 not detected — install from Microsoft if the app fails to open"
}

func registryOpen(path string) (string, error) {
	var k windows.Handle
	err := windows.RegOpenKeyEx(windows.HKEY_LOCAL_MACHINE, windows.StringToUTF16Ptr(path), 0, windows.KEY_READ|windows.KEY_WOW64_64KEY, &k)
	if err != nil {
		err = windows.RegOpenKeyEx(windows.HKEY_LOCAL_MACHINE, windows.StringToUTF16Ptr(path), 0, windows.KEY_READ|windows.KEY_WOW64_32KEY, &k)
	}
	if err != nil {
		return "", err
	}
	_ = windows.RegCloseKey(k)
	return path, nil
}

func windowsVersion() (string, bool) {
	if runtime.GOOS != "windows" {
		return runtime.GOOS, false
	}
	maj, min, build := windows.RtlGetNtVersionNumbers()
	detail := fmt.Sprintf("Windows %d.%d (build %d) %s", maj, min, build, runtime.GOARCH)
	ok := maj >= 10
	return detail, ok
}

func freeSpaceGB() (float64, bool) {
	var freeBytes, total, totalFree uint64
	path := paths.LocalAppData()
	p, _ := windows.UTF16PtrFromString(path)
	err := windows.GetDiskFreeSpaceEx(p, &freeBytes, &total, &totalFree)
	if err != nil {
		return 0, true
	}
	gb := float64(freeBytes) / (1024 * 1024 * 1024)
	return gb, gb >= 0.4
}

func fileExists(p string) bool {
	st, err := os.Stat(p)
	return err == nil && !st.IsDir()
}

func ternary(cond bool, a, b string) string {
	if cond {
		return a
	}
	return b
}

func init() {
	_ = syscall.Getpid()
}
