//go:build windows

package local

import (
	"os/exec"
	"syscall"
)

const createNewProcessGroup = 0x00000200
const createNoWindow = 0x08000000

func StartDetached(exe, dataDir string) (*exec.Cmd, error) {
	cmd := exec.Command(exe, "--dir", dataDir)
	cmd.SysProcAttr = &syscall.SysProcAttr{
		HideWindow:    true,
		CreationFlags: createNewProcessGroup | createNoWindow,
	}
	return cmd, cmd.Start()
}
