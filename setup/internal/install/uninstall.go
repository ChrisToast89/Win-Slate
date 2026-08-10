package install

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/ChrisToast89/Win-Slate/setup/internal/logx"
	"github.com/ChrisToast89/Win-Slate/setup/internal/manifest"
	"github.com/ChrisToast89/Win-Slate/setup/internal/paths"
)

// UninstallResult is returned after removing a Win-Slate install.
type UninstallResult struct {
	OK              bool   `json:"ok"`
	RemovedDir      string `json:"removedDir"`
	ProjectsDir     string `json:"projectsDir"`
	ProjectsPreserved bool `json:"projectsPreserved"`
	Summary         string `json:"summary"`
}

// Uninstall removes only a verified Win-Slate install directory and shortcuts.
// Never touches Documents\Slate projects or the npm/Electron Programs\Slate tree.
func Uninstall(progress ProgressFn) (UninstallResult, error) {
	if progress == nil {
		progress = func(string, string, int) {}
	}
	var res UninstallResult
	res.ProjectsDir = paths.ProjectsDir()
	res.ProjectsPreserved = true

	dir, _, ok := manifest.Discover()
	if !ok || dir == "" {
		return res, fmt.Errorf("no Win-Slate installation found")
	}
	if !manifest.IsWinSlateInstallDir(dir) {
		return res, fmt.Errorf("refusing to uninstall: %s is not a Win-Slate install", dir)
	}
	if paths.IsProtectedPath(dir) || paths.IsNpmSlateInstallDir(dir) {
		return res, fmt.Errorf("SAFETY STOP: install path is protected — not removed")
	}
	// Refuse anything that looks like Documents\Slate
	if strings.Contains(strings.ToLower(filepath.Clean(dir)), strings.ToLower(filepath.Join("documents", "slate"))) {
		return res, fmt.Errorf("SAFETY STOP: path looks like projects folder")
	}

	progress("Uninstall", "Removing Start Menu / Desktop shortcuts…", 30)
	_ = os.Remove(paths.StartMenuShortcut())
	_ = os.Remove(paths.DesktopShortcut())

	progress("Uninstall", "Removing program files…", 70)
	logx.Log("Uninstall removing %s", dir)
	if err := os.RemoveAll(dir); err != nil {
		return res, fmt.Errorf("could not remove install folder: %w", err)
	}

	// Clear last-install pointer in Win-Slate setup config
	_ = manifest.SaveConfig("")
	// Or rewrite empty config — SaveConfig with empty might still write; better remove config installDir
	cfgPath := paths.ConfigPath()
	if raw, err := os.ReadFile(cfgPath); err == nil && len(raw) > 0 {
		_ = os.WriteFile(cfgPath, []byte("{}\n"), 0o644)
	}

	res.OK = true
	res.RemovedDir = dir
	res.Summary = fmt.Sprintf(
		"Win-Slate was uninstalled from:\n%s\n\nYour projects were NOT removed:\n%s",
		dir, paths.ProjectsDir(),
	)
	progress("Uninstall", "Complete.", 100)
	return res, nil
}
