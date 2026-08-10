package install

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/ChrisToast89/Win-Slate/setup/internal/paths"
)

func TestInstallWritesExeAndManifest(t *testing.T) {
	dir := t.TempDir()
	// Fake payload (must be >= 1024 bytes)
	payload := make([]byte, 2048)
	for i := range payload {
		payload[i] = byte(i % 251)
	}
	res, err := Run(Options{
		InstallDir:      dir,
		DesktopShortcut: false,
		IsUpdate:        false,
		ReleaseTag:      "v-test",
		Payload:         payload,
	}, nil)
	if err != nil {
		t.Fatal(err)
	}
	if res.ExePath != paths.InstalledExe(dir) {
		t.Fatalf("exe path %s", res.ExePath)
	}
	st, err := os.Stat(res.ExePath)
	if err != nil || st.Size() != 2048 {
		t.Fatalf("exe missing or wrong size: %v %v", err, st)
	}
	if _, err := os.Stat(filepath.Join(dir, paths.ManifestName())); err != nil {
		t.Fatal("manifest missing", err)
	}
	// Safety: refuse projects dir
	_, err = Run(Options{
		InstallDir: paths.ProjectsDir(),
		Payload:    payload,
	}, nil)
	if err == nil {
		t.Fatal("expected safety error for projects dir")
	}
}

func TestRejectTinyPayload(t *testing.T) {
	_, err := Run(Options{
		InstallDir: t.TempDir(),
		Payload:    []byte("tiny"),
	}, nil)
	if err == nil {
		t.Fatal("expected error")
	}
}
