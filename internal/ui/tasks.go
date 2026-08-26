package ui

import (
	"fmt"
	"sort"
	"strings"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/alplix/lavandegrid/internal/app"
	"github.com/alplix/lavandegrid/internal/i18n"
)

func init() {
	registerScreen("nav.tasks", theme.ListIcon(), buildTasks, refreshTasks)
}

var tv struct {
	rows    []taskRow
	sel     map[string]bool
	search  *widget.Entry
	state   *widget.Select
	hostSel *widget.Select
	list    *widget.List
	count   *widget.Label
	bulkP   *widget.Button
	bulkR   *widget.Button
	bulkA   *widget.Button
}

type taskRow struct {
	hostID   string
	hostName string
	t        app.Task
}

const stateAll = "All"

func buildTasks(w fyne.Window) fyne.CanvasObject {
	tv.sel = map[string]bool{}
	tv.search = widget.NewEntry()
	tv.search.PlaceHolder = i18n.T("tasks.searchPh")
	tv.search.OnChanged = func(string) { refreshTasks() }

	states := []string{stateAll, "Running", "Paused", "Queued", "Ready", "Errors"}
	tv.state = widget.NewSelect(states, func(string) { refreshTasks() })
	tv.state.SetSelected(stateAll)

	hostNames := []string{stateAll}
	for _, h := range mgr.Store.List() {
		hostNames = append(hostNames, h.Name)
	}
	tv.hostSel = widget.NewSelect(hostNames, func(string) { refreshTasks() })
	tv.hostSel.SetSelected(stateAll)

	tv.count = widget.NewLabel("")

	mkBulk := func(label string, imp widget.Importance, op string) *widget.Button {
		b := widget.NewButtonWithIcon(label, iconFor(op), nil)
		b.Importance = imp
		b.Disable()
		b.OnTapped = func() { runBulk(op) }
		return b
	}
	tv.bulkR = mkBulk("Resume", widget.MediumImportance, "resume")
	tv.bulkP = mkBulk("Pause", widget.MediumImportance, "suspend")
	tv.bulkA = mkBulk("Abort", widget.DangerImportance, "abort")

	toolbar := container.NewVBox(
		container.NewHBox(
			container.NewBorder(nil, nil, widget.NewLabel("State"), nil, tv.state),
			container.NewBorder(nil, nil, widget.NewLabel("Server"), nil, tv.hostSel),
		),
		container.NewHBox(tv.bulkR, tv.bulkP, tv.bulkA,
			container.NewCenter(tv.count)),
	)

	tv.list = widget.NewList(
		func() int { return len(tv.rows) },
		func() fyne.CanvasObject { return newTaskItem() },
		func(i widget.ListItemID, o fyne.CanvasObject) {
			if i < len(tv.rows) {
				o.(*taskItem).set(tv.rows[i])
			}
		},
	)
	tv.list.OnSelected = func(i widget.ListItemID) {
		if i < len(tv.rows) {
			showTaskDetail(w, tv.rows[i])
		}
	}

	return container.NewBorder(toolbar, nil, nil, nil, container.NewVScroll(tv.list))
}

func iconFor(op string) fyne.Resource {
	switch op {
	case "resume":
		return theme.MediaPlayIcon()
	case "suspend":
		return theme.MediaPauseIcon()
	default:
		return theme.ContentClearIcon()
	}
}

func keyOf(r taskRow) string { return r.hostID + "\x00" + r.t.Name }

func runBulk(op string) {
	for k := range tv.sel {
		parts := strings.SplitN(k, "\x00", 2)
		if len(parts) == 2 {
			go func(h, n string) { _ = mgr.TaskOp(h, n, op) }(parts[0], parts[1])
		}
	}
	tv.sel = map[string]bool{}
	updateTaskUI()
	go func() {
		time.Sleep(400 * time.Millisecond)
		fireRefresh()
	}()
}

func updateTaskUI() {
	sel := len(tv.sel)
	tv.count.SetText(fmt.Sprintf("%s · %s", fmt.Sprintf(i18n.T("tasks.count"), len(tv.rows)), fmt.Sprintf(i18n.T("tasks.selected"), sel)))
	setBulkEnabled(sel > 0)
	if tv.list != nil {
		tv.list.Refresh()
	}
}

func setBulkEnabled(on bool) {
	if on {
		tv.bulkR.Enable()
		tv.bulkP.Enable()
		tv.bulkA.Enable()
	} else {
		tv.bulkR.Disable()
		tv.bulkP.Disable()
		tv.bulkA.Disable()
	}
}

