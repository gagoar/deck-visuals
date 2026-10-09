#!/usr/bin/env bash
# Cross-compile dv-tools into runtime-free binaries and write SHA256SUMS.
# Needs Go once, on the operator/CI machine — the client that runs a binary
# needs nothing installed. Mirrors onboarding/server/build.sh.
#
# The binaries are NOT committed. SHA256SUMS is: the dv-tools wrapper downloads
# the binary for this release and refuses any file that does not match it. Builds
# are reproducible (pinned toolchain in go.mod, -trimpath, no VCS stamp, empty
# build id), so rebuilding must reproduce the committed sums.
#
# Env: BIN_DIR   output dir (default ./bin)
#      TARGETS   space-separated os/arch list (default: all five)
# Setting either skips writing SHA256SUMS (it would describe a partial build).
set -euo pipefail

here="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
bin="${BIN_DIR:-$here/bin}"
targets="${TARGETS:-darwin/arm64 darwin/amd64 linux/amd64 linux/arm64 windows/amd64}"
mkdir -p "$bin"

# Build from inside the module dir so `go build` finds go.mod regardless of the
# caller's working directory.
cd "$here"

build() {
  local goos="${1%/*}" goarch="${1#*/}" ext=""
  [ "$goos" = windows ] && ext=".exe"
  echo "building $goos/$goarch..."
  GOOS="$goos" GOARCH="$goarch" CGO_ENABLED=0 \
    go build -trimpath -buildvcs=false -ldflags "-s -w -buildid=" \
    -o "$bin/dv-tools-$goos-$goarch$ext" .
}

for t in $targets; do build "$t"; done

if [ -z "${BIN_DIR:-}" ] && [ -z "${TARGETS:-}" ]; then
  hash() { if command -v sha256sum >/dev/null; then sha256sum "$@"; else shasum -a 256 "$@"; fi; }
  ( cd "$bin" && hash $(ls dv-tools-* | LC_ALL=C sort) ) > "$here/SHA256SUMS"
  echo "wrote $here/SHA256SUMS"
fi

echo "done — binaries in $bin"
ls -1 "$bin"
