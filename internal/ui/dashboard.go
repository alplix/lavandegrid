package ui

import (
	"fmt"
	"image/color"
	"sort"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/alplix/lavandegrid/internal/app"
	"github.com/alplix/lavandegrid/internal/i18n"
)

func init() {
	registerScreen("nav.dash", theme.HomeIcon(), buildDashboard, refreshDashboard)
}

var dashRefs struct {
	cards  []*widget.Label
	chart  *canvas.Image
	feed   *widget.List
	hostLs *widget.List
}

func statCard(title string) (fyne.CanvasObject, *widget.Label) {
	val := widget.NewLabel("–")
	val.TextStyle = fyne.TextStyle{Bold: true}
	card := widget.NewCard("", title, container.NewVBox(val))
	return card, val
}

type dashHostRow struct {
	cfg app.HostCfg
	snap *app.Snapshot
}

func buildDashboard(w fyne.Window) fyne.CanvasObject {
	dashRefs.cards = nil
	titles := []string{i18n.T("dash.running"), i18n.T("dash.queue"), i18n.T("st.error"), i18n.T("dash.fleetRac")}
	var cards []fyne.CanvasObject
	for _, t := range titles {
		c, lbl := statCard(t)
		dashRefs.cards = append(dashRefs.cards, lbl)
		cards = append(cards, c)
	}
	stats := container.NewGridWithColumns(4, cards...)

	dashRefs.chart = canvas.NewImageFromImage(renderChart(nil, 900, 220, !fyneApp.Preferences().Bool("light")))
	dashRefs.chart.FillMode = canvas.ImageFillStretch

	hostsTitle := widget.NewLabelWithStyle(i18n.T("dash.servers"), fyne.TextAlignLeading, fyne.TextStyle{Bold: true})

	rows := func() []dashHostRow {
		var out []dashHostRow
		for _, cfg := range mgr.Store.List() {
			out = append(out, dashHostRow{cfg: cfg, snap: mgr.Snap(cfg.ID)})
		}
		return out
	}

	dashRefs.hostLs = widget.NewList(
		func() int { return len(rows()) },
		func() fyne.CanvasObject {
			return newDashHostRow()
		},
		func(i widget.ListItemID, o fyne.CanvasObject) {
			r := rows()
			if i < len(r) {
				o.(*dashHostItem).set(r[i])
			}
		},
	)

	feedTitle := widget.NewLabelWithStyle(i18n.T("dash.recent"), fyne.TextAlignLeading, fyne.TextStyle{Bold: true})
	msgRows := func() []app.MsgLine {
		var all []app.MsgLine
		for _, s := range mgr.AllSnaps() {
			if s.Online {
				all = append(all, s.Messages...)
			}
		}
		sort.Slice(all, func(i, j int) bool { return all[i].Time.After(all[j].Time) })
		if len(all) > 14 {
			all = all[:14]
		}
		return all
	}
	dashRefs.feed = widget.NewList(
		func() int { return len(msgRows()) },
		func() fyne.CanvasObject {
			return container.NewHBox(widget.NewLabel(""), widget.NewLabel(""))
		},
		func(i widget.ListItemID, o fyne.CanvasObject) {
			mr := msgRows()
			if i >= len(mr) {
				return
			}
			m := mr[i]
			box := o.(*fyne.Container).Objects
			l1 := box[0].(*widget.Label)
			l2 := box[1].(*widget.Label)
			l1.SetText(app.FmtAgo(m.Time))
			l1.Resize(fyne.NewSize(70, l1.MinSize().Height))
			l2.Wrapping = fyne.TextWrapWord
			l2.SetText(m.Body)
		},
	)

	right := container.NewBorder(container.NewVBox(feedTitle), nil, nil, nil, dashRefs.feed)
	bottom := container.NewGridWithColumns(2,
		container.NewBorder(hostsTitle, nil, nil, nil, dashRefs.hostLs),
		right,
	)

	return container.NewBorder(stats, bottom, nil, nil,
		container.NewVScroll(container.NewVBox(
			widget.NewCard(i18n.T("dash.actT"), "", container.NewStack(dashRefs.chart)),
			widget.NewLabel(""),
		)),
	)
}

func refreshDashboard() {
	totalRun, queued, errs := 0, 0, 0
	rac := 0.0
	for _, s := range mgr.AllSnaps() {
		if !s.Online {
			continue
		}
		totalRun += s.Totals.Running
		queued += s.Totals.Queued
		errs += s.Totals.Errors
		rac += s.Totals.RAC
	}
	vals := []string{
		fmt.Sprintf("%d", totalRun),
		fmt.Sprintf("%d", queued),
		fmt.Sprintf("%d", errs),
		app.FmtNum(rac),
	}
	for i, v := range vals {
		if i < len(dashRefs.cards) {
			dashRefs.cards[i].SetText(v)
		}
	}
	light := fyneApp.Preferences().Bool("light")
	series := fleetSeries(light)
	dashRefs.chart.Image = renderChart(series, 900, 220, !light)
	dashRefs.chart.Refresh()

	if dashRefs.hostLs != nil {
		dashRefs.hostLs.Refresh()
	}
	if dashRefs.feed != nil {
		dashRefs.feed.Refresh()
	}
}

func fleetSeries(light bool) []chartSeries {
	var out []chartSeries
	for _, h := range mgr.Store.List() {
		hist := mgr.History(h.ID)
		if len(hist) < 2 {
			continue
		}
		col := projColorNRGBA(h.Name)
		if light {
			col = darken(col)
		}
		cs := chartSeries{Label: h.Name, Color: col}
		labels := make([]string, len(hist))
		now := time.Now()
		for i, p := range hist {
			cs.Points = append(cs.Points, chartPt{X: float64(i), Y: float64(p.Running)})
			labels[i] = now.Sub(p.T).Round(time.Minute).String() + " ago"
		}
		cs.Labels = labels
		out = append(out, cs)
	}
	return out
}

