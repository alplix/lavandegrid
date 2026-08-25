package ui

import (
	"fmt"
	"strconv"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/alplix/lavandegrid/internal/app"
)

type serverItem struct {
	fyne.Container
	name *widget.Label
	addr *widget.Label
	info *widget.Label
	dot  *canvas.Circle
	ver  *widget.Label
	cur  serverRow
	has  bool
}

func newServerItem(w fyne.Window) *serverItem {
	it := &serverItem{}
	it.dot = &canvas.Circle{FillColor: statusDot(false)}
	it.name = widget.NewLabel("")
	it.name.TextStyle = fyne.TextStyle{Bold: true}
	it.ver = widget.NewLabel("")
	it.addr = widget.NewLabel("")
	it.info = widget.NewLabel("")
	it.info.Wrapping = fyne.TextWrapWord

	edit := widget.NewButtonWithIcon("", theme.DocumentCreateIcon(), nil)
	test := widget.NewButtonWithIcon("", theme.QuestionIcon(), nil)
	del := widget.NewButtonWithIcon("", theme.DeleteIcon(), nil)

	edit.OnTapped = func() {
		if !it.has {
			return
		}
		base := it.cur.cfg
		showHostDialog(w, base, func(nc app.HostCfg) {
			nc.ID = base.ID
			mgr.Store.Upsert(nc)
			_ = mgr.Store.Save()
			mgr.Kick(nc.ID)
			fireRefresh()
		})
	}
	test.OnTapped = func() {
		if !it.has {
			return
		}
		cfg := it.cur.cfg
		go func() {
			ver, err := mgr.TestHost(cfg)
			msg := "Connected. BOINC version " + ver
			if err != nil {
				msg = "Connection failed:\n" + err.Error()
			}
			fyreDo(func() { dialog.NewInformation("Test connection", msg, w).Show() })
		}()
	}
	del.OnTapped = func() {
		if !it.has {
			return
		}
		cfg := it.cur.cfg
		dialog.NewConfirm("Remove server", "Remove "+cfg.Name+" from LavandeGrid?", func(ok bool) {
			if ok {
				mgr.Store.Remove(cfg.ID)
				_ = mgr.Store.Save()
				fireRefresh()
			}
		}, w).Show()
	}

	head := container.NewHBox(
		container.NewCenter(it.dot),
		it.name, it.ver, layoutSpace(), edit, test, del,
	)
	body := container.NewVBox(head, it.addr, it.info)
	it.Objects = []fyne.CanvasObject{container.NewPadded(body)}
	return it
}

func (it *serverItem) set(r serverRow) {
	it.cur = r
	it.has = true
	cfg, s := r.cfg, r.snap
	mode := ""
	if cfg.Demo {
		mode = "  [demo]"
	}
	it.name.SetText(cfg.Name + mode)
	addr := fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)
	online := s != nil && s.Online
	it.dot.FillColor = statusDot(online)
	canvas.Refresh(it.dot)

	if online {
		it.ver.SetText("BOINC " + s.Version)
		hostLine := s.HostInfo.OS + " · " + strconv.Itoa(s.HostInfo.Cores) + " cores · " +
			app.FmtBytes(s.HostInfo.Memory) + " RAM"
		for _, g := range s.HostInfo.GPUs {
			if len(g.Names) == 0 {
				continue
			}
			gname := g.Names[0]
			if g.Count > 1 {
				gname = fmt.Sprintf("%dx %s", g.Count, gname)
			}
			hostLine += "\nGPU: " + gname
			if g.Driver != "" {
				hostLine += " · driver " + g.Driver
			}
			if g.Cuda != "" {
				hostLine += " · CUDA " + g.Cuda
			}
			if g.VRAM > 0 {
				hostLine += " · " + app.FmtBytes(g.VRAM)
			}
		}
		spec := fmt.Sprintf("%d running · %d paused · %d queued · %d errors · RAC %s · credit %s",
			s.Totals.Running, s.Totals.Paused, s.Totals.Queued, s.Totals.Errors,
			app.FmtNum(s.Totals.RAC), app.FmtNum(s.Totals.Credit))
		it.addr.SetText(addr + " · " + hostLine)
		it.info.SetText(spec)
	} else {
		it.ver.SetText("")
		errTxt := "offline"
		if s != nil && s.Error != "" {
			errTxt += " · " + s.Error
		}
		it.addr.SetText(addr + " · " + errTxt)
		it.info.SetText("")
	}
}

func refreshServers() {
	if sv.list == nil {
		return
	}
	n := len(mgr.Store.List())
	online := 0
	for _, s := range mgr.AllSnaps() {
		if s.Online {
			online++
		}
	}
	sv.hint.SetText(fmt.Sprintf("%d servers · %d online", n, online))
	sv.list.Refresh()
}

func showHostDialog(w fyne.Window, base app.HostCfg, onSave func(app.HostCfg)) {
	title := "Add server"
	if base.ID != "" {
		title = "Edit server"
	}
	nameE := widget.NewEntry()
	nameE.SetText(base.Name)
	nameE.PlaceHolder = "My laptop"
	hostE := widget.NewEntry()
	hostE.SetText(base.Host)
	hostE.PlaceHolder = "192.168.1.20"
	portE := widget.NewEntry()
	portE.SetText(strconv.Itoa(orPort(base.Port)))
	passE := widget.NewPasswordEntry()
	passE.SetText(base.Password)
	demoC := widget.NewCheck("Demo mode (simulated data)", nil)
	demoC.SetChecked(base.Demo)

	help := widget.NewLabel("The RPC password is in gui_rpc_auth.cfg inside the BOINC data directory of the remote machine.")
	help.Wrapping = fyne.TextWrapWord
	items := []*widget.FormItem{
		widget.NewFormItem("Name", nameE),
		widget.NewFormItem("Host / IP", hostE),
		widget.NewFormItem("GUI RPC port", portE),
		widget.NewFormItem("RPC password", passE),
		widget.NewFormItem("", demoC),
		widget.NewFormItem("", help),
	}

	dlg := dialog.NewForm(title, "Save", "Cancel", items, func(ok bool) {
		if !ok || strings.TrimSpace(hostE.Text) == "" {
			return
		}
		port, _ := strconv.Atoi(strings.TrimSpace(portE.Text))
		cfg := app.HostCfg{
			Name:     fallback(nameE.Text, hostE.Text),
			Host:     strings.TrimSpace(hostE.Text),
			Port:     clampPort(port),
			Password: passE.Text,
			Demo:     demoC.Checked,
		}
		if onSave != nil {
			onSave(cfg)
		}
	}, w)
	dlg.Resize(fyne.NewSize(480, 440))
	dlg.Show()
}

func orPort(p int) int {
	if p > 0 {
		return p
	}
	return 31416
}

func clampPort(p int) int {
	if p <= 0 || p > 65535 {
		return 31416
	}
	return p
}

func fallback(a, b string) string {
	if a != "" {
		return a
	}
	return b
}

func layoutSpace() fyne.CanvasObject {
	sp := widget.NewLabel("")
	sp.Hide()
	return sp
}
