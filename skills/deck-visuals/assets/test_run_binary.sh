#!/usr/bin/env bash
# test_run_binary.sh — offline tests for run-binary.sh (download + checksum pin,
# cache, tamper rejection, and the build-from-source fallback). Uses file:// as the
# release base and stub binaries, so it needs no network and no Go.
set -euo pipefail

src="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
T="$(mktemp -d)"; trap 'rm -rf "$T"' EXIT
fails=0
ok()   { echo "ok   $*"; }
bad()  { echo "FAIL $*" >&2; fails=$((fails + 1)); }

case "$(uname -s)" in Darwin) os=darwin ;; *) os=linux ;; esac
case "$(uname -m)" in arm64|aarch64) arch=arm64 ;; *) arch=amd64 ;; esac
asset="dv-test-$os-$arch"
hash() { if command -v sha256sum >/dev/null; then sha256sum "$1" | cut -d' ' -f1; else shasum -a 256 "$1" | cut -d' ' -f1; fi; }

# A copy of the assets layout with a stub release.
mkdir -p "$T/assets/tools" "$T/rel"
cp "$src/run-binary.sh" "$T/assets/"; echo 9.9.9 > "$T/assets/VERSION"
printf '#!/bin/sh\necho stub-ok "$@"\n' > "$T/rel/$asset"
good="$(hash "$T/rel/$asset")"
printf '%s  %s\n' "$good" "$asset" > "$T/assets/tools/SHA256SUMS"
printf '#!/bin/sh\nmkdir -p "$BIN_DIR"\nprintf "#!/bin/sh\\necho stub-built \\"\\$@\\"\\n" > "$BIN_DIR/%s"\nchmod +x "$BIN_DIR/%s"\n' "$asset" "$asset" > "$T/assets/tools/build.sh"
chmod +x "$T/assets/tools/build.sh"
mkdir -p "$T/fakego"; printf '#!/bin/sh\necho go version go9.9.9 test/test\n' > "$T/fakego/go"; chmod +x "$T/fakego/go"

# A PATH holding only the tools run-binary.sh needs, so "no Go" holds on any machine
# (CI images ship a go in /usr/bin).
mkdir -p "$T/nogo"
for tool in awk cut tr uname dirname mkdir rm mv chmod curl sha256sum shasum; do
  found="$(command -v "$tool" 2>/dev/null || true)"
  case "$found" in /*) ln -s "$found" "$T/nogo/$tool" ;; esac
done
NOGO_PATH="$T/nogo"
run() { # run <cache> <base> <path> args...
  local cache="$1" base="$2" path="$3"; shift 3
  CLAUDE_PLUGIN_DATA="$cache" DV_RELEASE_BASE="$base" PATH="$path" \
    "$T/assets/run-binary.sh" dv-test "$T/assets/tools" "$T/assets/tools/build.sh" "$@"
}

# 1. downloads, verifies, runs
out="$(run "$T/c1" "file://$T/rel" "$NOGO_PATH" a b 2>&1)" && [ "$out" = "stub-ok a b" ] \
  && ok "download verified and ran" || bad "download: got '$out'"
[ -x "$T/c1/9.9.9/$asset" ] && ok "binary cached under the version" || bad "binary not cached"

# 2. cache hit works with the network gone
out="$(run "$T/c1" "file://$T/missing" "$NOGO_PATH" x 2>&1)" && [ "$out" = "stub-ok x" ] \
  && ok "cache hit needs no network" || bad "cache hit: got '$out'"

# 3. tampered download is rejected, deleted, and never run
printf '%s  %s\n' "$(printf "%064d" 0)" "$asset" > "$T/assets/tools/SHA256SUMS"
if out="$(run "$T/c3" "file://$T/rel" "$NOGO_PATH" 2>&1)"; then bad "tampered binary was run: '$out'"; else
  case "$out" in *"checksum mismatch"*) ok "tampered download rejected" ;; *) bad "tamper: wrong message '$out'" ;; esac
fi
ls "$T/c3/9.9.9/" 2>/dev/null | grep -q . && bad "tampered file left in cache: $(ls "$T/c3/9.9.9")" || ok "nothing left in cache after rejection"

# 4. no verified binary + Go available -> builds from source
out="$(run "$T/c4" "file://$T/missing" "$T/fakego:$NOGO_PATH" z 2>/dev/null)" && [ "$out" = "stub-built z" ] \
  && ok "falls back to a source build" || bad "fallback: got '$out'"

# 5. no verified binary and no Go -> clear failure
if out="$(run "$T/c5" "file://$T/missing" "$NOGO_PATH" 2>&1)"; then bad "ran with nothing available"; else
  case "$out" in *"no Go toolchain"*) ok "fails clearly with no binary and no Go" ;; *) bad "no-go: wrong message '$out'" ;; esac
fi

[ "$fails" -eq 0 ] && echo "all run-binary tests passed" || { echo "$fails failed" >&2; exit 1; }
