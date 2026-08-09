// Package stills — Stills Library (Circle Take discovery + ffmpeg extraction).
// Parity with slate-0.3.2/src/main/stills.ts
package stills

import (
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/wassermanproductions/slate-windows/internal/types"
)

const maxStills = 8

type circleTakePaths struct {
	RecentsFile string
}

func defaultCircleTakePaths() circleTakePaths {
	home, _ := os.UserHomeDir()
	// Original mac path — also probe Windows suite locations if present.
	mac := filepath.Join(home, "Library", "Application Support", "circle-take", "recents.json")
	win := filepath.Join(home, "AppData", "Roaming", "circle-take", "recents.json")
	if _, err := os.Stat(win); err == nil {
		return circleTakePaths{RecentsFile: win}
	}
	return circleTakePaths{RecentsFile: mac}
}

// DiscoverCircledTakes scans Circle Take recents for circled takes.
func DiscoverCircledTakes() ([]types.CircledTake, error) {
	return discover(defaultCircleTakePaths())
}

func discover(paths circleTakePaths) ([]types.CircledTake, error) {
	raw, err := os.ReadFile(paths.RecentsFile)
	if err != nil {
		return []types.CircledTake{}, nil
	}
	var recents []struct {
		Path     string `json:"path"`
		OpenedAt string `json:"openedAt"`
	}
	if err := json.Unmarshal(raw, &recents); err != nil {
		return []types.CircledTake{}, nil
	}
	seen := map[string]bool{}
	var out []types.CircledTake
	for _, r := range recents {
		if r.Path == "" || seen[r.Path] {
			continue
		}
		if _, err := os.Stat(r.Path); err != nil {
			continue
		}
		seen[r.Path] = true
		docRaw, err := os.ReadFile(filepath.Join(r.Path, "project.json"))
		if err != nil {
			continue
		}
		var doc struct {
			Name  string `json:"name"`
			Shots []struct {
				ID   string `json:"id"`
				Name string `json:"name"`
			} `json:"shots"`
			Takes []struct {
				ID        string   `json:"id"`
				ShotID    string   `json:"shotId"`
				MediaPath string   `json:"mediaPath"`
				FileName  string   `json:"fileName"`
				Circled   bool     `json:"circled"`
				Rating    float64  `json:"rating"`
				InSec     *float64 `json:"inSec"`
				OutSec    *float64 `json:"outSec"`
			} `json:"takes"`
		}
		if err := json.Unmarshal(docRaw, &doc); err != nil {
			continue
		}
		shotName := map[string]string{}
		for _, s := range doc.Shots {
			if s.ID != "" {
				shotName[s.ID] = s.Name
			}
		}
		for _, t := range doc.Takes {
			if !t.Circled || t.MediaPath == "" {
				continue
			}
			media := t.MediaPath
			if !filepath.IsAbs(media) {
				media = filepath.Join(r.Path, media)
			}
			if _, err := os.Stat(media); err != nil {
				continue
			}
			var shot *string
			if t.ShotID != "" {
				if n, ok := shotName[t.ShotID]; ok {
					shot = &n
				}
			}
			fn := t.FileName
			if fn == "" {
				fn = filepath.Base(media)
			}
			proj := doc.Name
			if proj == "" {
				proj = strings.TrimSuffix(filepath.Base(r.Path), ".ctake")
			}
			out = append(out, types.CircledTake{
				Project:   proj,
				Shot:      shot,
				MediaPath: media,
				FileName:  fn,
				Rating:    t.Rating,
				InSec:     t.InSec,
				OutSec:    t.OutSec,
			})
		}
	}
	if out == nil {
		out = []types.CircledTake{}
	}
	return out, nil
}

func runFFmpeg(args ...string) error {
	cmd := exec.Command("ffmpeg", args...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		msg := string(out)
		if len(msg) > 400 {
			msg = msg[:400]
		}
		return fmt.Errorf("%s", msg)
	}
	return nil
}

// ExtractStills extracts up to maxStills frames into cacheRoot/stills.
func ExtractStills(cacheRoot, mediaPath string, inSec, outSec *float64) ([]string, error) {
	key := fmt.Sprintf("%s:%v:%v", mediaPath, inSec, outSec)
	h := sha1.Sum([]byte(key))
	hash := hex.EncodeToString(h[:])[:12]
	outDir := filepath.Join(cacheRoot, "stills", hash)
	if err := os.MkdirAll(outDir, 0o755); err != nil {
		return nil, err
	}
	if existing := listJPGs(outDir); len(existing) > 0 {
		return existing, nil
	}

	var windowArgs []string
	if inSec != nil && *inSec > 0 {
		windowArgs = append(windowArgs, "-ss", fmt.Sprintf("%v", *inSec))
	}
	if outSec != nil && *outSec > 0 {
		windowArgs = append(windowArgs, "-to", fmt.Sprintf("%v", *outSec))
	}

	args1 := append([]string{}, windowArgs...)
	args1 = append(args1,
		"-i", mediaPath,
		"-vf", "select='gt(scene,0.24)',scale=768:-2",
		"-vsync", "vfr",
		"-frames:v", fmt.Sprintf("%d", maxStills),
		"-q:v", "3",
		filepath.Join(outDir, "scene-%02d.jpg"),
	)
	_ = runFFmpeg(args1...)

	frames := listJPGs(outDir)
	if len(frames) < 4 {
		args2 := append([]string{}, windowArgs...)
		args2 = append(args2,
			"-i", mediaPath,
			"-vf", "fps=1,scale=768:-2",
			"-frames:v", fmt.Sprintf("%d", maxStills-len(frames)),
			"-q:v", "3",
			filepath.Join(outDir, "tick-%02d.jpg"),
		)
		_ = runFFmpeg(args2...)
		frames = listJPGs(outDir)
	}
	if len(frames) == 0 {
		return nil, fmt.Errorf("No frames could be extracted — is ffmpeg installed and on PATH?")
	}
	return frames, nil
}

func listJPGs(dir string) []string {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil
	}
	var out []string
	for _, e := range entries {
		if strings.HasSuffix(strings.ToLower(e.Name()), ".jpg") {
			out = append(out, filepath.Join(dir, e.Name()))
		}
	}
	for i := 0; i < len(out); i++ {
		for j := i + 1; j < len(out); j++ {
			if out[j] < out[i] {
				out[i], out[j] = out[j], out[i]
			}
		}
	}
	return out
}
