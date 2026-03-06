package recon

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"slices"
	"testing"
)

func TestParseStatusCodes(t *testing.T) {
	set, err := ParseStatusCodes("200, 401")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !set[200] || !set[401] {
		t.Fatalf("expected 200 and 401 in set")
	}
}

func TestShouldIncludeStatus(t *testing.T) {
	show := map[int]bool{200: true}
	exclude := map[int]bool{401: true}
	if !ShouldIncludeStatus(200, show, exclude) {
		t.Fatalf("expected status 200 to be included by show set")
	}
	if ShouldIncludeStatus(401, show, exclude) {
		t.Fatalf("show set should have priority and exclude 401 when not present in show")
	}
}

func TestLoadSpecFromFile(t *testing.T) {
	scanner := NewAPIScanner(1, nil)
	dir := t.TempDir()
	specPath := filepath.Join(dir, "openapi.json")
	data := `{"paths":{"/health":{"get":{}}}}`
	if err := os.WriteFile(specPath, []byte(data), 0o644); err != nil {
		t.Fatalf("write spec file: %v", err)
	}

	spec, err := scanner.LoadSpec(specPath)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(spec.Paths) != 1 {
		t.Fatalf("expected 1 path, got %d", len(spec.Paths))
	}
}

func TestScanAgainstHTTPServer(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/health":
			if r.Method == http.MethodGet {
				w.WriteHeader(http.StatusOK)
				return
			}
		case "/items":
			if r.Method == http.MethodPost {
				w.WriteHeader(http.StatusCreated)
				return
			}
		}
		w.WriteHeader(http.StatusNotFound)
	}))
	defer server.Close()

	spec := OpenAPISpec{
		Paths: map[string]map[string]json.RawMessage{
			"/health": {"get": {}},
			"/items":  {"post": {}, "head": {}},
		},
	}
	scanner := NewAPIScanner(2, nil)
	results := scanner.Scan(server.URL, spec)
	if len(results) != 2 {
		t.Fatalf("expected 2 results for allowed methods, got %d", len(results))
	}

	statuses := []int{results[0].StatusCode, results[1].StatusCode}
	slices.Sort(statuses)
	if !slices.Equal(statuses, []int{200, 201}) {
		t.Fatalf("unexpected statuses: %v", statuses)
	}
}

func TestBuildAuthHeaders(t *testing.T) {
	headers, err := BuildAuthHeaders("jwt-token", "X-API-Key", "secret")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if headers["Authorization"] != "Bearer jwt-token" {
		t.Fatalf("unexpected Authorization header: %q", headers["Authorization"])
	}
	if headers["X-API-Key"] != "secret" {
		t.Fatalf("unexpected api key header value")
	}
}

func TestScanSendsAuthHeaders(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") != "Bearer jwt-token" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		if r.Header.Get("X-API-Key") != "secret" {
			w.WriteHeader(http.StatusForbidden)
			return
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	spec := OpenAPISpec{
		Paths: map[string]map[string]json.RawMessage{
			"/secure": {"get": {}},
		},
	}
	headers := map[string]string{
		"Authorization": "Bearer jwt-token",
		"X-API-Key":     "secret",
	}
	scanner := NewAPIScanner(2, headers)
	results := scanner.Scan(server.URL, spec)
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	if results[0].StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", results[0].StatusCode)
	}
}

func TestLoadSpecFromURLWithAuthHeaders(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") != "Bearer jwt-token" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"paths":{"/health":{"get":{}}}}`))
	}))
	defer server.Close()

	scanner := NewAPIScanner(2, map[string]string{"Authorization": "Bearer jwt-token"})
	spec, err := scanner.LoadSpec(server.URL)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(spec.Paths) != 1 {
		t.Fatalf("expected 1 path, got %d", len(spec.Paths))
	}
}

func TestScanResultOrderIsDeterministic(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	spec := OpenAPISpec{
		Paths: map[string]map[string]json.RawMessage{
			"/zeta":  {"post": {}, "get": {}},
			"/alpha": {"put": {}, "get": {}, "head": {}},
		},
	}
	scanner := NewAPIScanner(2, nil)
	results := scanner.Scan(server.URL, spec)
	if len(results) != 4 {
		t.Fatalf("expected 4 results, got %d", len(results))
	}

	got := []string{
		results[0].Method + " " + results[0].Path,
		results[1].Method + " " + results[1].Path,
		results[2].Method + " " + results[2].Path,
		results[3].Method + " " + results[3].Path,
	}
	want := []string{
		"GET /alpha",
		"PUT /alpha",
		"GET /zeta",
		"POST /zeta",
	}
	if !slices.Equal(got, want) {
		t.Fatalf("unexpected order: got=%v want=%v", got, want)
	}
}
