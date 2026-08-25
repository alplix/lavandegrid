package ui

import (
	"sort"
	"strconv"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

func init() {
	registerScreen("Preferences", theme.SettingsIcon(), buildPrefs, refreshPrefs)
}

var prefFields = []string{
	"max_cpus", "max_ncpus_pct", "cpu_usage_limit",
	"max_bytes_sec_up", "max_bytes_sec_down",
	"disk_max_used_gb", "disk_min_free_gb", "daily_xfer_limit_mb", "daily_xfer_period_days",
	"work_buf_min_days", "work_buf_additional_days",
	"suspend_cpu_usage", "suspend_if_user_active", "suspend_gpu_while_user_active",
	"cpu_scheduling_period_minutes", "end_on_battery",
}

var prefLabels = map[string]string{
	"max_cpus":                    "Max CPU cores used",
	"max_ncpus_pct":               "Max CPU usage (% of all cores)",
	"cpu_usage_limit":             "CPU limit per task (%)",
	"max_bytes_sec_up":            "Max upload rate (bytes/s)",
	"max_bytes_sec_down":          "Max download rate (bytes/s)",
	"disk_max_used_gb":            "Max disk used (GB)",
	"daily_xfer_limit_mb":         "Daily transfer limit (MB)",
	"daily_xfer_period_days":      "Transfer limit period (days)",
	"disk_min_free_gb":            "Keep free disk space (GB)",
	"work_buf_min_days":           "Work buffer min (days)",
	"work_buf_additional_days":    "Work buffer extra (days)",
	"suspend_cpu_usage":           "Suspend when CPU above (%)",
	"suspend_if_user_active":      "Suspend when user active (0/1)",
	"suspend_gpu_while_user_active": "Suspend GPU when user active (0/1)",
	"cpu_scheduling_period_minutes": "CPU scheduling period (min)",
	"end_on_battery":              "Stop on battery (0/1)",
}

var pfv struct {
	hostSel  *widget.Select
	fields   map[string]*widget.Entry
	status   *widget.Label
	runRow   *fyne.Container
	netRow   *fyne.Container
	cur      string
}

func buildPrefs(w fyne.Window) fyne.CanvasObject {
	pfv.fields = map[string]*widget.Entry{}
	pfv.hostSel = widget.NewSelect(hostChoices(), func(s string) {
		pfv.cur = s
		go loadPrefs()
	})
	pfv.hostSel.PlaceHolder = "Choose a server"

	modeBtn := func(label string, op, mode string) *widget.Button {
		b := widget.NewButton(label, nil)
		b.OnTapped = func() {
			if pfv.cur == "" {
				return
			}
			id := hostIDByName(pfv.cur)
			if id == "" {
				return
			}
			go func() { _ = mgr.ClientOp(id, op, mode) }()
			later(fireRefresh)
		}
		return b
	}
	pfv.runRow = container.NewHBox(
		widget.NewLabel("Compute mode:"),
		modeBtn("Run always", "setRunMode", "always"),
		modeBtn("Auto", "setRunMode", "auto"),
		modeBtn("Suspend", "setRunMode", "never"),
	)
	pfv.netRow = container.NewHBox(
		widget.NewLabel("Network mode:"),
		modeBtn("Always", "setNetworkMode", "always"),
		modeBtn("Auto", "setNetworkMode", "auto"),
		modeBtn("Never", "setNetworkMode", "never"),
	)
	bench := widget.NewButtonWithIcon("Run CPU benchmarks", theme.MediaFastForwardIcon(), func() {
		if pfv.cur == "" {
			return
		}
		id := hostIDByName(pfv.cur)
		go func() { _ = mgr.ClientOp(id, "benchmarks", "") }()
	})
	clearOvr := widget.NewButton("Clear local overrides", nil)
	clearOvr.Importance = widget.DangerImportance
	clearOvr.OnTapped = func() {
		if pfv.cur == "" {
			return
		}
		dialog.NewConfirm("Clear overrides",
			"Remove ALL locally overridden preferences on this server?\nServer defaults will apply again.",
			func(ok bool) {
				if ok {
					go func() { _ = mgr.PrefsClear(hostIDByName(pfv.cur)) }()
					go loadPrefs()
				}
			}, w).Show()
	}

	saveBtn := widget.NewButtonWithIcon("Save overrides", theme.ConfirmIcon(), savePrefs)
	saveBtn.Importance = widget.HighImportance
	pfv.status = widget.NewLabel("")

	var leftCol, rightCol []fyne.CanvasObject
	for i, f := range prefFields {
		e := widget.NewEntry()
		e.PlaceHolder = "default"
		pfv.fields[f] = e
		item := container.NewBorder(nil, nil, widget.NewLabel(prefLabels[f]), nil, e)
		if i%2 == 0 {
			leftCol = append(leftCol, item)
		} else {
			rightCol = append(rightCol, item)
		}
	}
	grid := container.NewGridWithColumns(2,
		container.NewVBox(leftCol...), container.NewVBox(rightCol...))

	top := container.NewVBox(
		container.NewHBox(
			container.NewBorder(nil, nil, widget.NewLabel("Server"), nil, pfv.hostSel),
			saveBtn,
		),
		container.NewPadded(pfv.status),
	)

	body := container.NewVScroll(container.NewVBox(
		widget.NewCard("Client modes", "", container.NewVBox(pfv.runRow, pfv.netRow, bench)),
		widget.NewCard("Local preference overrides", "Empty = use project defaults. Values here override the server's global prefs.", grid),
		container.NewHBox(clearOvr),
	))

	return container.NewBorder(top, nil, nil, nil, body)
}

func hostIDByName(name string) string {
	for _, cfg := range mgr.Store.List() {
		if cfg.Name == name {
			return cfg.ID
		}
	}
	return ""
}

func refreshPrefs() {}

func loadPrefs() {
	if pfv.hostSel == nil || pfv.cur == "" || len(pfv.fields) == 0 {
		return
	}
	id := hostIDByName(pfv.cur)
	if id == "" {
		return
	}
	vals, err := mgr.PrefsGet(id)
	fyreDo(func() {
		if err != nil {
			pfv.status.SetText("Failed to load prefs: " + err.Error())
			return
		}
		keys := make([]string, 0, len(vals))
		for k := range vals {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		for _, f := range prefFields {
			pfv.fields[f].SetText(vals[f])
		}
		extra := 0
		known := map[string]bool{}
		for _, f := range prefFields {
			known[f] = true
		}
		for _, k := range keys {
			if !known[k] && vals[k] != "" {
				extra++
			}
		}
		if extra > 0 {
			pfv.status.SetText("Loaded. (" + strconv.Itoa(extra) + " advanced overrides not shown are kept.)")
		} else {
			pfv.status.SetText("Loaded.")
		}
	})
}

func savePrefs() {
	if pfv.cur == "" {
		return
	}
	var pairs [][2]string
	for _, f := range prefFields {
		v := strings.TrimSpace(pfv.fields[f].Text)
		if v != "" {
			pairs = append(pairs, [2]string{f, v})
		}
	}
	err := mgr.PrefsSet(hostIDByName(pfv.cur), pairs)
	if err != nil {
		pfv.status.SetText("Save failed: " + err.Error())
		return
	}
	pfv.status.SetText("Saved.")
}
