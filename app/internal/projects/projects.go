// Package projects — one folder per project under ~/Documents/Slate.
// Parity with slate-0.3.2/src/main/projects.ts
package projects

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sort"
	"time"

	"github.com/google/uuid"
	"github.com/wassermanproductions/slate-windows/internal/types"
)

func Root() string {
	if d := os.Getenv("SLATE_DATA_DIR"); d != "" {
		return d
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return filepath.Join(".", "Slate")
	}
	// Match original: Documents/Slate
	return filepath.Join(home, "Documents", "Slate")
}

func projectDir(id string) string {
	return filepath.Join(Root(), id)
}

func projectFile(id string) string {
	return filepath.Join(projectDir(id), "project.json")
}

func CacheDir(id string) string {
	return filepath.Join(projectDir(id), ".cache")
}

func ensureRoot() error {
	return os.MkdirAll(Root(), 0o755)
}

func List() ([]types.ProjectMeta, error) {
	if err := ensureRoot(); err != nil {
		return nil, err
	}
	entries, err := os.ReadDir(Root())
	if err != nil {
		return nil, err
	}
	var metas []types.ProjectMeta
	for _, e := range entries {
		if !e.IsDir() {
			continue
		}
		file := projectFile(e.Name())
		if _, err := os.Stat(file); err != nil {
			continue
		}
		raw, err := os.ReadFile(file)
		if err != nil {
			continue
		}
		var p types.Project
		if err := json.Unmarshal(raw, &p); err != nil {
			continue
		}
		shots := 0
		for _, s := range p.Scenes {
			shots += len(s.Shots)
		}
		metas = append(metas, types.ProjectMeta{
			ID:         p.ID,
			Name:       p.Name,
			Logline:    p.Logline,
			Path:       projectDir(p.ID),
			UpdatedAt:  p.UpdatedAt,
			SceneCount: len(p.Scenes),
			ShotCount:  shots,
		})
	}
	sort.Slice(metas, func(i, j int) bool {
		return metas[i].UpdatedAt > metas[j].UpdatedAt
	})
	if metas == nil {
		metas = []types.ProjectMeta{}
	}
	return metas, nil
}

func NewProject(name string) types.Project {
	now := time.Now().UTC().Format(time.RFC3339Nano)
	return types.Project{
		ID:      uuid.NewString(),
		Name:    name,
		Logline: "",
		World:   "",
		Defaults: types.ProjectDefaults{
			AspectRatio: "16:9",
			FPS:         24,
			DurationSec: 8,
			TargetModel: "seedance-2",
			Brain:       "claude",
		},
		Scenes:     []types.Scene{},
		Characters: []types.CharacterSheet{},
		ArtDept:    []types.ArtDeptSheet{},
		Locations:  []types.LocationSheet{},
		Lookbook:   []types.StyleProfile{},
		References: []types.Reference{},
		MySetups:   []types.CustomSetup{},
		Music:      []types.MusicCue{},
		Voices:     []types.VoiceSheet{},
		CreatedAt:  now,
		UpdatedAt:  now,
	}
}

func Create(name string) (*types.Project, error) {
	if err := ensureRoot(); err != nil {
		return nil, err
	}
	p := NewProject(name)
	if err := os.MkdirAll(projectDir(p.ID), 0o755); err != nil {
		return nil, err
	}
	if err := Save(&p); err != nil {
		return nil, err
	}
	return &p, nil
}

func Open(id string) (*types.Project, error) {
	raw, err := os.ReadFile(projectFile(id))
	if err != nil {
		return nil, nil // parity: null when missing
	}
	var p types.Project
	if err := json.Unmarshal(raw, &p); err != nil {
		return nil, err
	}
	return &p, nil
}

func Save(project *types.Project) error {
	project.UpdatedAt = time.Now().UTC().Format(time.RFC3339Nano)
	dir := projectDir(project.ID)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return err
	}
	raw, err := json.MarshalIndent(project, "", "  ")
	if err != nil {
		return err
	}
	tmp := filepath.Join(dir, ".project."+time.Now().Format("20060102150405.000")+".tmp")
	if err := os.WriteFile(tmp, raw, 0o644); err != nil {
		return err
	}
	dest := projectFile(project.ID)
	// Windows: rename over existing can fail — remove first.
	_ = os.Remove(dest)
	return os.Rename(tmp, dest)
}

func Delete(id string) error {
	return os.RemoveAll(projectDir(id))
}

func ProjectPath(id string) string {
	return projectFile(id)
}
