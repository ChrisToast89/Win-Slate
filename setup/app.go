package main

import (
	"context"
	"fmt"
	"os/exec"
	"path/filepath"
	"sync"

	"github.com/wailsapp/wails/v2/pkg/runtime"

	"github.com/ChrisToast89/slate-for-windows/setup/internal/audit"
	"github.com/ChrisToast89/slate-for-windows/setup/internal/deps"
	"github.com/ChrisToast89/slate-for-windows/setup/internal/install"
	"github.com/ChrisToast89/slate-for-windows/setup/internal/logx"
	"github.com/ChrisToast89/slate-for-windows/setup/internal/manifest"
	"github.com/ChrisToast89/slate-for-windows/setup/internal/paths"
	"github.com/ChrisToast89/slate-for-windows/setup/internal/update"
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
	logx.Log("Slate for Windows Setup v%s started", paths.SetupVersion)
}

func (a *App) emitProgress(step, detail string, percent int) {
	if a.ctx == nil {
		return
	}
	runtime.EventsEmit(a.ctx, "install:progress", map[string]interface{}{
		"step": step, "detail": detail, "percent": percent,
	})
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
	return audit.Run()
}

func (a *App) GetInstallStatus() map[string]interface{} {
	dir, m, ok := manifest.Discover()
	out := map[string]interface{}{
		"installed": ok,
		"installDir": dir,
		"version":    "",
		"exePath":    "",
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
	return update.Check()
}

// PickInstallFolder opens a directory picker; returns chosen path or empty if cancelled.
func (a *App) PickInstallFolder(current string) string {
	if current == "" {
		current = paths.DefaultInstallDir()
	}
	path, err := runtime.OpenDirectoryDialog(a.ctx, runtime.OpenDialogOptions{
		Title:                "Choose install folder for Slate for Windows",
		DefaultDirectory:     current,
		CanCreateDirectories: true,
	})
	if err != nil || path == "" {
		return ""
	}
	return path
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
	rep := audit.Run()
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
		return fmt.Errorf("Slate is not installed")
	}
	exe := paths.InstalledExe(dir)
	cmd := exec.Command(exe)
	cmd.Dir = filepath.Dir(exe)
	return cmd.Start()
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
