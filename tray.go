package main

import (
	_ "embed"

	"github.com/energye/systray"
	wailsRuntime "github.com/wailsapp/wails/v2/pkg/runtime"
)

//go:embed build/windows/icon.ico
var trayIcon []byte

func (a *App) initTray() {
	systray.Run(a.onTrayReady, func() {})
}

func (a *App) onTrayReady() {
	systray.SetIcon(trayIcon)

	a.trayShow = systray.AddMenuItem("", "")
	systray.AddSeparator()
	a.trayQuit = systray.AddMenuItem("", "")

	a.updateTrayLocale()

	a.trayShow.Click(func() {
		if a.ctx != nil {
			wailsRuntime.WindowShow(a.ctx)
		}
	})

	a.trayQuit.Click(func() {
		a.quitting = true
		if a.ctx != nil {
			wailsRuntime.Quit(a.ctx)
		}
	})
}

func (a *App) updateTrayLocale() {
	locale := a.i18n.GetLocale(a.config.Language)
	if a.trayShow != nil {
		if s, ok := locale["tray.show"].(string); ok {
			a.trayShow.SetTitle(s)
		}
	}
	if a.trayQuit != nil {
		if s, ok := locale["tray.exit"].(string); ok {
			a.trayQuit.SetTitle(s)
		}
	}
	if name, ok := locale["app.title"].(string); ok {
		systray.SetTooltip(name)
	}
}
