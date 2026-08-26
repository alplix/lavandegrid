package ui

import (
	"fmt"
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/alplix/lavandegrid/internal/i18n"
	"github.com/alplix/lavandegrid/internal/local"
)

func init() {
	registerScreen("nav.settings", theme.MoreHorizontalIcon(), buildSettings, func() {})
}

func buildSettings(w fyne.Window) fyne.CanvasObject {
	prefs := fyneApp.Preferences()

	langNames := map[string]string{
		"en": "English", "tr": "Turkce", "de": "Deutsch", "fr": "Francais",
		"es": "Espanol", "it": "Italiano", "pt": "Portugues", "ru": "Русский", "ja": "日本語",
	}
	var langCodes []string
	var langLabels []string
	for _, c := range i18n.Available() {
		langCodes = append(langCodes, c)
		langLabels = append(langLabels, langNames[c])
	}

	langRadio := widget.NewRadioGroup(langLabels, func(val string) {
		for i, l := range langLabels {
			if l == val {
				i18n.SetLang(langCodes[i])
				RestartWindow()
				return
			}
		}
	})
	curIdx := 0
	for i, c := range langCodes {
		if c == i18n.CurrentCode() {
			curIdx = i
			break
		}
	}
	if curIdx < len(langLabels) {
		langRadio.SetSelected(langLabels[curIdx])
	}
	langCard := widget.NewCard(i18n.T("set.lang"), "", container.NewHBox(widget.NewLabel(i18n.T("set.lang") + ":"), langRadio))

	themeRadio := widget.NewRadioGroup([]string{i18n.T("cmd.darkTheme"), i18n.T("cmd.lightTheme")}, func(val string) {
		light := val == i18n.T("cmd.lightTheme")
		prefs.SetBool("light", light)
		ApplyTheme(fyneApp, light)
	})
	if prefs.Bool("light") {
		themeRadio.SetSelected(i18n.T("cmd.lightTheme"))
	} else {
		themeRadio.SetSelected(i18n.T("cmd.darkTheme"))
	}
	themeCard := widget.NewCard(i18n.T("set.appearance"), "", container.NewHBox(widget.NewLabel(i18n.T("set.theme") + ":"), themeRadio))

	notifCheck := widget.NewCheck(i18n.T("notif.enable"), nil)
	notifCheck.SetChecked(prefs.Bool("notifications"))
	notifCheck.OnChanged = func(on bool) { prefs.SetBool("notifications", on) }
	notifCard := widget.NewCard(i18n.T("set.notifT"), "", notifCheck)

	info := localInfoText()
	d := getDaemon()
	status := d.Status()
	runningColor := color.NRGBA{R: 0x34, G: 0xd3, B: 0x99, A: 0xff}
	stoppedColor := color.NRGBA{R: 0xfb, G: 0x71, B: 0x85, A: 0xff}
	statusLabel := i18n.T("set.stopped")
	dotColor := stoppedColor
	if status == local.DaemonRunning {
		dotColor = runningColor
		statusLabel = fmt.Sprintf("%s (PID %d)", i18n.T("set.running"), d.PID())
	}
	statusDot := canvas.NewText("●", dotColor)
	statusDot.TextSize = 14

	startBtn := widget.NewButtonWithIcon(i18n.T("set.startLocal"), theme.MediaPlayIcon(), func() {
		go func() {
			err := startLocalDaemon()
			msg := i18n.T("set.localStarted")
			if err != nil {
				msg = i18n.T("set.localFailed") + ": " + err.Error()
			}
			fyreDo(func() { dialog.NewInformation(i18n.T("nav.hosts"), msg, w).Show() })
		}()
	})
	stopBtn := widget.NewButtonWithIcon(i18n.T("set.stopLocal"), theme.MediaStopIcon(), func() {
		go func() {
			err := stopLocalDaemon()
			msg := i18n.T("set.localStopped")
			if err != nil {
				msg = i18n.T("set.localFailed") + ": " + err.Error()
			}
			fyreDo(func() { dialog.NewInformation(i18n.T("nav.hosts"), msg, w).Show() })
		}()
	})
	startBtn.Disable()
	stopBtn.Disable()
	if status == local.DaemonRunning {
		stopBtn.Enable()
	} else if status == local.DaemonStopped && d.Info.Found {
		startBtn.Enable()
	}

	statusRow := container.NewHBox(statusDot, widget.NewLabel(statusLabel))
	daemonBtns := container.NewHBox(startBtn, stopBtn)

	verText := ""
	if status == local.DaemonRunning {
		verText = fmt.Sprintf("LavandeGrid sends user_agent=LavandeGrid/%s to projects", Version)
	}

	daemonCard := widget.NewCard(i18n.T("set.localTitle"), info, container.NewVBox(statusRow, daemonBtns, widget.NewLabel(verText)))

	aboutBody := container.NewVBox(
		widget.NewLabelWithStyle("LavandeGrid v"+Version, fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		widget.NewLabel(i18n.T("about.desc")),
		widget.NewLabel(i18n.T("about.built")),
		widget.NewLabel(i18n.T("about.license")),
	)
	aboutCard := widget.NewCard(i18n.T("about.title"), "", aboutBody)

	return container.NewVScroll(container.NewVBox(
		container.NewPadded(langCard),
		container.NewPadded(themeCard),
		container.NewPadded(notifCard),
		container.NewPadded(daemonCard),
		container.NewPadded(aboutCard),
	))
}
