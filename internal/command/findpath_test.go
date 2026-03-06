package command

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
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
		got := normalizeBaseURL(tt.in)
		if got != tt.want {
			t.Fatalf("normalizeBaseURL(%q) = %q, want %q", tt.in, got, tt.want)
		}
	}
}

func TestRequestWithFallbackToHTTP(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = io.WriteString(w, "ok")
	}))
	defer server.Close()

	httpsTarget := strings.Replace(server.URL, "http://", "https://", 1) + "/admin"
	resp, actualURL, err := requestWithFallback(&http.Client{}, httpsTarget, 1, true)
	if err != nil {
		t.Fatalf("requestWithFallback returned error: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("status = %d, want 200", resp.StatusCode)
	}
	if !strings.HasPrefix(actualURL, "http://") {
		t.Fatalf("expected fallback to http, got %q", actualURL)
	}
}

func TestJoinURL(t *testing.T) {
	got := joinURL("https://example.com/", "/admin")
	want := "https://example.com/admin"
	if got != want {
		t.Fatalf("joinURL() = %q, want %q", got, want)
	}
}
