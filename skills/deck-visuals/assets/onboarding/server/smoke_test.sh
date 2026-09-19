#!/usr/bin/env bash
# smoke_test.sh — build the dv-onboard server and prove both pickers still work end
# to end: each serves its page (with data injected and the shared stylesheet), takes
# one POST, writes a well-formed results.json, and shuts itself down. Run it locally
# after changing the server, and in CI on every change
# (see .github/workflows/onboarding-server.yml).
#
# Needs: go, curl, python3. Exits non-zero on any failed check.
set -euo pipefail

here="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
work="$(mktemp -d)"
embedded=(picker.html levity-picker.html picker.css)
cleanup() {
  for f in "${embedded[@]}"; do rm -f "$here/$f"; done
  [ -n "${srvpid:-}" ] && kill "$srvpid" 2>/dev/null || true
  rm -rf "$work"
}
trap cleanup EXIT

fail() { echo "SMOKE FAIL: $*" >&2; exit 1; }

# go:embed can't reach a parent dir, so copy the canonical assets in before any tool
# that compiles the package.
for f in "${embedded[@]}"; do cp "$here/../$f" "$here/$f"; done

echo "==> go vet"
( cd "$here" && go vet ./... )

echo "==> build (host platform)"
( cd "$here" && CGO_ENABLED=0 go build -trimpath -o "$work/dv-onboard" . )

# run_form <form> <data-file> <post-body> <grep-in-page> <python-asserts>
run_form() {
  local form="$1" data="$2" body="$3" needle="$4" pyfile="$5"
  echo "==> [$form] launch on a random loopback port"
  local out="$work/results-$form.json"
  "$work/dv-onboard" --form "$form" --data "$data" --out "$out" --port 0 --timeout 20 \
    >"$work/url-$form.txt" 2>"$work/err-$form.txt" &
  srvpid=$!

  local url=""
  for _ in $(seq 1 100); do
    url="$(cat "$work/url-$form.txt" 2>/dev/null || true)"
    [ -n "$url" ] && break
    sleep 0.1
  done
  [ -n "$url" ] || fail "[$form] server never printed a URL (stderr: $(cat "$work/err-$form.txt"))"
  local base="${url%/}"
  echo "    serving at $url"
  case "$url" in http://127.0.0.1:*) : ;; *) fail "[$form] did not bind loopback: $url" ;; esac

  echo "==> [$form] GET / — data injected, token gone"
  local page; page="$(curl -fsS "$base/")"
  echo "$page" | grep -q "$needle" || fail "[$form] expected content '$needle' not in page"
  echo "$page" | grep -q "REPLACE_DATA_JSON" && fail "[$form] injection token still present"

  echo "==> [$form] GET /picker.css — shared stylesheet served"
  curl -fsS "$base/picker.css" | grep -q "." || fail "[$form] picker.css empty or missing"

  echo "==> [$form] POST /submit — accepted"
  local resp; resp="$(curl -fsS -X POST "$base/submit" -H 'Content-Type: application/json' -d "$body")"
  echo "$resp" | grep -q '"ok":true' || fail "[$form] submit did not return ok:true (got: $resp)"

  echo "==> [$form] server shuts down after one submit (exit 0)"
  wait "$srvpid"; local code=$?; srvpid=""
  [ "$code" -eq 0 ] || fail "[$form] server exited $code, expected 0"

  echo "==> [$form] results.json is well-formed"
  python3 "$pyfile" "$out"
}

# ---- onboard form ----
cat > "$work/cards.json" <<'JSON'
{"cards":[
  {"id":"the-office","title":"The Office","domain":"TV comedy","image_url":null,"gif_url":null,"degraded":true},
  {"id":"breaking-bad","title":"Breaking Bad","domain":"TV drama","image_url":null,"gif_url":null,"degraded":true}
]}
JSON
cat > "$work/assert_onboard.py" <<'PY'
import json, sys
r = json.load(open(sys.argv[1]))
assert r.get("submitted_at"), "missing submitted_at"
assert r["talk"]["duration_minutes"] == 30, "duration not recorded"
assert r["audience"]["description"] == "internal eng team", "audience not verbatim"
assert r["recognized_ids"] == ["the-office"], "recognized_ids wrong"
print("    onboard results OK")
PY
run_form onboard "$work/cards.json" \
  '{"talk":{"duration_minutes":30},"audience":{"description":"internal eng team"},"recognized_ids":["the-office"],"not_recognized_ids":["breaking-bad"],"skipped_ids":[],"suggestions":[]}' \
  "Breaking Bad" "$work/assert_onboard.py"

# ---- levity form ----
cat > "$work/slots.json" <<'JSON'
{"deck":"demo","slots":[
  {"id":"s3","slide":"Feature comparison","claim":"Too many features isn't always fantastic","shape":"Excess / feature bloat","candidates":[
    {"label":"The Simpsons — The Homer","gif_url":null,"caption":"every feature bolted on"},
    {"label":"Rube Goldberg machine","gif_url":null,"caption":"ten steps for a one-step job"}
  ]}
]}
JSON
cat > "$work/assert_levity.py" <<'PY'
import json, sys
r = json.load(open(sys.argv[1]))
assert r.get("submitted_at"), "missing submitted_at"
sel = r["selections"][0]
assert sel["slot_id"] == "s3", "slot_id wrong"
assert sel["chosen_index"] == 0, "chosen_index wrong"
print("    levity results OK")
PY
run_form levity "$work/slots.json" \
  '{"selections":[{"slot_id":"s3","chosen_index":0,"chosen_label":"The Simpsons — The Homer","chosen_gif_url":null,"skipped":false}]}' \
  "Feature comparison" "$work/assert_levity.py"

echo "SMOKE PASS"
