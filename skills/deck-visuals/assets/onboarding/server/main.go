// dv-onboard — the runtime-free local server behind the deck-visuals pickers.
//
// One static Go binary drives both web pickers this plugin uses:
//
//	--form onboard  the show-recognition swipe deck (rebuilds the profile)
//	--form levity   choose one of 2-3 candidate GIFs per slide, at deck-build time
//
// It serves the chosen page (and one shared stylesheet) on 127.0.0.1, accepts one
// POST of the user's answers, writes them to a results JSON file, and shuts down.
// Nothing leaves the machine: loopback only, no outbound calls, no interpreter or
// runtime needed on the machine that runs it. The picker HTML/CSS are compiled in
// via go:embed, so the shipped binary needs no companion files. build.sh copies the
// canonical assets next to this source before building (embed can't reach a parent
// directory).
//
// Usage:
//
//	dv-onboard --form onboard --data <cards.json>  --out <results.json> [--port 0] [--timeout 1800]
//	dv-onboard --form levity  --data <slots.json>  --out <results.json> [--port 0] [--timeout 1800]
package main

import (
	"context"
	_ "embed"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"net"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"
)

//go:embed picker.html
var onboardHTML string

//go:embed levity-picker.html
var levityHTML string

//go:embed picker.css
var pickerCSS string

const (
	dataToken    = "REPLACE_DATA_JSON"
	maxBodyBytes = 512 * 1024
)

func main() {
	form := flag.String("form", "onboard", "which picker to serve: onboard | levity")
	dataPath := flag.String("data", "", "path to the JSON injected into the page (required)")
	outPath := flag.String("out", "", "path to write results.json (required)")
	port := flag.Int("port", 0, "TCP port on 127.0.0.1; 0 lets the OS choose")
	timeout := flag.Int("timeout", 1800, "seconds to wait for a submit before giving up")
	flag.Parse()

	var tmpl string
	switch *form {
	case "onboard":
		tmpl = onboardHTML
	case "levity":
		tmpl = levityHTML
	default:
		fmt.Fprintf(os.Stderr, "error: --form must be onboard or levity, got %q\n", *form)
		os.Exit(2)
	}
	if *dataPath == "" || *outPath == "" {
		fmt.Fprintln(os.Stderr, "error: --data and --out are required")
		flag.Usage()
		os.Exit(2)
	}

	dataBytes, err := os.ReadFile(*dataPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: cannot read --data %q: %v\n", *dataPath, err)
		os.Exit(2)
	}
	if !json.Valid(dataBytes) {
		fmt.Fprintf(os.Stderr, "error: --data %q is not valid JSON\n", *dataPath)
		os.Exit(2)
	}
	page := strings.ReplaceAll(tmpl, dataToken, string(dataBytes))

	// Result of the run: 0 = submitted, 2 = timed out / abandoned.
	exitCh := make(chan int, 1)
	var once sync.Once
	finish := func(code int) { once.Do(func() { exitCh <- code }) }

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		_, _ = w.Write([]byte(page))
	})
	mux.HandleFunc("/picker.css", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/css; charset=utf-8")
		_, _ = w.Write([]byte(pickerCSS))
	})
	mux.HandleFunc("/submit", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		r.Body = http.MaxBytesReader(w, r.Body, maxBodyBytes)
		// Accept any well-formed JSON object; the local page owns the shape. We only
		// stamp the server-authoritative submitted_at and persist it.
		var payload map[string]any
		dec := json.NewDecoder(r.Body)
		if err := dec.Decode(&payload); err != nil {
			http.Error(w, "bad request", http.StatusBadRequest)
			return
		}
		payload["submitted_at"] = time.Now().UTC().Format(time.RFC3339)

		if err := writeAtomic(*outPath, payload); err != nil {
			fmt.Fprintf(os.Stderr, "error: writing results: %v\n", err)
			http.Error(w, "could not save", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"ok":true}`))
		finish(0) // one submit is enough; schedule shutdown
	})

	ln, err := net.Listen("tcp", fmt.Sprintf("127.0.0.1:%d", *port))
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: cannot bind 127.0.0.1:%d: %v\n", *port, err)
		os.Exit(2)
	}
	fmt.Printf("http://%s/\n", ln.Addr().String())

	srv := &http.Server{Handler: mux}
	go func() {
		if serveErr := srv.Serve(ln); serveErr != nil && !errors.Is(serveErr, http.ErrServerClosed) {
			fmt.Fprintf(os.Stderr, "error: serve: %v\n", serveErr)
			finish(2)
		}
	}()

	// Watchdog: give up if nobody submits (an abandoned tab).
	timer := time.AfterFunc(time.Duration(*timeout)*time.Second, func() {
		fmt.Fprintln(os.Stderr, "timed out waiting for a submit")
		finish(2)
	})
	defer timer.Stop()

	code := <-exitCh
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_ = srv.Shutdown(shutdownCtx)
	os.Exit(code)
}

// writeAtomic marshals payload and writes it via a temp file + rename so a reader
// never sees a half-written file.
func writeAtomic(path string, payload any) error {
	data, err := json.MarshalIndent(payload, "", "  ")
	if err != nil {
		return err
	}
	tmp := path + ".tmp"
	if err := os.WriteFile(tmp, data, 0o644); err != nil {
		return err
	}
	return os.Rename(tmp, path)
}
