#!/usr/bin/env bash
# Regenerate the plugin icon, social card and README proof image from their HTML
# sources — doc tooling only, not a runtime dependency. Requires Google Chrome (or
# set CHROME to another Chromium-based browser's binary).
set -euo pipefail
cd "$(dirname "$0")"

CHROME="${CHROME:-/Applications/Google Chrome.app/Contents/MacOS/Google Chrome}"
if [ ! -x "$CHROME" ]; then
  echo "error: Chrome not found at $CHROME (set \$CHROME to another Chromium browser)" >&2
  exit 1
fi

shot() { # shot <html> <out.png> <width> <height> [scale]
  "$CHROME" --headless=new --disable-gpu --hide-scrollbars --force-prefers-reduced-motion \
    --force-device-scale-factor="${5:-1}" --virtual-time-budget=5000 \
    --window-size="${3},${4}" --screenshot="$2" "file://$PWD/$1" 2>/dev/null
  echo "wrote $2"
}

shot assets/icon.html ../.claude-plugin/icon.png 1024 1024
shot assets/og.html assets/og.png 1200 630
shot assets/gallery/beforeafter.html assets/before-after.png 800 414 2
echo "done. Inspect the PNGs before committing."
