#!/bin/sh
# run-binary.sh — run a deck-visuals Go tool without committing its binaries.
#
#   run-binary.sh <name> <wrapper-dir> <build-script> [args...]
#
# <wrapper-dir>/SHA256SUMS pins the expected hash of every released binary
# (<name>-<os>-<arch>). The binary is taken, in order, from:
#   1. the cache, if already there;
#   2. this release's GitHub Release assets, accepted ONLY if its SHA-256 matches
#      the pin (a mismatch is deleted and reported, never run);
#   3. a local build from the committed Go source, if `go` is on PATH.
# Nothing here reads credentials; the only network call is one unauthenticated
# HTTPS GET for the release asset.
#
# Env: DV_RELEASE_BASE  override the download base URL (used by the tests)
#      CLAUDE_PLUGIN_DATA / XDG_CACHE_HOME  where binaries are cached
set -eu

name="$1"; wrapper_dir="$2"; build_script="$3"; shift 3
assets="$(cd "$(dirname "$0")" && pwd)"
version="$(tr -d '[:space:]' < "$assets/VERSION")"
repo_url="https://github.com/gagoar/deck-visuals/releases/download"
base="${DV_RELEASE_BASE:-$repo_url/v$version}"

case "$(uname -s)" in
  Darwin) os=darwin ;;
  Linux) os=linux ;;
  *) echo "$name: unsupported OS $(uname -s)" >&2; exit 1 ;;
esac
case "$(uname -m)" in
  arm64|aarch64) arch=arm64 ;;
  x86_64|amd64) arch=amd64 ;;
  *) echo "$name: unsupported architecture $(uname -m)" >&2; exit 1 ;;
esac

asset="$name-$os-$arch"
cache="${CLAUDE_PLUGIN_DATA:-${XDG_CACHE_HOME:-$HOME/.cache}/deck-visuals}/$version"
bin="$cache/$asset"

sha256() {
  if command -v sha256sum >/dev/null 2>&1; then sha256sum "$1" | cut -d' ' -f1
  else shasum -a 256 "$1" | cut -d' ' -f1; fi
}

fetch_verified() {
  want="$(awk -v a="$asset" '$2 == a { print $1 }' "$wrapper_dir/SHA256SUMS")"
  if [ -z "$want" ]; then echo "$name: no pinned hash for $asset" >&2; return 1; fi
  tmp="$bin.part.$$"
  if ! curl -fsSL --max-time 120 -o "$tmp" "$base/$asset" 2>/dev/null; then
    rm -f "$tmp"; echo "$name: could not download $base/$asset" >&2; return 1
  fi
  got="$(sha256 "$tmp")"
  if [ "$got" != "$want" ]; then
    rm -f "$tmp"
    echo "$name: checksum mismatch for $asset (expected $want, got $got); discarded" >&2
    return 1
  fi
  chmod +x "$tmp" && mv "$tmp" "$bin"
}

if [ ! -x "$bin" ]; then
  mkdir -p "$cache"
  if ! fetch_verified; then
    if command -v go >/dev/null 2>&1; then
      echo "$name: building from source with $(go version | cut -d' ' -f3)..." >&2
      BIN_DIR="$cache" TARGETS="$os/$arch" "$build_script" >&2 || { echo "$name: build failed" >&2; exit 1; }
    else
      echo "$name: no verified binary and no Go toolchain. Install Go and retry, or fetch $asset from $base" >&2
      exit 1
    fi
  fi
fi

exec "$bin" "$@"
