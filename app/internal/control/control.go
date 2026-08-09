// Package control — localhost MCP control server with bearer token.
// Parity with slate-0.3.2/src/main/control.ts
package control

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/wassermanproductions/slate-windows/internal/brain"
	"github.com/wassermanproductions/slate-windows/internal/projects"
	"github.com/wassermanproductions/slate-windows/internal/types"
)

type Notify func()

var token string

func init() {
	b := make([]byte, 24)
	_, _ = rand.Read(b)
	token = hex.EncodeToString(b)
}

func descriptorPath() string {
	if base := os.Getenv("APPDATA"); base != "" {
		return filepath.Join(base, "slate", "control.json")
	}
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".config", "slate", "control.json")
}

func findShot(p *types.Project, shotID string) (*types.Scene, *types.Shot) {
	for si := range p.Scenes {
		for shi := range p.Scenes[si].Shots {
			if p.Scenes[si].Shots[shi].ID == shotID {
				return &p.Scenes[si], &p.Scenes[si].Shots[shi]
			}
		}
	}
	return nil, nil
}

func invokeTool(tool string, args map[string]interface{}, notify Notify) (interface{}, error) {
	str := func(k string) string {
		if v, ok := args[k]; ok && v != nil {
			return fmt.Sprint(v)
		}
		return ""
	}

	switch tool {
	case "list_projects":
		return projects.List()
	case "get_project":
		p, err := projects.Open(str("projectId"))
		if err != nil {
			return nil, err
		}
		if p == nil {
			return nil, fmt.Errorf("Project not found")
		}
		return p, nil
	case "create_project":
		name := str("name")
		if name == "" {
			name = "Untitled"
		}
		p, err := projects.Create(name)
		if err != nil {
			return nil, err
		}
		if notify != nil {
			notify()
		}
		return map[string]string{"id": p.ID, "name": p.Name}, nil
	case "list_shots":
		p, err := projects.Open(str("projectId"))
		if err != nil {
			return nil, err
		}
		if p == nil {
			return nil, fmt.Errorf("Project not found")
		}
		var out []map[string]interface{}
		for _, sc := range p.Scenes {
			shots := make([]map[string]interface{}, 0, len(sc.Shots))
			for _, s := range sc.Shots {
				shots = append(shots, map[string]interface{}{
					"id": s.ID, "name": s.Name, "intent": s.Intent, "spec": s.Spec,
				})
			}
			out = append(out, map[string]interface{}{
				"sceneId": sc.ID, "scene": sc.Name, "shots": shots,
			})
		}
		return out, nil
	case "get_shot_prompt":
		p, err := projects.Open(str("projectId"))
		if err != nil {
			return nil, err
		}
		if p == nil {
			return nil, fmt.Errorf("Project not found")
		}
		_, shot := findShot(p, str("shotId"))
		if shot == nil {
			return nil, fmt.Errorf("Shot not found")
		}
		return map[string]interface{}{
			"prompt": shot.Prompt, "spec": shot.Spec, "beatSheet": shot.BeatSheet,
		}, nil
	case "set_shot_prompt":
		p, err := projects.Open(str("projectId"))
		if err != nil {
			return nil, err
		}
		if p == nil {
			return nil, fmt.Errorf("Project not found")
		}
		_, shot := findShot(p, str("shotId"))
		if shot == nil {
			return nil, fmt.Errorf("Shot not found")
		}
		shot.History = append([]types.PromptVersion{{
			ID:      fmt.Sprintf("v%d", time.Now().UnixMilli()),
			SavedAt: time.Now().UTC().Format(time.RFC3339Nano),
			Label:   "external edit (MCP)",
			Prompt:  shot.Prompt,
		}}, shot.History...)
		shot.Prompt = str("prompt")
		shot.UpdatedAt = time.Now().UTC().Format(time.RFC3339Nano)
		if err := projects.Save(p); err != nil {
			return nil, err
		}
		if notify != nil {
			notify()
		}
		return map[string]bool{"ok": true}, nil
	case "add_scene":
		p, err := projects.Open(str("projectId"))
		if err != nil {
			return nil, err
		}
		if p == nil {
			return nil, fmt.Errorf("Project not found")
		}
		name := str("name")
		if name == "" {
			name = fmt.Sprintf("Scene %d", len(p.Scenes)+1)
		}
		scene := types.Scene{
			ID:       fmt.Sprintf("scene-%d", time.Now().UnixMilli()),
			Name:     name,
			Synopsis: str("synopsis"),
			Shots:    []types.Shot{},
		}
		p.Scenes = append(p.Scenes, scene)
		if err := projects.Save(p); err != nil {
			return nil, err
		}
		if notify != nil {
			notify()
		}
		return map[string]string{"sceneId": scene.ID}, nil
	case "list_characters":
		p, err := projects.Open(str("projectId"))
		if err != nil {
			return nil, err
		}
		if p == nil {
			return nil, fmt.Errorf("Project not found")
		}
		return p.Characters, nil
	case "list_locations":
		p, err := projects.Open(str("projectId"))
		if err != nil {
			return nil, err
		}
		if p == nil {
			return nil, fmt.Errorf("Project not found")
		}
		return p.Locations, nil
	case "list_lookbook":
		p, err := projects.Open(str("projectId"))
		if err != nil {
			return nil, err
		}
		if p == nil {
			return nil, fmt.Errorf("Project not found")
		}
		return p.Lookbook, nil
	case "brain_status":
		return brain.Status(""), nil
	case "list_local_models":
		ep := str("endpoint")
		return brain.DetectLocal(ep), nil
	case "set_brain":
		p, err := projects.Open(str("projectId"))
		if err != nil {
			return nil, err
		}
		if p == nil {
			return nil, fmt.Errorf("Project not found")
		}
		b := str("brain")
		if b != "claude" && b != "codex" && b != "local" {
			return nil, fmt.Errorf("brain must be one of: claude, codex, local")
		}
		p.Defaults.Brain = b
		if _, ok := args["endpoint"]; ok {
			p.Defaults.LocalEndpoint = str("endpoint")
		}
		if _, ok := args["model"]; ok {
			p.Defaults.LocalModel = str("model")
		}
		if err := projects.Save(p); err != nil {
			return nil, err
		}
		if notify != nil {
			notify()
		}
		return map[string]interface{}{"ok": true, "brain": p.Defaults.Brain}, nil
	default:
		return nil, fmt.Errorf("Unknown tool: %s", tool)
	}
}

