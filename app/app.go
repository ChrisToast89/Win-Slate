package main

import (
	"context"
	"encoding/base64"
	"fmt"
	"mime"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/wailsapp/wails/v2/pkg/runtime"
	"github.com/wassermanproductions/slate-windows/internal/audio"
	"github.com/wassermanproductions/slate-windows/internal/brain"
	"github.com/wassermanproductions/slate-windows/internal/control"
	"github.com/wassermanproductions/slate-windows/internal/media"
	"github.com/wassermanproductions/slate-windows/internal/projects"
	"github.com/wassermanproductions/slate-windows/internal/stills"
	"github.com/wassermanproductions/slate-windows/internal/types"
)

// App is the Wails-bound host API, mirroring Electron IPC / window.slate.
type App struct {
	ctx context.Context
}

func NewApp() *App {
	return &App{}
}

func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
	_ = control.Start(func() {
		if a.ctx != nil {
			runtime.EventsEmit(a.ctx, "projects:changed")
		}
	})
}

func (a *App) ListProjects() ([]types.ProjectMeta, error) {
	return projects.List()
}

func (a *App) CreateProject(name string) (*types.Project, error) {
	p, err := projects.Create(name)
	if err != nil {
		return nil, err
	}
	if a.ctx != nil {
		runtime.EventsEmit(a.ctx, "projects:changed")
	}
	return p, nil
}

func (a *App) OpenProject(id string) (*types.Project, error) {
	return projects.Open(id)
}

func (a *App) SaveProject(project types.Project) error {
	return projects.Save(&project)
}

func (a *App) DeleteProject(id string) error {
	err := projects.Delete(id)
	if err == nil && a.ctx != nil {
		runtime.EventsEmit(a.ctx, "projects:changed")
	}
	return err
}

func (a *App) RevealProject(id string) error {
	path := projects.ProjectPath(id)
	// Windows Explorer: select the project.json file
	cmd := exec.Command("explorer", "/select,", path)
	_ = cmd.Start()
	return nil
}

func (a *App) BrainStatus(localEndpoint string) types.BrainStatus {
	return brain.Status(localEndpoint)
}

func (a *App) LocalModels(endpoint string) types.LocalModelsResult {
	return brain.DetectLocal(endpoint)
}

func (a *App) BrainRun(req types.BrainRequest) types.BrainResult {
	backend := req.Backend
	if backend == "" {
		backend = "claude"
	}
	return brain.Run(req, backend)
}

func (a *App) BrainCancel(id string) {
	brain.Cancel(id)
}

func (a *App) BrainTest(backend string, local types.LocalOpts) types.BrainResult {
	return brain.Test(backend, local)
}

func (a *App) StillsDiscover() ([]types.CircledTake, error) {
	return stills.DiscoverCircledTakes()
}

func (a *App) StillsExtract(projectID, mediaPath string, inSec, outSec *float64) ([]string, error) {
	return stills.ExtractStills(projects.CacheDir(projectID), mediaPath, inSec, outSec)
}

func (a *App) PickMedia() ([]string, error) {
	paths, err := runtime.OpenMultipleFilesDialog(a.ctx, runtime.OpenDialogOptions{
		Title: "Select media",
		Filters: []runtime.FileFilter{
			{DisplayName: "Media", Pattern: "*.png;*.jpg;*.jpeg;*.webp;*.gif;*.mp4;*.mov;*.m4v;*.webm;*.mkv"},
		},
	})
	if err != nil {
		return []string{}, nil
	}
	if paths == nil {
		return []string{}, nil
	}
	return paths, nil
}

func (a *App) PickAudio() ([]string, error) {
	path, err := runtime.OpenFileDialog(a.ctx, runtime.OpenDialogOptions{
		Title: "Select audio",
		Filters: []runtime.FileFilter{
			{DisplayName: "Audio & Video", Pattern: "*.mp3;*.wav;*.m4a;*.aac;*.flac;*.ogg;*.aif;*.aiff;*.mp4;*.mov;*.m4v;*.webm;*.mkv"},
		},
	})
	if err != nil || path == "" {
		return []string{}, nil
	}
	return []string{path}, nil
}

func (a *App) IngestMedia(projectID, path string) (types.MediaIngestResult, error) {
	kind := media.Kind(path)
	if kind == "" {
		return types.MediaIngestResult{}, fmt.Errorf("Unsupported media type")
	}
	if kind == "image" {
		return types.MediaIngestResult{Kind: kind, Frames: []string{path}}, nil
	}
	frames, err := media.ExtractFrames(projectID, path)
	if err != nil {
		return types.MediaIngestResult{}, err
	}
	return types.MediaIngestResult{Kind: kind, Frames: frames}, nil
}

func (a *App) AnalyzeAudio(path string) (types.AudioFingerprint, error) {
	return audio.Analyze(path)
}

func (a *App) CopyText(text string) {
	runtime.ClipboardSetText(a.ctx, text)
}

// FileAsDataURL reads a local image/video still and returns a data URL for WebView2.
// Replaces Electron file:// image loading in the original UI.
func (a *App) FileAsDataURL(path string) (string, error) {
	raw, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	ext := strings.ToLower(filepath.Ext(path))
	mimeType := mime.TypeByExtension(ext)
	if mimeType == "" {
		switch ext {
		case ".jpg", ".jpeg":
			mimeType = "image/jpeg"
		case ".png":
			mimeType = "image/png"
		case ".webp":
			mimeType = "image/webp"
		case ".gif":
			mimeType = "image/gif"
		default:
			mimeType = "application/octet-stream"
		}
	}
	return "data:" + mimeType + ";base64," + base64.StdEncoding.EncodeToString(raw), nil
}

