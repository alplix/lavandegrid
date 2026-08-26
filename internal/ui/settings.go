package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/alplix/lavandegrid/internal/i18n"
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
	localStart := widget.NewButtonWithIcon(i18n.T("set.startLocal"), theme.MediaPlayIcon(), func() {
		go func() {
			err := startLocalClient()
			msg := i18n.T("set.localStarted")
			if err != nil {
				msg = i18n.T("set.localFailed") + ": " + err.Error()
			}
			fyreDo(func() { dialog.NewInformation(i18n.T("nav.hosts"), msg, w).Show() })
		}()
	})
	localCard := widget.NewCard(i18n.T("set.localTitle"), info, localStart)

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
		container.NewPadded(localCard),
		container.NewPadded(aboutCard),
	))
}
