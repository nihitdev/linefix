# linefix

`linefix` is a small, dependency-free CLI that converts text-file line endings
in place. It exists to make line-ending cleanup explicit, safe, and easy to use
in scripts without bringing in a formatter or runtime.

## Features

- Converts LF and CRLF safely without producing accidental `\r\r\n` sequences
- Detects LF, CRLF, mixed, and newline-free files
- Preserves file permissions and final-newline status
- Skips writes when the requested format is already present
- Refuses likely binary input
- Ships as one executable for Linux, macOS, and Windows

## Installation

### Linux and macOS

```sh
curl -fsSL https://raw.githubusercontent.com/nihitdev/linefix/main/install.sh | sh
```

The installer detects Linux or macOS and x86-64 or ARM64, verifies the archive
checksum, and installs to `~/.local/bin` by default. Override the destination:

```sh
curl -fsSL https://raw.githubusercontent.com/nihitdev/linefix/main/install.sh | INSTALL_DIR=/usr/local/bin sh
```

The default works without root. A system destination may require running the
downloaded script with appropriate permissions. Apple Silicon (`arm64`) and
Intel (`amd64`) Macs are both supported.

### Windows

Run in PowerShell:

```powershell
irm https://raw.githubusercontent.com/nihitdev/linefix/main/install.ps1 | iex
```

This installs to `%LOCALAPPDATA%\Programs\linefix\bin`, verifies the checksum,
and adds that directory to the current user's `PATH`. It does not require
Administrator privileges and is safe to run again when upgrading.

### Manual installation

Download the archive for your OS and architecture from [GitHub Releases](https://github.com/nihitdev/linefix/releases),
verify it against `SHA256SUMS`, extract it, and place `linefix` (or
`linefix.exe`) in a directory on your `PATH`.

With Go 1.22 or newer, you can instead build the current development version:

```sh
go install github.com/nihitdev/linefix@latest
```

## Usage

```text
linefix lf <file>
linefix crlf <file>
linefix check <file>
linefix --version
```

- `lf` converts CRLF (`\r\n`) endings to LF (`\n`).
- `crlf` converts LF endings to CRLF. Existing CRLF endings are normalized first.
- `check` prints `LF`, `CRLF`, `Mixed`, or `No line endings`.

Examples:

```sh
linefix check README.md
linefix lf script.sh
linefix crlf notes.txt
```

Conversions preserve a file's final-newline status and permissions. A file that
already has the requested endings is left untouched. Empty files and files with
no newline are valid and are not changed.

## Safety and conversion behavior

Before conversion, `linefix` samples the input for NUL bytes and an unusually
high proportion of non-text control bytes. Likely binary files are rejected
instead of modified. The converted content is written to a temporary file in the
same directory and then moved over the original.

When converting to CRLF, existing CRLF sequences are first reduced to LF. Every
LF is then expanded exactly once, preventing accidental `\r\r\n` sequences from
the common naive replacement approach.

`linefix` deliberately treats only LF and CRLF as line endings. A bare carriage
return is preserved as ordinary input.

## Building from source

```sh
git clone https://github.com/nihitdev/linefix.git
cd linefix
go build -o linefix .
```

Release builds inject their version at link time:

```sh
go build -ldflags "-s -w -X main.version=0.1.0" -o linefix .
```

## Development

```sh
gofmt -w .
go test ./...
go vet ./...
```

The project uses only the Go standard library, so there are no dependencies to
download beyond the Go toolchain.

To reproduce all release archives locally, run:

```sh
./scripts/build-release.sh 0.1.0
```

## Release process

Create and push an annotated `v*` tag. The release workflow cross-compiles all
supported targets, packages each executable with the README and license,
generates `SHA256SUMS`, and publishes a GitHub Release.

## Supported platforms

Official release archives cover Linux, macOS, and Windows on `amd64` and
`arm64`.

## License

linefix is licensed under the GNU General Public License v3.0 or later
(`GPL-3.0-or-later`). See [LICENSE](LICENSE).
