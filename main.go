package main

import (
	"context"
	"embed"
	"os"
	"runtime"

	"github.com/getlantern/systray"
	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	"github.com/wailsapp/wails/v2/pkg/options/mac"
	"github.com/wailsapp/wails/v2/pkg/options/windows"
)

//go:embed all:frontend/dist
var assets embed.FS

//go:embed build/icon.png
var iconPNG []byte

var wailsApp *App

func main() {
	runtime.LockOSThread()
	systray.Run(onTrayReady, onTrayExit)
}

func onTrayReady() {
	systray.SetIcon(iconPNG)
	systray.SetTitle("LavandeGrid")
	systray.SetTooltip("LavandeGrid - BOINC Manager")

	mShow := systray.AddMenuItem("Open LavandeGrid", "Show main window")
	mRefresh := systray.AddMenuItem("Refresh all servers", "Poll every configured server now")
	mHide := systray.AddMenuItem("Hide to tray", "Hide the main window")
	systray.AddSeparator()
	mQuit := systray.AddMenuItem("Quit", "Quit LavandeGrid")

	mShow.Disable()
	mHide.Disable()

	go func() {
		for {
			select {
			case <-mShow.ClickedCh:
				if wailsApp != nil && wailsApp.ctx != nil {
					wailsApp.showWindow()
				}
			case <-mHide.ClickedCh:
				if wailsApp != nil && wailsApp.ctx != nil {
					wailsApp.hideWindow()
				}
			case <-mRefresh.ClickedCh:
				if wailsApp != nil {
					wailsApp.RefreshAll()
				}
			case <-mQuit.ClickedCh:
				systray.Quit()
				os.Exit(0)
			}
		}
	}()

	go startWails(mShow, mHide)
}

func onTrayExit() {
}

func startWails(mShow, mHide *systray.MenuItem) {
	wailsApp = NewApp()

	err := wails.Run(&options.App{
		Title:     "LavandeGrid",
		Width:     1280,
		Height:    800,
		MinWidth:  960,
		MinHeight: 640,
		BackgroundColour: &options.RGBA{
			R: 17, G: 16, B: 26, A: 255,
		},
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		SingleInstanceLock: &options.SingleInstanceLock{
			UniqueId: "dev.alplix.lavandegrid",
			OnSecondInstanceLaunch: func(_ options.SecondInstanceData) {
				if wailsApp != nil && wailsApp.ctx != nil {
					wailsApp.showWindow()
				}
			},
		},
		OnStartup:  wailsApp.startup,
		OnShutdown: wailsApp.shutdown,
		OnBeforeClose: func(ctx context.Context) bool {
			if wailsApp != nil {
				wailsApp.hideWindow()
			}
			return true
		},
		Windows: &windows.Options{
			WebviewIsTransparent: true,
			WindowIsTranslucent:  true,
			BackdropType:         windows.Acrylic,
			Theme:                windows.SystemDefault,
		},
		Mac: &mac.Options{
			About: &mac.AboutInfo{
				Title:   "LavandeGrid " + appVersion,
				Message: "A lavender-themed desktop manager for your BOINC fleet.\n\n© 2026 Alperen Yavuz\n" + appRepo,
				Icon:    iconPNG,
			},
		},
		Bind: []interface{}{
			wailsApp,
		},
	})
	if err != nil {
		println("Error:", err.Error())
		os.Exit(1)
	}
	mShow.Enable()
	mHide.Enable()
}