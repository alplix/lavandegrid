//go:build !windows

package local

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"syscall"
	"time"
)

func statusByPID(pid int) DaemonStatus {
	err := syscall.Kill(pid, 0)
	if err != nil {
		return DaemonStopped
	}
	return DaemonRunning
}

func StartDetached(exe, dataDir string) (*exec.Cmd, error) {
	cmd := exec.Command(exe, "--dir", dataDir)
	cmd.SysProcAttr = &syscall.SysProcAttr{Setsid: true}
	if err := cmd.Start(); err != nil {
		return nil, fmt.Errorf("start failed: %w", err)
	}
	return cmd, nil
}

func StartDaemon(d *Daemon, version string) error {
	if !d.Info.Found {
		return fmt.Errorf("BOINC client not found")
	}
	if d.Status() == DaemonRunning {
		return nil
	}
	if err := d.WriteConfig(version); err != nil {
		return fmt.Errorf("write config: %w", err)
	}
	cmd, err := StartDetached(d.Info.Exe, d.Info.DataDir)
	if err != nil {
		return err
	}
	if cmd.Process != nil {
		writePIDFile(d.PIDFile, cmd.Process.Pid)
	}
	time.Sleep(800 * time.Millisecond)
	return nil
}

func StopDaemon(d *Daemon) error {
	pid := d.PID()
	if pid <= 0 {
		return nil
	}
	err := syscall.Kill(pid, syscall.SIGTERM)
	removePIDFile(d.PIDFile)
	if err != nil {
		return fmt.Errorf("stop failed: %w", err)
	}
	return nil
}

func writePIDFile(path string, pid int) {
	if path == "" {
		return
	}
	dir := path[:len(path)-len("/boinc.pid")]
	os.MkdirAll(dir, 0o755)
	os.WriteFile(path, []byte(strconv.Itoa(pid)), 0o644)
}

func removePIDFile(path string) {
	if path != "" {
		os.Remove(path)
	}
}
