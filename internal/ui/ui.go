package ui

import (
	"sync"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/alplix/lavandegrid/internal/app"
	"github.com/alplix/lavandegrid/internal/local"
)

const Version = "1.0.0"

var (
	mgr     *app.Manager
	mainWin fyne.Window
	fyneApp fyne.App

	screenMu   sync.Mutex
	curRefresh func()
)

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
	label   string
	icon    fyne.Resource
	build   func(fyne.Window) fyne.CanvasObject
	refresh func()
}

var navItems []navItem

func registerScreen(label string, icon fyne.Resource, build func(fyne.Window) fyne.CanvasObject, refresh func()) {
	navItems = append(navItems, navItem{label: label, icon: icon, build: build, refresh: refresh})
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

func Run(a fyne.App) {
	fyneApp = a
	light := a.Preferences().Bool("light")
	ApplyTheme(a, light)
	notifsOn := a.Preferences().Bool("notifications")

	mgr = app.NewManager()
	mgr.OnNotice(func(n app.Notice) {
		if notifsOn {
			a.SendNotification(&fyne.Notification{Title: n.Title, Content: n.Body})
		}
	})
	mgr.Start()

	w := a.NewWindow("LavandeGrid v" + Version)
	w.Resize(fyne.NewSize(1200, 780))
	mainWin = w

	nav := container.NewVBox()
	content := container.NewStack()

	for _, it := range navItems {
		item := it
		btn := widget.NewButtonWithIcon(item.label, item.icon, nil)
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

	root := container.NewBorder(nil, nil, container.NewVScroll(nav), nil, content)
	w.SetContent(root)

	go func() {
		t := time.NewTicker(2 * time.Second)
		defer t.Stop()
		for range t.C {
			fireRefresh()
		}
	}()

	uiQueue = make(chan func(), 64)
	go func() {
		for fn := range uiQueue {
			func() {
				defer func() { recover() }()
				fn()
			}()
		}
	}()

	if len(nav.Objects) > 0 {
		if b, ok := nav.Objects[0].(*widget.Button); ok {
			b.OnTapped()
		}
	}

	w.ShowAndRun()
}

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
			return "Bundled BOINC client found:\n" + info.Exe
		}
		return "System BOINC client found:\n" + info.Exe
	}
	return "No local BOINC client detected.\n" + info.Hint
}
