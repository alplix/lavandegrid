package app

import (
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/alplix/lavandegrid/internal/boinc"
)

type TaskStatus string

const (
	StatusRunning    TaskStatus = "running"
	StatusPaused     TaskStatus = "paused"
	StatusQueued     TaskStatus = "queued"
	StatusDownloading TaskStatus = "downloading"
	StatusUploading  TaskStatus = "uploading"
	StatusError      TaskStatus = "error"
	StatusReady      TaskStatus = "ready"
)

type GPU struct {
	Vendor string
	Count  int
	Names  []string
	Driver string
	Cuda   string
	VRAM   int64
}

type HostSpec struct {
	OS        string
	OSVersion string
	CPU       string
	Cores     int
	Flops     float64
	Memory    int64
	DiskFree  int64
	DiskTotal int64
	CPID      string
	GPUs      []GPU
}

type ProjectInfo struct {
	Name       string
	URL        string
	Venue      string
	UserName   string
	TeamName   string
	UserCredit float64
	RAC        float64
	HostCredit float64
	HostRAC    float64
	Share      float64
	Suspended  bool
	NoMoreWork bool
	Pending    bool
	Ended      bool
}

type Task struct {
	Name        string
	Wu          string
	URL         string
	ProjectName string
	Status      TaskStatus
	Progress    float64
	Elapsed     float64
	CPUTime     float64
	ETA         float64
	Deadline    int64
	CPTime      float64
	Exit        int
	Mem         int64
	Resources   string
	Slot        int
	AppVersion  string
	Active      bool
	Suspended   bool
	Ready       bool
}

type Transfer struct {
	Name        string
	URL         string
	ProjectName string
	Upload      bool
	Total       int64
	Done        int64
	Progress    float64
	Paused      bool
	Finished    bool
}

type MsgLine struct {
	Seq     int
	Pri     int
	Time    time.Time
	Body    string
	Project string
}

type Totals struct {
	Running   int
	Paused    int
	Queued    int
	Errors    int
	Downloads int
	Uploads   int
	Memory    int64
	Credit    float64
	RAC       float64
}

type Snapshot struct {
	HostID    string
	Demo      bool
	Online    bool
	Error     string
	Version   string
	TS        time.Time
	HostInfo  HostSpec
	Projects  []ProjectInfo
	Tasks     []Task
	Transfers []Transfer
	Messages  []MsgLine
	Totals    Totals
	TaskMode  string
	NetMode   string
}

func statusOf(r boinc.Result) TaskStatus {
	exit := r.ExitStatus.I()
	state := r.State.I()
	if r.ActiveTask.B() && !r.SuspendedViaGUI.B() {
		return StatusRunning
	}
	if r.ActiveTask.B() && r.SuspendedViaGUI.B() {
		return StatusPaused
	}
	if r.SuspendedViaGUI.B() {
		return StatusPaused
	}
	if exit != 0 || state == 6 || state == 3 {
		return StatusError
	}
	if r.ReadyToReport.B() || state == 7 {
		return StatusReady
	}
	if state == 1 || state == 2 {
		return StatusDownloading
	}
	if state == 4 || state == 5 {
		return StatusUploading
	}
	return StatusQueued
}

func progressOf(r boinc.Result) float64 {
	f := r.FractionDone.F()
	if f <= 0 {
		cpu := r.CurrentCPUTime.F()
		eta := r.EstimatedCPUTimeRemaining.F()
		if cpu > 0 && eta > 0 {
			f = cpu / (cpu + eta)
		}
	}
	if f < 0 {
		f = 0
	}
	if f > 1 {
		f = 1
	}
	return f
}

func cudaStr(v float64) string {
	if v <= 0 {
		return ""
	}
	if v >= 1000 {
		maj := int(v) / 1000
		min := float64(int(v)%1000) / 10
		s := strconvFormat(min)
		return itoa(maj) + "." + s
	}
	return itoa(int(v))
}

