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
)

func init() {
	registerScreen("nav.stats", theme.MediaPhotoIcon(), buildStats, refreshStats)
}

var stv struct {
	hostSel *widget.Select
	days    *widget.Select
	mode    *widget.Select
	chart   *canvas.Image
	sumLbl  *widget.Label
	cur     string
	data    []statSeriesUI
	light   bool
}

type statSeriesUI struct {
	name   string
	colorN color.NRGBA
	points []chartPt
	labels []string
}

func buildStats(w fyne.Window) fyne.CanvasObject {
	stv.hostSel = widget.NewSelect(hostChoices(), func(s string) {
		stv.cur = s
		go loadStats()
	})
	stv.hostSel.PlaceHolder = "Choose a server"

	stv.days = widget.NewSelect([]string{"30 days", "90 days", "180 days", "1 year"}, func(string) { go loadStats() })
	stv.days.SetSelected("90 days")

	stv.mode = widget.NewSelect([]string{"Total credit", "Daily average (RAC)"}, func(string) {
		drawStats()
	})
	stv.mode.SetSelected("Total credit")

	stv.sumLbl = widget.NewLabel("")
	stv.chart = canvas.NewImageFromImage(renderChart(nil, 900, 320, true))
	stv.chart.FillMode = canvas.ImageFillStretch

	top := container.NewHBox(
		container.NewBorder(nil, nil, widget.NewLabel("Server"), nil, stv.hostSel),
		container.NewBorder(nil, nil, widget.NewLabel("Range"), nil, stv.days),
		container.NewBorder(nil, nil, widget.NewLabel("Metric"), nil, stv.mode),
	)

	body := container.NewVBox(
		widget.NewCard("Credit history per project", "", container.NewStack(stv.chart)),
		container.NewPadded(stv.sumLbl),
	)

	return container.NewBorder(top, nil, nil, nil, container.NewVScroll(body))
}

func refreshStats() {}

func loadStats() {
	hostID := ""
	for _, cfg := range mgr.Store.List() {
		if cfg.Name == stv.cur {
			hostID = cfg.ID
			break
		}
	}
	if hostID == "" {
		fyreDo(func() { stv.sumLbl.SetText("Choose a server to load statistics.") })
		return
	}
	stats, err := mgr.Stats(hostID)
	if err != nil {
		fyreDo(func() { stv.sumLbl.SetText("Failed to load stats: " + err.Error()) })
		return
	}

	days := 90
	switch stv.days.Selected {
	case "30 days":
		days = 30
	case "180 days":
		days = 180
	case "1 year":
		days = 365
	}
	cutoff := time.Now().AddDate(0, 0, -days).Format("20060102")

	var out []statSeriesUI
	for _, ss := range stats {
		ui := statSeriesUI{name: ss.Name}
		if ui.name == "" {
			ui.name = ss.URL
		}
		ui.colorN = projColorNRGBA(ui.name)
		startIdx := 0
		for i, d := range ss.Daily {
			if d.Day < cutoff {
				startIdx = i + 1
				continue
			}
			break
		}
		base := 0.0
		if startIdx > 0 && startIdx < len(ss.Daily) {
			base = ss.Daily[startIdx-1].HostCredit
		}
		for _, d := range ss.Daily[startIdx:] {
			dayT, _ := time.Parse("20060102", d.Day)
			ui.points = append(ui.points, chartPt{X: float64(len(ui.points)), Y: d.HostCredit - base})
			ui.labels = append(ui.labels, dayT.Format("Jan 02"))
		}
		out = append(out, ui)
	}
	sort.Slice(out, func(i, j int) bool {
		a, b := 0.0, 0.0
		if n := len(out[i].points); n > 0 {
			a = out[i].points[n-1].Y
		}
		if n := len(out[j].points); n > 0 {
			b = out[j].points[n-1].Y
		}
		return a > b
	})

	fyreDo(func() {
		stv.data = out
		stv.light = fyneApp.Preferences().Bool("light")
		drawStats()
	})
}

func drawStats() {
	mode := stv.mode.Selected
	var series []chartSeries
	totalGain := 0.0
	for _, s := range stv.data {
		cs := chartSeries{Label: s.name, Color: s.colorN, Labels: s.labels}
		for i, p := range s.points {
			y := p.Y
			if mode == "Daily average (RAC)" && i > 0 {
				y = p.Y - s.points[i-1].Y
			}
			cs.Points = append(cs.Points, chartPt{X: p.X, Y: y})
		}
		if n := len(cs.Points); n > 0 {
			totalGain += cs.Points[n-1].Y
		}
		series = append(series, cs)
	}
	stv.chart.Image = renderChart(series, 980, 340, !stv.light)
	stv.chart.Refresh()
	if mode == "Total credit" {
		stv.sumLbl.SetText(sumText(series, totalGain))
	} else {
		stv.sumLbl.SetText(legendText(series))
	}
}

func legendText(series []chartSeries) string {
	s := ""
	for i, c := range series {
		if i > 0 {
			s += " · "
		}
		s += c.Label
		if len(c.Points) > 0 {
			last := c.Points[len(c.Points)-1]
			prev := last.Y
			if len(c.Points) > 2 {
				prev = c.Points[len(c.Points)-3].Y
			}
			s += fmtNumShort(last.Y - prev)
		}
	}
	if s == "" {
		return "No data yet."
	}
	return s + " (credit/day, last 3 days)"
}

func sumText(series []chartSeries, gain float64) string {
	if len(series) == 0 {
		return "No data yet."
	}
	names := ""
	for i, c := range series {
		if i >= 4 {
			names += fmt.Sprintf(" · +%d more", len(series)-i)
			break
		}
		if names != "" && !endsWithMore(names) {
			names += " · "
		}
		last := ""
		if len(c.Points) > 0 {
			last = fmtNumShort(c.Points[len(c.Points)-1].Y)
		}
		names += c.Label + " " + last
	}
	return names + " credits in range"
}

func endsWithMore(s string) bool {
	n := len(s)
	return n > 5 && s[n-1] == 'e' && s[n-2] == 'r' && s[n-3] == 'o' && s[n-4] == 'm'
}

func fmtNumShort(v float64) string {
	switch {
	case v >= 1e9:
		return trimF(v/1e9) + "B"
	case v >= 1e6:
		return trimF(v/1e6) + "M"
	case v >= 1e3:
		return trimF(v/1e3) + "k"
	default:
		return trimF(v)
	}
}
