package ui

import (
	"image/color"
	"sync"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/alplix/lavandegrid/internal/app"
	"github.com/alplix/lavandegrid/internal/i18n"
	"github.com/alplix/lavandegrid/internal/local"
)

const Version = "1.0.0"

var (
	mgr        *app.Manager
	mainWin    fyne.Window
	fyneApp    fyne.App
	wantRestart bool

	screenMu   sync.Mutex
	curRefresh func()

	screenDefs []screenDef
)

type screenDef struct {
	i18nKey string
	icon    fyne.Resource
	build   func(fyne.Window) fyne.CanvasObject
	refresh func()
}

func setRefresh(fn func()) {
	screenMu.Lock()
	curRefresh = fn
	screenMu.Unlock()
}

func fireRefresh() {
	screenMu.Lock()
	fn := curRefresh
	screenMu.Unlock()
	if fn != nil {
		func() {
			defer func() { recover() }()
			fn()
		}()
	}
}

type navItem struct {
	i18nKey string
	icon    fyne.Resource
	build   func(fyne.Window) fyne.CanvasObject
	refresh func()
}

func (n navItem) label() string { return i18n.T(n.i18nKey) }

var navItems []navItem

func registerScreen(i18nKey string, icon fyne.Resource, build func(fyne.Window) fyne.CanvasObject, refresh func()) {
	screenDefs = append(screenDefs, screenDef{i18nKey: i18nKey, icon: icon, build: build, refresh: refresh})
}

func rebuildNav() {
	navItems = navItems[:0]
	for _, d := range screenDefs {
		navItems = append(navItems, navItem{i18nKey: d.i18nKey, icon: d.icon, build: d.build, refresh: d.refresh})
	}
}

var uiQueue chan func()

func fyreDo(fn func()) {
	if uiQueue != nil {
		select {
		case uiQueue <- fn:
		default:
			go fn()
		}
	} else {
		go fn()
	}
}

func buildTopBar() fyne.CanvasObject {
	logoLabel := canvas.NewText("LavandeGrid", color.NRGBA{R: 0xc4, G: 0x95, B: 0xf0, A: 0xff})
	logoLabel.TextSize = 20
	logoLabel.TextStyle = fyne.TextStyle{Bold: true}

	verLabel := canvas.NewText("v"+Version, color.NRGBA{R: 0x86, G: 0x78, B: 0xa3, A: 0xff})
	verLabel.TextSize = 12

	flower1 := canvas.NewText("~", color.NRGBA{R: 0xc4, G: 0x95, B: 0xf0, A: 0x60})
	flower1.TextSize = 24

	flower2 := canvas.NewText("~", color.NRGBA{R: 0xa7, G: 0x8b, B: 0xfa, A: 0x40})
	flower2.TextSize = 18

	subtitle := canvas.NewText(i18n.T("page.dash.sub"), color.NRGBA{R: 0x86, G: 0x78, B: 0xa3, A: 0xff})
	subtitle.TextSize = 11

	left := container.NewHBox(flower1, logoLabel, verLabel)
	right := container.NewHBox(subtitle, flower2)

	return container.NewBorder(nil, nil, left, right)
}

func Run(a fyne.App) {
	fyneApp = a
	prefs := a.Preferences()
	i18n.RegisterPrefs(func() string { return prefs.String("lang") }, func(s string) { prefs.SetString("lang", s) })
	i18n.Init()
	light := prefs.Bool("light")
	ApplyTheme(a, light)
	notifsOn := prefs.Bool("notifications")

	if mgr == nil {
		mgr = app.NewManager()
		mgr.OnNotice(func(n app.Notice) {
			if notifsOn {
				a.SendNotification(&fyne.Notification{Title: n.Title, Content: n.Body})
			}
		})
		mgr.Start()
	}

	rebuildNav()

	w := a.NewWindow(i18n.T("app.name") + " v" + Version)
	w.Resize(fyne.NewSize(1200, 780))
	mainWin = w

	nav := container.NewVBox()
	content := container.NewStack()

	for _, it := range navItems {
		item := it
		btn := widget.NewButtonWithIcon(item.label(), item.icon, nil)
		btn.Alignment = widget.ButtonAlignLeading
		btn.OnTapped = func() {
			setRefresh(item.refresh)
			content.Objects = []fyne.CanvasObject{item.build(w)}
			content.Refresh()
			highlightNav(nav, btn)
			fireRefresh()
		}
		nav.Add(btn)
	}

	petalBar := canvas.NewRectangle(color.NRGBA{R: 0x2b, G: 0x20, B: 0x40, A: 0x30})
	petalBar.Resize(fyne.NewSize(200, 2))

	navWrap := container.NewBorder(nil, nil, nil, nil,
		container.NewVBox(
			container.NewPadded(buildTopBar()),
			petalBar,
			nav,
			layout.NewSpacer(),
			widget.NewLabel(""),
		),
	)
	navScroll := container.NewVScroll(navWrap)

	root := container.NewBorder(nil, nil, navScroll, nil, content)
	w.SetContent(root)

	go func() {
		t := time.NewTicker(2 * time.Second)
		defer t.Stop()
		for range t.C {
			fireRefresh()
		}
	}()

	if uiQueue == nil {
		uiQueue = make(chan func(), 64)
		go func() {
			for fn := range uiQueue {
				func() {
					defer func() { recover() }()
					fn()
				}()
			}
		}()
	}

	if len(nav.Objects) > 0 {
		if b, ok := nav.Objects[0].(*widget.Button); ok {
			b.OnTapped()
		}
	}

	w.ShowAndRun()
}

func RestartWindow() {
	wantRestart = true
	if mainWin != nil {
		w := mainWin
		mainWin = nil
		w.Close()
	}
}

func NeedsRestart() bool { return wantRestart }

func highlightNav(nav *fyne.Container, active *widget.Button) {
	for _, obj := range nav.Objects {
		if b, ok := obj.(*widget.Button); ok {
			if b == active {
				b.Importance = widget.HighImportance
			} else {
				b.Importance = widget.MediumImportance
			}
			b.Refresh()
		}
	}
}

func localInfoText() string {
	info := local.Detect()
	if info.Found {
		if info.Bundled {
			return i18n.T("local.bundled") + ":\n" + info.Exe
		}
		return i18n.T("local.system") + ":\n" + info.Exe
	}
	return i18n.T("local.none") + "\n" + info.Hint
}