// AppVersion for About panel.
func (a *App) AppVersion() string {
	return "0.3.2-win.1"
}

// OpenExternal opens a URL in the default browser.
func (a *App) OpenExternal(url string) {
	runtime.BrowserOpenURL(a.ctx, url)
}

// EmitAbout / EmitHelp are callable from Go menus.
func (a *App) EmitAbout() {
	if a.ctx != nil {
		runtime.EventsEmit(a.ctx, "about:open")
	}
}

func (a *App) EmitHelp() {
	if a.ctx != nil {
		runtime.EventsEmit(a.ctx, "help:open")
	}
}

// EmitBrainRefresh asks the UI to re-probe brain availability.
func (a *App) EmitBrainRefresh() {
	if a.ctx != nil {
		runtime.EventsEmit(a.ctx, "brain:refresh")
	}
}

// ConnectClaudeAccount launches Claude Code CLI login (no keys stored in Slate).
func (a *App) ConnectClaudeAccount() (string, error) {
	msg, err := brain.LaunchClaudeAuthLogin()
	a.EmitBrainRefresh()
	return msg, err
}

// ConnectCodexAccount launches Codex CLI login (optional; same no-key model).
func (a *App) ConnectCodexAccount() (string, error) {
	msg, err := brain.LaunchCodexAuthLogin()
	a.EmitBrainRefresh()
	return msg, err
}

// TestClaudeBrain runs the tiny connectivity check against Claude Code.
func (a *App) TestClaudeBrain() types.BrainResult {
	res := brain.Test("claude", types.LocalOpts{})
	a.EmitBrainRefresh()
	return res
}

// ClaudeCLIInstalled reports whether the claude binary is resolvable.
func (a *App) ClaudeCLIInstalled() bool {
	return brain.CLIAvailable("claude")
}

// GetPlatform for UI tweaks.
func (a *App) GetPlatform() string {
	return "windows"
}

// dialog helpers used by the native Brain menu.
func (a *App) showInfo(title, message string) {
	if a.ctx == nil {
		return
	}
	_, _ = runtime.MessageDialog(a.ctx, runtime.MessageDialogOptions{
		Type:    runtime.InfoDialog,
		Title:   title,
		Message: message,
	})
}

func (a *App) showError(title, message string) {
	if a.ctx == nil {
		return
	}
	_, _ = runtime.MessageDialog(a.ctx, runtime.MessageDialogOptions{
		Type:    runtime.ErrorDialog,
		Title:   title,
		Message: message,
	})
}

// MenuConnectClaude is invoked from the native menu (handles dialogs).
func (a *App) MenuConnectClaude() {
	msg, err := a.ConnectClaudeAccount()
	if err != nil {
		a.showError("Connect Claude account", err.Error())
		return
	}
	a.showInfo("Connect Claude account", msg)
}

// MenuConnectCodex is invoked from the native menu.
func (a *App) MenuConnectCodex() {
	msg, err := a.ConnectCodexAccount()
	if err != nil {
		a.showError("Connect Codex account", err.Error())
		return
	}
	a.showInfo("Connect Codex account", msg)
}

// MenuTestClaude is invoked from the native menu.
func (a *App) MenuTestClaude() {
	if !brain.CLIAvailable("claude") {
		a.showError(
			"Test Claude brain",
			"Claude Code CLI was not found on PATH.\n\nUse Brain → Connect Claude account… after installing Claude Code.",
		)
		return
	}
	// Single dialog after the call (MessageDialog is modal and would block the test).
	res := a.TestClaudeBrain()
	if res.OK && strings.Contains(strings.ToLower(res.Text), "ready") {
		a.showInfo(
			"Test Claude brain",
			fmt.Sprintf("Brain online — Claude Code replied in %.1fs.\n\nReply: %s", float64(res.ElapsedMs)/1000, strings.TrimSpace(res.Text)),
		)
		return
	}
	errMsg := res.Error
	if errMsg == "" {
		errMsg = strings.TrimSpace(res.Text)
	}
	if errMsg == "" {
		errMsg = "Unexpected empty response. Try Brain → Connect Claude account… then test again."
	}
	a.showError("Test Claude brain", errMsg)
}

// MenuRefreshBrain re-probes CLI / local server availability.
func (a *App) MenuRefreshBrain() {
	st := brain.Status("")
	var lines []string
	if st.Claude.Available {
		v := "available"
		if st.Claude.Version != nil {
			v = *st.Claude.Version
		}
		lines = append(lines, "Claude Code: yes — "+v)
	} else {
		lines = append(lines, "Claude Code: not found (install CLI + Brain → Connect Claude account…)")
	}
	if st.Codex.Available {
		v := "available"
		if st.Codex.Version != nil {
			v = *st.Codex.Version
		}
		lines = append(lines, "Codex: yes — "+v)
	} else {
		lines = append(lines, "Codex: not found")
	}
	if st.Local.Available {
		v := "yes"
		if st.Local.Version != nil {
			v = *st.Local.Version
		}
		lines = append(lines, "Local model server: "+v)
	} else {
		lines = append(lines, "Local model server: none detected")
	}
	a.EmitBrainRefresh()
	a.showInfo("Brain status", strings.Join(lines, "\n"))
}
