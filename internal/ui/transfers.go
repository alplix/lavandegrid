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
)

func init() {
	registerScreen("nav.transfers", theme.DownloadIcon(), buildTransfers, refreshTransfers)
}

var xf struct {
	hostSel *widget.Select
	list    *widget.List
	rows    []xfRow
	cur     string
}

type xfRow struct {
	hostID   string
	hostName string
	t        struct {
		Name, Project string
		Upload        bool
		Done          int64
		Total         int64
		Progress      float64
		Paused        bool
		Finished      bool
	}
}

func buildTransfers(w fyne.Window) fyne.CanvasObject {
	xf.hostSel = widget.NewSelect(hostChoices(), func(s string) {
		xf.cur = s
		refreshTransfers()
	})
	xf.hostSel.PlaceHolder = "All servers"

	top := container.NewHBox(
		container.NewBorder(nil, nil, widget.NewLabel("Server"), nil, xf.hostSel),
	)

	xf.list = widget.NewList(
		func() int { return len(xf.rows) },
		func() fyne.CanvasObject { return newXfItem(w) },
		func(i widget.ListItemID, o fyne.CanvasObject) {
			if i < len(xf.rows) {
				o.(*xfItem).set(xf.rows[i])
			}
		},
	)

	return container.NewBorder(top, nil, nil, nil, container.NewVScroll(xf.list))
}

func refreshTransfers() {
	if xf.list == nil {
		return
	}
	xf.rows = xf.rows[:0]
	snaps := mgr.AllSnaps()
	for _, cfg := range mgr.Store.List() {
		if xf.cur != "" && xf.cur != "All servers" && cfg.Name != xf.cur {
			continue
		}
		s := snaps[cfg.ID]
		if s == nil || !s.Online {
			continue
		}
		for _, t := range s.Transfers {
			row := xfRow{hostID: cfg.ID, hostName: cfg.Name}
			row.t.Name = t.Name
			row.t.Project = t.ProjectName
			row.t.Upload = t.Upload
			row.t.Done = t.Done
			row.t.Total = t.Total
			row.t.Progress = t.Progress
			row.t.Paused = t.Paused
			row.t.Finished = t.Finished
			xf.rows = append(xf.rows, row)
		}
	}
	sort.SliceStable(xf.rows, func(i, j int) bool {
		a, b := xf.rows[i].t, xf.rows[j].t
		if a.Finished != b.Finished {
			return !a.Finished
		}
		return a.Name < b.Name
	})
	xf.list.Refresh()
}

type xfItem struct {
	fyne.Container
	name  *widget.Label
	sub   *widget.Label
	bar   *widget.ProgressBar
	state *widget.Label
	cur   xfRow
	has   bool
}

func newXfItem(w fyne.Window) *xfItem {
	it := &xfItem{}
	it.name = widget.NewLabel("")
	it.name.TextStyle = fyne.TextStyle{Bold: true}
	it.name.Wrapping = fyne.TextWrapWord
	it.sub = widget.NewLabel("")
	it.bar = widget.NewProgressBar()
	it.state = widget.NewLabel("")

	retry := widget.NewButtonWithIcon("", theme.ViewRefreshIcon(), nil)
	retry.OnTapped = func() {
		if it.has {
			go func() { _ = mgr.TransferOp(it.cur.hostID, it.cur.t.Name, "retry") }()
			go later(fireRefresh)
		}
	}
	abort := widget.NewButtonWithIcon("", theme.DeleteIcon(), nil)
	abort.Importance = widget.DangerImportance
	abort.OnTapped = func() {
		if !it.has {
			return
		}
		row := it.cur
		dialog.NewConfirm("Abort transfer", "Abort transfer of\n"+row.t.Name+"?", func(ok bool) {
			if ok {
				go func() { _ = mgr.TransferOp(row.hostID, row.t.Name, "abort") }()
				go later(fireRefresh)
			}
		}, w).Show()
	}

	head := container.NewHBox(it.name, layoutSpace(), retry, abort)
	body := container.NewVBox(head, it.sub, container.NewHBox(it.bar, container.NewCenter(it.state)))
	it.Objects = []fyne.CanvasObject{container.NewPadded(body)}
	return it
}

func (it *xfItem) set(r xfRow) {
	it.cur = r
	it.has = true
	t := r.t
	dir := "↓ download"
	if t.Upload {
		dir = "↑ upload"
	}
	state := ""
	switch {
	case t.Finished:
		state = "[DONE]"
	case t.Paused:
		state = "[PAUSED]"
	}
	it.name.SetText(t.Name)
	it.sub.SetText(fmt.Sprintf("%s · %s · %s · %s / %s", dir, t.Project, r.hostName,
		fmtBytes(t.Done), fmtBytes(t.Total)))
	it.bar.Value = clamp01(t.Progress)
	it.bar.Refresh()
	it.state.SetText(state + strings.Repeat(" ", 0))
}

func fmtBytes(b int64) string {
	const kb, mb, gb = 1024, 1048576, 1073741824
	switch {
	case b >= gb:
		return fmt.Sprintf("%.1f GB", float64(b)/gb)
	case b >= mb:
		return fmt.Sprintf("%.1f MB", float64(b)/mb)
	case b >= kb:
		return fmt.Sprintf("%.1f KB", float64(b)/kb)
	default:
		return fmt.Sprintf("%d B", b)
	}
}

func later(fn func()) {
	go func() {
		time.Sleep(400 * time.Millisecond)
		fyreDo(fn)
	}()
}
