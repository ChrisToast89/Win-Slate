package logx

import (
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/ChrisToast89/Win-Slate/setup/internal/paths"
)

var (
	mu     sync.Mutex
	file   *os.File
	inited bool
)

func Init() error {
	mu.Lock()
	defer mu.Unlock()
	if inited && file != nil {
		return nil
	}
	f, err := os.OpenFile(paths.TempLog(), os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o644)
	if err != nil {
		inited = true
		return err
	}
	file = f
	inited = true
	// Do NOT call Log() here — it needs the same mutex (would deadlock).
	line := time.Now().Format("15:04:05") + " ——— log open " + time.Now().Format(time.RFC3339) + " ———\n"
	_, _ = file.WriteString(line)
	return nil
}

func Path() string { return paths.TempLog() }

func Log(format string, args ...interface{}) {
	mu.Lock()
	defer mu.Unlock()
	line := fmt.Sprintf(time.Now().Format("15:04:05")+" "+format+"\n", args...)
	if file != nil {
		_, _ = file.WriteString(line)
		_ = file.Sync() // best-effort; ignore errors
	}
}
