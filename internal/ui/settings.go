package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

func init() {
	registerScreen("Settings", theme.MoreHorizontalIcon(), buildSettings, func() {})
}

func buildSettings(w fyne.Window) fyne.CanvasObject {
	prefs := fyneApp.Preferences()

	themeRadio := widget.NewRadioGroup([]string{"Dark", "Light"}, func(val string) {
		light := val == "Light"
		prefs.SetBool("light", light)
		ApplyTheme(fyneApp, light)
	})
	if prefs.Bool("light") {
		themeRadio.SetSelected("Light")
	} else {
		themeRadio.SetSelected("Dark")
	}
	themeCard := widget.NewCard("Appearance", "", container.NewHBox(widget.NewLabel("Color theme:"), themeRadio))

	notifCheck := widget.NewCheck("Desktop notifications for task errors and hosts going offline", nil)
	notifCheck.SetChecked(prefs.Bool("notifications"))
	notifCheck.OnChanged = func(on bool) { prefs.SetBool("notifications", on) }
	notifCard := widget.NewCard("Notifications", "", notifCheck)

	info := localInfoText()
	localStart := widget.NewButtonWithIcon("Start bundled BOINC client", theme.MediaPlayIcon(), func() {
		go func() {
			err := startLocalClient()
			msg := "BOINC client started. It may take a few seconds before it accepts connections."
			if err != nil {
				msg = "Failed to start: " + err.Error()
			}
			fyreDo(func() { dialog.NewInformation("Local BOINC", msg, w).Show() })
		}()
	})
	localCard := widget.NewCard("Local BOINC client", info, localStart)

	aboutBody := container.NewVBox(
		widget.NewLabelWithStyle("LavandeGrid v"+Version, fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		widget.NewLabel("A lavender-themed manager for your BOINC fleet."),
		widget.NewLabel("Built with Go + Fyne. Connects to BOINC clients via GUI RPC."),
		widget.NewLabel("© 2026 Alperen Yavuz · MIT License"),
	)
	aboutCard := widget.NewCard("About", "", aboutBody)

	return container.NewVScroll(container.NewVBox(
		container.NewPadded(aboutCard),
		container.NewPadded(themeCard),
		container.NewPadded(notifCard),
		container.NewPadded(localCard),
	))
}
