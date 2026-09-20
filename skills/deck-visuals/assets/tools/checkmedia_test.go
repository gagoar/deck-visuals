package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestIsPassVerdict(t *testing.T) {
	cases := []struct {
		name        string
		status      int
		contentType string
		size        int
		want        bool
	}{
		{"good image", 200, "image/gif", 13321, true},
		{"partial content honored range", 206, "video/mp4", 100000, true},
		{"404 status fails regardless of type/size", 404, "image/gif", 5000, false},
		{"html content-type fails", 200, "text/html", 5000, false},
		{"too small fails", 200, "image/gif", 35, false},
		{"exactly at floor passes", 200, "image/gif", 1024, true},
		{"one byte under floor fails", 200, "image/gif", 1023, false},
		{"video prefix passes", 200, "video/webm", 2048, true},
		{"empty content-type fails", 200, "", 5000, false},
		{"500 status fails", 500, "image/gif", 5000, false},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := isPassVerdict(tc.status, tc.contentType, tc.size)
			if got != tc.want {
				t.Errorf("isPassVerdict(%d, %q, %d) = %v, want %v", tc.status, tc.contentType, tc.size, got, tc.want)
			}
		})
	}
}

func TestCheckOneURL_GoodMedia(t *testing.T) {
	body := strings.Repeat("x", 2048)
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "image/gif")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(body))
	}))
	defer srv.Close()

	res := checkOneURL(srv.URL)
	if !res.ok {
		t.Fatalf("expected ok=true, got result=%+v", res)
	}
	if res.status != 200 {
		t.Errorf("expected status 200, got %d", res.status)
	}
	if res.contentType != "image/gif" {
		t.Errorf("expected content-type image/gif, got %q", res.contentType)
	}
}

func TestCheckOneURL_ContentTypeWithParameters(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body := strings.Repeat("y", 2048)
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(body))
	}))
	defer srv.Close()

	res := checkOneURL(srv.URL)
	if res.ok {
		t.Fatalf("expected ok=false for html, got result=%+v", res)
	}
	if res.contentType != "text/html" {
		t.Errorf("expected content-type trimmed to text/html, got %q", res.contentType)
	}
}

func TestCheckOneURL_NotFound(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.NotFound(w, r)
	}))
	defer srv.Close()

	res := checkOneURL(srv.URL + "/missing")
	if res.ok {
		t.Fatalf("expected ok=false for 404, got result=%+v", res)
	}
	if res.status != 404 {
		t.Errorf("expected status 404, got %d", res.status)
	}
}

func TestCheckOneURL_ConnectionFailure(t *testing.T) {
	// Port 0 on loopback with no listener guarantees a connection failure.
	res := checkOneURL("http://127.0.0.1:1/nope")
	if res.ok {
		t.Fatalf("expected ok=false for connection failure, got result=%+v", res)
	}
	if res.hasStatus {
		t.Errorf("expected hasStatus=false for connection failure, got true (status=%d)", res.status)
	}
	if !strings.HasPrefix(res.note, "request failed: ") {
		t.Errorf("expected note to start with %q, got %q", "request failed: ", res.note)
	}
}

func TestRunCheckMediaURL_NoArgsExitsTwo(t *testing.T) {
	if code := runCheckMediaURL(nil); code != 2 {
		t.Errorf("expected exit 2 for no args, got %d", code)
	}
}
