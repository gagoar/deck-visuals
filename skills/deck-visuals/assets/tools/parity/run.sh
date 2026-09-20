#!/usr/bin/env bash
# run.sh — language-neutral golden characterization harness for the three
# deck-visuals scripts being ported from Python to Go.
#
# Usage:
#   bash run.sh <python|gotools> [--capture]
#
#   python   - run the existing Python scripts (the source of truth today).
#   gotools  - run ../dv-tools <subcommand> (does not exist yet in stage 1;
#              this mode is for stage 2 once the Go port lands).
#   --capture - only valid with the python runner. (Re)writes golden/<case>.out
#              and golden/<case>.exit from what the python scripts actually do
#              right now, instead of diffing against them.
#
# For every case file in cases/*.case this:
#   1. resolves the executable for the case's subcommand under the requested
#      runner,
#   2. substitutes the {BASE} placeholder (the check-media-url fixture
#      server's base URL) into the case's args,
#   3. runs it, capturing stdout and the exit code,
#   4. normalizes both the captured output and (unless --capture) the golden
#      output the same way, then diffs them.
#
# Normalization (applied to both captured and golden output before diffing):
#   - 127.0.0.1:<port>              -> 127.0.0.1:PORT
#   - "<N> bytes read/declared"     -> "<N> bytes read/declared" literal marker
# Nothing else is normalized — a real behavior difference should still show up
# as a diff.
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
cd "$SCRIPT_DIR"

MODE="${1:-}"
CAPTURE=0
if [[ "${2:-}" == "--capture" ]]; then
  CAPTURE=1
fi

if [[ "$MODE" != "python" && "$MODE" != "gotools" ]]; then
  echo "usage: run.sh <python|gotools> [--capture]" >&2
  exit 2
fi

if [[ "$CAPTURE" == 1 && "$MODE" != "python" ]]; then
  echo "--capture only regenerates goldens from the python runner; run: run.sh python --capture" >&2
  exit 2
fi

CASES_DIR="$SCRIPT_DIR/cases"
GOLDEN_DIR="$SCRIPT_DIR/golden"
mkdir -p "$GOLDEN_DIR"

normalize() {
  # Reads stdin, writes normalized text to stdout.
  sed -E \
    -e 's/127\.0\.0\.1:[0-9]+/127.0.0.1:PORT/g' \
    -e 's/[0-9]+ bytes read\/declared/<N> bytes read\/declared/g'
}

# --- start the Go fixture HTTP server -----------------------------------
# Its own module lives under fixtures/server so it never depends on the
# future dv-tools module. It prints "PORT=<n>" to stdout as soon as it is
# ready, then serves until we kill it.
SERVER_DIR="$SCRIPT_DIR/fixtures/server"
SERVER_LOG="$(mktemp)"

(cd "$SERVER_DIR" && exec go run .) >"$SERVER_LOG" 2>&1 &
SERVER_PID=$!

cleanup() {
  kill "$SERVER_PID" >/dev/null 2>&1 || true
  wait "$SERVER_PID" >/dev/null 2>&1 || true
  rm -f "$SERVER_LOG"
}
trap cleanup EXIT

PORT=""
for _ in $(seq 1 200); do
  if ! kill -0 "$SERVER_PID" >/dev/null 2>&1; then
    echo "fixture server exited early:" >&2
    cat "$SERVER_LOG" >&2
    exit 1
  fi
  if grep -q '^PORT=' "$SERVER_LOG" 2>/dev/null; then
    PORT="$(grep '^PORT=' "$SERVER_LOG" | head -n1 | cut -d= -f2)"
    break
  fi
  sleep 0.1
done

if [[ -z "$PORT" ]]; then
  echo "fixture server did not report a port in time:" >&2
  cat "$SERVER_LOG" >&2
  exit 1
fi
BASE="http://127.0.0.1:${PORT}"

