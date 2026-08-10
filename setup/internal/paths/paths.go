package paths

import (
	"os"
	"path/filepath"
	"strings"
)

// Product identity — Windows packaging of Sam Wasserman's Slate.
const (
	ProductName       = "Win-Slate"
	AppExeName        = "Slate.exe"
	SetupVersion      = "0.3.2-win.1"
	AppVersion        = "0.3.2-win.1"
	GitHubOwner       = "ChrisToast89"
	GitHubRepo        = "Win-Slate"
	UpstreamSlateURL  = "https://github.com/wassermanproductions/slate"
	UpstreamAuthor    = "Sam Wasserman"
	GitHubReleasesAPI = "https://api.github.com/repos/" + GitHubOwner + "/" + GitHubRepo + "/releases/latest"
	GitHubRepoURL     = "https://github.com/" + GitHubOwner + "/" + GitHubRepo
)

// ConfigDir stores last chosen install path (per-user).
func ConfigDir() string {
	return filepath.Join(LocalAppData(), "Win-Slate")
}

func ConfigPath() string {
	return filepath.Join(ConfigDir(), "setup-config.json")
}

func LocalAppData() string {
	if d := os.Getenv("LOCALAPPDATA"); d != "" {
		return d
	}
	home, _ := os.UserHomeDir()
	return filepath.Join(home, "AppData", "Local")
}

// DefaultInstallDir is the suggested per-user install location.
func DefaultInstallDir() string {
	return filepath.Join(LocalAppData(), "Programs", "Win-Slate")
}

func ManifestName() string { return "install-manifest.json" }

func ManifestPath(installDir string) string {
	return filepath.Join(installDir, ManifestName())
}

func InstalledExe(installDir string) string {
	return filepath.Join(installDir, AppExeName)
}

// ProjectsDir is where the app stores user projects — NEVER touch.
func ProjectsDir() string {
	if d := os.Getenv("SLATE_DATA_DIR"); d != "" {
		return d
	}
	home, _ := os.UserHomeDir()
	return filepath.Join(home, "Documents", "Slate")
}

func TempLog() string {
	return filepath.Join(os.TempDir(), "win-slate-setup.log")
}

func StartMenuShortcut() string {
	programs := filepath.Join(os.Getenv("APPDATA"), "Microsoft", "Windows", "Start Menu", "Programs")
	return filepath.Join(programs, "Win-Slate.lnk")
}

func DesktopShortcut() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, "Desktop", "Win-Slate.lnk")
}

func IsProtectedPath(p string) bool {
	p = filepath.Clean(p)
	protected := []string{
		ProjectsDir(),
		filepath.Join(os.Getenv("USERPROFILE"), "Documents", "Slate"),
	}
	for _, root := range protected {
		root = filepath.Clean(root)
		if root == "" || root == "." {
			continue
		}
		if strings.EqualFold(p, root) {
			return true
		}
		rel, err := filepath.Rel(root, p)
		if err == nil && rel != ".." && !strings.HasPrefix(rel, ".."+string(os.PathSeparator)) {
			return true
		}
	}
	return false
}

func AssertSafeInstallDir(installDir string) error {
	inst := filepath.Clean(installDir)
	proj := filepath.Clean(ProjectsDir())
	if inst == "" || inst == "." {
		return errProtected("install directory is empty")
	}
	if strings.EqualFold(inst, proj) {
		return errProtected("install directory equals projects directory")
	}
	if IsProtectedPath(inst) {
		return errProtected("install directory is under protected projects path")
	}
	if rel, err := filepath.Rel(inst, proj); err == nil && rel != ".." && !strings.HasPrefix(rel, ".."+string(os.PathSeparator)) {
		return errProtected("projects appear to live under install directory")
	}
	// Refuse system roots
	for _, bad := range []string{`C:\`, `C:\Windows`, `C:\Windows\System32`} {
		if strings.EqualFold(inst, filepath.Clean(bad)) {
			return errProtected("refusing to install into system path")
		}
	}
	return nil
}

type protectError string

func (e protectError) Error() string { return string(e) }

func errProtected(msg string) error {
	return protectError("SAFETY STOP: " + msg + ". User project files were not touched.")
}
