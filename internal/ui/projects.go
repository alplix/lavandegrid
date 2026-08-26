package ui

import (
	"fmt"
	"sort"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/alplix/lavandegrid/internal/app"
)

func init() {
	registerScreen("nav.projects", theme.StorageIcon(), buildProjects, refreshProjects)
}

var pv struct {
	hostSel *widget.Select
	list    *widget.List
	rows    []projRow
	cur     string
}

type projRow struct {
	hostID   string
	hostName string
	p        app.ProjectInfo
}

func buildProjects(w fyne.Window) fyne.CanvasObject {
	pv.hostSel = widget.NewSelect(hostChoices(), func(s string) {
		pv.cur = s
		refreshProjects()
	})
	pv.hostSel.PlaceHolder = "Choose a server"

	attach := widget.NewButtonWithIcon("Attach to project", theme.ContentAddIcon(), func() {
		showAttachDialog(w)
	})
	attach.Importance = widget.HighImportance

	top := container.NewHBox(
		container.NewBorder(nil, nil, widget.NewLabel("Server"), nil, pv.hostSel),
		attach,
	)

	pv.list = widget.NewList(
		func() int { return len(pv.rows) },
		func() fyne.CanvasObject { return newProjItem(w) },
		func(i widget.ListItemID, o fyne.CanvasObject) {
			if i < len(pv.rows) {
				o.(*projItem).set(pv.rows[i])
			}
		},
	)

	return container.NewBorder(top, nil, nil, nil, container.NewVScroll(pv.list))
}

func hostChoices() []string {
	var out []string
	for _, h := range mgr.Store.List() {
		out = append(out, h.Name)
	}
	sort.Strings(out)
	return out
}

func refreshProjects() {
	if pv.list == nil {
		return
	}
	pv.rows = pv.rows[:0]
	if pv.cur == "" && len(pv.hostSel.Options) > 0 {
		pv.cur = pv.hostSel.Options[0]
		pv.hostSel.SetSelected(pv.cur)
	}
	for _, cfg := range mgr.Store.List() {
		if pv.cur != "" && cfg.Name != pv.cur {
			continue
		}
		s := mgr.Snap(cfg.ID)
		if s == nil || !s.Online {
			continue
		}
		for _, p := range s.Projects {
			pv.rows = append(pv.rows, projRow{hostID: cfg.ID, hostName: cfg.Name, p: p})
		}
	}
	sort.SliceStable(pv.rows, func(i, j int) bool { return pv.rows[i].p.Name < pv.rows[j].p.Name })
	pv.list.Refresh()
}

type projItem struct {
	fyne.Container
	name  *widget.Label
	url   *widget.Label
	stats *widget.Label
	flags *widget.Label
	cur   projRow
	has   bool
}

func newProjItem(w fyne.Window) *projItem {
	it := &projItem{}
	it.name = widget.NewLabel("")
	it.name.TextStyle = fyne.TextStyle{Bold: true}
	it.url = widget.NewLabel("")
	it.stats = widget.NewLabel("")
	it.flags = widget.NewLabel("")

	mkAct := func(label string, icon fyne.Resource, op string) *widget.Button {
		b := widget.NewButtonWithIcon(label, icon, nil)
		b.OnTapped = func() {
			if !it.has {
				return
			}
			go func() { _ = mgr.ProjectOp(it.cur.hostID, it.cur.p.URL, op) }()
			go func() {
				time.Sleep(400 * time.Millisecond)
				fireRefresh()
			}()
		}
		return b
	}
	toggleRun := widget.NewButtonWithIcon("", theme.MediaPauseIcon(), nil)
	toggleRun.OnTapped = func() {
		if !it.has {
			return
		}
		op := "suspend"
		if it.cur.p.Suspended {
			op = "resume"
		}
		go func() { _ = mgr.ProjectOp(it.cur.hostID, it.cur.p.URL, op) }()
		go func() {
			time.Sleep(400 * time.Millisecond)
			fireRefresh()
		}()
	}
	toggleWork := widget.NewButtonWithIcon("", theme.CancelIcon(), nil)
	toggleWork.OnTapped = func() {
		if !it.has {
			return
		}
		op := "nomorework"
		if it.cur.p.NoMoreWork {
			op = "allowmorework"
		}
		go func() { _ = mgr.ProjectOp(it.cur.hostID, it.cur.p.URL, op) }()
		go func() {
			time.Sleep(400 * time.Millisecond)
			fireRefresh()
		}()
	}
	updateBtn := mkAct("", theme.ViewRefreshIcon(), "update")
	resetBtn := mkAct("", theme.HistoryIcon(), "reset")
	detachBtn := mkAct("", theme.DeleteIcon(), "detach")
	detachBtn.Importance = widget.DangerImportance
	detachBtn.OnTapped = func() {
		if !it.has {
			return
		}
		row := it.cur
		dialog.NewConfirm("Detach project",
			"Detach from "+row.p.Name+"?\nAll tasks and credit for this project on this host will be removed.",
			func(ok bool) {
				if ok {
					go func() { _ = mgr.ProjectOp(row.hostID, row.p.URL, "detach") }()
					go func() {
						time.Sleep(500 * time.Millisecond)
						fireRefresh()
					}()
				}
			}, w).Show()
	}

	head := container.NewHBox(it.name, layoutSpace(),
		toggleRun, toggleWork, updateBtn, resetBtn, detachBtn)
	body := container.NewVBox(head, it.url, it.stats, it.flags)
	it.Objects = []fyne.CanvasObject{container.NewPadded(body)}
	return it
}

func (it *projItem) set(r projRow) {
	it.cur = r
	it.has = true
	p := r.p
	flag := ""
	if p.Pending {
		flag = "[syncing…] "
	}
	if p.Suspended {
		flag += "[paused] "
	}
	if p.NoMoreWork {
		flag += "[no new work] "
	}
	if p.Ended {
		flag += "[ended] "
	}
	it.name.SetText(flag + p.Name)
	it.url.SetText(p.URL)
	it.stats.SetText(fmt.Sprintf(
		"share %.0f%% · host RAC %s (total %s) · user RAC %s (total %s)",
		p.Share, app.FmtNum(p.HostRAC), app.FmtNum(p.HostCredit),
		app.FmtNum(p.RAC), app.FmtNum(p.UserCredit)))
	acct := ""
	if p.UserName != "" {
		acct += "user: " + p.UserName
	}
	if p.TeamName != "" {
		acct += " · team: " + p.TeamName
	}
	if p.Venue != "" {
		acct += " · venue: " + p.Venue
	}
	it.flags.SetText(acct)
}
