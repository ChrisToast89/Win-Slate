//go:build !windows

package brain

import "os/exec"

func configureHiddenProcess(cmd *exec.Cmd) {
	// No-op on Unix; console flash is a Windows GUI-host issue.
}
