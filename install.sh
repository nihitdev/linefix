#!/bin/sh
set -eu

repo="nihitdev/linefix"
install_dir="${INSTALL_DIR:-${HOME}/.local/bin}"

os=$(uname -s | tr '[:upper:]' '[:lower:]')
case "$os" in
  linux|darwin) ;;
  *) echo "linefix: unsupported operating system: $os" >&2; exit 1 ;;
esac

machine=$(uname -m)
case "$machine" in
  x86_64|amd64) arch="amd64" ;;
  arm64|aarch64) arch="arm64" ;;
  *) echo "linefix: unsupported architecture: $machine" >&2; exit 1 ;;
esac

latest_url=$(curl -fsSL -o /dev/null -w '%{url_effective}' "https://github.com/${repo}/releases/latest")
tag=${latest_url##*/}
case "$tag" in
  v[0-9]*) ;;
  *) echo "linefix: could not determine latest release" >&2; exit 1 ;;
esac
version=${tag#v}
archive="linefix_${version}_${os}_${arch}.tar.gz"
base_url="https://github.com/${repo}/releases/download/${tag}"

tmp_dir=$(mktemp -d 2>/dev/null || mktemp -d -t linefix)
trap 'rm -rf "$tmp_dir"' EXIT HUP INT TERM

curl -fsSL "${base_url}/${archive}" -o "${tmp_dir}/${archive}"
curl -fsSL "${base_url}/SHA256SUMS" -o "${tmp_dir}/SHA256SUMS"
expected=$(awk -v name="$archive" '$2 == name { print $1 }' "${tmp_dir}/SHA256SUMS")
[ -n "$expected" ] || { echo "linefix: checksum not found for ${archive}" >&2; exit 1; }

if command -v sha256sum >/dev/null 2>&1; then
  actual=$(sha256sum "${tmp_dir}/${archive}" | awk '{print $1}')
elif command -v shasum >/dev/null 2>&1; then
  actual=$(shasum -a 256 "${tmp_dir}/${archive}" | awk '{print $1}')
else
  echo "linefix: sha256sum or shasum is required" >&2
  exit 1
fi
[ "$actual" = "$expected" ] || { echo "linefix: checksum verification failed" >&2; exit 1; }

tar -xzf "${tmp_dir}/${archive}" -C "$tmp_dir" linefix
mkdir -p "$install_dir"
install -m 0755 "${tmp_dir}/linefix" "${install_dir}/linefix"
echo "linefix ${version} installed to ${install_dir}/linefix"

# Releases created before the manual page was introduced do not contain it.
if tar -tzf "${tmp_dir}/${archive}" | grep -qx 'linefix\.1'; then
  tar -xzf "${tmp_dir}/${archive}" -C "$tmp_dir" linefix.1
  man_dir="${MAN_DIR:-${install_dir%/bin}/share/man/man1}"
  mkdir -p "$man_dir"
  install -m 0644 "${tmp_dir}/linefix.1" "${man_dir}/linefix.1"
  echo "Manual installed to ${man_dir}/linefix.1"
fi
case ":${PATH}:" in
  *:"${install_dir}":*) ;;
  *) echo "Add ${install_dir} to your PATH." ;;
esac
