#!/usr/bin/env python3
"""
check_media_url.py — verify a levity GIF/video URL resolves to real media
before it ships in a deck.

Why this exists: a made-up or link-rotted Giphy/Tenor URL renders as a
broken image in a delivered deck. A page URL (giphy.com/gifs/<slug>,
tenor.com/view/...) dropped into an <img src> also renders broken — it
returns an HTML document, not image bytes. This script is the mechanical
half of the verification step in references/levity.md: run it on every
direct media URL before pinning it into a fragment.

Checks per URL, via a real HTTP GET (redirects followed):
  - status code must be 200
  - Content-Type must start with "image/" or "video/"
  - body must be at least MIN_BYTES (1 KB) — catches error pages and
    empty/placeholder responses served with a 200 and an image-ish header

The GET is streamed/ranged rather than fully downloaded: only enough of
the body is read to confirm it clears MIN_BYTES, so this stays cheap even
against a large .mp4.

Usage:
  python3 check_media_url.py <url> [<url> ...]

Exit code: 1 if any URL FAILs, else 0. Uses only the Python standard
library (urllib) — no third-party dependencies.
"""

import sys
import urllib.error
import urllib.request

MIN_BYTES = 1024
TIMEOUT_SECONDS = 15
READ_CAP_BYTES = 65536  # only read enough to confirm size/content past MIN_BYTES
USER_AGENT = "Mozilla/5.0 (compatible; check_media_url.py; +deck-visuals)"


def check_one(url):
    """Return (ok, status, content_type, size_note) for a single URL."""
    req = urllib.request.Request(
        url,
        headers={
            "User-Agent": USER_AGENT,
            # Ask for a small range; not all hosts honor it, but when they
            # do this avoids pulling a whole video over the wire.
            "Range": f"bytes=0-{READ_CAP_BYTES - 1}",
        },
    )
    try:
        with urllib.request.urlopen(req, timeout=TIMEOUT_SECONDS) as resp:
            status = resp.status
            content_type = resp.headers.get("Content-Type", "").split(";")[0].strip()
            declared_length = resp.headers.get("Content-Length")
            body = resp.read(READ_CAP_BYTES)
    except urllib.error.HTTPError as e:
        # A 403/404/etc still carries a status and headers worth reporting.
        status = e.code
        content_type = e.headers.get("Content-Type", "").split(";")[0].strip() if e.headers else ""
        declared_length = e.headers.get("Content-Length") if e.headers else None
        body = b""
    except (urllib.error.URLError, TimeoutError, OSError) as e:
        return False, None, "", f"request failed: {e}"

    # Prefer the declared Content-Length when present and larger than what
    # we actually read (e.g. a ranged/partial read on a big file).
    try:
        size = max(len(body), int(declared_length)) if declared_length else len(body)
    except ValueError:
        size = len(body)

    is_media_type = content_type.startswith("image/") or content_type.startswith("video/")
    # A 200 (full body) or 206 (partial content honoring our Range header)
    # both indicate the resource itself resolved.
    status_ok = status in (200, 206)
    size_ok = size >= MIN_BYTES

    ok = status_ok and is_media_type and size_ok
    note = f"{size} bytes read/declared"
    return ok, status, content_type or "(none)", note


def main(argv):
    urls = argv[1:]
    if not urls:
        print("usage: python3 check_media_url.py <url> [<url> ...]", file=sys.stderr)
        return 2

    any_fail = False
    for url in urls:
        ok, status, content_type, note = check_one(url)
        verdict = "PASS" if ok else "FAIL"
        if not ok:
            any_fail = True
        status_str = status if status is not None else "n/a"
        print(f"[{verdict}] {url}")
        print(f"        status={status_str} content-type={content_type} {note}")

    return 1 if any_fail else 0


if __name__ == "__main__":
    sys.exit(main(sys.argv))
