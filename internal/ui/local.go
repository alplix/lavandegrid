package ui

import (
	"github.com/alplix/lavandegrid/internal/local"
)

func getDaemon() *local.Daemon {
	return local.NewDaemon(local.Detect())
}

func startLocalDaemon() error {
	d := getDaemon()
	if !d.Info.Found {
		return errNoClient
	}
	return local.StartDaemon(d, Version)
}

func stopLocalDaemon() error {
	d := getDaemon()
	return local.StopDaemon(d)
}

func daemonStatus() local.DaemonStatus {
	return getDaemon().Status()
}

type strErr string

func (s strErr) Error() string { return string(s) }

const errNoClient = strErr("no BOINC client found on this machine")
