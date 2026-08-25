//go:build !windows

package local

import (
	"os/exec"
	"syscall"
)

func StartDetached(exe, dataDir string) (*exec.Cmd, error) {
	cmd := exec.Command(exe, "--dir", dataDir)
	cmd.SysProcAttr = &syscall.SysProcAttr{Setsid: true}
	return cmd, cmd.Start()
}
