package app

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestNormalizeStatuses(t *testing.T) {
	cases := []struct {
		st   TaskStatus
		want TaskStatus
	}{
		{StatusRunning, StatusRunning},
		{StatusPaused, StatusPaused},
		{StatusError, StatusError},
		{StatusQueued, StatusQueued},
		{StatusReady, StatusReady},
	}
	for _, c := range cases {
		m := NewMock(HostCfg{ID: "t", Name: "Test"})
		if len(m.tasks) == 0 {
			t.Fatal("mock created no tasks")
		}
		mt := m.tasks[0]
		mt.t.Status = c.st
		mt.t.Active = c.st == StatusRunning || c.st == StatusPaused
		mt.t.Suspended = c.st == StatusPaused
		mt.t.Ready = c.st == StatusReady
		snap := m.Snapshot()
		if len(snap.Tasks) == 0 {
			t.Fatalf("expected mock tasks, got none")
		}
		if snap.Tasks[0].Status != c.want {
			t.Fatalf("status %q -> got %q, want %q", c.st, snap.Tasks[0].Status, c.want)
		}
	}
}

func TestMockSnapshotHasData(t *testing.T) {
	m := NewMock(HostCfg{ID: "demo", Name: "Demo"})
	snap := m.Snapshot()
	if !snap.Online {
		t.Fatal("demo mock should report online")
	}
	if len(snap.Projects) == 0 {
		t.Error("demo mock should expose projects")
	}
	if len(snap.Tasks) == 0 {
		t.Error("demo mock should expose tasks")
	}
	if len(snap.Transfers) == 0 {
		t.Error("demo mock should expose at least one transfer")
	}
	if len(snap.Messages) == 0 {
		t.Error("demo mock should expose messages")
	}
	if snap.Totals.Running == 0 {
		t.Error("demo mock should have running tasks")
	}
	if snap.Totals.Credit <= 0 {
		t.Error("demo mock should aggregate credit")
	}
	if len(snap.HostInfo.GPUs) == 0 {
		t.Error("demo mock should report a GPU")
	}
	if snap.Version == "" {
		t.Error("demo mock should report a version")
	}
}

func TestMockOps(t *testing.T) {
	m := NewMock(HostCfg{ID: "demo", Name: "Demo"})
	snap := m.Snapshot()
	if len(snap.Tasks) == 0 {
		t.Fatal("no tasks to operate on")
	}
	name := snap.Tasks[0].Name
	if err := m.ResultOp(name, "suspend"); err != nil {
		t.Fatalf("suspend: %v", err)
	}
	if err := m.ResultOp(name, "resume"); err != nil {
		t.Fatalf("resume: %v", err)
	}
	if len(snap.Projects) == 0 {
		t.Fatal("no projects to operate on")
	}
	purl := snap.Projects[0].URL
	if err := m.ProjectOp(purl, "suspend"); err != nil {
		t.Fatalf("project suspend: %v", err)
	}
	if err := m.ProjectOp(purl, "resume"); err != nil {
		t.Fatalf("project resume: %v", err)
	}
	if err := m.ProjectOp(purl, "detach"); err != nil {
		t.Fatalf("project detach: %v", err)
	}
	if st, err := m.Stats(); err != nil || len(st) == 0 {
		t.Fatalf("stats: %v, n=%d", err, len(st))
	}
}

func TestMsgTimeIsUnix(t *testing.T) {
	m := NewMock(HostCfg{ID: "demo", Name: "Demo"})
	snap := m.Snapshot()
	for _, msg := range snap.Messages {
		if msg.Time <= 0 {
			t.Fatalf("message time should be a unix timestamp, got %v", msg.Time)
		}
	}
}

func TestStoreRoundTrip(t *testing.T) {
	dir := t.TempDir()
	s := &Store{path: filepath.Join(dir, "hosts.json")}
	h := HostCfg{Name: "A", Host: "localhost", Port: 31416, Password: "pw"}
	h = s.Upsert(h)
	if err := s.Save(); err != nil {
		t.Fatalf("save: %v", err)
	}
	s2 := &Store{path: filepath.Join(dir, "hosts.json")}
	if data, err := os.ReadFile(s2.path); err == nil {
		_ = json.Unmarshal(data, &s2.hosts)
	}
	list := s2.List()
	if len(list) != 1 || list[0].ID != h.ID || list[0].Password != "pw" {
		t.Fatalf("round trip mismatch: %+v", list)
	}
}

func TestFmtHelpers(t *testing.T) {
	if FmtBytes(1024) != "1.0 KB" {
		t.Errorf("FmtBytes: %s", FmtBytes(1024))
	}
	if FmtDuration(90) != "1m 30s" {
		t.Errorf("FmtDuration: %s", FmtDuration(90))
	}
	if FmtNum(5e6) != "5.00 M" {
		t.Errorf("FmtNum: %s", FmtNum(5e6))
	}
}