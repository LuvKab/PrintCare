# PrintCare

<p align="center">
  <strong>Keep your printhead healthy with scheduled maintenance printing.</strong>
</p>

<p align="center">
  <a href="#features">Features</a> &bull;
  <a href="#download">Download</a> &bull;
  <a href="#screenshots">Screenshots</a> &bull;
  <a href="#usage">Usage</a> &bull;
  <a href="#build">Build</a> &bull;
  <a href="#license">License</a>
</p>

---

Inkjet printers suffer from clogged nozzles when left idle for days. **PrintCare** solves this by automatically printing a maintenance image at configurable intervals — keeping your printhead clean and ready.

## Features

- **Scheduled Printing** — Automatically prints at user-defined intervals (e.g., every 3 days)
- **Multi-format Support** — Print BMP, PNG, and JPEG images
- **Modern UI** — Clean, responsive interface built with Svelte + Tailwind CSS
- **Multi-language** — English, 简体中文, 繁體中文, 日本語
- **System Tray** — Runs quietly in the background with tray icon and right-click menu
- **First-run Wizard** — Guided setup on first launch for a smooth onboarding experience
- **Dark Mode** — Full dark theme support
- **Portable** — Single `.exe`, no installation required. Configuration stored in `%APPDATA%`

## Download

Go to [Releases](../../releases) and download the latest `PrintCare.exe`.

**Requirements:** Windows 10/11 with WebView2 runtime (included in most Windows 10/11 installations).

## Usage

1. Run `PrintCare.exe`
2. On first launch, a setup wizard will guide you through:
   - Selecting your language
   - Choosing a printer
   - Selecting a maintenance image (a nozzle-check pattern works best)
   - Setting the print interval
3. The app minimizes to the system tray and prints on schedule
4. Right-click the tray icon to show the window or exit

### Configuration

Settings are stored in:
```
%APPDATA%\PrintheadMaintainer\config.json
```

Logs are stored in:
```
%APPDATA%\PrintheadMaintainer\app.log
```

## Screenshots

> _Screenshots coming soon._

## Tech Stack

| Component | Technology |
|-----------|-----------|
| Backend | Go + [Wails v2](https://wails.io/) |
| Frontend | Svelte + TypeScript + Tailwind CSS |
| Printing | Windows GDI API via `syscall` |
| System Tray | [energye/systray](https://github.com/energye/systray) |
| Config | JSON file |
| i18n | Embedded JSON locale files |

## Build

### Prerequisites

- Go 1.21+
- Node.js 16+
- [Wails CLI v2](https://wails.io/docs/gettingstarted/installation)

### Steps

```bash
# Install Wails CLI
go install github.com/wailsapp/wails/v2/cmd/wails@latest

# Build
cd printhead-go
wails build
```

The output binary will be at `build/bin/PrintCare.exe`.

## Project Structure

```
├── main.go              # Wails entry point
├── app.go               # Core app struct + Wails bindings
├── config.go            # JSON config read/write
├── scheduler.go         # Scheduled print timer
├── printer_windows.go   # Windows GDI print API
├── tray.go              # System tray integration
├── i18n.go              # Multi-language manager
├── logger.go            # File-based logging
├── locales/             # Language JSON files
│   ├── en.json
│   ├── zh-CN.json
│   ├── zh-TW.json
│   └── ja.json
└── frontend/            # Svelte + Tailwind UI
    └── src/
        ├── pages/       # Dashboard, Settings, PrintNow, Logs
        └── components/  # Sidebar, SetupWizard, StatusCard
```

## License

MIT
