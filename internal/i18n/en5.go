package i18n

func init() {
	reg("en", map[string]string{
		"set.lang": "Language", "set.appearance": "Appearance", "set.theme": "Color theme",
		"set.startLocal": "Start bundled BOINC client", "set.localStarted": "BOINC client started. It may take a few seconds before it accepts connections.",
		"set.localFailed": "Failed to start", "set.localTitle": "Local BOINC client",
		"set.stopped": "Stopped", "set.running": "Running", "set.stopLocal": "Stop client",
		"set.localStopped": "BOINC client stopped",
		"about.title": "About", "about.desc": "A lavender-themed manager for your BOINC fleet.",
		"about.built": "Built with Go + Fyne. Connects to BOINC clients via GUI RPC.",
		"about.license": "2026 Alperen Yavuz. MIT License.",
	})
}