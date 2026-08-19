# Troubleshooting

## `linefix: command not found`

The install directory is not on your PATH. Unix installations default to `~/.local/bin`; add it to your shell configuration:

```bash
export PATH="$HOME/.local/bin:$PATH"
```

On Windows, open a new terminal after installation so it picks up the updated user PATH.

## Permission denied during installation

Use the default per-user destination, or set a writable directory:

```bash
curl -fsSL https://raw.githubusercontent.com/nihitdev/linefix/main/install.sh | INSTALL_DIR="$HOME/bin" sh
```

## Binary-file rejection

linefix refuses files containing NUL bytes or a high proportion of control bytes. This protects binary formats from modification. If a legitimate text encoding triggers the safeguard, convert it to a byte-compatible text encoding such as UTF-8 first.

## `No line endings`

This is expected for an empty file or a file containing no LF or CRLF sequence. Conversion leaves such a file untouched.

## Verify a download manually

On Linux:

```bash
sha256sum -c SHA256SUMS
```

On macOS:

```bash
shasum -a 256 linefix_*.tar.gz
```

On Windows PowerShell:

```powershell
Get-FileHash .\linefix_0.1.0_windows_amd64.zip -Algorithm SHA256
```
