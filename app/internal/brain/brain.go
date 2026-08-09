// Package brain runs Claude Code / Codex CLIs or any OpenAI-compatible local server.
// Parity with slate-0.3.2/src/main/brain.ts — no API keys stored.
package brain

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/wassermanproductions/slate-windows/internal/types"
)

var (
	mu            sync.Mutex
	running       = map[string]*exec.Cmd{}
	runningCancel = map[string]context.CancelFunc{}
)

var localCandidates = []string{
	"http://localhost:11434/v1",
	"http://localhost:1234/v1",
	"http://localhost:8000/v1",
	"http://localhost:8080/v1",
}

var imageMIME = map[string]string{
	".png":  "image/png",
	".jpg":  "image/jpeg",
	".jpeg": "image/jpeg",
	".webp": "image/webp",
	".gif":  "image/gif",
}

var claudeTierModel = map[string]string{
	"fast":     "haiku",
	"standard": "sonnet",
	// "top": user's default
}

// resolveBatchShim follows npm-style .cmd/.bat wrappers to a real .exe when possible.
// Spawning the wrapper goes through cmd.exe and flashes a console under a GUI host
// (and adds a full shell process per First AD turn).
func resolveBatchShim(path string) string {
	ext := strings.ToLower(filepath.Ext(path))
	if ext != ".cmd" && ext != ".bat" {
		return path
	}
	raw, err := os.ReadFile(path)
	if err != nil {
		return path
	}
	content := string(raw)
	dir := filepath.Dir(path)
	patterns := []*regexp.Regexp{
		// "%dp0%\node_modules\…\claude.exe"
		regexp.MustCompile(`(?i)"%dp0%\\([^"]+\.exe)"`),
		// "%~dp0node_modules\…\claude.exe"
		regexp.MustCompile(`(?i)"%~dp0([^"]+\.exe)"`),
		// unquoted %dp0%\…\tool.exe
		regexp.MustCompile(`(?i)%dp0%\\(\S+\.exe)`),
	}
	for _, re := range patterns {
		if m := re.FindStringSubmatch(content); len(m) == 2 {
			cand := filepath.Clean(filepath.Join(dir, m[1]))
			if st, err := os.Stat(cand); err == nil && !st.IsDir() {
				return cand
			}
		}
	}
	return path
}

func preferRealExe(path string) string {
	if path == "" {
		return path
	}
	return resolveBatchShim(path)
}

