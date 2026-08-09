//go:build windows

package brain

import (
	"os/exec"
	"syscall"
)

// CREATE_NO_WINDOW — hide console windows when a GUI host (Wails) spawns CLIs.
// Without this, every `claude -p` / version probe flashes a DOS box.
const createNoWindow = 0x08000000

// configureHiddenProcess prevents console window flash for non-interactive brain jobs.
// Auth login still uses launchCLIAuth, which intentionally shows a window.
func configureHiddenProcess(cmd *exec.Cmd) {
	if cmd == nil {
		return
	}
	cmd.SysProcAttr = &syscall.SysProcAttr{
		HideWindow:    true,
		CreationFlags: createNoWindow,
	}
}
