#!/usr/bin/env bash
# Cross-compile the dv-onboard server into runtime-free binaries and drop them in
# ../bin. Needs Go once, on the operator/CI machine — the client that runs a binary
# needs nothing installed.
#
# go:embed cannot reach a parent directory, so we copy the canonical embedded assets
# (picker.html, levity-picker.html, palette-picker.html, picker.css, bg.jpg) next to
# main.go before building, then clean them up.
set -euo pipefail

here="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
bin="$here/../bin"
mkdir -p "$bin"

embedded=(picker.html levity-picker.html palette-picker.html picker.css bg.jpg)
for f in "${embedded[@]}"; do cp "$here/../$f" "$here/$f"; done
trap 'for f in "${embedded[@]}"; do rm -f "$here/$f"; done' EXIT

# Build from inside the module dir so `go build` finds go.mod regardless of the
# caller's working directory (CI runs this script from the repo root).
cd "$here"

build() {
  local goos="$1" goarch="$2" ext="${3:-}"
  echo "building $goos/$goarch..."
  GOOS="$goos" GOARCH="$goarch" CGO_ENABLED=0 \
    go build -trimpath -ldflags "-s -w" \
    -o "$bin/dv-onboard-$goos-$goarch$ext" .
}

build darwin arm64
build darwin amd64
build linux amd64
build linux arm64
build windows amd64 .exe

echo "done — binaries in $bin"
ls -1 "$bin"
