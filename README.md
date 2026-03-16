# PrintCare — 打印养护

<p align="center">
  <strong>定时打印，保持喷头畅通，防止堵塞。</strong>
</p>

<p align="center">
  <a href="#功能特性">功能特性</a> &bull;
  <a href="#下载">下载</a> &bull;
  <a href="#使用方法">使用方法</a> &bull;
  <a href="#构建">构建</a> &bull;
  <a href="#english">English</a>
</p>

---

喷墨打印机长时间不使用，喷头容易堵塞干涸。**PrintCare** 通过定时自动打印维护图片，让喷头始终保持通畅。

## 功能特性

- **定时打印** — 按自定义间隔自动打印（如每 3 天一次）
- **多格式支持** — 支持 BMP、PNG、JPEG 图片
- **现代化界面** — 基于 Svelte + Tailwind CSS 构建的简洁 UI
- **多语言** — 简体中文、繁體中文、English、日本語
- **系统托盘** — 后台静默运行，右键菜单支持多语言
- **首次引导** — 首次启动自动弹出设置向导，快速上手
- **深色模式** — 完整的暗色主题支持
- **绿色便携** — 单文件 `.exe`，无需安装，配置存储在 `%APPDATA%`

## 下载

前往 [Releases](../../releases) 下载最新版 `PrintCare.exe`。

**系统要求：** Windows 10/11（需要 WebView2 运行时，大多数 Win10/11 已自带）。

## 使用方法

1. 运行 `PrintCare.exe`
2. 首次启动会弹出设置向导：
   - 选择语言
   - 选择打印机
   - 选择维护打印图片（推荐使用喷嘴检测图案）
   - 设置打印间隔天数
3. 设置完成后程序最小化到系统托盘，按计划自动打印
4. 右键托盘图标可以显示主窗口或退出程序

### 配置文件位置

```
%APPDATA%\PrintheadMaintainer\config.json
```

### 日志文件位置

```
%APPDATA%\PrintheadMaintainer\app.log
```

## 技术栈

| 组件 | 技术 |
|------|------|
| 后端 | Go + [Wails v2](https://wails.io/) |
| 前端 | Svelte + TypeScript + Tailwind CSS |
| 打印 | Windows GDI API（`syscall` 直接调用） |
| 系统托盘 | [energye/systray](https://github.com/energye/systray) |
| 配置 | JSON 文件 |
| 多语言 | 内嵌 JSON 语言资源文件 |

## 构建

### 前置条件

- Go 1.21+
- Node.js 16+
- [Wails CLI v2](https://wails.io/docs/gettingstarted/installation)

### 步骤

```bash
go install github.com/wailsapp/wails/v2/cmd/wails@latest

wails build
```

编译产物位于 `build/bin/PrintCare.exe`。

## 项目结构

```
├── main.go              # Wails 入口
├── app.go               # 核心 App 结构体 + Wails 绑定
├── config.go            # JSON 配置读写
├── scheduler.go         # 定时打印调度器
├── printer_windows.go   # Windows GDI 打印封装
├── tray.go              # 系统托盘
├── i18n.go              # 多语言管理
├── logger.go            # 文件日志
├── locales/             # 语言资源
│   ├── en.json
│   ├── zh-CN.json
│   ├── zh-TW.json
│   └── ja.json
└── frontend/            # Svelte + Tailwind 前端
    └── src/
        ├── pages/       # 仪表盘、设置、立即打印、日志
        └── components/  # 侧边栏、设置向导、状态卡片
```

---

## English

**PrintCare** keeps your inkjet printhead healthy by automatically printing a maintenance image at configurable intervals.

**Features:** Scheduled printing, BMP/PNG/JPEG support, modern UI (Svelte + Tailwind), multi-language (en/zh-CN/zh-TW/ja), system tray, first-run wizard, dark mode, single portable EXE.

**Usage:** Run `PrintCare.exe` → complete the setup wizard → the app runs in the system tray and prints on schedule.

## License

MIT
