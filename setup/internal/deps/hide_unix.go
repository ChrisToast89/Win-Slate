//go:build !windows

package deps

import "os/exec"

func hideConsole(cmd *exec.Cmd) {}
