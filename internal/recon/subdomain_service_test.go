package recon

import (
	"slices"
	"testing"
)

func TestCollectSubDomainsDeduplicatesWordlistAndResults(t *testing.T) {
	wordlist := []string{"api", "api", "dev"}
	resolver := func(host string) bool {
		return host == "api.example.com" || host == "dev.example.com"
	}

	scanner := NewSubdomainScanner(resolver)
	got := scanner.Collect("example.com", wordlist, false, 0, 4)
	want := []string{"api.example.com", "dev.example.com"}
	if !slices.Equal(got, want) {
		t.Fatalf("Collect() = %v, want %v", got, want)
	}
}

func TestCollectSubDomainsRespectsMaxDepth(t *testing.T) {
	wordlist := []string{"a"}
	resolver := func(host string) bool {
		return host == "a.example.com" || host == "a.a.example.com"
	}

	scanner := NewSubdomainScanner(resolver)
	withoutRecursion := scanner.Collect("example.com", wordlist, true, 0, 2)
	wantWithoutRecursion := []string{"a.example.com"}
	if !slices.Equal(withoutRecursion, wantWithoutRecursion) {
		t.Fatalf("max-depth=0 result = %v, want %v", withoutRecursion, wantWithoutRecursion)
	}

	withDepthOne := scanner.Collect("example.com", wordlist, true, 1, 2)
	wantWithDepthOne := []string{"a.a.example.com", "a.example.com"}
	if !slices.Equal(withDepthOne, wantWithDepthOne) {
		t.Fatalf("max-depth=1 result = %v, want %v", withDepthOne, wantWithDepthOne)
	}
}

func TestCollectSubDomainsNormalizesConcurrency(t *testing.T) {
	wordlist := []string{"api", "dev"}
	resolver := func(host string) bool {
		return host == "api.example.com" || host == "dev.example.com"
	}

	scanner := NewSubdomainScanner(resolver)
	got := scanner.Collect("example.com", wordlist, false, 0, 0)
	want := []string{"api.example.com", "dev.example.com"}
	if !slices.Equal(got, want) {
		t.Fatalf("Collect() = %v, want %v", got, want)
	}
}
