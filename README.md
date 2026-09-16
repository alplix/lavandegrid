# LavandeGrid

A lavender-themed desktop manager for your **BOINC / Camellia fleet** — one window, every machine.

LavandeGrid connects to BOINC clients over GUI RPC (port 31416 by default), monitors them live, and
lets you control tasks, projects, transfers and preferences across all of them from a single native
desktop app. A built-in **demo mode** lets you explore the whole interface with zero real clients.

Built with **Go + [Wails v2](https://wails.io)**. Native WebView on every platform — no browser, no Electron.

## Features

| Area | What you can do |
|---|---|
| **Fleet dashboard** | Live totals (running / paused / queued), RAC & credit, per-host activity sparkline, recent client messages |
| **Tasks** | Progress, elapsed/CPU time, ETA, deadlines, GPU/CPU resources; pause / resume / abort per task; status filters + search |
| **Projects** | Attach (with email/password authenticator lookup), detach, suspend/resume, update, no-more-work / allow-work |
| **Transfers** | Live upload/download progress with retry & abort |
| **Messages** | Severity-highlighted client log across all servers |
| **Statistics** | Per-project credit-history charts (365 days), 30-day transfer history, per-project disk usage |
| **Preferences** | Remote global-pref overrides (key=value editor), per-host & **fleet-wide** run/network modes, benchmarks |
| **Hardware** | OS, CPU, cores, FLOPS, RAM, disk, GPU list with VRAM per host |
| **Notifications** | Desktop alerts for deadlines, task errors and offline hosts; tray menu with refresh/hide/show |
| **Local client** | Auto-detect a local Camellia/BOINC client, start/stop it from Settings |

## Platforms

| OS | Architectures | Packaging | Notes |
|---|---|---|---|
| Windows 10/11 | amd64 | NSIS installer + portable zip | WebView2 bundled; optional BOINC installer in release |
| macOS | amd64, arm64 | `.app` bundle (zip) | Optional BOINC `.dmg` in release |
| Linux | amd64, arm64 | portable `.tar.gz` | GTK3 + WebKitGTK 4.1 required; uses system `boinc` client |

The app detects a local BOINC/Camellia client automatically (bundled copy first, then the system
installation) and registers it as the **Local Camellia** server.

## Getting started

1. On each machine running BOINC, find its RPC password in `gui_rpc_auth.cfg` inside the BOINC data directory.
2. Allow GUI RPC — start the client with `-allow_remote_gui_rpc`, or configure `remote_hosts` for remote management.
3. In LavandeGrid open **Servers → Add Server** and enter host, port (default `31416`) and password.
4. First launch without any servers creates a **Demo Server** with simulated data so you can look around safely.

Everything is stored in the OS config directory (`~/.config/lavandegrid` on Linux, `%APPDATA%\lavandegrid`
on Windows, `~/Library/Application Support/lavandegrid` on macOS). Passwords are saved in `hosts.json`
(0600) with the rest of the configuration.

## Building from source

Prerequisites:

- Go 1.26+
- [Wails CLI](https://wails.io/docs/gettingstarted/installation) (`go install github.com/wailsapp/wails/v2/cmd/wails@latest`)
- Node.js 20+ (frontend tooling)

On Linux you also need the [Wails WebKit/GTK dependencies](https://wails.io/docs/guides/linux)
(`libgtk-3-dev libwebkit2gtk-4.1-dev libayatana-appindicator3-dev`). On Windows, [WebView2](https://developer.microsoft.com/microsoft-edge/webview2/) is required at runtime (bundled by the installer).

```
wails dev     # live development build
wails build   # production build for the current platform
```

Cross-compiling for release targets is handled by the CI pipelines:

- `.github/workflows/ci.yml` — vet, tests, frontend build and a Linux smoke build on every PR/push
- `.github/workflows/release.yml` — tagged releases (`vX.Y.Z`) for Windows, macOS and Linux with
  checksums; trigger with `git tag v1.1.0 && git push origin v1.1.0`

### Manual cross-build example (Linux host)

```bash
# Windows amd64 (+ NSIS installer)
wails build -platform windows/amd64 -nsis -webview2 download

# Linux arm64
wails build -platform linux/arm64
```

## Architecture

```
frontend/                 Vite + vanilla JS/CSS single-page UI
  src/main.js             state, rendering, Wails bindings, notifications
  src/style.css           lavender glassmorphism theme (light / dark)
  wailsjs/                generated TS bindings to the Go API

internal/app/             pure-Go application logic (unit-tested)
  types.go                snapshot model + JSON contract (camelCase)
  mock.go                 simulated fleet for demo mode
  manager.go              polling loop, host operations, notices, stats
  store.go                atomic-persisted host configuration
internal/boinc/           BOINC GUI-RPC client (xmlrpc-over-TCP)
internal/local/           local daemon detection & lifecycle management
internal/i18n/            embedded UI translations (foundation)
app.go                    Wails-bound API surface (main namespace)
main.go                   tray + window setup, single-instance lock
```

The Go API is bound under `window.go.main.App` and talks to the UI through plain method calls plus
a single event channel (`notice`) for deadline/error/offline notifications.

## Release checklist

1. Bump `appVersion` in `app.go`, `productVersion` in `wails.json`, and the README version tag.
2. `git tag v1.1.0 && git push origin v1.1.0`
3. Download the release assets; each job uploads `SHA256SUMS.txt` alongside the archives.

## Roadmap

- Frontend i18n wiring (translations are already embedded as Go strings)
- Windows arm64 build target
- .deb / .rpm / AppImage packaging for Linux
- Per-host task scheduling & project share editing

## License

MIT — see [LICENSE](LICENSE). BOINC is LGPL; release builds optionally bundle official BOINC
installers downloaded at build time.