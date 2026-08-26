package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/alplix/lavandegrid/internal/app"
	"github.com/alplix/lavandegrid/internal/i18n"
)

func init() {
	registerScreen("nav.hosts", theme.ComputerIcon(), buildServers, refreshServers)
}

var sv struct {
	list *widget.List
	hint *widget.Label
}

type serverRow struct {
	cfg  app.HostCfg
	snap *app.Snapshot
}

func buildServers(w fyne.Window) fyne.CanvasObject {
	addBtn := widget.NewButtonWithIcon(i18n.T("hosts.add"), theme.ContentAddIcon(), func() {
		showHostDialog(w, app.HostCfg{}, func(nc app.HostCfg) {
			saved := mgr.Store.Upsert(nc)
			_ = mgr.Store.Save()
			mgr.Kick(saved.ID)
			fireRefresh()
		})
	})
	addBtn.Importance = widget.HighImportance
	sv.hint = widget.NewLabel("")

	sv.list = widget.NewList(
		func() int { return len(serverRows()) },
		func() fyne.CanvasObject { return newServerItem(w) },
		func(i widget.ListItemID, o fyne.CanvasObject) {
			rows := serverRows()
			if i < len(rows) {
				o.(*serverItem).set(rows[i])
			}
		},
	)

	return container.NewBorder(
		container.NewHBox(addBtn, sv.hint),
		nil, nil, nil,
		container.NewVScroll(sv.list),
	)
}

func serverRows() []serverRow {
	var rows []serverRow
	for _, cfg := range mgr.Store.List() {
		rows = append(rows, serverRow{cfg: cfg, snap: mgr.Snap(cfg.ID)})
	}
	return rows
}