func matchVram(props []boinc.OpenCLProp, names []string) int64 {
	var total int64
	for _, n := range names {
		ln := strings.ToLower(n)
		for _, p := range props {
			lk := strings.ToLower(p.Name)
			if lk != "" && (strings.Contains(lk, ln) || strings.Contains(ln, lk)) {
				total += p.GlobalMem.I64()
				break
			}
		}
	}
	return total
}

func gpusOf(hi boinc.HostInfo, ocl []boinc.OpenCLProp) []GPU {
	if hi.Coprocs.Count.I() <= 0 {
		return nil
	}
	var out []GPU
	nv := nonEmpty(hi.Coprocs.NvidiaDeviceNames)
	if nv > 0 && len(hi.Coprocs.NvidiaDeviceNames) > 0 {
		out = append(out, GPU{Vendor: "NVIDIA", Count: nv, Names: hi.Coprocs.NvidiaDeviceNames,
			Driver: hi.Coprocs.NvidiaDriverVersion, Cuda: cudaStr(hi.Coprocs.CudaVersion.F()), VRAM: matchVram(ocl, hi.Coprocs.NvidiaDeviceNames)})
	}
	at := nonEmpty(hi.Coprocs.AtiDeviceNames)
	drv := hi.Coprocs.AmdDriverVersion
	if drv == "" {
		drv = hi.Coprocs.AtiDriverVersion
	}
	if at > 0 && len(hi.Coprocs.AtiDeviceNames) > 0 {
		out = append(out, GPU{Vendor: "AMD", Count: at, Names: hi.Coprocs.AtiDeviceNames, Driver: drv, VRAM: matchVram(ocl, hi.Coprocs.AtiDeviceNames)})
	}
	it := nonEmpty(hi.Coprocs.IntelGpuDeviceNames)
	if it > 0 && len(hi.Coprocs.IntelGpuDeviceNames) > 0 {
		out = append(out, GPU{Vendor: "Intel", Count: it, Names: hi.Coprocs.IntelGpuDeviceNames, VRAM: matchVram(ocl, hi.Coprocs.IntelGpuDeviceNames)})
	}
	return out
}

func nonEmpty(names []string) int {
	c := 0
	for _, n := range names {
		if strings.TrimSpace(n) != "" {
			c++
		}
	}
	if c == 0 {
		return len(names)
	}
	return c
}

