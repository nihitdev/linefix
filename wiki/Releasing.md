# Releasing

Releases are automated by `.github/workflows/release.yml`.

## Process

1. Update `CHANGELOG.md`.
2. Ensure `go test ./...` and `go vet ./...` pass.
3. Push the intended release commit to `main` and wait for CI.
4. Create and push an annotated semantic-version tag:

   ```bash
   git tag -a v0.2.0 -m "linefix v0.2.0"
   git push origin v0.2.0
   ```

5. Verify the generated GitHub Release and all checksums.

## Published targets

- Linux: AMD64 and ARM64 (`.tar.gz`)
- macOS: AMD64 and ARM64 (`.tar.gz`)
- Windows: AMD64 and ARM64 (`.zip`)

Each archive contains the executable, `README.md`, and `LICENSE`. `SHA256SUMS` covers every archive. The version is injected at build time with Go linker flags.
