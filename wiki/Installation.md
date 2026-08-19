# Installation

Official builds support Linux, macOS, and Windows on AMD64 and ARM64.

## Linux and macOS

```bash
curl -fsSL https://raw.githubusercontent.com/nihitdev/linefix/main/install.sh | sh
```

The installer detects the operating system and CPU, downloads the latest matching release, verifies its SHA-256 checksum, and installs to `~/.local/bin`.

Choose another destination with `INSTALL_DIR`:

```bash
curl -fsSL https://raw.githubusercontent.com/nihitdev/linefix/main/install.sh | INSTALL_DIR=/usr/local/bin sh
```

The default location does not need root. A system-wide destination may require elevated permissions.

## Windows

Run in PowerShell:

```powershell
irm https://raw.githubusercontent.com/nihitdev/linefix/main/install.ps1 | iex
```

The installer downloads and verifies the appropriate ZIP, installs `linefix.exe` under `%LOCALAPPDATA%\Programs\linefix\bin`, and adds that directory to the user PATH when needed. Administrator privileges are not required.

## Manual installation

1. Download your platform archive from [GitHub Releases](https://github.com/nihitdev/linefix/releases/latest).
2. Verify its digest using the published `SHA256SUMS` file.
3. Extract the archive.
4. Move `linefix` or `linefix.exe` into a directory on your PATH.

## Build with Go

Go 1.22 or newer is required:

```bash
go install github.com/nihitdev/linefix@latest
```

This produces a development build whose `--version` output is `linefix dev`.
