package ui

import (
	"sort"
	"strings"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	"github.com/alplix/lavandegrid/internal/i18n"
)

func init() {
	registerScreen("nav.messages", theme.MailComposeIcon(), buildMessages, refreshMessages)
}

var mv struct {
	search *widget.Entry
	sev    *widget.Select
	host   *widget.Select
	list   *widget.List
	rows   []msgRow
	curH   string
}

type msgRow struct {
	host string
	m    struct {
		Time  int64
		Pri   int
		Body  string
		Proj  string
	}
}

func buildMessages(w fyne.Window) fyne.CanvasObject {
	mv.search = widget.NewEntry()
	mv.search.PlaceHolder = i18n.T("msg.searchPh")
	mv.search.OnChanged = func(string) { refreshMessages() }

	mv.sev = widget.NewSelect([]string{"All", "Info", "Warnings", "Errors"}, func(string) { refreshMessages() })
	mv.sev.SetSelected("All")

	mv.host = widget.NewSelect(hostChoices(), func(s string) {
		mv.curH = s
		refreshMessages()
	})
	mv.host.PlaceHolder = "All servers"

	top := container.NewVBox(
		container.NewHBox(
			container.NewGridWithColumns(3,
				container.NewBorder(nil, nil, widget.NewLabel("Level"), nil, mv.sev),
				container.NewBorder(nil, nil, widget.NewLabel("Server"), nil, mv.host),
				container.NewBorder(nil, nil, widget.NewLabel("Search"), nil, mv.search),
			),
		),
	)

	mv.list = widget.NewList(
		func() int { return len(mv.rows) },
		func() fyne.CanvasObject {
			return container.NewVBox(
				container.NewHBox(widget.NewLabel(""), widget.NewLabel("")),
				widget.NewLabel(""),
			)
		},
		func(i widget.ListItemID, o fyne.CanvasObject) {
			if i >= len(mv.rows) {
				return
			}
			r := mv.rows[i]
			vbox := o.(*fyne.Container).Objects
			hbox := vbox[0].(*fyne.Container).Objects
			meta := hbox[0].(*widget.Label)
			host := hbox[1].(*widget.Label)
			body := vbox[1].(*widget.Label)
			meta.SetText(msgMeta(r.m.Pri, r.m.Time))
			host.SetText(r.host)
			body.Wrapping = fyne.TextWrapWord
			body.TextStyle = fyne.TextStyle{Monospace: true}
			body.SetText(r.m.Body)
		},
	)

	return container.NewBorder(top, nil, nil, nil, container.NewVScroll(mv.list))
}

func msgMeta(pri int, ts int64) string {
	tag := "INFO"
	switch {
	case pri >= 3:
		tag = "ERROR"
	case pri == 2:
		tag = "WARN"
	}
	return "[" + tag + "] " + fmtTimeShort(ts)
}

func fmtTimeShort(ts int64) string {
	return unixToStamp(ts)
}

func refreshMessages() {
	if mv.list == nil || mv.sev.Selected == "" {
		return
	}
	q := strings.ToLower(mv.search.Text)
	lvl := mv.sev.Selected
	rows := []msgRow{}
	snaps := mgr.AllSnaps()
	for _, cfg := range mgr.Store.List() {
		if mv.curH != "" && mv.curH != "All servers" && cfg.Name != mv.curH {
			continue
		}
		s := snaps[cfg.ID]
		if s == nil || !s.Online {
			continue
		}
		for i := len(s.Messages) - 1; i >= 0; i-- {
			m := s.Messages[i]
			ok := true
			switch lvl {
			case "Warnings":
				ok = m.Pri == 2
			case "Errors":
				ok = m.Pri >= 3
			}
			if ok && q != "" && !strings.Contains(strings.ToLower(m.Body), q) {
				ok = false
			}
			if !ok {
				continue
			}
			var r msgRow
			r.host = cfg.Name
			r.m.Time = m.Time.Unix()
			r.m.Pri = m.Pri
			r.m.Body = m.Body
			r.m.Proj = m.Project
			rows = append(rows, r)
		}
	}
	sort.SliceStable(rows, func(i, j int) bool { return rows[i].m.Time > rows[j].m.Time })
	mv.rows = rows
	mv.list.Refresh()
}

func unixToStamp(ts int64) string {
	return timeUnix(ts).Format("Jan 02 15:04:05")
}

func timeUnix(ts int64) time.Time {
	if ts <= 0 {
		return time.Unix(0, 0)
	}
	return time.Unix(ts, 0)
}
