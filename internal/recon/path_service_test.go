package recon

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestNormalizeBaseURL(t *testing.T) {
	tests := []struct {
		in   string
		want string
	}{
		{in: "example.com", want: "https://example.com"},
		{in: "example.com/", want: "https://example.com"},
		{in: "http://example.com/", want: "http://example.com"},
		{in: "https://example.com", want: "https://example.com"},
	}

	for _, tt := range tests {
		got := NormalizeBaseURL(tt.in)
		if got != tt.want {
			t.Fatalf("NormalizeBaseURL(%q) = %q, want %q", tt.in, got, tt.want)
		}
	}
}

func TestJoinURL(t *testing.T) {
	got := JoinURL("https://example.com/", "/admin")
	want := "https://example.com/admin"
	if got != want {
		t.Fatalf("JoinURL() = %q, want %q", got, want)
	}
}

func TestHasHTTPPrefixCaseInsensitive(t *testing.T) {
	if !HasHTTPPrefix("HTTPS://example.com") {
		t.Fatalf("expected HTTPS:// prefix to be recognized")
	}
}

func TestPathScannerFallbackToHTTP(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = io.WriteString(w, "ok")
	}))
	defer server.Close()

	scanner := NewPathScanner(1, 0, nil, true)
	baseURL := strings.Replace(server.URL, "http://", "https://", 1)
	res := scanner.Scan(baseURL, []string{"admin"})

	if len(res) != 1 {
		t.Fatalf("expected 1 result, got %d", len(res))
	}
	if res[0].Err != nil {
		t.Fatalf("unexpected error: %v", res[0].Err)
	}
	if res[0].StatusCode != http.StatusOK {
		t.Fatalf("status = %d, want 200", res[0].StatusCode)
	}
	if !strings.HasPrefix(res[0].URL, "http://") {
		t.Fatalf("expected fallback url to use http, got %q", res[0].URL)
	}
}

func TestPathScannerAppliesDelayOnErrors(t *testing.T) {
	scanner := NewPathScanner(1, 1, nil, false)

	start := time.Now()
	res := scanner.Scan("https://%", []string{"admin"})
	elapsed := time.Since(start)

	if len(res) != 1 {
		t.Fatalf("expected 1 result, got %d", len(res))
	}
	if res[0].Err == nil {
		t.Fatalf("expected an error result")
	}
	if elapsed < time.Second {
		t.Fatalf("expected delay to apply on error, elapsed=%v", elapsed)
	}
}
