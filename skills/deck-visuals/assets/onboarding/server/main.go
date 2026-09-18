// dv-onboard — the runtime-free local server for the deck-visuals onboarding picker.
//
// It serves one page (the embedded picker.html) on 127.0.0.1, accepts one POST of
// the presenter's answers, writes them to a results JSON file, and shuts down.
// Nothing leaves the machine: it binds loopback only, makes no outbound calls, and
// needs no interpreter or runtime on the machine that runs it (a single static Go
// binary). Only the server is a binary; media sourcing and coverage stay Python and
// run where Claude runs, never on the client.
//
// picker.html is compiled in via go:embed, so the shipped binary needs no companion
// file. build.sh copies the canonical assets/onboarding/picker.html next to this
// source before building (embed cannot reach a parent directory).
//
// Usage:
//   dv-onboard --cards <cards.json> --out <results.json> [--port 0] [--timeout 1800]
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
var pickerHTML string

const (
	cardsToken     = "REPLACE_CARDS_JSON"
	maxBodyBytes   = 256 * 1024
	maxSuggestions = 100
)

type suggestion struct {
	Title string `json:"title"`
	Note  string `json:"note,omitempty"`
}

type results struct {
	SubmittedAt string `json:"submitted_at"`
	Talk        struct {
		DurationMinutes *int `json:"duration_minutes"`
	} `json:"talk"`
	Audience struct {
		Description string `json:"description"`
	} `json:"audience"`
	RecognizedIDs    []string     `json:"recognized_ids"`
	NotRecognizedIDs []string     `json:"not_recognized_ids"`
	SkippedIDs       []string     `json:"skipped_ids"`
	Suggestions      []suggestion `json:"suggestions"`
}

func main() {
	cardsPath := flag.String("cards", "", "path to the cards.json served into the picker (required)")
	outPath := flag.String("out", "", "path to write results.json (required)")
	port := flag.Int("port", 0, "TCP port on 127.0.0.1; 0 lets the OS choose")
	timeout := flag.Int("timeout", 1800, "seconds to wait for a submit before giving up")
	flag.Parse()

	if *cardsPath == "" || *outPath == "" {
		fmt.Fprintln(os.Stderr, "error: --cards and --out are required")
		flag.Usage()
		os.Exit(2)
	}

	cardsBytes, err := os.ReadFile(*cardsPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: cannot read --cards %q: %v\n", *cardsPath, err)
		os.Exit(2)
	}
	// Fail fast on a malformed cards file rather than shipping a broken page.
	if !json.Valid(cardsBytes) {
		fmt.Fprintf(os.Stderr, "error: --cards %q is not valid JSON\n", *cardsPath)
		os.Exit(2)
	}
	page := strings.ReplaceAll(pickerHTML, cardsToken, string(cardsBytes))

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
	mux.HandleFunc("/submit", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		r.Body = http.MaxBytesReader(w, r.Body, maxBodyBytes)
		var res results
		dec := json.NewDecoder(r.Body)
		dec.DisallowUnknownFields()
		if err := dec.Decode(&res); err != nil {
			http.Error(w, "bad request", http.StatusBadRequest)
			return
		}
		if len(res.Suggestions) > maxSuggestions {
			res.Suggestions = res.Suggestions[:maxSuggestions]
		}
		res.SubmittedAt = time.Now().UTC().Format(time.RFC3339)

		if err := writeAtomic(*outPath, res); err != nil {
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

	// Watchdog: give up if nobody submits (an abandoned tab) so the caller isn't
	// left with a server running forever.
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

// writeAtomic marshals res and writes it via a temp file + rename so a reader never
// sees a half-written file.
func writeAtomic(path string, res results) error {
	data, err := json.MarshalIndent(res, "", "  ")
	if err != nil {
		return err
	}
	tmp := path + ".tmp"
	if err := os.WriteFile(tmp, data, 0o644); err != nil {
		return err
	}
	return os.Rename(tmp, path)
}
