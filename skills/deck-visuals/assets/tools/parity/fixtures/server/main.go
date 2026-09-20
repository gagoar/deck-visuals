// Command fixtureserver is a tiny, deterministic static HTTP server used only
// by the deck-visuals parity harness (skills/deck-visuals/assets/tools/parity).
//
// It exists so the check-media-url parity cases don't depend on the real
// network (flaky, non-deterministic, and outside this repo's control) or on
// Python (which will eventually be removed once the Go port lands). It is
// intentionally its own Go module so it has no dependency on the future
// dv-tools module.
//
// It binds to 127.0.0.1 on a random free port, prints that port to stdout as
// "PORT=<n>\n" as soon as it is ready to accept connections, and then serves
// forever until killed by its caller (run.sh).
//
// Routes:
//
//	/good.gif   - a valid GIF, >= ~2KB, Content-Type: image/gif  (PASS case)
//	/tiny.gif   - a valid GIF, < 1024 bytes, Content-Type: image/gif (FAIL: size)
//	/page.html  - an HTML document, 200 OK (FAIL: content-type)
//	(anything else, e.g. /missing) - 404 (FAIL: status)
package main

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"image/gif"
	"log"
	"math/rand"
	"net"
	"net/http"
	"os"
	"strconv"
)

// noisePalette returns a 256-entry palette with no runs of identical
// neighboring colors, so a noise image encodes with no cheap LZW
// compression and reliably lands well above the 2KB floor.
func noisePalette() color.Palette {
	pal := make(color.Palette, 256)
	for i := 0; i < 256; i++ {
		pal[i] = color.RGBA{
			R: uint8(i),
			G: uint8((i * 3) % 256),
			B: uint8((i * 7) % 256),
			A: 255,
		}
	}
	return pal
}

// makeNoiseGIF renders a w×h image of pseudo-random pixel indices (seeded,
// so the fixture is reproducible run to run) and GIF-encodes it. Random
// noise is close to incompressible under LZW, so the encoded size scales
// predictably with pixel count.
func makeNoiseGIF(w, h int, seed int64) []byte {
	r := rand.New(rand.NewSource(seed))
	pal := noisePalette()
	img := image.NewPaletted(image.Rect(0, 0, w, h), pal)
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			img.SetColorIndex(x, y, uint8(r.Intn(256)))
		}
	}
	var buf bytes.Buffer
	if err := gif.Encode(&buf, img, nil); err != nil {
		log.Fatalf("encode good.gif fixture: %v", err)
	}
	return buf.Bytes()
}

// makeTinyGIF renders a minimal 2x2 GIF — a genuinely valid image, just a
// tiny one — comfortably under the 1024-byte MIN_BYTES floor.
func makeTinyGIF() []byte {
	pal := color.Palette{
		color.RGBA{255, 0, 0, 255},
		color.RGBA{0, 255, 0, 255},
	}
	img := image.NewPaletted(image.Rect(0, 0, 2, 2), pal)
	var buf bytes.Buffer
	if err := gif.Encode(&buf, img, nil); err != nil {
		log.Fatalf("encode tiny.gif fixture: %v", err)
	}
	return buf.Bytes()
}

const minGoodBytes = 2048

func buildGoodGIF() []byte {
	size := 96
	for {
		b := makeNoiseGIF(size, size, 42)
		if len(b) >= minGoodBytes {
			return b
		}
		size += 32 // grow and retry if the encoder ever surprises us
	}
}

func main() {
	goodGIF := buildGoodGIF()
	tinyGIF := makeTinyGIF()
	if len(tinyGIF) >= 1024 {
		log.Fatalf("tiny.gif fixture is %d bytes, expected < 1024", len(tinyGIF))
	}

	mux := http.NewServeMux()

	mux.HandleFunc("/good.gif", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "image/gif")
		w.Header().Set("Content-Length", strconv.Itoa(len(goodGIF)))
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(goodGIF)
	})

	mux.HandleFunc("/tiny.gif", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "image/gif")
		w.Header().Set("Content-Length", strconv.Itoa(len(tinyGIF)))
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(tinyGIF)
	})

	mux.HandleFunc("/page.html", func(w http.ResponseWriter, r *http.Request) {
		body := []byte("<!doctype html><html><body><h1>not media</h1></body></html>")
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Header().Set("Content-Length", strconv.Itoa(len(body)))
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(body)
	})

	// Everything else (e.g. /missing) falls through to ServeMux's default
	// 404 behavior, which is what the check-media-url "404 FAIL" case wants.

	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		log.Fatalf("listen: %v", err)
	}
	port := ln.Addr().(*net.TCPAddr).Port

	// The harness polls stdout for this line to learn the port; flush
	// immediately and unbuffered so the poll doesn't stall on it.
	fmt.Fprintf(os.Stdout, "PORT=%d\n", port)
	_ = os.Stdout.Sync()

	srv := &http.Server{Handler: mux}
	log.Fatal(srv.Serve(ln))
}
