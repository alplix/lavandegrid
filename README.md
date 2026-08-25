# LavandeGrid

A lavender-themed desktop manager for your BOINC fleet - one window, every machine.

LavandeGrid connects to BOINC clients over GUI RPC and lets you monitor and control
all of them from a single native desktop app. A built-in demo mode lets you explore
the full interface without any real clients.

Built with Go + [Fyne](https://fyne.io). No browser, no Electron.

## Features

- **Multi-server fleet** - add any number of hosts and watch them all live
- **Tasks** - progress, CPU time, ETA, deadlines, GPU/CPU resources, bulk ops
- **Projects** - attach/detach, suspend/resume, update, reset, authenticator lookup
- **Transfers** - live downloads/uploads with retry/abort
- **Messages** - severity-filtered client log across all servers
- **Statistics** - credit history charts per project (30/90/180/365 days)
- **Preferences** - remote global-pref overrides, run/network modes, benchmarks
- **Notifications** - desktop alerts for task errors and offline hosts

## Platforms

| OS | Architectures | BOINC client |
|---|---|---|
| Windows | amd64 | bundled installer in release |
| macOS | amd64, arm64 | bundled installer in release |
| Linux | amd64, arm64 | use system client (`boinc` / `boinc_client`) |
| FreeBSD | amd64 | use system client |

The app detects a local BOINC client automatically (bundled copy first, then the
system installation) and can start it from Settings.

## Getting started

1. On each machine running BOINC, find its RPC password in
   `gui_rpc_auth.cfg` inside the BOINC data directory.
2. Make sure GUI RPC is allowed (BOINC client option `-allow_remote_gui_rpc`
   or `remote_hosts` configured for remote access).
3. In LavandeGrid open **Servers -> Add server** and enter host, port
   (default 31416) and password.
4. First launch without any servers starts a **Demo Server** with simulated
   data so you can look around safely.

## Building from source

Requires Go 1.26+, a C compiler and OpenGL/X11 headers (see
[Fyne docs](https://docs.fyne.io)). Then:

```
go build -o lavandegrid .
```

Cross-compilation for all release targets is handled by
`.github/workflows/release.yml` using [fyne-cross](https://github.com/fyne-io/fyne-cross).

## License

MIT - see [LICENSE](LICENSE). BOINC is LGPL; releases bundle official client
installers downloaded at build time.