func resolveCLI(name string) string {
	// Prefer a real .exe over npm .cmd/.ps1 shims so we never shell out via cmd.exe.
	if p, err := exec.LookPath(name + ".exe"); err == nil {
		return preferRealExe(p)
	}
	if p, err := exec.LookPath(name); err == nil {
		return preferRealExe(p)
	}
	home, _ := os.UserHomeDir()
	npm := filepath.Join(home, "AppData", "Roaming", "npm")
	var candidates []string
	// Known npm package layouts (name-specific — never share paths across CLIs).
	switch name {
	case "claude":
		candidates = append(candidates,
			filepath.Join(npm, "node_modules", "@anthropic-ai", "claude-code", "bin", "claude.exe"),
		)
	case "codex":
		candidates = append(candidates,
			filepath.Join(npm, "node_modules", "@openai", "codex", "bin", "codex.exe"),
			filepath.Join(npm, "node_modules", "@openai", "codex", "vendor", "x86_64-pc-windows-msvc", "codex", "codex.exe"),
		)
	}
	candidates = append(candidates,
		filepath.Join(npm, name+".exe"),
		filepath.Join(npm, name+".cmd"),
		filepath.Join(home, ".local", "bin", name+".exe"),
		filepath.Join(home, ".local", "bin", name),
		`C:\Program Files\nodejs\`+name+".cmd",
	)
	for _, c := range candidates {
		if st, err := os.Stat(c); err == nil && !st.IsDir() {
			return preferRealExe(c)
		}
	}
	// Last resort: bare name (PATH at spawn time).
	return name
}

// ResolveCLI returns the absolute path to a brain CLI when found, or the bare name.
func ResolveCLI(name string) string {
	return resolveCLI(name)
}

// CLIAvailable reports whether the named CLI can be executed.
func CLIAvailable(name string) bool {
	p := resolveCLI(name)
	if p == name {
		// Bare name only counts if LookPath finds it.
		if _, err := exec.LookPath(name); err != nil {
			if _, err2 := exec.LookPath(name + ".exe"); err2 != nil {
				return false
			}
		}
		return true
	}
	if st, err := os.Stat(p); err == nil && !st.IsDir() {
		return true
	}
	return false
}

// LaunchClaudeAuthLogin starts `claude auth login` in a new console (browser OAuth).
// Slate never stores credentials — Claude Code owns the session, matching upstream.
func LaunchClaudeAuthLogin() (string, error) {
	if !CLIAvailable("claude") {
		return "", fmt.Errorf(
			"Claude Code CLI was not found on PATH.\n\n" +
				"1. Install Claude Code from https://claude.com/claude-code\n" +
				"2. Open a terminal and confirm: claude --version\n" +
				"3. Restart Slate, then try again.\n\n" +
				"No API keys are stored in this app.",
		)
	}
	path := resolveCLI("claude")
	if err := launchCLIAuth(path, []string{"auth", "login"}, "Claude Code Login"); err != nil {
		return "", fmt.Errorf("could not start Claude login: %w", err)
	}
	return "Opened Claude Code login in a new window.\n\n" +
		"1. Complete sign-in in the browser.\n" +
		"2. Close the login window when finished.\n" +
		"3. Use Brain → Test Claude brain… (or the title-bar Brain pill) to verify.\n\n" +
		"Slate uses your Claude Code session only — no API keys are saved here.", nil
}

// LaunchCodexAuthLogin starts Codex CLI login in a new console (optional second brain).
func LaunchCodexAuthLogin() (string, error) {
	if !CLIAvailable("codex") {
		return "", fmt.Errorf(
			"Codex CLI was not found on PATH.\n\n" +
				"Install the Codex CLI, run: codex login\n" +
				"Then restart Slate. No API keys are stored in this app.",
		)
	}
	path := resolveCLI("codex")
	if err := launchCLIAuth(path, []string{"login"}, "Codex Login"); err != nil {
		return "", fmt.Errorf("could not start Codex login: %w", err)
	}
	return "Opened Codex login in a new window.\n\n" +
		"Finish sign-in, then use Brain → Test… or the title-bar Brain pill.\n" +
		"No API keys are stored in Slate.", nil
}

// launchCLIAuth opens an interactive CLI login in a separate process (Windows-friendly).
func launchCLIAuth(cliPath string, args []string, windowTitle string) error {
	// Prefer PowerShell Start-Process so .cmd shims and spaces in paths work.
	argList := make([]string, 0, len(args))
	for _, a := range args {
		argList = append(argList, "'"+strings.ReplaceAll(a, "'", "''")+"'")
	}
	ps := fmt.Sprintf(
		`Start-Process -FilePath %q -ArgumentList %s`,
		cliPath,
		strings.Join(argList, ","),
	)
	cmd := exec.Command("powershell", "-NoProfile", "-ExecutionPolicy", "Bypass", "-Command", ps)
	if err := cmd.Start(); err != nil {
		// Fallback: new cmd window
		all := append([]string{"/c", "start", windowTitle, cliPath}, args...)
		cmd2 := exec.Command("cmd", all...)
		return cmd2.Start()
	}
	return nil
}

func normalizeEndpoint(url string) string {
	u := strings.TrimSpace(url)
	u = strings.TrimRight(u, "/")
	if !strings.HasPrefix(u, "http://") && !strings.HasPrefix(u, "https://") {
		u = "http://" + u
	}
	if !strings.HasSuffix(u, "/v1") {
		u = u + "/v1"
	}
	return u
}

func probeLocal(endpoint string) []types.LocalModelInfo {
	ctx, cancel := context.WithTimeout(context.Background(), 1500*time.Millisecond)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint+"/models", nil)
	if err != nil {
		return nil
	}
	req.Header.Set("Authorization", "Bearer slate")
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		return nil
	}
	var body struct {
		Data []struct {
			ID string `json:"id"`
		} `json:"data"`
	}
	if err := json.NewDecoder(res.Body).Decode(&body); err != nil {
		return nil
	}
	var models []types.LocalModelInfo
	for _, m := range body.Data {
		if m.ID != "" {
			models = append(models, types.LocalModelInfo{ID: m.ID})
		}
	}
	return models
}

// DetectLocal finds a live local server.
func DetectLocal(preferred string) types.LocalModelsResult {
	var candidates []string
	if preferred != "" {
		candidates = []string{normalizeEndpoint(preferred)}
	} else {
		candidates = localCandidates
	}
	for _, ep := range candidates {
		models := probeLocal(ep)
		if models != nil {
			epCopy := ep
			return types.LocalModelsResult{Endpoint: &epCopy, Models: models}
		}
	}
	return types.LocalModelsResult{Endpoint: nil, Models: []types.LocalModelInfo{}}
}

func whichVersion(cmd string, args ...string) *string {
	c := exec.Command(resolveCLI(cmd), args...)
	configureHiddenProcess(c)
	var out bytes.Buffer
	c.Stdout = &out
	c.Stderr = &out
	if err := c.Run(); err != nil {
		return nil
	}
	line := strings.TrimSpace(strings.Split(out.String(), "\n")[0])
	if line == "" {
		line = "available"
	}
	return &line
}

// Status reports available brains.
func Status(localEndpoint string) types.BrainStatus {
	claudeV := whichVersion("claude", "--version")
	codexV := whichVersion("codex", "--version")
	local := DetectLocal(localEndpoint)
	st := types.BrainStatus{
		Claude: types.BackendStatus{Available: claudeV != nil, Version: claudeV},
		Codex:  types.BackendStatus{Available: codexV != nil, Version: codexV},
	}
	if local.Endpoint != nil {
		v := fmt.Sprintf("%d model(s) @ %s", len(local.Models), strings.TrimPrefix(strings.TrimPrefix(*local.Endpoint, "https://"), "http://"))
		st.Local = types.LocalBackendStatus{Available: true, Version: &v, Endpoint: local.Endpoint}
	} else {
		st.Local = types.LocalBackendStatus{Available: false, Version: nil, Endpoint: nil}
	}
	return st
}

// ExtractJSON finds the first balanced JSON object or array in text.
func ExtractJSON(text string) (interface{}, error) {
	cleaned := strings.TrimSpace(regexp.MustCompile("```(?:json)?").ReplaceAllString(text, ""))
	// Prefer a full parse when the response is already pure JSON.
	if cleaned != "" && (cleaned[0] == '{' || cleaned[0] == '[') {
		var whole interface{}
		if err := json.Unmarshal([]byte(cleaned), &whole); err == nil {
			return whole, nil
		}
	}
	// Scan for the earliest JSON start so arrays win over nested objects.
	objIdx := strings.IndexByte(cleaned, '{')
	arrIdx := strings.IndexByte(cleaned, '[')
	order := [][2]string{{"{", "}"}, {"[", "]"}}
	if arrIdx >= 0 && (objIdx < 0 || arrIdx < objIdx) {
		order = [][2]string{{"[", "]"}, {"{", "}"}}
	}
	for _, pair := range order {
		open, close := pair[0][0], pair[1][0]
		i := strings.IndexByte(cleaned, open)
		if i < 0 {
			continue
		}
		depth := 0
		inStr := false
		esc := false
		for j := i; j < len(cleaned); j++ {
			ch := cleaned[j]
			if esc {
				esc = false
				continue
			}
			if ch == '\\' {
				esc = true
				continue
			}
			if ch == '"' {
				inStr = !inStr
				continue
			}
			if inStr {
				continue
			}
			if ch == open {
				depth++
			} else if ch == close {
				depth--
				if depth == 0 {
					var v interface{}
					if err := json.Unmarshal([]byte(cleaned[i:j+1]), &v); err == nil {
						return v, nil
					}
					break
				}
			}
		}
	}
	return nil, errors.New("No valid JSON found in response")
}

type chatMsg struct {
	Role    string      `json:"role"`
	Content interface{} `json:"content"`
}

func localMessages(req types.BrainRequest) []chatMsg {
	var user interface{} = req.Prompt
	if len(req.Images) > 0 {
		parts := []map[string]interface{}{
			{"type": "text", "text": req.Prompt},
		}
		for _, img := range req.Images {
			ext := strings.ToLower(filepath.Ext(img))
			mime, ok := imageMIME[ext]
			if !ok {
				continue
			}
			raw, err := os.ReadFile(img)
			if err != nil {
				continue
			}
			b64 := base64.StdEncoding.EncodeToString(raw)
			parts = append(parts, map[string]interface{}{
				"type": "image_url",
				"image_url": map[string]string{
					"url": "data:" + mime + ";base64," + b64,
				},
			})
		}
		user = parts
	}
	return []chatMsg{
		{Role: "system", Content: req.System},
		{Role: "user", Content: user},
	}
}

func runLocalOnce(ctx context.Context, req types.BrainRequest, endpoint, model, extraNudge string) (string, error) {
	messages := localMessages(req)
	if extraNudge != "" {
		messages = append(messages, chatMsg{Role: "user", Content: extraNudge})
	}
	body, _ := json.Marshal(map[string]interface{}{
		"model":    model,
		"messages": messages,
		"stream":   false,
	})
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint+"/chat/completions", bytes.NewReader(body))
	if err != nil {
		return "", err
	}
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer slate")
	res, err := http.DefaultClient.Do(httpReq)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()
	raw, _ := io.ReadAll(res.Body)
	if res.StatusCode != 200 {
		detail := string(raw)
		if len(detail) > 300 {
			detail = detail[:300]
		}
		return "", fmt.Errorf("Local model server responded %d at %s. %s", res.StatusCode, endpoint, detail)
	}
	var parsed struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
		Error *struct {
			Message string `json:"message"`
		} `json:"error"`
	}
	if err := json.Unmarshal(raw, &parsed); err != nil {
		return "", err
	}
	if parsed.Error != nil && parsed.Error.Message != "" {
		return "", errors.New(parsed.Error.Message)
	}
	if len(parsed.Choices) == 0 || strings.TrimSpace(parsed.Choices[0].Message.Content) == "" {
		return "", errors.New("Local model returned an empty response.")
	}
	return strings.TrimSpace(parsed.Choices[0].Message.Content), nil
}

