package manifest

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"

	"github.com/ChrisToast89/Win-Slate/setup/internal/paths"
)

// Manifest records what Win-Slate Setup installed (for audit + update checks).
// Kind + Product identify this Wails port so we never confuse it with the
// separate npm/Electron install of Sam Wasserman's Slate.
type Manifest struct {
	Product        string `json:"product"`
	Kind           string `json:"kind"` // paths.InstallKind
	AppVersion     string `json:"appVersion"`
	SetupVersion   string `json:"setupVersion"`
	ReleaseTag     string `json:"releaseTag,omitempty"`
	InstallDir     string `json:"installDir"`
	ExePath        string `json:"exePath"`
	InstalledAt    string `json:"installedAt"`
	SmokeOK        bool   `json:"smokeOk"`
	UpstreamCredit string `json:"upstreamCredit"`
}

func Write(m Manifest) error {
	if err := paths.AssertSafeInstallDir(m.InstallDir); err != nil {
		return err
	}
	m.Product = paths.ProductName
	m.Kind = paths.InstallKind
	m.UpstreamCredit = paths.UpstreamAuthor + " — " + paths.UpstreamSlateURL
	raw, err := json.MarshalIndent(m, "", "  ")
	if err != nil {
		return err
	}
	if err := os.WriteFile(paths.ManifestPath(m.InstallDir), raw, 0o644); err != nil {
		return err
	}
	// Marker file: presence alone is not enough for Discover, but helps humans.
	_ = os.WriteFile(filepath.Join(m.InstallDir, paths.MarkerFile()), []byte(paths.InstallKind+"\n"), 0o644)
	return nil
}

func Read(installDir string) (*Manifest, error) {
	raw, err := os.ReadFile(paths.ManifestPath(installDir))
	if err != nil {
		return nil, err
	}
	var m Manifest
	if err := json.Unmarshal(raw, &m); err != nil {
		return nil, err
	}
	return &m, nil
}

// IsOurs reports whether a manifest belongs to Win-Slate (not npm Slate).
func IsOurs(m *Manifest) bool {
	if m == nil {
		return false
	}
	if m.Kind == paths.InstallKind {
		return true
	}
	// Legacy manifests written before Kind existed: product name only.
	// Never accept empty product — that would match random folders.
	p := strings.TrimSpace(m.Product)
	return p == paths.ProductName || p == "Slate for Windows"
}

// HasWinSlateExe returns true if our binary (or transitional names) is present.
func HasWinSlateExe(installDir string) bool {
	cands := []string{
		paths.InstalledExe(installDir),
		// Transitional: early Win-Slate builds still named the exe Slate.exe
		// but only count when our manifest is also present (checked by caller).
		filepath.Join(installDir, "Slate.exe"),
	}
	for _, c := range cands {
		if st, err := os.Stat(c); err == nil && !st.IsDir() {
			return true
		}
	}
	return false
}

// IsWinSlateInstallDir is true only when this folder is a real Win-Slate install.
func IsWinSlateInstallDir(installDir string) bool {
	if installDir == "" || paths.IsNpmSlateInstallDir(installDir) {
		return false
	}
	m, err := Read(installDir)
	if err != nil || !IsOurs(m) {
		return false
	}
	return HasWinSlateExe(installDir)
}

// Discover finds an existing *Win-Slate* install only (never npm/Electron Slate).
func Discover() (installDir string, m *Manifest, ok bool) {
	// 1) Last path from Win-Slate setup config
	if cfg, err := readConfig(); err == nil && cfg.InstallDir != "" {
		if IsWinSlateInstallDir(cfg.InstallDir) {
			mm, _ := Read(cfg.InstallDir)
			return cfg.InstallDir, mm, true
		}
	}
	// 2) Default Win-Slate directory
	def := paths.DefaultInstallDir()
	if IsWinSlateInstallDir(def) {
		mm, _ := Read(def)
		return def, mm, true
	}
	// 3) Previous Win-Slate product folder name only (not Programs\Slate)
	legacy := filepath.Join(paths.LocalAppData(), "Programs", "Slate for Windows")
	if IsWinSlateInstallDir(legacy) {
		mm, _ := Read(legacy)
		return legacy, mm, true
	}
	return "", nil, false
}

type setupConfig struct {
	InstallDir string `json:"installDir"`
}

func readConfig() (setupConfig, error) {
	var c setupConfig
	raw, err := os.ReadFile(paths.ConfigPath())
	if err != nil {
		return c, err
	}
	err = json.Unmarshal(raw, &c)
	return c, err
}

func SaveConfig(installDir string) error {
	_ = os.MkdirAll(paths.ConfigDir(), 0o755)
	if strings.TrimSpace(installDir) == "" {
		// Clear last install pointer after uninstall.
		return os.WriteFile(paths.ConfigPath(), []byte("{}\n"), 0o644)
	}
	raw, _ := json.MarshalIndent(setupConfig{InstallDir: installDir}, "", "  ")
	return os.WriteFile(paths.ConfigPath(), raw, 0o644)
}

func VersionLooksNewer(remote, local string) bool {
	r := strings.TrimPrefix(strings.TrimSpace(remote), "v")
	l := strings.TrimPrefix(strings.TrimSpace(local), "v")
	if r == "" || l == "" {
		return r != "" && r != l
	}
	return r != l
}
