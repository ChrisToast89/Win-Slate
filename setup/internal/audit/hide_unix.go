//go:build !windows

package audit

import "os/exec"

func hideConsole(cmd *exec.Cmd) {}
