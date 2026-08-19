# Development

linefix requires Go 1.22 or newer and uses only the Go standard library.

## Clone and verify

```bash
git clone https://github.com/nihitdev/linefix.git
cd linefix
gofmt -w .
go vet ./...
go test ./...
go build .
```

## Release-style local build

```bash
go build -trimpath -ldflags "-s -w -X main.version=0.1.0" -o linefix .
./linefix --version
```

## Build all release archives

```bash
./scripts/build-release.sh 0.1.0
```

The generated archives and `SHA256SUMS` are placed under `dist/`.

## Design guidelines

- Keep dependencies at zero unless a clear requirement justifies one.
- Keep CLI output concise and suitable for shell use.
- Preserve input data outside the requested line-ending change.
- Add tests for every behavior change and edge case.
- Keep Linux, macOS, and Windows behavior aligned.

Bug reports and focused pull requests are welcome in the [repository](https://github.com/nihitdev/linefix).
