package recon

import (
	"strings"
	"testing"
)

func TestGenerateLinksAllModes(t *testing.T) {
	svc := NewGoogleLinksService()
	links, err := svc.GenerateLinks("example.com", "all")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(links) != 17 {
		t.Fatalf("expected 17 links, got %d", len(links))
	}
	if !strings.Contains(links[0], "example.com") {
		t.Fatalf("expected target replacement in first link")
	}
}

func TestGenerateLinksSpecificModes(t *testing.T) {
	svc := NewGoogleLinksService()
	links, err := svc.GenerateLinks("example.com", "1,5,5")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(links) != 2 {
		t.Fatalf("expected deduplicated links, got %d", len(links))
	}
}

func TestGenerateLinksInvalidMode(t *testing.T) {
	svc := NewGoogleLinksService()
	_, err := svc.GenerateLinks("example.com", "0")
	if err == nil {
		t.Fatalf("expected error for invalid mode")
	}
}
