package ui

import (
	"github.com/alplix/lavandegrid/internal/local"
)

func startLocalClient() error {
	info := local.Detect()
	if !info.Found {
		return errNoClient
	}
	_, err := local.StartDetached(info.Exe, info.DataDir)
	return err
}

type strErr string

func (s strErr) Error() string { return string(s) }

const errNoClient = strErr("no BOINC client found on this machine")
