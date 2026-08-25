package ui

import (
	"fmt"
	"strings"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

type projCatalogEntry struct {
	Name string
	URL  string
}

var projCatalog = []projCatalogEntry{
	{"Einstein@Home", "https://einsteinathome.org/"},
	{"Rosetta@home", "https://boinc.bakerlab.org/rosetta/"},
	{"Milkyway@home", "https://milkyway.cs.rpi.edu/milkyway/"},
	{"PrimeGrid", "https://www.primegrid.org/"},
	{"World Community Grid", "https://www.worldcommunitygrid.org/"},
	{"GPUGRID.net", "https://www.gpugrid.net/"},
	{"Amicable Numbers", "https://sech.me/boinc/AMicable/"},
	{"NumberFields@home", "https://numberfields.asu.edu/NumberFields/"},
	{"Cosmology@Home", "https://www.cosmologyathome.org/"},
	{"Asteroids@home", "http://asteroidsathome.net/boinc/"},
	{"Enigma@Home", "https://enigmaathome.net/"},
	{"yoyo@home", "https://www.rechenkraft.net/yoyo/"},
	{"Moo! Wrapper", "https://moowrap.net/"},
	{"SRBase", "https://srbase.my-firefox.org/boinc/"},
}

func catalogNames() []string {
	out := make([]string, len(projCatalog))
	for i, p := range projCatalog {
		out[i] = p.Name
	}
	return out
}

func catalogURL(name string) string {
	for _, p := range projCatalog {
		if p.Name == name {
			return p.URL
		}
	}
	return ""
}

func showAttachDialog(w fyne.Window) {
	if pv.cur == "" {
		dialog.NewInformation("Attach to project", "Choose a server first.", w).Show()
		return
	}
	hostID := ""
	for _, cfg := range mgr.Store.List() {
		if cfg.Name == pv.cur {
			hostID = cfg.ID
			break
		}
	}
	if hostID == "" {
		dialog.NewInformation("Attach to project", "Server not found.", w).Show()
		return
	}

	catalogSel := widget.NewSelect(catalogNames(), nil)
	urlE := widget.NewEntry()
	urlE.PlaceHolder = "https://project.example.edu/"
	catalogSel.OnChanged = func(name string) {
		if u := catalogURL(name); u != "" && urlE.Text == "" {
			urlE.SetText(u)
		} else if u := catalogURL(name); u != "" {
			urlE.SetText(u)
		}
	}

	authE := widget.NewPasswordEntry()
	authE.PlaceHolder = "32-character authenticator"
	emailE := widget.NewEntry()
	emailE.PlaceHolder = "you@example.com"
	passE := widget.NewPasswordEntry()
	passE.PlaceHolder = "Project account password"

	lookupBtn := widget.NewButton("Look up authenticator with account", nil)
	statusLbl := widget.NewLabel("")

	lookupBtn.OnTapped = func() {
		base := strings.TrimSpace(urlE.Text)
		if base == "" || emailE.Text == "" || passE.Text == "" {
			statusLbl.SetText("Fill URL, e-mail and password first.")
			return
		}
		lookupBtn.Disable()
		statusLbl.SetText("Contacting project…")
		go func() {
			auth, err := mgr.LookupAccount(base, strings.TrimSpace(emailE.Text), passE.Text)
			fyreDo(func() {
				lookupBtn.Enable()
				if err != nil {
					statusLbl.SetText("Lookup failed: " + err.Error())
					return
				}
				authE.SetText(auth)
				statusLbl.SetText("Authenticator retrieved.")
			})
		}()
	}

	items := []*widget.FormItem{
		widget.NewFormItem("Project", catalogSel),
		widget.NewFormItem("Project URL", urlE),
		widget.NewFormItem("", widget.NewLabel("Option A - paste your weak account key:")),
		widget.NewFormItem("Authenticator", authE),
		widget.NewFormItem("", widget.NewLabel("Option B - fetch it using your project account:")),
		widget.NewFormItem("E-mail", emailE),
		widget.NewFormItem("Password", passE),
		widget.NewFormItem("", container.NewVBox(lookupBtn, statusLbl)),
	}

	dlg := dialog.NewForm("Attach "+pv.cur+" to a project", "Attach", "Cancel", items, func(ok bool) {
		if !ok {
			return
		}
		u := strings.TrimSpace(urlE.Text)
		a := strings.TrimSpace(authE.Text)
		if u == "" || a == "" {
			dialog.NewInformation("Attach to project", "URL and authenticator are required.", w).Show()
			return
		}
		name := catalogSel.Selected
		go func() {
			err := mgr.Attach(hostID, u, a, name)
			if err != nil {
				fyreDo(func() {
					dialog.NewError(fmt.Errorf("attach failed: %v", err), w).Show()
				})
				return
			}
			time.Sleep(600 * time.Millisecond)
			fyreDo(fireRefresh)
		}()
	}, w)
	dlg.Resize(fyne.NewSize(560, 620))
	dlg.Show()
}
