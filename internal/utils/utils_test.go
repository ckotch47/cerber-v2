package utils

import (
	"os"
	"path/filepath"
	"testing"
)

func TestReadFileSuccess(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "wordlist.txt")
	if err := os.WriteFile(path, []byte("app\n#comment\nportal\n"), 0o644); err != nil {
		t.Fatalf("write temp file: %v", err)
	}

	lines, err := ReadFile(path)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(lines) != 2 || lines[0] != "app" || lines[1] != "portal" {
		t.Fatalf("unexpected lines: %v", lines)
	}
}

func TestReadFileReturnsErrorOnEmptyFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "empty.txt")
	if err := os.WriteFile(path, []byte("\n#only-comment\n"), 0o644); err != nil {
		t.Fatalf("write temp file: %v", err)
	}

	_, err := ReadFile(path)
	if err == nil {
		t.Fatalf("expected error for empty file")
	}
}

func TestReadFileReturnsErrorWhenNotFound(t *testing.T) {
	_, err := ReadFile("/tmp/definitely-missing-cerber-wordlist.txt")
	if err == nil {
		t.Fatalf("expected error for missing file")
	}
}
