package i18n

func init() {
	reg("en", map[string]string{
		"app.name": "LavandeGrid",
		"local.bundled": "Bundled BOINC client found", "local.system": "System BOINC client found",
		"local.none": "No local BOINC client detected.",
	})
}