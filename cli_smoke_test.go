package main

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
	"time"
)

var ansiPattern = regexp.MustCompile(`\x1b\[[0-9;]*m`)

func TestCLISmokeVersion(t *testing.T) {
	t.Parallel()

	out, err := runCerber(t, "--lang", "en", "version")
	if err != nil {
		t.Fatalf("version command failed: %v\noutput:\n%s", err, out)
	}
	if !strings.Contains(stripANSI(out), "Version: v0.0.1a") {
		t.Fatalf("unexpected version output:\n%s", out)
	}
}

func TestCLISmokeFindPath(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/admin" {
			w.WriteHeader(http.StatusOK)
			return
		}
		w.WriteHeader(http.StatusNotFound)
	}))
	defer srv.Close()

	wordlist := filepath.Join(t.TempDir(), "paths.txt")
	if err := os.WriteFile(wordlist, []byte("admin\n"), 0o600); err != nil {
		t.Fatalf("write wordlist: %v", err)
	}

	out, err := runCerber(
		t,
		"--lang", "en",
		"find", "path", srv.URL,
		"-w", wordlist,
		"-t", "0",
		"--request-timeout", "2",
	)
	if err != nil {
		t.Fatalf("find path command failed: %v\noutput:\n%s", err, out)
	}
	if !strings.Contains(stripANSI(out), "admin : 200") {
		t.Fatalf("expected successful admin result, got:\n%s", out)
	}
}

func TestCLISmokeAPIScan(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet && r.URL.Path == "/secure" {
			w.WriteHeader(http.StatusOK)
			return
		}
		w.WriteHeader(http.StatusNotFound)
	}))
	defer srv.Close()

	spec := `{"paths":{"/secure":{"get":{}}}}`
	specPath := filepath.Join(t.TempDir(), "openapi.json")
	if err := os.WriteFile(specPath, []byte(spec), 0o600); err != nil {
		t.Fatalf("write spec: %v", err)
	}

	out, err := runCerber(
		t,
		"--lang", "en",
		"api", "scan",
		"--spec", specPath,
		"--host", srv.URL,
		"--show", "200",
		"--request-timeout", "2",
	)
	if err != nil {
		t.Fatalf("api scan command failed: %v\noutput:\n%s", err, out)
	}
	if !strings.Contains(stripANSI(out), "[GET] /secure : 200") {
		t.Fatalf("expected scan result [GET] /secure : 200, got:\n%s", out)
	}
}

func runCerber(t *testing.T, args ...string) (string, error) {
	t.Helper()

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, "go", append([]string{"run", "."}, args...)...)
	cmd.Dir = "/Users/blant/GoLangProject/lessons/cerber"
	cmd.Env = append(os.Environ(), "GOCACHE=/tmp/.gocache")
	out, err := cmd.CombinedOutput()
	return string(out), err
}

func stripANSI(s string) string {
	return ansiPattern.ReplaceAllString(s, "")
}
