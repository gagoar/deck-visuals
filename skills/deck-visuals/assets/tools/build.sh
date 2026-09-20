#!/usr/bin/env bash
# Cross-compile the dv-tools binary into runtime-free binaries and drop them
# in ./bin. Needs Go once, on the operator/CI machine — the client that runs
# a binary needs nothing installed. Mirrors
# skills/deck-visuals/assets/onboarding/server/build.sh.
set -euo pipefail

here="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
bin="$here/bin"
mkdir -p "$bin"

# Build from inside the module dir so `go build` finds go.mod regardless of the
# caller's working directory.
cd "$here"

build() {
  local goos="$1" goarch="$2" ext="${3:-}"
  echo "building $goos/$goarch..."
  GOOS="$goos" GOARCH="$goarch" CGO_ENABLED=0 \
    go build -trimpath -ldflags "-s -w" \
    -o "$bin/dv-tools-$goos-$goarch$ext" .
}

build darwin arm64
build darwin amd64
build linux amd64
build linux arm64
build windows amd64 .exe

echo "done — binaries in $bin"
ls -1 "$bin"
