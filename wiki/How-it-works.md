# How it works

linefix reads a regular file, checks whether it looks like text, converts its bytes in memory, and only writes when the result differs.

## Detection

Every LF byte is classified as:

- CRLF when immediately preceded by a carriage return
- LF otherwise

Seeing both produces `Mixed`. Seeing neither produces `No line endings`. Bare carriage returns are not treated as line endings.

## Binary safeguard

Before conversion, linefix samples up to the first 8 KiB. It rejects input containing a NUL byte or an unusually high proportion of non-text control bytes. This is intentionally a conservative heuristic, not a file-format detector.

`check` applies the same safeguard so binary data is never presented as meaningful text-line-ending output.

## In-place replacement

When a change is needed, linefix:

1. creates a temporary file in the original file's directory;
2. applies the original permission bits;
3. writes and syncs the converted content;
4. replaces the original path.

Using the same directory keeps replacement on the same filesystem. Files already using the requested format are never rewritten.

## Edge cases

- Empty files are accepted and remain empty.
- Files without a final newline keep that state.
- Files with a final newline keep one in the requested format.
- Mixed files are normalized to the selected format.
- Existing CRLF is normalized before CRLF conversion.