func runLocal(req types.BrainRequest, started time.Time) types.BrainResult {
	detected := DetectLocal(req.LocalEndpoint)
	if detected.Endpoint == nil {
		return types.BrainResult{
			ID: req.ID, OK: false, Text: "",
			Error:     "No local model server found. Start Ollama, LM Studio, vLLM, or llama.cpp (or set a custom endpoint in Project Settings → Brain), then retry.",
			ElapsedMs: time.Since(started).Milliseconds(),
		}
	}
	model := strings.TrimSpace(req.LocalModel)
	if model == "" && len(detected.Models) > 0 {
		model = detected.Models[0].ID
	}
	if model == "" {
		return types.BrainResult{
			ID: req.ID, OK: false, Text: "",
			Error:     fmt.Sprintf("Local server at %s has no models loaded. Pull or load a model, then retry.", *detected.Endpoint),
			ElapsedMs: time.Since(started).Milliseconds(),
		}
	}
	ctx, cancel := context.WithCancel(context.Background())
	mu.Lock()
	runningCancel[req.ID] = cancel
	mu.Unlock()
	defer func() {
		mu.Lock()
		delete(runningCancel, req.ID)
		mu.Unlock()
		cancel()
	}()

	text, err := runLocalOnce(ctx, req, *detected.Endpoint, model, "")
	if err != nil {
		return types.BrainResult{ID: req.ID, OK: false, Text: "", Error: err.Error(), ElapsedMs: time.Since(started).Milliseconds()}
	}
	var js interface{}
	if req.ExpectJSON {
		js, err = ExtractJSON(text)
		if err != nil {
			text, err = runLocalOnce(ctx, req, *detected.Endpoint, model, "IMPORTANT: Respond with ONLY the requested JSON. No prose, no code fences.")
			if err != nil {
				return types.BrainResult{ID: req.ID, OK: false, Text: "", Error: err.Error(), ElapsedMs: time.Since(started).Milliseconds()}
			}
			js, err = ExtractJSON(text)
			if err != nil {
				return types.BrainResult{ID: req.ID, OK: false, Text: text, Error: err.Error(), ElapsedMs: time.Since(started).Milliseconds()}
			}
		}
	}
	return types.BrainResult{ID: req.ID, OK: true, Text: text, JSON: js, ElapsedMs: time.Since(started).Milliseconds()}
}