func Normalize(hostID string, demo bool, st *boinc.ClientState, transfers []boinc.FileTransfer, cc *boinc.CcStatus, msgs []boinc.Msg, version string) *Snapshot {
	nameByURL := map[string]string{}
	for _, p := range st.Projects {
		nameByURL[p.MasterURL] = p.Name
	}
	snap := &Snapshot{
		HostID:  hostID,
		Demo:    demo,
		Online:  true,
		Version: version,
		TS:      time.Now(),
		Projects: make([]ProjectInfo, 0, len(st.Projects)),
		Tasks:   make([]Task, 0, len(st.Results)),
		Transfers: make([]Transfer, 0, len(transfers)),
	}
	hi := &st.HostInfo
	snap.HostInfo = HostSpec{
		OS:        hi.OSName,
		OSVersion: hi.OSVersion,
		CPU:       strings.TrimSpace(hi.PVendor + " " + hi.PModel),
		Cores:     hi.PNcpus.I(),
		Flops:     hi.PFlops.F(),
		Memory:    hi.MNbytes.I64(),
		DiskFree:  hi.DFree.I64(),
		DiskTotal: hi.DTotal.I64(),
		CPID:      hi.HostCPID,
		GPUs:      gpusOf(*hi, st.OpenCLGpuProps),
	}
	for _, p := range st.Projects {
		snap.Projects = append(snap.Projects, ProjectInfo{
			Name: p.Name, URL: p.MasterURL, Venue: p.Venue,
			UserName: p.UserName, TeamName: p.TeamName,
			UserCredit: p.UserTotalCredit.F(), RAC: p.UserExpavgCredit.F(),
			HostCredit: p.HostTotalCredit.F(), HostRAC: p.HostExpavgCredit.F(),
			Share: p.ResourceShare.F(), Suspended: p.SuspendedViaGUI.B(),
			NoMoreWork: p.DontRequestMoreWork.B(), Pending: p.SchedRPCPending.B(), Ended: p.Ended.B(),
		})
	}
	for _, r := range st.Results {
		snap.Tasks = append(snap.Tasks, Task{
			Name: r.Name, Wu: r.WuName, URL: r.ProjectURL,
			ProjectName: nameByURL[r.ProjectURL],
			Status:      statusOf(r), Progress: progressOf(r),
			Elapsed: r.ElapsedTime.F(), CPUTime: r.CurrentCPUTime.F(),
			ETA: r.EstimatedCPUTimeRemaining.F(), Deadline: r.ReportDeadline.I64(),
			CPTime: r.CheckpointCPUTime.F(), Exit: r.ExitStatus.I(),
			Mem: r.WorkingSetSize.I64(), Resources: strings.TrimSpace(r.Resources),
			Slot: slotOf(r.Slot.I()), AppVersion: appVer(r.VersionNum.I()),
			Active: r.ActiveTask.B(), Suspended: r.SuspendedViaGUI.B(), Ready: r.ReadyToReport.B(),
		})
	}
	for _, t := range transfers {
		up := t.IsUpload.B()
		total := t.Nbytes.I64()
		done := t.BytesXferred.I64()
		pr := 0.0
		if total > 0 {
			pr = float64(done) / float64(total)
		}
		fin := t.IsValid.B() && total > 0 && done >= total
		snap.Transfers = append(snap.Transfers, Transfer{
			Name: t.Name, URL: t.ProjectURL,
			ProjectName: nameByURL[t.ProjectURL], Upload: up,
			Total: total, Done: done, Progress: pr,
			Paused: t.Paused.B(), Finished: fin,
		})
	}
	sort.Slice(msgs, func(i, j int) bool { return msgs[i].Seqno.F() < msgs[j].Seqno.F() })
	for i, m := range msgs {
		if i < len(msgs)-80 {
			continue
		}
		snap.Messages = append(snap.Messages, MsgLine{
			Seq: m.Seqno.I(), Pri: m.Pri.I(), Time: time.Unix(m.Time.I64(), 0),
			Body: m.Body, Project: m.Project,
		})
	}
	t := &snap.Totals
	for _, x := range snap.Tasks {
		switch x.Status {
		case StatusRunning:
			t.Running++
		case StatusPaused:
			t.Paused++
		case StatusError:
			t.Errors++
		case StatusQueued, StatusDownloading, StatusUploading:
			t.Queued++
		}
		if x.Status == StatusRunning || x.Status == StatusPaused {
			t.Memory += x.Mem
		}
	}
	for _, tr := range snap.Transfers {
		switch {
		case !tr.Upload && !tr.Finished:
			t.Downloads++
		case tr.Upload && !tr.Finished:
			t.Uploads++
		}
	}
	for _, p := range snap.Projects {
		t.Credit += p.HostCredit
		t.RAC += p.HostRAC
	}
	mode := map[float64]string{1: "always", 2: "auto", 3: "never"}
	if cc != nil {
		snap.TaskMode = mode[cc.TaskMode.F()]
		snap.NetMode = mode[cc.NetworkMode.F()]
	}
	return snap
}

func slotOf(n int) int {
	if n >= 0 {
		return n
	}
	return -1
}

func appVer(v int) string {
	if v <= 0 {
		return ""
	}
	return itoa(v)
}

func itoa(i int) string {
	if i == 0 {
		return "0"
	}
	neg := i < 0
	if neg {
		i = -i
	}
	var b [20]byte
	pos := len(b)
	for i > 0 {
		pos--
		b[pos] = byte('0' + i%10)
		i /= 10
	}
	if neg {
		pos--
		b[pos] = '-'
	}
	return string(b[pos:])
}

func strconvFormat(f float64) string {
	s := fmt.Sprintf("%.1f", f)
	s = strings.TrimSuffix(s, ".0")
	return s
}
