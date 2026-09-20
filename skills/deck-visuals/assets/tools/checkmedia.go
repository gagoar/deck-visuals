// checkmedia.go — port of check_media_url.py: verify a levity GIF/video URL
// resolves to real media before it ships in a deck.
//
// Checks per URL, via a real HTTP GET (redirects followed):
//   - status code must be 200 or 206
//   - Content-Type must start with "image/" or "video/"
//   - body must be at least minBytes (1 KB)
//
// Only enough of the body is read to confirm size (capped at readCapBytes),
// so this stays cheap even against a large video file.
package main

import (
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

const (
	minBytes       = 1024
	timeoutSeconds = 15
	readCapBytes   = 65536
	userAgent      = "Mozilla/5.0 (compatible; check_media_url.py; +deck-visuals)"
)

// checkResult mirrors the (ok, status, content_type, note) tuple returned by
// check_one() in check_media_url.py. hasStatus is false only when the request
// itself failed (DNS, connection refused, timeout, ...) before any response
// was received — the Python equivalent of status being None.
type checkResult struct {
	ok          bool
	hasStatus   bool
	status      int
	contentType string
	note        string
}

// isPassVerdict is the pure PASS/FAIL decision, split out from checkOneURL
// so it can be unit-tested without a network round trip: PASS iff status is
// 200 or 206, Content-Type (already stripped of any ";..." parameter) starts
// with "image/" or "video/", and size clears minBytes.
func isPassVerdict(status int, contentType string, size int) bool {
	isMediaType := strings.HasPrefix(contentType, "image/") || strings.HasPrefix(contentType, "video/")
	statusOk := status == 200 || status == 206
	sizeOk := size >= minBytes
	return statusOk && isMediaType && sizeOk
}

func checkOneURL(url string) checkResult {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return checkResult{ok: false, note: fmt.Sprintf("request failed: %v", err)}
	}
	req.Header.Set("User-Agent", userAgent)
	// Ask for a small range; not all hosts honor it, but when they do this
	// avoids pulling a whole video over the wire.
	req.Header.Set("Range", fmt.Sprintf("bytes=0-%d", readCapBytes-1))

	client := &http.Client{Timeout: timeoutSeconds * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return checkResult{ok: false, note: fmt.Sprintf("request failed: %v", err)}
	}
	defer resp.Body.Close()

	status := resp.StatusCode
	contentType := resp.Header.Get("Content-Type")
	if idx := strings.Index(contentType, ";"); idx >= 0 {
		contentType = contentType[:idx]
	}
	contentType = strings.TrimSpace(contentType)
	declaredLength := resp.Header.Get("Content-Length")

	// Read up to readCapBytes of the body — enough to confirm size past
	// minBytes without pulling a whole large file over the wire.
	buf := make([]byte, readCapBytes)
	nRead := 0
	for nRead < len(buf) {
		n, rerr := resp.Body.Read(buf[nRead:])
		nRead += n
		if rerr != nil {
			break
		}
	}

	// Prefer the declared Content-Length when present and larger than what
	// we actually read (e.g. a ranged/partial read on a big file).
	size := nRead
	if declaredLength != "" {
		if dl, derr := strconv.Atoi(declaredLength); derr == nil && dl > size {
			size = dl
		}
	}

	ok := isPassVerdict(status, contentType, size)

	ctOut := contentType
	if ctOut == "" {
		ctOut = "(none)"
	}

	return checkResult{
		ok:          ok,
		hasStatus:   true,
		status:      status,
		contentType: ctOut,
		note:        fmt.Sprintf("%d bytes read/declared", size),
	}
}

func runCheckMediaURL(args []string) int {
	fs := flag.NewFlagSet("check-media-url", flag.ContinueOnError)
	fs.SetOutput(io.Discard)
	if err := fs.Parse(args); err != nil {
		fmt.Fprintln(os.Stderr, "usage: dv-tools check-media-url <url> [<url> ...]")
		return 2
	}
	urls := fs.Args()
	if len(urls) == 0 {
		fmt.Fprintln(os.Stderr, "usage: dv-tools check-media-url <url> [<url> ...]")
		return 2
	}

	anyFail := false
	for _, url := range urls {
		res := checkOneURL(url)
		verdict := "PASS"
		if !res.ok {
			verdict = "FAIL"
			anyFail = true
		}
		statusStr := "n/a"
		if res.hasStatus {
			statusStr = strconv.Itoa(res.status)
		}
		fmt.Printf("[%s] %s\n", verdict, url)
		fmt.Printf("        status=%s content-type=%s %s\n", statusStr, res.contentType, res.note)
	}

	if anyFail {
		return 1
	}
	return 0
}
