package main

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/wailsapp/wails/v2/pkg/runtime"

	"github.com/ChrisToast89/Win-Slate/setup/internal/audit"
	"github.com/ChrisToast89/Win-Slate/setup/internal/deps"
	"github.com/ChrisToast89/Win-Slate/setup/internal/install"
	"github.com/ChrisToast89/Win-Slate/setup/internal/logx"
	"github.com/ChrisToast89/Win-Slate/setup/internal/manifest"
	"github.com/ChrisToast89/Win-Slate/setup/internal/paths"
	"github.com/ChrisToast89/Win-Slate/setup/internal/update"
)

// App is the Setup wizard host.
type App struct {
	ctx     context.Context
	payload []byte
	mu      sync.Mutex
	running bool
	lastRes install.Result
}

func NewApp(payload []byte) *App {
	return &App{payload: payload}
}

func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
	_ = logx.Init()
	logx.Log("Win-Slate Setup v%s started", paths.SetupVersion)
}

func (a *App) emitProgress(step, detail string, percent int) {
	if a.ctx == nil {
		return
	}
	payload := map[string]interface{}{
		"step": step, "detail": detail, "percent": percent,
	}
	// Emit async so a bound method (RunAudit) cannot deadlock waiting on the UI.
	go func() {
		runtime.EventsEmit(a.ctx, "install:progress", payload)
		runtime.EventsEmit(a.ctx, "audit:progress", payload)
	}()
}

// GetPaths exposes constants and default destinations for the UI.
func (a *App) GetPaths() map[string]string {
	dir, _, ok := manifest.Discover()
	installDir := paths.DefaultInstallDir()
	if ok && dir != "" {
		installDir = dir
	}
	return map[string]string{
		"productName":      paths.ProductName,
		"defaultInstall":   paths.DefaultInstallDir(),
		"installDir":       installDir,
		"projectsDir":      paths.ProjectsDir(),
		"logPath":          paths.TempLog(),
		"setupVersion":     paths.SetupVersion,
		"appVersion":       paths.AppVersion,
		"repoURL":          paths.GitHubRepoURL,
		"upstreamURL":      paths.UpstreamSlateURL,
		"upstreamAuthor":   paths.UpstreamAuthor,
	}
}

func (a *App) RunAudit() audit.Report {
	// Emit immediately so the UI leaves "Preparing…" even before first check.
	a.emitProgress("Starting", "Running system checks…", 2)
	return audit.Run(func(step, detail string, percent int) {
		// Progress first (non-blocking), log second — never block UI on log I/O.
		a.emitProgress(step, detail, percent)
		logx.Log("audit [%d%%] %s — %s", percent, step, detail)
	})
}

func (a *App) GetInstallStatus() map[string]interface{} {
	dir, m, ok := manifest.Discover()
	out := map[string]interface{}{
		"installed":  ok,
		"installDir": dir,
		"version":    "",
		"exePath":    "",
		"product":    paths.ProductName,
	}
	if ok {
		out["exePath"] = paths.InstalledExe(dir)
		if m != nil {
			out["version"] = m.AppVersion
			if m.ReleaseTag != "" {
				out["version"] = m.ReleaseTag
			}
			out["manifest"] = m
		}
	}
	return out
}

func (a *App) CheckForUpdates() update.CheckResult {
	// Best-effort only; never required for install. Keep timeout short inside update.Check.
	a.emitProgress("Updates", "Optional: checking GitHub (max ~5s)…", 96)
	res := update.Check()
	if res.OK {
		a.emitProgress("Updates", res.Message, 100)
	} else {
		a.emitProgress("Updates", "Update check skipped — you can still install offline.", 100)
	}
	return res
}

// PickInstallFolder opens a directory picker; returns chosen path or empty if cancelled.
func (a *App) PickInstallFolder(current string) string {
	if current == "" {
		current = paths.DefaultInstallDir()
	}
	// Dialogs fail on some systems if DefaultDirectory does not exist yet.
	start := current
	if st, err := os.Stat(start); err != nil || !st.IsDir() {
		start = filepath.Dir(start)
	}
	if st, err := os.Stat(start); err != nil || !st.IsDir() {
		start = paths.LocalAppData()
	}
	_ = os.MkdirAll(start, 0o755)

	title := "Choose install folder for Win-Slate"
	// Prefer native WinForms picker — OpenDirectoryDialog is unreliable in some Wails builds.
	if p := pickFolderPowerShell(title, start); p != "" {
		return p
	}
	if a.ctx != nil {
		path, err := runtime.OpenDirectoryDialog(a.ctx, runtime.OpenDialogOptions{
			Title:                title,
			DefaultDirectory:     start,
			CanCreateDirectories: true,
		})
		if err != nil {
			logx.Log("OpenDirectoryDialog error: %v", err)
			return ""
		}
		return path
	}
	return ""
}

