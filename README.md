# barcomic

![Build Status](https://img.shields.io/github/actions/workflow/status/thegraydot/barcomic_server/tag.yml?style=flat)
![Test Status](https://img.shields.io/github/actions/workflow/status/thegraydot/barcomic_server/tag.yml?style=flat&label=test)

![Release Version](https://img.shields.io/github/v/release/thegraydot/barcomic_server?style=flat)
![Release downloads](https://img.shields.io/github/downloads/thegraydot/barcomic_server/total?label=downloads)

![Go Version](https://img.shields.io/github/go-mod/go-version/thegraydot/barcomic_server)
![Code Coverage](https://img.shields.io/badge/coverage-XX%25-blue)

An HTTP API for receiving comic book barcodes from the Barcomic Android application

## Barcomic App

The Barcomic application for Android and iOS (which leverages this HTTP API) is in active development and not currently publicly available.

## Quick Start

- Download [latest release](https://github.com/thegraydot/barcomic_server/releases/latest/) from GitHub releases page
- Double click to run the program
- This should automatically open and start the server in interactive mode, if it doesn't, open a terminal and run (e.g., `./barcomic-linux` on Linux)
- Pick an IP address from the list, usually you Ethernet or Wi-Fi adapter, so the Barcomic Android app can connect to the server
- Connect the Barcomic Android app using the QR code

## Command Arguments

| Flag | Description |
|------|-------------|
| `-a` | IP address or hostname to listen on (default `0.0.0.0`) |
| `-p` | Port to listen on (default `9999`) |
| `-k` | Enable HTTPS with a self-signed certificate |
| `-s` | Disable keystroke injection (useful with `-v` for logging only) |
| `-i` | Run interactive network interface selection (default `true`) |
| `-v` | Print verbose request logging |

> NOTE: For any of the examples provided below, change `barcomic-linux-amd64` to the correct release name you have downloaded. For example, `barcomic-windows-amd64.exe` or `barcomic-darwin-arm64`.

### Start server with IP and port specified

Use this if you have already configured your Barcomic Android app and want to start the server using known network configuration.

```
./barcomic-linux-amd64 -a 192.168.1.100 -p 9876
```

### Start server without keystrokes enabled

Use this if you don't want to have the server "type" the barcode out. Good when used in verbose or logging mode.

```
./barcomic-linux-amd64 -s -v
```

## Build Project

Compiled binaries are provided in GitHub releases for this project. However, the following instructions provide some general guidance on building the project. The barcomic server has the following requirements:

- Go (>= 1.24)
- Make

To build a snapshot for all platforms:

```
make build
```

Built binaries are placed in the `bin/` folder. Releases are published automatically by goreleaser when a `v*` tag is pushed.

## Platform Prerequisites

Barcomic injects keystrokes into the currently focused window, the same model as a USB barcode scanner. Focus your target application (e.g. a spreadsheet or web form) before scanning.

| Platform | Requirement |
|----------|-------------|
| Windows | None, PowerShell is built in |
| macOS | None, `osascript` is built in. Grant Accessibility permission on first run via **System Settings → Privacy & Security → Accessibility** |
| Linux X11 | Install `xdotool`: `apt install xdotool` / `dnf install xdotool` |
| Linux Wayland | Install `ydotool` and start the daemon: `systemctl --user start ydotoold`. Add your user to the `input` group: `sudo usermod -aG input $USER` |
