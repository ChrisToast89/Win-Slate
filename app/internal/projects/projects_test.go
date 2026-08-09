package projects

import (
	"os"
	"path/filepath"
	"testing"
)

func TestProjectCRUD(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("SLATE_DATA_DIR", dir)

	p, err := Create("Test Film")
	if err != nil {
		t.Fatal(err)
	}
	if p.Name != "Test Film" {
		t.Fatalf("name: %s", p.Name)
	}
	if p.Defaults.TargetModel != "seedance-2" {
		t.Fatalf("default model: %s", p.Defaults.TargetModel)
	}

	metas, err := List()
	if err != nil {
		t.Fatal(err)
	}
	if len(metas) != 1 {
		t.Fatalf("want 1 project, got %d", len(metas))
	}

	p.Logline = "A neon noir"
	if err := Save(p); err != nil {
		t.Fatal(err)
	}

	opened, err := Open(p.ID)
	if err != nil {
		t.Fatal(err)
	}
	if opened == nil || opened.Logline != "A neon noir" {
		t.Fatalf("open/save failed: %+v", opened)
	}

	// Atomic file exists
	if _, err := os.Stat(filepath.Join(dir, p.ID, "project.json")); err != nil {
		t.Fatal(err)
	}

	if err := Delete(p.ID); err != nil {
		t.Fatal(err)
	}
	metas, _ = List()
	if len(metas) != 0 {
		t.Fatalf("after delete want 0, got %d", len(metas))
	}
}