// pickFolderPowerShell uses WinForms FolderBrowserDialog when Wails dialog is unavailable.
func pickFolderPowerShell(title, start string) string {
	// Escape single quotes for PowerShell single-quoted strings.
	esc := func(s string) string { return strings.ReplaceAll(s, "'", "''") }
	ps := fmt.Sprintf(`
Add-Type -AssemblyName System.Windows.Forms | Out-Null
$f = New-Object System.Windows.Forms.FolderBrowserDialog
$f.Description = '%s'
$f.ShowNewFolderButton = $true
if (Test-Path -LiteralPath '%s') { $f.SelectedPath = '%s' }
$r = $f.ShowDialog()
if ($r -eq [System.Windows.Forms.DialogResult]::OK) { Write-Output $f.SelectedPath }
`, esc(title), esc(start), esc(start))
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()
	cmd := exec.CommandContext(ctx, "powershell", "-NoProfile", "-STA", "-ExecutionPolicy", "Bypass", "-Command", ps)
	var out, errb bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errb
	if err := cmd.Run(); err != nil {
		logx.Log("FolderBrowserDialog failed: %v %s", err, errb.String())
		return ""
	}
	p := strings.TrimSpace(out.String())
	// PowerShell may emit trailing newlines only
	if p == "" {
		return ""
	}
	// First line only
	if i := strings.IndexAny(p, "\r\n"); i >= 0 {
		p = p[:i]
	}
	return strings.TrimSpace(p)
}

// StartInstall installs or updates into installDir using the embedded app binary.
func (a *App) StartInstall(installDir string, desktopShortcut bool, isUpdate bool) (map[string]interface{}, error) {
	a.mu.Lock()
	if a.running {
		a.mu.Unlock()
		return nil, fmt.Errorf("install already running")
	}
	a.running = true
	a.mu.Unlock()
	defer func() {
		a.mu.Lock()
		a.running = false
		a.mu.Unlock()
	}()

	progress := func(step, detail string, percent int) {
		logx.Log("[%d%%] %s — %s", percent, step, detail)
		a.emitProgress(step, detail, percent)
	}

	progress("Safety", "Confirming project files will not be touched…", 2)
	if installDir == "" {
		installDir = paths.DefaultInstallDir()
	}
	if err := paths.AssertSafeInstallDir(installDir); err != nil {
		return nil, err
	}

	progress("Check", "Re-checking this PC…", 8)
	rep := audit.Run(func(step, detail string, percent int) {
		a.emitProgress(step, detail, percent)
	})
	if !rep.CanProceed {
		return nil, fmt.Errorf("%s", rep.Summary)
	}

	progress("Dependencies", "Checking optional tools…", 25)
	notes := deps.EnsureOptional(progress)

	// Prefer latest tag from GitHub for the manifest label
	upd := update.Check()
	tag := "v" + paths.AppVersion
	if upd.LatestTag != "" {
		// If this setup's bundled version matches or is the install source, record bundled tag
		tag = "v" + paths.AppVersion
		if upd.LatestTag != "" && !upd.UpdateAvailable {
			tag = upd.LatestTag
		}
	}

	res, err := install.Run(install.Options{
		InstallDir:      installDir,
		DesktopShortcut: desktopShortcut,
		IsUpdate:        isUpdate,
		ReleaseTag:      tag,
		Payload:         a.payload,
	}, progress)
	if err != nil {
		return nil, err
	}
	a.lastRes = res

	// Soft Claude attempt if missing
	progress("Claude Code", "Checking Claude Code…", 96)
	claudeOK, claudeMsg := false, ""
	if !rep.ClaudeOK {
		claudeOK, claudeMsg = deps.EnsureClaude(progress)
	} else {
		claudeOK, claudeMsg = true, "Claude Code already available."
	}

	return map[string]interface{}{
		"result":     res,
		"depNotes":   notes,
		"claudeOk":   claudeOK,
		"claudeMsg":  claudeMsg,
		"projectsDir": paths.ProjectsDir(),
		"upstream":   paths.UpstreamSlateURL,
	}, nil
}

func (a *App) InstallClaude() map[string]interface{} {
	ok, msg := deps.EnsureClaude(func(step, detail string, percent int) {
		a.emitProgress(step, detail, percent)
	})
	return map[string]interface{}{"ok": ok, "message": msg}
}

func (a *App) LaunchClaudeLogin() (string, error) {
	return deps.LaunchClaudeAuth()
}

func (a *App) LaunchApp() error {
	dir, _, ok := manifest.Discover()
	if !ok {
		return fmt.Errorf("Win-Slate is not installed")
	}
	exe := paths.InstalledExe(dir)
	cmd := exec.Command(exe)
	cmd.Dir = filepath.Dir(exe)
	return cmd.Start()
}

// Uninstall removes a verified Win-Slate install only (never Documents\Slate projects).
func (a *App) Uninstall() (map[string]interface{}, error) {
	a.mu.Lock()
	if a.running {
		a.mu.Unlock()
		return nil, fmt.Errorf("another operation is already running")
	}
	a.running = true
	a.mu.Unlock()
	defer func() {
		a.mu.Lock()
		a.running = false
		a.mu.Unlock()
	}()

	res, err := install.Uninstall(func(step, detail string, percent int) {
		logx.Log("uninstall [%d%%] %s — %s", percent, step, detail)
		a.emitProgress(step, detail, percent)
	})
	if err != nil {
		return nil, err
	}
	return map[string]interface{}{
		"ok":                res.OK,
		"removedDir":        res.RemovedDir,
		"projectsDir":       res.ProjectsDir,
		"projectsPreserved": res.ProjectsPreserved,
		"summary":           res.Summary,
	}, nil
}

func (a *App) OpenExternal(url string) {
	if a.ctx != nil {
		runtime.BrowserOpenURL(a.ctx, url)
	}
}

func (a *App) OpenProjectsFolder() {
	p := paths.ProjectsDir()
	runtime.BrowserOpenURL(a.ctx, "file:///"+filepath.ToSlash(p))
}

func (a *App) OpenLogFolder() {
	runtime.BrowserOpenURL(a.ctx, "file:///"+filepath.ToSlash(filepath.Dir(paths.TempLog())))
}

func (a *App) LogPath() string {
	return logx.Path()
}
