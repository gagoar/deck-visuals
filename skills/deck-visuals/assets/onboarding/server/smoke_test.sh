#!/usr/bin/env bash
# smoke_test.sh — build the dv-onboard server and prove both pickers still work end
# to end: each serves its page (with data injected and the shared stylesheet), takes
# one POST, writes a well-formed results.json, and shuts itself down. Run it locally
# after changing the server, and in CI on every change
# (see .github/workflows/onboarding-server.yml).
#
# Needs: go, curl. Exits non-zero on any failed check.
set -euo pipefail

here="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
work="$(mktemp -d)"
embedded=(picker.html levity-picker.html palette-picker.html picker.css)
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

# assert_grep <file> <pattern> <message> — fail unless pattern appears in file.
assert_grep() {
  local file="$1" pattern="$2" msg="$3"
  grep -q -- "$pattern" "$file" || fail "$msg"
}

# run_form <form> <data-file> <post-body> <grep-in-page>
run_form() {
  local form="$1" data="$2" body="$3" needle="$4"
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
}

# ---- onboard form ----
cat > "$work/cards.json" <<'JSON'
{"cards":[
  {"id":"the-office","title":"The Office","domain":"TV comedy","image_url":null,"gif_url":null,"degraded":true},
  {"id":"breaking-bad","title":"Breaking Bad","domain":"TV drama","image_url":null,"gif_url":null,"degraded":true}
]}
JSON
run_form onboard "$work/cards.json" \
  '{"talk":{"duration_minutes":30},"audience":{"description":"internal eng team"},"recognized_ids":["the-office"],"not_recognized_ids":["breaking-bad"],"skipped_ids":[],"suggestions":[]}' \
  "Breaking Bad"

echo "==> [onboard] results.json is well-formed"
out_onboard="$work/results-onboard.json"
assert_grep "$out_onboard" '"submitted_at"' "[onboard] missing submitted_at"
assert_grep "$out_onboard" '"duration_minutes": 30' "[onboard] duration not recorded"
assert_grep "$out_onboard" '"description": "internal eng team"' "[onboard] audience not verbatim"
assert_grep "$out_onboard" '"the-office"' "[onboard] recognized_ids wrong"
echo "    onboard results OK"

# ---- levity form ----
cat > "$work/slots.json" <<'JSON'
{"deck":"demo","slots":[
  {"id":"s3","slide":"Feature comparison","claim":"Too many features isn't always fantastic","shape":"Excess / feature bloat","candidates":[
    {"label":"The Simpsons — The Homer","gif_url":null,"caption":"every feature bolted on"},
    {"label":"Rube Goldberg machine","gif_url":null,"caption":"ten steps for a one-step job"}
  ]}
]}
JSON
run_form levity "$work/slots.json" \
  '{"selections":[{"slot_id":"s3","chosen_index":0,"chosen_label":"The Simpsons — The Homer","chosen_gif_url":null,"skipped":false}]}' \
  "Feature comparison"

echo "==> [levity] results.json is well-formed"
out_levity="$work/results-levity.json"
assert_grep "$out_levity" '"submitted_at"' "[levity] missing submitted_at"
assert_grep "$out_levity" '"slot_id": "s3"' "[levity] slot_id wrong"
assert_grep "$out_levity" '"chosen_index": 0' "[levity] chosen_index wrong"
echo "    levity results OK"

# ---- palette form ----
cat > "$work/presets.json" <<'JSON'
{"presets":[
  {"id":"house","subject":"brand default","accent":"#3ea6ff","bg":"#0a1428","data_slots":["#2894eb","#d95926","#199e70","#c98500","#d55181","#008300","#9085e9","#e66767"]},
  {"id":"space","subject":"space / night sky","accent":"#6a87f0","bg":"#0a1428","data_slots":["#2894eb","#d95926","#199e70","#c98500","#d55181","#008300","#9085e9","#e66767"]},
  {"id":"architecture","subject":"architecture","accent":"#00a4a4","bg":"#0a1428","data_slots":["#2894eb","#d95926","#199e70","#c98500","#d55181","#008300","#9085e9","#e66767"]},
  {"id":"sunrise","subject":"warm / keynote","accent":"#ffb020","bg":"#0a1428","data_slots":["#2894eb","#d95926","#199e70","#c98500","#d55181","#008300","#9085e9","#e66767"]}
]}
JSON
run_form palette "$work/presets.json" \
  '{"chosen_preset_id":"house"}' \
  "house"

echo "==> [palette] results.json is well-formed"
out_palette="$work/results-palette.json"
assert_grep "$out_palette" '"submitted_at"' "[palette] missing submitted_at"
assert_grep "$out_palette" '"chosen_preset_id": "house"' "[palette] chosen_preset_id wrong"
echo "    palette results OK"

echo "SMOKE PASS"