func refreshTasks() {
	if mgr == nil || tv.list == nil || tv.state.Selected == "" {
		return
	}
	q := strings.ToLower(tv.search.Text)
	state := tv.state.Selected
	hsel := tv.hostSel.Selected

	var rows []taskRow
	snaps := mgr.AllSnaps()
	for _, cfg := range mgr.Store.List() {
		if hsel != stateAll && cfg.Name != hsel {
			continue
		}
		s := snaps[cfg.ID]
		if s == nil || !s.Online {
			continue
		}
		for _, t := range s.Tasks {
			switch state {
			case "Running":
				if t.Status != app.StatusRunning {
					continue
				}
			case "Paused":
				if t.Status != app.StatusPaused {
					continue
				}
			case "Queued":
				if t.Status != app.StatusQueued && t.Status != app.StatusDownloading && t.Status != app.StatusUploading {
					continue
				}
			case "Ready":
				if !t.Ready {
					continue
				}
			case "Errors":
				if t.Status != app.StatusError {
					continue
				}
			}
			name := strings.ToLower(t.Name + " " + t.ProjectName + " " + t.Resources)
			if q != "" && !strings.Contains(name, q) {
				continue
			}
			rows = append(rows, taskRow{hostID: cfg.ID, hostName: cfg.Name, t: t})
		}
	}
	sort.SliceStable(rows, func(i, j int) bool { return rows[i].t.Name < rows[j].t.Name })
	tv.rows = rows
	updateTaskUI()
}

type taskItem struct {
	fyne.Container
	check *widget.Check
	name  *widget.Label
	sub   *widget.Label
	bar   *widget.ProgressBar
	state *widget.Label
	cur   taskRow
	has   bool
}

func newTaskItem() *taskItem {
	ti := &taskItem{}
	ti.check = widget.NewCheck("", nil)
	ti.check.OnChanged = func(on bool) {
		if !ti.has {
			return
		}
		k := keyOf(ti.cur)
		if on {
			tv.sel[k] = true
		} else {
			delete(tv.sel, k)
		}
		tv.count.SetText(fmt.Sprintf("%d tasks · %d selected", len(tv.rows), len(tv.sel)))
		setBulkEnabled(len(tv.sel) > 0)
	}
	ti.name = widget.NewLabel("")
	ti.name.TextStyle = fyne.TextStyle{Bold: true}
	ti.name.Wrapping = fyne.TextWrapWord
	ti.sub = widget.NewLabel("")
	ti.bar = widget.NewProgressBar()
	ti.bar.Resize(fyne.NewSize(160, 18))
	ti.state = widget.NewLabel("")

	left := container.NewCenter(ti.check)
	mid := container.NewVBox(ti.name, ti.sub, ti.bar)
	right := container.NewVBox(container.NewCenter(ti.state))
	body := container.NewBorder(nil, nil, left, right, mid)
	ti.Objects = []fyne.CanvasObject{container.NewPadded(body)}
	return ti
}

func (ti *taskItem) set(r taskRow) {
	ti.cur = r
	ti.has = true
	t := r.t
	ti.name.SetText(t.Name)
	ti.sub.SetText(fmt.Sprintf("%s · %s", t.ProjectName, r.hostName))
	ti.bar.Value = clamp01(t.Progress)
	ti.bar.Refresh()
	ti.state.SetText(statusLabel(t.Status, t.Ready))
	ti.check.SetChecked(tv.sel[keyOf(r)])
}

func clamp01(f float64) float64 {
	if f < 0 {
		return 0
	}
	if f > 1 {
		return 1
	}
	return f
}

func statusLabel(st app.TaskStatus, ready bool) string {
	if ready {
		return "[READY]"
	}
	return "[" + strings.ToUpper(string(st)) + "]"
}

func showTaskDetail(w fyne.Window, r taskRow) {
	t := r.t
	form := widget.NewForm()

	kvPairs := [][2]string{
		{"Task", t.Name},
		{"Project", t.ProjectName},
		{"Application version", t.AppVersion},
		{"Host", r.hostName},
		{"Status", statusLabel(t.Status, t.Ready)},
		{"Progress", fmt.Sprintf("%.1f%%", t.Progress*100)},
		{"Elapsed", app.FmtDuration(t.Elapsed)},
		{"CPU time", app.FmtDuration(t.CPUTime)},
		{"Remaining (est.)", app.FmtDuration(t.ETA)},
		{"Deadline", app.FmtTime(time.Unix(t.Deadline, 0))},
		{"Memory", app.FmtBytes(t.Mem)},
		{"Resources", t.Resources},
		{"Slot", fmt.Sprintf("%d", t.Slot)},
	}
	for _, p := range kvPairs {
		e := widget.NewLabel(p[1])
		e.Wrapping = fyne.TextWrapWord
		form.Append(p[0], e)
	}

	doOp := func(op string) {
		go func() { _ = mgr.TaskOp(r.hostID, t.Name, op) }()
		time.Sleep(300 * time.Millisecond)
		fireRefresh()
	}

	btnPause := widget.NewButtonWithIcon("Pause", theme.MediaPauseIcon(), func() { doOp("suspend") })
	btnResume := widget.NewButtonWithIcon("Resume", theme.MediaPlayIcon(), func() { doOp("resume") })
	btnAbort := widget.NewButtonWithIcon("Abort", theme.DeleteIcon(), func() {
		dialog.NewConfirm("Abort task", "Abort "+t.Name+"?", func(ok bool) {
			if ok {
				doOp("abort")
			}
		}, w).Show()
	})
	rowBtns := container.NewHBox(btnPause, btnResume, btnAbort)

	card := widget.NewCard(t.ProjectName, "", container.NewVBox(form, rowBtns))
	dlg := dialog.NewCustom("Task details", "Close", card, w)
	dlg.Resize(fyne.NewSize(560, 600))
	dlg.Show()
	tv.list.UnselectAll()
}