func toolCatalog() []map[string]string {
	return []map[string]string{
		{"name": "list_projects", "description": "List Slate projects with scene/shot counts"},
		{"name": "get_project", "description": "Get a full Slate project (args: projectId)"},
		{"name": "create_project", "description": "Create a new Slate project (args: name)"},
		{"name": "list_shots", "description": "List scenes and shots (args: projectId)"},
		{"name": "get_shot_prompt", "description": "Get a shot prompt + spec (args: projectId, shotId)"},
		{"name": "set_shot_prompt", "description": "Set a shot prompt, versioning the old one (args: projectId, shotId, prompt)"},
		{"name": "add_scene", "description": "Add a scene (args: projectId, name, synopsis)"},
		{"name": "list_characters", "description": "List character sheets (args: projectId)"},
		{"name": "list_locations", "description": "List location sheets (args: projectId)"},
		{"name": "list_lookbook", "description": "List style profiles (args: projectId)"},
		{"name": "brain_status", "description": "Report which brains are available: Claude Code, Codex, and any local model server"},
		{"name": "list_local_models", "description": "List models on the local OpenAI-compatible server (args: endpoint?)"},
		{"name": "set_brain", "description": "Set a project brain: claude | codex | local (args: projectId, brain, endpoint?, model?)"},
	}
}

// Start launches the control server and writes control.json for the MCP bridge.
func Start(notify Notify) error {
	mux := http.NewServeMux()
	mux.HandleFunc("/tools", func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") != "Bearer "+token {
			http.Error(w, `{"error":"unauthorized"}`, http.StatusUnauthorized)
			return
		}
		if r.Method != http.MethodGet {
			http.Error(w, `{"error":"method"}`, http.StatusMethodNotAllowed)
			return
		}
		_ = json.NewEncoder(w).Encode(map[string]interface{}{"tools": toolCatalog()})
	})
	mux.HandleFunc("/invoke", func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") != "Bearer "+token {
			http.Error(w, `{"error":"unauthorized"}`, http.StatusUnauthorized)
			return
		}
		if r.Method != http.MethodPost {
			http.Error(w, `{"error":"method"}`, http.StatusMethodNotAllowed)
			return
		}
		body, _ := io.ReadAll(r.Body)
		var req struct {
			Tool string                 `json:"tool"`
			Args map[string]interface{} `json:"args"`
		}
		if err := json.Unmarshal(body, &req); err != nil {
			w.WriteHeader(400)
			_ = json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}
		if req.Args == nil {
			req.Args = map[string]interface{}{}
		}
		result, err := invokeTool(req.Tool, req.Args, notify)
		if err != nil {
			w.WriteHeader(400)
			_ = json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}
		_ = json.NewEncoder(w).Encode(map[string]interface{}{"result": result})
	})

	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		return err
	}
	port := ln.Addr().(*net.TCPAddr).Port

	desc := descriptorPath()
	if err := os.MkdirAll(filepath.Dir(desc), 0o755); err != nil {
		return err
	}
	meta := map[string]interface{}{
		"v": 1, "app": "slate", "port": port, "token": token, "pid": os.Getpid(),
	}
	raw, _ := json.MarshalIndent(meta, "", "  ")
	if err := os.WriteFile(desc, raw, 0o600); err != nil {
		return err
	}

	go func() {
		_ = http.Serve(ln, mux)
	}()
	return nil
}
