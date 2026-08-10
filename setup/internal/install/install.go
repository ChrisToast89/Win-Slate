package install

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/ChrisToast89/Win-Slate/setup/internal/logx"
	"github.com/ChrisToast89/Win-Slate/setup/internal/manifest"
	"github.com/ChrisToast89/Win-Slate/setup/internal/paths"
)

type ProgressFn func(step, detail string, percent int)

type Result struct {
	InstallDir        string `json:"installDir"`
	ExePath           string `json:"exePath"`
	StartMenuOK       bool   `json:"startMenuOk"`
	DesktopOK         bool   `json:"desktopOk"`
	SmokeOK           bool   `json:"smokeOk"`
	SmokeDetail       string `json:"smokeDetail"`
	Summary           string `json:"summary"`
	ProjectsDir       string `json:"projectsDir"`
	ProjectsPreserved bool   `json:"projectsPreserved"`
	IsUpdate          bool   `json:"isUpdate"`
	AppVersion        string `json:"appVersion"`
	ReleaseTag        string `json:"releaseTag"`
}

type Options struct {
	InstallDir      string
	DesktopShortcut bool
	IsUpdate        bool
	ReleaseTag      string
	Payload         []byte // embedded Slate.exe
}

func Run(opts Options, progress ProgressFn) (Result, error) {
	if progress == nil {
		progress = func(string, string, int) {}
	}
	var res Result
	res.ProjectsDir = paths.ProjectsDir()
	res.ProjectsPreserved = true
	res.IsUpdate = opts.IsUpdate
	res.AppVersion = paths.AppVersion
	res.ReleaseTag = opts.ReleaseTag
	if res.ReleaseTag == "" {
		res.ReleaseTag = "v" + paths.AppVersion
	}

	dest := strings.TrimSpace(opts.InstallDir)
	if dest == "" {
		dest = paths.DefaultInstallDir()
	}
	if err := paths.AssertSafeInstallDir(dest); err != nil {
		return res, err
	}
	if len(opts.Payload) < 1024 {
		return res, fmt.Errorf("embedded Win-Slate.exe payload is missing or too small — rebuild Setup with payload")
	}

	label := "Installing"
	if opts.IsUpdate {
		label = "Updating"
	}
	progress(label, "Writing program files (projects folder is not touched)…", 80)
	logx.Log("%s -> %s (%d bytes payload)", label, dest, len(opts.Payload))

	if err := os.MkdirAll(dest, 0o755); err != nil {
		return res, fmt.Errorf("create install dir: %w", err)
	}

	exePath := paths.InstalledExe(dest)
	// If existing exe is running, overwrite may fail — try replace carefully
	tmp := exePath + ".new"
	if err := os.WriteFile(tmp, opts.Payload, 0o755); err != nil {
		return res, fmt.Errorf("write app: %w", err)
	}
	_ = os.Remove(exePath)
	if err := os.Rename(tmp, exePath); err != nil {
		// fallback copy
		if err2 := os.WriteFile(exePath, opts.Payload, 0o755); err2 != nil {
			return res, fmt.Errorf("place %s (is it running?): %w", paths.AppExeName, err2)
		}
		_ = os.Remove(tmp)
	}

	// LICENSE / NOTICE next to app when present beside payload (optional)
	// Write a short credit file always.
	credit := fmt.Sprintf(
		"Win-Slate %s\n\nSlate by Sam Wasserman (Apache-2.0)\n%s\n\nThis Windows build is a derivative packaging.\nProjects live in: %s\n",
		paths.AppVersion, paths.UpstreamSlateURL, paths.ProjectsDir(),
	)
	_ = os.WriteFile(filepath.Join(dest, "README-CREDIT.txt"), []byte(credit), 0o644)

	res.InstallDir = dest
	res.ExePath = exePath

	progress("Shortcuts", "Creating Start Menu shortcut…", 90)
	res.StartMenuOK = createShortcut(paths.StartMenuShortcut(), exePath, dest) == nil
	if opts.DesktopShortcut {
		res.DesktopOK = createShortcut(paths.DesktopShortcut(), exePath, dest) == nil
	}

	progress("Smoke test", "Starting Win-Slate briefly to verify…", 94)
	res.SmokeOK, res.SmokeDetail = smokeTest(exePath)

	_ = manifest.Write(manifest.Manifest{
		Product:      paths.ProductName,
		Kind:         paths.InstallKind,
		AppVersion:   paths.AppVersion,
		SetupVersion: paths.SetupVersion,
		ReleaseTag:   res.ReleaseTag,
		InstallDir:   dest,
		ExePath:      exePath,
		InstalledAt:  time.Now().UTC().Format(time.RFC3339),
		SmokeOK:      res.SmokeOK,
	})
	_ = manifest.SaveConfig(dest)

	lines := []string{}
	if opts.IsUpdate {
		lines = append(lines, "Win-Slate was updated.")
	} else {
		lines = append(lines, "Win-Slate is installed.")
	}
	lines = append(lines, "Location: "+exePath)
	lines = append(lines, "Version: "+res.ReleaseTag)
	if res.StartMenuOK {
		lines = append(lines, "Start Menu: Win-Slate")
	}
	if res.SmokeOK {
		lines = append(lines, "Startup check: passed ("+res.SmokeDetail+")")
	} else {
		lines = append(lines, "Startup check: "+res.SmokeDetail)
	}
	lines = append(lines, "Projects NOT modified: "+paths.ProjectsDir())
	lines = append(lines, "Original Slate by Sam Wasserman — "+paths.UpstreamSlateURL)
	res.Summary = strings.Join(lines, "\n")
	progress(label, "Complete.", 100)
	return res, nil
}

func smokeTest(exe string) (bool, string) {
	cmd := exec.Command(exe)
	cmd.Dir = filepath.Dir(exe)
	if err := cmd.Start(); err != nil {
		return false, "could not start: " + err.Error()
	}
	done := make(chan error, 1)
	go func() { done <- cmd.Wait() }()
	select {
	case err := <-done:
		if err != nil {
			// Wails apps that exit immediately with tag error would fail here
			return false, "exited early: " + err.Error()
		}
		return true, "exited cleanly"
	case <-time.After(6 * time.Second):
		if cmd.Process != nil {
			_ = cmd.Process.Kill()
		}
		return true, "still running after 6s"
	}
}

func createShortcut(lnk, target, workDir string) error {
	_ = os.MkdirAll(filepath.Dir(lnk), 0o755)
	// PowerShell WScript shortcut
	ps := fmt.Sprintf(`
$ws = New-Object -ComObject WScript.Shell
$s = $ws.CreateShortcut(%q)
$s.TargetPath = %q
$s.WorkingDirectory = %q
$s.Description = "Win-Slate — by Sam Wasserman (Windows port)"
$s.Save()
`, lnk, target, workDir)
	cmd := exec.Command("powershell", "-NoProfile", "-ExecutionPolicy", "Bypass", "-Command", ps)
	out, err := cmd.CombinedOutput()
	if err != nil {
		logx.Log("shortcut %s: %s err=%v", lnk, string(out), err)
		return err
	}
	return nil
}
