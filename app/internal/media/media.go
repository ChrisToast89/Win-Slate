// Package media — ffmpeg frame extraction and media kind detection.
// Parity with slate-0.3.2/src/main/ingest.ts
package media

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/wassermanproductions/slate-windows/internal/projects"
)

var imageExt = map[string]bool{
	".png": true, ".jpg": true, ".jpeg": true, ".webp": true, ".gif": true,
	".tif": true, ".tiff": true, ".bmp": true, ".heic": true,
}
var videoExt = map[string]bool{
	".mp4": true, ".mov": true, ".m4v": true, ".webm": true, ".mkv": true, ".avi": true, ".mxf": true,
}

const maxFrames = 16

func runFFmpeg(args ...string) error {
	cmd := exec.Command("ffmpeg", args...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		msg := string(out)
		if len(msg) > 400 {
			msg = msg[len(msg)-400:]
		}
		return fmt.Errorf("%s", msg)
	}
	return nil
}

func Available() bool {
	cmd := exec.Command("ffmpeg", "-version")
	return cmd.Run() == nil
}

func Kind(path string) string {
	ext := strings.ToLower(filepath.Ext(path))
	if imageExt[ext] {
		return "image"
	}
	if videoExt[ext] {
		return "video"
	}
	return ""
}

// ExtractFrames pulls representative stills from a video into the project cache.
func ExtractFrames(projectID, videoPath string) ([]string, error) {
	h := sha1.Sum([]byte(videoPath))
	hash := hex.EncodeToString(h[:])[:12]
	outDir := filepath.Join(projects.CacheDir(projectID), "frames", hash)
	if err := os.MkdirAll(outDir, 0o755); err != nil {
		return nil, err
	}
	if existing := listJPGs(outDir); len(existing) > 0 {
		return existing, nil
	}

	_ = runFFmpeg(
		"-i", videoPath,
		"-vf", "select='gt(scene,0.28)',scale=768:-2",
		"-vsync", "vfr",
		"-frames:v", fmt.Sprintf("%d", maxFrames),
		"-q:v", "4",
		filepath.Join(outDir, "cut-%03d.jpg"),
	)

	frames := listJPGs(outDir)
	if len(frames) < 4 {
		_ = runFFmpeg(
			"-i", videoPath,
			"-vf", "fps=1/2,scale=768:-2",
			"-frames:v", fmt.Sprintf("%d", maxFrames-len(frames)),
			"-q:v", "4",
			filepath.Join(outDir, "sample-%03d.jpg"),
		)
		frames = listJPGs(outDir)
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
	// sort by name
	for i := 0; i < len(out); i++ {
		for j := i + 1; j < len(out); j++ {
			if out[j] < out[i] {
				out[i], out[j] = out[j], out[i]
			}
		}
	}
	return out
}

// ResolveFFmpeg finds ffmpeg on PATH (used by callers for diagnostics).
func ResolveFFmpeg() string {
	if p, err := exec.LookPath("ffmpeg"); err == nil {
		return p
	}
	return "ffmpeg"
}

// Touch keeps go compiler happy about unused imports in edge builds.
var _ = time.Now
