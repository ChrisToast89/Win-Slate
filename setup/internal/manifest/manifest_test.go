package manifest

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/ChrisToast89/Win-Slate/setup/internal/paths"
)

func TestIsOurs(t *testing.T) {
	if IsOurs(nil) {
		t.Fatal("nil")
	}
	if IsOurs(&Manifest{Product: "Slate"}) {
		t.Fatal("plain Slate product must not count as Win-Slate")
	}
	if !IsOurs(&Manifest{Kind: paths.InstallKind}) {
		t.Fatal("kind")
	}
	if !IsOurs(&Manifest{Product: "Win-Slate"}) {
		t.Fatal("product Win-Slate")
	}
	if !IsOurs(&Manifest{Product: "Slate for Windows"}) {
		t.Fatal("legacy product name")
	}
}

func TestDiscoverIgnoresNpmSlateLayout(t *testing.T) {
	// Simulate Programs\Slate with only Slate.exe and electron-style manifest — not ours.
	dir := t.TempDir()
	// Point config at a fake npm layout under temp by writing files and calling IsWinSlateInstallDir
	npmLike := filepath.Join(dir, "Programs", "Slate")
	if err := os.MkdirAll(npmLike, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(npmLike, "Slate.exe"), []byte("fake-npm"), 0o644); err != nil {
		t.Fatal(err)
	}
	// electron helper-style manifest without our kind
	raw, _ := json.Marshal(map[string]string{"sourceRef": "main", "exePath": filepath.Join(npmLike, "Slate.exe")})
	_ = os.WriteFile(filepath.Join(npmLike, "slate-install-manifest.json"), raw, 0o644)

	if IsWinSlateInstallDir(npmLike) {
		t.Fatal("npm/Electron Slate tree must not be treated as Win-Slate")
	}
}

func TestWinSlateInstallRequiresOurManifest(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, paths.AppExeName), []byte("bin"), 0o644); err != nil {
		t.Fatal(err)
	}
	if IsWinSlateInstallDir(dir) {
		t.Fatal("exe alone is not enough")
	}
	if err := Write(Manifest{
		InstallDir: dir,
		ExePath:    paths.InstalledExe(dir),
		AppVersion: "test",
	}); err != nil {
		t.Fatal(err)
	}
	if !IsWinSlateInstallDir(dir) {
		t.Fatal("expected ours after Write")
	}
	m, err := Read(dir)
	if err != nil || !IsOurs(m) || m.Kind != paths.InstallKind {
		t.Fatalf("manifest %#v err %v", m, err)
	}
}
