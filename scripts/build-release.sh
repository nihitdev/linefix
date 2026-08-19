#!/bin/sh
set -eu

version=${1:-}
case "$version" in
  [0-9]*.[0-9]*.[0-9]*) ;;
  *) echo "usage: $0 MAJOR.MINOR.PATCH" >&2; exit 2 ;;
esac

project_dir=$(CDPATH= cd -- "$(dirname -- "$0")/.." && pwd)
dist_dir="${project_dir}/dist"
rm -rf "$dist_dir"
mkdir -p "$dist_dir"

for target in linux/amd64 linux/arm64 darwin/amd64 darwin/arm64 windows/amd64 windows/arm64; do
  os=${target%/*}
  arch=${target#*/}
  name="linefix_${version}_${os}_${arch}"
  stage=$(mktemp -d)
  binary="linefix"
  [ "$os" = windows ] && binary="linefix.exe"

  (
    cd "$project_dir"
    GOOS="$os" GOARCH="$arch" CGO_ENABLED=0 go build \
      -trimpath -ldflags "-s -w -X main.version=${version}" \
      -o "${stage}/${binary}" .
  )
  cp "${project_dir}/README.md" "${project_dir}/LICENSE" "$stage/"
  if [ "$os" = windows ]; then
    if command -v zip >/dev/null 2>&1; then
      (cd "$stage" && zip -q "${dist_dir}/${name}.zip" "$binary" README.md LICENSE)
    elif command -v 7z >/dev/null 2>&1; then
      (cd "$stage" && 7z a -bd -y "${dist_dir}/${name}.zip" "$binary" README.md LICENSE >/dev/null)
    else
      echo "linefix release: zip or 7z is required to package Windows" >&2
      exit 1
    fi
  else
    cp "${project_dir}/man/linefix.1" "$stage/"
    tar -C "$stage" -czf "${dist_dir}/${name}.tar.gz" "$binary" linefix.1 README.md LICENSE
  fi
  rm -rf "$stage"
done

(cd "$dist_dir" && sha256sum linefix_* > SHA256SUMS)
