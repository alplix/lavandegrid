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
	mQuit := systray.AddMenuItem("Quit", "Quit LavandeGrid")

	go func() {
		for {
			select {
			case <-mShow.ClickedCh:
				if wailsApp != nil && wailsApp.ctx != nil {
					wailsApp.showWindow()
				}
			case <-mQuit.ClickedCh:
				systray.Quit()
				os.Exit(0)
			}
		}
	}()

	go startWails()
}

func onTrayExit() {
}

func startWails() {
	wailsApp = NewApp()

	err := wails.Run(&options.App{
		Title:     "LavandeGrid",
		Width:     1280,
		Height:    800,
		MinWidth:  960,
		MinHeight: 640,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		OnStartup:  wailsApp.startup,
		OnShutdown: wailsApp.shutdown,
		OnBeforeClose: func(ctx context.Context) bool {
			if wailsApp != nil {
				wailsApp.hideWindow()
			}
			return true
		},
		Bind: []interface{}{
			wailsApp,
		},
	})
	if err != nil {
		println("Error:", err.Error())
		os.Exit(1)
	}
}