func darken(c color.NRGBA) color.NRGBA {
	return color.NRGBA{
		R: uint8(float64(c.R) * 0.62),
		G: uint8(float64(c.G) * 0.62),
		B: uint8(float64(c.B) * 0.62),
		A: 0xff,
	}
}

// ---- host row widget ----

type dashHostItem struct {
	fyne.Container
	name    *widget.Label
	sub     *widget.Label
	gpu     *widget.Label
	dot     *canvas.Circle
	version *widget.Label
}

func newDashHostRow() *dashHostItem {
	it := &dashHostItem{}
	it.dot = &canvas.Circle{FillColor: color.NRGBA{R: 0x80, G: 0x80, B: 0x80, A: 0xff}}
	it.dot.Resize(fyne.NewSize(10, 10))
	it.name = widget.NewLabel("")
	it.name.TextStyle = fyne.TextStyle{Bold: true}
	it.version = widget.NewLabel("")
	it.sub = widget.NewLabel("")
	it.gpu = widget.NewLabel("")

	head := container.NewHBox(container.NewCenter(it.dot), it.name, it.version)
	body := container.NewVBox(head, it.sub, it.gpu)
	it.Objects = []fyne.CanvasObject{container.NewPadded(body)}
	return it
}

func textMutedColor() color.Color {
	return color.NRGBA{R: 0xa8, G: 0x9c, B: 0xc6, A: 0xff}
}

func statusDot(online bool) color.Color {
	if online {
		return color.NRGBA{R: 0x34, G: 0xd3, B: 0x99, A: 0xff}
	}
	return color.NRGBA{R: 0xfb, G: 0x71, B: 0x85, A: 0xff}
}

func (it *dashHostItem) set(r dashHostRow) {
	it.dot.FillColor = statusDot(r.snap != nil && r.snap.Online)
	canvas.Refresh(it.dot)
	mode := ""
	if r.cfg.Demo {
		mode = "  [demo]"
	}
	it.name.SetText(r.cfg.Name + mode)
	addr := r.cfg.Host + ":" + fmt.Sprintf("%d", r.cfg.Port)
	if r.snap == nil || !r.snap.Online {
		it.version.SetText("")
		errTxt := "offline"
		if r.snap != nil && r.snap.Error != "" {
			errTxt = "offline · " + r.snap.Error
		}
		it.sub.SetText(addr + " · " + errTxt)
		it.gpu.SetText("")
		return
	}
	s := r.snap
	it.version.SetText("BOINC " + s.Version)
	it.sub.SetText(fmt.Sprintf("%s · %d running · %d paused · %d queued · RAC %s · credit %s",
		addr, s.Totals.Running, s.Totals.Paused, s.Totals.Queued, app.FmtNum(s.Totals.RAC), app.FmtNum(s.Totals.Credit)))
	gpuLine := ""
	for _, g := range s.HostInfo.GPUs {
		name := ""
		if len(g.Names) > 0 {
			name = g.Names[0]
		}
		if g.Count > 1 {
			name = fmt.Sprintf("%dx %s", g.Count, name)
		}
		line := name
		if g.VRAM > 0 {
			line += " · " + app.FmtBytes(g.VRAM)
		}
		if g.Cuda != "" {
			line += " · CUDA " + g.Cuda
		}
		gpuLine = line
		break
	}
	cpu := s.HostInfo.CPU
	if cpu == "" {
		cpu = "CPU n/a"
	}
	spec := fmt.Sprintf("%s · %d cores · %s RAM", cpu, s.HostInfo.Cores, app.FmtBytes(s.HostInfo.Memory))
	if s.HostInfo.Flops > 0 {
		spec += fmt.Sprintf(" · %.0f GFLOPS", s.HostInfo.Flops/1e9)
	}
	if gpuLine != "" {
		spec += "\n" + gpuLine
	}
	it.gpu.SetText(spec)
}

func projColorNRGBA(key string) color.NRGBA {
	var h uint32
	for _, r := range key {
		h = h*31 + uint32(r)
	}
	hue := float64(h%360) / 360
	return hslToNRGBA(hue, 0.62, 0.62)
}

func hslToNRGBA(h, s, l float64) color.NRGBA {
	var r, g, b float64
	if s == 0 {
		r, g, b = l, l, l
	} else {
		q := l + s - l*s
		if l >= 0.5 {
			q = l + s - l*s
		}
		p := 2*l - q
		r = hue2rgb(p, q, h+1.0/3.0)
		g = hue2rgb(p, q, h)
		b = hue2rgb(p, q, h-1.0/3.0)
	}
	to8 := func(v float64) uint8 {
		if v < 0 {
			v = 0
		}
		if v > 1 {
			v = 1
		}
		return uint8(v * 255)
	}
	return color.NRGBA{R: to8(r), G: to8(g), B: to8(b), A: 0xff}
}

func hue2rgb(p, q, t float64) float64 {
	if t < 0 {
		t++
	}
	if t > 1 {
		t--
	}
	switch {
	case t < 1.0/6.0:
		return p + (q-p)*6*t
	case t < 0.5:
		return q
	case t < 2.0/3.0:
		return p + (q-p)*(2.0/3.0-t)*6
	default:
		return p
	}
}

var _ = time.Now
