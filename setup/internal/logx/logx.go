package logx

import (
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/ChrisToast89/Win-Slate/setup/internal/paths"
)

var (
	mu   sync.Mutex
	file *os.File
)

func Init() error {
	mu.Lock()
	defer mu.Unlock()
	if file != nil {
		return nil
	}
	f, err := os.OpenFile(paths.TempLog(), os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o644)
	if err != nil {
		return err
	}
	file = f
	Log("——— log open %s ———", time.Now().Format(time.RFC3339))
	return nil
}

func Path() string { return paths.TempLog() }

func Log(format string, args ...interface{}) {
	mu.Lock()
	defer mu.Unlock()
	line := fmt.Sprintf(time.Now().Format("15:04:05")+" "+format+"\n", args...)
	if file != nil {
		_, _ = file.WriteString(line)
	}
}