func parseClaudeOutput(raw string) (string, error) {
	var parsed map[string]interface{}
	if err := json.Unmarshal([]byte(raw), &parsed); err != nil {
		return strings.TrimSpace(raw), nil
	}
	if isErr, _ := parsed["is_error"].(bool); isErr {
		msg, _ := parsed["result"].(string)
		if msg == "" {
			msg = "Claude Code returned an error."
		}
		re := regexp.MustCompile(`(?i)authenticat|oauth|401|logged? ?in|revoked`)
		if re.MatchString(msg) {
			return "", fmt.Errorf("Claude Code's sign-in has expired or been revoked. Open a terminal, run: claude auth login  — approve in the browser, then retry. (%s)", msg)
		}
		return "", errors.New(msg)
	}
	if result, ok := parsed["result"].(string); ok {
		return result, nil
	}
	return strings.TrimSpace(raw), nil
}

// Run executes a brain request against the chosen backend.
func Run(req types.BrainRequest, backend string) types.BrainResult {
	started := time.Now()
	if backend == "" {
		backend = req.Backend
	}
	if backend == "" {
		backend = "claude"
	}

	// Demo mock mode
	if dir := os.Getenv("SLATE_BRAIN_MOCK"); dir != "" {
		keys := []string{"first-ad", "reference-analysis", "score-compile", "voice-compile", "compile", "directors-note"}
		key := "default"
		for _, k := range keys {
			if strings.HasPrefix(req.Task, k) {
				key = k
				break
			}
		}
		file := filepath.Join(dir, key+".json")
		if raw, err := os.ReadFile(file); err == nil {
			time.Sleep(1200 * time.Millisecond)
			var canned interface{}
			if err := json.Unmarshal(raw, &canned); err == nil {
				text := string(raw)
				if s, ok := canned.(string); ok {
					text = s
					return types.BrainResult{ID: req.ID, OK: true, Text: text, ElapsedMs: time.Since(started).Milliseconds()}
				}
				return types.BrainResult{ID: req.ID, OK: true, Text: text, JSON: canned, ElapsedMs: time.Since(started).Milliseconds()}
			}
		}
	}

	if backend == "local" {
		return runLocal(req, started)
	}

	tmpDir := os.TempDir()
	lastMessageFile := filepath.Join(tmpDir, fmt.Sprintf("slate-codex-%s.txt", req.ID))

	var cmdPath string
	var args []string
	var stdinData string

	if backend == "claude" {
		cmdPath = resolveCLI("claude")
		args = []string{"-p", "--output-format", "json"}
		if m, ok := claudeTierModel[req.Tier]; ok && m != "" {
			args = append(args, "--model", m)
		}
		if len(req.Images) > 0 {
			args = append(args, "--allowedTools", "Read")
		}
		args = append(args, "--append-system-prompt", req.System)
		prompt := req.Prompt
		if len(req.Images) > 0 {
			prompt += "\n\nReference media frames to view (use the Read tool on each before answering):\n"
			for _, p := range req.Images {
				prompt += "- " + p + "\n"
			}
		}
		stdinData = prompt
	} else {
		// codex
		cmdPath = resolveCLI("codex")
		args = []string{"exec", "--skip-git-repo-check", "--output-last-message", lastMessageFile}
		for _, img := range req.Images {
			args = append(args, "-i", img)
		}
		args = append(args, "-")
		stdinData = req.System + "\n\n---\n\n" + req.Prompt
	}

	runOnce := func(extraNudge string) (string, error) {
		ctx, cancel := context.WithCancel(context.Background())
		cmd := exec.CommandContext(ctx, cmdPath, args...)
		// Hide DOS/console windows for Claude/Codex one-shots (First AD, notes, etc.).
		configureHiddenProcess(cmd)
		var stdout, stderr bytes.Buffer
		cmd.Stdout = &stdout
		cmd.Stderr = &stderr
		in := stdinData
		if extraNudge != "" {
			in += "\n\n" + extraNudge
		}
		cmd.Stdin = strings.NewReader(in)
		mu.Lock()
		running[req.ID] = cmd
		runningCancel[req.ID] = cancel
		mu.Unlock()
		err := cmd.Run()
		mu.Lock()
		delete(running, req.ID)
		delete(runningCancel, req.ID)
		mu.Unlock()
		out := stdout.String()
		if err != nil && strings.TrimSpace(out) == "" {
			msg := strings.TrimSpace(stderr.String())
			if msg == "" {
				msg = err.Error()
			}
			return "", errors.New(msg)
		}
		if backend == "claude" {
			return parseClaudeOutput(out)
		}
		// codex: prefer last-message file
		if msg, err := os.ReadFile(lastMessageFile); err == nil {
			_ = os.Remove(lastMessageFile)
			if t := strings.TrimSpace(string(msg)); t != "" {
				return t, nil
			}
		}
		return strings.TrimSpace(out), nil
	}

	text, err := runOnce("")
	if err != nil {
		return types.BrainResult{ID: req.ID, OK: false, Text: "", Error: err.Error(), ElapsedMs: time.Since(started).Milliseconds()}
	}
	var js interface{}
	if req.ExpectJSON {
		js, err = ExtractJSON(text)
		if err != nil {
			text, err = runOnce("IMPORTANT: Respond with ONLY the requested JSON. No prose, no code fences.")
			if err != nil {
				return types.BrainResult{ID: req.ID, OK: false, Text: "", Error: err.Error(), ElapsedMs: time.Since(started).Milliseconds()}
			}
			js, err = ExtractJSON(text)
			if err != nil {
				return types.BrainResult{ID: req.ID, OK: false, Text: text, Error: err.Error(), ElapsedMs: time.Since(started).Milliseconds()}
			}
		}
	}
	return types.BrainResult{ID: req.ID, OK: true, Text: text, JSON: js, ElapsedMs: time.Since(started).Milliseconds()}
}

// Cancel aborts a running brain job.
func Cancel(id string) {
	mu.Lock()
	defer mu.Unlock()
	if cancel, ok := runningCancel[id]; ok {
		cancel()
		delete(runningCancel, id)
	}
	if cmd, ok := running[id]; ok {
		if cmd.Process != nil {
			_ = cmd.Process.Kill()
		}
		delete(running, id)
	}
}

// Test runs a tiny connectivity check.
func Test(backend string, local types.LocalOpts) types.BrainResult {
	return Run(types.BrainRequest{
		ID:            fmt.Sprintf("test-%d", time.Now().UnixMilli()),
		Task:          "self-test",
		System:        "You are a connectivity check. Reply with exactly one word.",
		Prompt:        "Reply with exactly: READY",
		Tier:          "fast",
		LocalEndpoint: local.Endpoint,
		LocalModel:    local.Model,
	}, backend)
}
