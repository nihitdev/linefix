# Command reference

## Convert to LF

```bash
linefix lf <file>
```

Converts every CRLF (`\r\n`) ending to LF (`\n`). Bare carriage returns are preserved.

## Convert to CRLF

```bash
linefix crlf <file>
```

Normalizes existing CRLF sequences to LF, then expands each LF exactly once to CRLF. This prevents the common accidental `\r\r\n` result.

## Check a file

```bash
linefix check <file>
```

Prints exactly one of:

| Output | Meaning |
| --- | --- |
| `LF` | All detected line endings are LF |
| `CRLF` | All detected line endings are CRLF |
| `Mixed` | The file contains both LF and CRLF |
| `No line endings` | The file is empty or contains no LF/CRLF ending |

## Global options

```bash
linefix --help
linefix --version
```

Release builds print their injected release version, such as `linefix 0.1.0`. Source builds default to `linefix dev`.

## Exit codes

- `0`: command completed successfully
- `1`: file or conversion error, including likely binary input
- `2`: invalid command-line arguments

Unknown commands, missing paths, extra arguments, and unsupported flags are rejected rather than ignored.
