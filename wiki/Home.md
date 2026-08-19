# linefix wiki

`linefix` is a lightweight, cross-platform CLI for converting text-file line endings between LF and CRLF.

## Start here

- [Installation](Installation) — Linux, macOS, Windows, and manual setup
- [Command reference](Command-reference) — every command, output, and exit behavior
- [How it works](How-it-works) — conversion, detection, and file-safety details
- [Development](Development) — build, test, and contribute
- [Releasing](Releasing) — release automation and artifacts
- [Troubleshooting](Troubleshooting) — common installation and usage issues

## Quick example

```bash
linefix check README.md
linefix lf notes.txt
linefix crlf windows-file.txt
```

`check` reports `LF`, `CRLF`, `Mixed`, or `No line endings`. Conversion happens in place while preserving the file's existing permissions and trailing-newline state.

## Project links

- [Website](https://nihitdev.github.io/linefix/)
- [Source repository](https://github.com/nihitdev/linefix)
- [Latest release](https://github.com/nihitdev/linefix/releases/latest)
- [Issue tracker](https://github.com/nihitdev/linefix/issues)

linefix is licensed under the GNU General Public License v3.0 or later (`GPL-3.0-or-later`).
