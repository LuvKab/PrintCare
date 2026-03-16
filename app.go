package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"time"

	"github.com/energye/systray"
	wailsRuntime "github.com/wailsapp/wails/v2/pkg/runtime"
)

type Status struct {
	Enabled     bool   `json:"enabled"`
	PrinterOK   bool   `json:"printerReady"`
	LastPrint   string `json:"lastPrint"`
	NextPrint   string `json:"nextPrint"`
	StatusText  string `json:"statusText"`
	StatusLevel string `json:"statusLevel"`
}

type ConfigInput struct {
	Enabled      bool   `json:"enabled"`
	PrinterName  string `json:"printerName"`
	PaperSource  int    `json:"paperSource"`
	IntervalDays int    `json:"intervalDays"`
	ImagePath    string `json:"imagePath"`
	Language     string `json:"language"`
}

type App struct {
	ctx      context.Context
	config   *Config
	scheduler *Scheduler
	logger   *Logger
	i18n     *I18nManager
	quitting bool
	trayShow *systray.MenuItem
	trayQuit *systray.MenuItem
}

func NewApp() *App {
	logger := NewLogger()
	cfg := LoadConfig(logger)
	i18n := NewI18nManager()

	app := &App{
		config: cfg,
		logger: logger,
		i18n:   i18n,
	}
	app.scheduler = NewScheduler(app)
	return app
}

func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
	a.logger.Info("Application started")
	a.scheduler.Start()
	go a.initTray()
}

func (a *App) shutdown(ctx context.Context) {
	a.scheduler.Stop()
	a.logger.Info("Application shutdown")
}

func (a *App) beforeClose(ctx context.Context) (prevent bool) {
	if a.quitting {
		return false
	}
	wailsRuntime.WindowHide(ctx)
	return true
}

func (a *App) IsFirstRun() bool {
	return a.config.PrinterName == ""
}

func (a *App) GetStatus() Status {
	s := Status{Enabled: a.config.Enabled}

	printers, _ := EnumPrinters()
	for _, p := range printers {
		if p == a.config.PrinterName {
			s.PrinterOK = true
			break
		}
	}

	if !a.config.LastPrint.IsZero() {
		s.LastPrint = a.config.LastPrint.Format("2006-01-02 15:04:05")
	}

	if a.config.Enabled && a.config.IntervalDays > 0 && !a.config.LastPrint.IsZero() {
		next := a.config.LastPrint.Add(time.Duration(a.config.IntervalDays) * 24 * time.Hour)
		s.NextPrint = next.Format("2006-01-02 15:04:05")
	}

	switch {
	case !a.config.Enabled:
		s.StatusText, s.StatusLevel = "disabled", "warning"
	case a.config.PrinterName == "":
		s.StatusText, s.StatusLevel = "no_printer", "error"
	case a.config.ImagePath == "":
		s.StatusText, s.StatusLevel = "no_image", "error"
	case !s.PrinterOK:
		s.StatusText, s.StatusLevel = "printer_offline", "error"
	default:
		s.StatusText, s.StatusLevel = "ready", "ok"
	}

	return s
}

func (a *App) GetConfig() *Config {
	return a.config
}

func (a *App) SaveConfig(input ConfigInput) error {
	oldLang := a.config.Language

	a.config.Enabled = input.Enabled
	a.config.PrinterName = input.PrinterName
	a.config.PaperSource = input.PaperSource
	a.config.IntervalDays = input.IntervalDays
	a.config.ImagePath = input.ImagePath
	a.config.Language = input.Language

	if err := saveConfig(a.config, a.logger); err != nil {
		return err
	}
	a.scheduler.Reset()
	a.logger.Info("Configuration saved")

	if input.Language != oldLang {
		a.updateTrayLocale()
		if title, ok := a.i18n.GetLocale(input.Language)["app.title"].(string); ok && a.ctx != nil {
			wailsRuntime.WindowSetTitle(a.ctx, title)
		}
	}

	return nil
}

func (a *App) PrintNow() error {
	if a.config.PrinterName == "" {
		return fmt.Errorf("no printer configured")
	}
	if a.config.ImagePath == "" {
		return fmt.Errorf("no image configured")
	}

	a.logger.Info("Manual print triggered")
	if err := PrintImage(a.config.PrinterName, a.config.ImagePath, a.config.PaperSource); err != nil {
		a.logger.Error("Print failed: %v", err)
		return err
	}

	a.config.LastPrint = time.Now()
	saveConfig(a.config, a.logger)
	a.logger.Info("Print completed successfully")
	return nil
}

func (a *App) GetPrinters() []string {
	printers, err := EnumPrinters()
	if err != nil {
		a.logger.Error("Failed to enumerate printers: %v", err)
		return []string{}
	}
	return printers
}

func (a *App) SelectImage() (string, error) {
	return wailsRuntime.OpenFileDialog(a.ctx, wailsRuntime.OpenDialogOptions{
		Title: "Select Image",
		Filters: []wailsRuntime.FileFilter{
			{DisplayName: "Images (*.bmp;*.png;*.jpg;*.jpeg)", Pattern: "*.bmp;*.png;*.jpg;*.jpeg"},
			{DisplayName: "All Files (*.*)", Pattern: "*.*"},
		},
	})
}

func (a *App) GetLogs() string {
	return a.logger.ReadAll()
}

func (a *App) ClearLogs() {
	a.logger.Clear()
}

func (a *App) GetLocale(lang string) map[string]interface{} {
	return a.i18n.GetLocale(lang)
}

func (a *App) GetAvailableLocales() []LocaleInfo {
	return a.i18n.GetAvailableLocales()
}

func (a *App) GetImagePreview() string {
	if a.config.ImagePath == "" {
		return ""
	}
	if _, err := os.Stat(a.config.ImagePath); os.IsNotExist(err) {
		return ""
	}
	return a.config.ImagePath
}

func getAppDataDir() string {
	if runtime.GOOS == "windows" {
		dir := filepath.Join(os.Getenv("APPDATA"), "PrintheadMaintainer")
		os.MkdirAll(dir, 0755)
		return dir
	}
	home, _ := os.UserHomeDir()
	dir := filepath.Join(home, ".printhead-maintainer")
	os.MkdirAll(dir, 0755)
	return dir
}
