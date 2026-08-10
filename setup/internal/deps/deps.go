package deps

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/ChrisToast89/Win-Slate/setup/internal/audit"
	"github.com/ChrisToast89/Win-Slate/setup/internal/logx"
)

type ProgressFn func(step, detail string, percent int)

// EnsureOptional installs missing helpful tools (ffmpeg via winget). Never fails the install hard.
func EnsureOptional(progress ProgressFn) []string {
	if progress == nil {
		progress = func(string, string, int) {}
	}
	var notes []string
	r := audit.Run()

	if !r.FFmpegOK {
		progress("ffmpeg", "Trying to install ffmpeg…", 40)
		if err := installFFmpeg(); err != nil {
			logx.Log("ffmpeg install: %v", err)
			notes = append(notes, "ffmpeg was not installed automatically — install later for video/audio features.")
			progress("ffmpeg", "Could not auto-install ffmpeg (optional).", 50)
		} else {
			notes = append(notes, "ffmpeg installed.")
			progress("ffmpeg", "ffmpeg is ready.", 50)
		}
	} else {
		progress("ffmpeg", "Already available.", 50)
	}
	return notes
}

func installFFmpeg() error {
	if which("winget") == "" {
		return fmt.Errorf("winget not available")
	}
	cmd := exec.Command("winget", "install", "-e", "--id", "Gyan.FFmpeg",
		"--accept-package-agreements", "--accept-source-agreements", "--disable-interactivity")
	out, err := cmd.CombinedOutput()
	logx.Log("winget ffmpeg: %s err=%v", truncate(string(out), 1500), err)
	if err != nil {
		// alternate id
		cmd = exec.Command("winget", "install", "-e", "--id", "ffmpeg",
			"--accept-package-agreements", "--accept-source-agreements", "--disable-interactivity")
		out, err = cmd.CombinedOutput()
		logx.Log("winget ffmpeg alt: %s err=%v", truncate(string(out), 1500), err)
	}
	time.Sleep(2 * time.Second)
	refreshPath()
	if which("ffmpeg") == "" {
		return fmt.Errorf("ffmpeg still not on PATH")
	}
	return nil
}

// EnsureClaude tries to make Claude Code available. Returns human-readable result.
func EnsureClaude(progress ProgressFn) (ok bool, message string) {
	if progress == nil {
		progress = func(string, string, int) {}
	}
	if p := resolveClaude(); p != "" {
		return true, "Claude Code already available: " + p
	}
	progress("Claude Code", "Claude Code not found — attempting install…", 70)
	if which("npm") == "" && which("node") == "" {
		return false, "Node.js/npm not found. Install Node LTS from https://nodejs.org then re-run “Install Claude Code”, or install Claude Code manually from https://claude.com/claude-code"
	}
	progress("Claude Code", "Running: npm install -g @anthropic-ai/claude-code …", 75)
	cmd := exec.Command("npm", "install", "-g", "@anthropic-ai/claude-code")
	// hide window
	hideConsole(cmd)
	out, err := cmd.CombinedOutput()
	logx.Log("npm claude: %s err=%v", truncate(string(out), 2000), err)
	refreshPath()
	if p := resolveClaude(); p != "" {
		return true, "Claude Code installed: " + p + "\n\nNext: open a terminal and run:\n  claude auth login\nThen in Slate set Brain → Claude Code."
	}
	if err != nil {
		return false, "Could not install Claude Code automatically.\n\n" + truncate(string(out), 400) + "\n\nManual: npm install -g @anthropic-ai/claude-code\nThen: claude auth login"
	}
	return false, "Install finished but claude is not on PATH yet. Close Setup, open a new terminal, run: claude --version"
}

func LaunchClaudeAuth() (string, error) {
	cli := resolveClaude()
	if cli == "" {
		return "", fmt.Errorf("Claude Code not found")
	}
	// Prefer real exe if .cmd
	args := []string{"auth", "login"}
	var cmd *exec.Cmd
	if strings.EqualFold(filepath.Ext(cli), ".cmd") || strings.EqualFold(filepath.Ext(cli), ".bat") {
		cmd = exec.Command("cmd", append([]string{"/c", "start", "Claude Code Login", cli}, args...)...)
	} else {
		ps := fmt.Sprintf(`Start-Process -FilePath %q -ArgumentList 'auth','login'`, cli)
		cmd = exec.Command("powershell", "-NoProfile", "-ExecutionPolicy", "Bypass", "-Command", ps)
	}
	if err := cmd.Start(); err != nil {
		return "", err
	}
	return "Opened Claude Code login. Finish sign-in in the browser, then test the Brain pill in Slate.", nil
}

func resolveClaude() string {
	if p := which("claude"); p != "" {
		return p
	}
	cands := []string{
		filepath.Join(os.Getenv("APPDATA"), "npm", "claude.cmd"),
		filepath.Join(os.Getenv("APPDATA"), "npm", "node_modules", "@anthropic-ai", "claude-code", "bin", "claude.exe"),
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

func truncate(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n] + "…"
}

func refreshPath() {
	// Best-effort: merge Machine+User PATH into this process for child lookups.
	// (Limited; full fix often needs new process.)
	cmd := exec.Command("powershell", "-NoProfile", "-Command",
		`[Environment]::GetEnvironmentVariable('Path','Machine') + ';' + [Environment]::GetEnvironmentVariable('Path','User')`)
	hideConsole(cmd)
	out, err := cmd.CombinedOutput()
	if err == nil {
		p := strings.TrimSpace(string(out))
		if p != "" {
			_ = os.Setenv("Path", p)
		}
	}
}
