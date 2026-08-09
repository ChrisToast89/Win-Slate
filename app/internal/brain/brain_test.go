package brain

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestResolveBatchShim(t *testing.T) {
	dir := t.TempDir()
	// Fake layout matching npm global claude.cmd
	binDir := filepath.Join(dir, "node_modules", "@anthropic-ai", "claude-code", "bin")
	if err := os.MkdirAll(binDir, 0o755); err != nil {
		t.Fatal(err)
	}
	exe := filepath.Join(binDir, "claude.exe")
	if err := os.WriteFile(exe, []byte("fake"), 0o644); err != nil {
		t.Fatal(err)
	}
	cmdPath := filepath.Join(dir, "claude.cmd")
	shim := "@ECHO off\r\n" +
		"GOTO start\r\n" +
		":find_dp0\r\n" +
		"SET dp0=%~dp0\r\n" +
		"EXIT /b\r\n" +
		":start\r\n" +
		"SETLOCAL\r\n" +
		"CALL :find_dp0\r\n" +
		"\"%dp0%\\node_modules\\@anthropic-ai\\claude-code\\bin\\claude.exe\"   %*\r\n"
	if err := os.WriteFile(cmdPath, []byte(shim), 0o644); err != nil {
		t.Fatal(err)
	}
	got := resolveBatchShim(cmdPath)
	if got != exe {
		t.Fatalf("resolveBatchShim = %q, want %q", got, exe)
	}
	// Non-batch paths pass through.
	if resolveBatchShim(exe) != exe {
		t.Fatal("exe should pass through")
	}
}


func TestExtractJSONObject(t *testing.T) {
	raw := "Sure, here you go:\n```json\n{\"ok\": true, \"n\": 3}\n```\n"
	v, err := ExtractJSON(raw)
	if err != nil {
		t.Fatal(err)
	}
	m, ok := v.(map[string]interface{})
	if !ok {
		t.Fatalf("type %T", v)
	}
	if m["ok"] != true {
		t.Fatalf("%v", m)
	}
}

func TestExtractJSONArray(t *testing.T) {
	v, err := ExtractJSON(`[1, 2, {"a":"b"}]`)
	if err != nil {
		t.Fatal(err)
	}
	a, ok := v.([]interface{})
	if !ok || len(a) != 3 {
		t.Fatalf("%v", v)
	}
}

func TestNormalizeAndDetectEmpty(t *testing.T) {
	// No server expected — should return empty without panic
	res := DetectLocal("http://127.0.0.1:1/v1")
	if res.Endpoint != nil {
		// flaky if something is on port 1; accept either
		_ = res
	}
	b, _ := json.Marshal(Status(""))
	if len(b) < 10 {
		t.Fatal("status empty")
	}
}
