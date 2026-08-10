package manifest

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"

	"github.com/ChrisToast89/Win-Slate/setup/internal/paths"
)

// Manifest records what Setup installed (for audit + update checks).
type Manifest struct {
	Product        string `json:"product"`
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
	m.UpstreamCredit = paths.UpstreamAuthor + " — " + paths.UpstreamSlateURL
	raw, err := json.MarshalIndent(m, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(paths.ManifestPath(m.InstallDir), raw, 0o644)
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

// Discover finds an existing install via config, default path, or common locations.
func Discover() (installDir string, m *Manifest, ok bool) {
	// 1) Last path from setup config
	if cfg, err := readConfig(); err == nil && cfg.InstallDir != "" {
		if st, err := os.Stat(paths.InstalledExe(cfg.InstallDir)); err == nil && !st.IsDir() {
			mm, _ := Read(cfg.InstallDir)
			return cfg.InstallDir, mm, true
		}
	}
	// 2) Default (Win-Slate)
	def := paths.DefaultInstallDir()
	if st, err := os.Stat(paths.InstalledExe(def)); err == nil && !st.IsDir() {
		mm, _ := Read(def)
		return def, mm, true
	}
	// 3) Previous product folder name
	for _, legacyName := range []string{"Slate for Windows", "Slate"} {
		legacy := filepath.Join(paths.LocalAppData(), "Programs", legacyName)
		if st, err := os.Stat(paths.InstalledExe(legacy)); err == nil && !st.IsDir() {
			mm, _ := Read(legacy)
			return legacy, mm, true
		}
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
	raw, _ := json.MarshalIndent(setupConfig{InstallDir: installDir}, "", "  ")
	return os.WriteFile(paths.ConfigPath(), raw, 0o644)
}

func VersionLooksNewer(remote, local string) bool {
	// Strip leading v; simple string compare on semver-ish tags after normalize.
	r := strings.TrimPrefix(strings.TrimSpace(remote), "v")
	l := strings.TrimPrefix(strings.TrimSpace(local), "v")
	if r == "" || l == "" {
		return r != "" && r != l
	}
	return r != l
}