# --- resolve the executable for a subcommand under the requested runner --
resolve_cmd() {
  local subcommand="$1"
  case "$subcommand" in
    check-media-url)
      if [[ "$MODE" == "python" ]]; then
        printf '%s\n' "python3 ../../scripts/check_media_url.py"
      else
        printf '%s\n' "../dv-tools check-media-url"
      fi
      ;;
    coverage)
      if [[ "$MODE" == "python" ]]; then
        printf '%s\n' "python3 ../../scripts/coverage_report.py"
      else
        printf '%s\n' "../dv-tools coverage"
      fi
      ;;
    validate-palette)
      if [[ "$MODE" == "python" ]]; then
        printf '%s\n' "python3 ../../scripts/validate_palette.py"
      else
        printf '%s\n' "../dv-tools validate-palette"
      fi
      ;;
    *)
      echo "unknown subcommand: $subcommand" >&2
      exit 2
      ;;
  esac
}

# --- run every case -------------------------------------------------------
TOTAL=0
FAILED=0
FAILED_NAMES=()

shopt -s nullglob
CASE_FILES=("$CASES_DIR"/*.case)
shopt -u nullglob

if [[ ${#CASE_FILES[@]} -eq 0 ]]; then
  echo "no case files found under $CASES_DIR" >&2
  exit 2
fi

for case_file in "${CASE_FILES[@]}"; do
  name="$(basename "$case_file" .case)"
  TOTAL=$((TOTAL + 1))

  # Each case file sets SUBCOMMAND and ARGS (a bash array). Reset both so a
  # case that forgets ARGS doesn't inherit the previous case's.
  SUBCOMMAND=""
  ARGS=()
  # shellcheck disable=SC1090
  source "$case_file"

  resolved="$(resolve_cmd "$SUBCOMMAND")"
  # shellcheck disable=SC2206
  cmd_words=($resolved)

  run_args=()
  for a in "${ARGS[@]+"${ARGS[@]}"}"; do
    run_args+=("${a//\{BASE\}/$BASE}")
  done

  set +e
  actual_out="$("${cmd_words[@]}" "${run_args[@]+"${run_args[@]}"}" 2>/dev/null)"
  actual_exit=$?
  set -e

  norm_actual="$(printf '%s' "$actual_out" | normalize)"

  golden_out_file="$GOLDEN_DIR/$name.out"
  golden_exit_file="$GOLDEN_DIR/$name.exit"

  if [[ "$CAPTURE" == 1 ]]; then
    printf '%s' "$actual_out" >"$golden_out_file"
    printf '%s\n' "" >>"$golden_out_file"  # trailing newline, tidy for diffing by hand
    printf '%s\n' "$actual_exit" >"$golden_exit_file"
    echo "captured $name (exit=$actual_exit)"
    continue
  fi

  if [[ ! -f "$golden_out_file" || ! -f "$golden_exit_file" ]]; then
    echo "FAIL $name — missing golden (run: run.sh python --capture)"
    FAILED=$((FAILED + 1))
    FAILED_NAMES+=("$name")
    continue
  fi

  golden_out="$(cat "$golden_out_file")"
  golden_exit="$(cat "$golden_exit_file")"
  norm_golden="$(printf '%s' "$golden_out" | normalize)"

  ok=1
  if [[ "$norm_actual" != "$norm_golden" ]]; then
    ok=0
  fi
  if [[ "$actual_exit" != "$golden_exit" ]]; then
    ok=0
  fi

  if [[ "$ok" == 1 ]]; then
    echo "PASS $name"
  else
    echo "FAIL $name (exit actual=$actual_exit golden=$golden_exit)"
    diff -u <(printf '%s\n' "$norm_golden") <(printf '%s\n' "$norm_actual") | sed 's/^/    /' || true
    FAILED=$((FAILED + 1))
    FAILED_NAMES+=("$name")
  fi
done

echo
if [[ "$CAPTURE" == 1 ]]; then
  echo "captured $TOTAL case(s) under $GOLDEN_DIR"
  exit 0
fi

PASSED=$((TOTAL - FAILED))
echo "$PASSED/$TOTAL cases passed"
if [[ "$FAILED" -gt 0 ]]; then
  echo "FAILED: ${FAILED_NAMES[*]}"
  exit 1
fi
exit 0
