package command

import (
	"slices"
	"testing"
)

func TestCollectSubDomainsDeduplicatesWordlistAndResults(t *testing.T) {
	wordlist := []string{"api", "api", "dev"}
	resolver := func(host string) bool {
		return host == "api.example.com" || host == "dev.example.com"
	}

	got := collectSubDomains("example.com", wordlist, false, 0, 4, resolver)
	want := []string{"api.example.com", "dev.example.com"}
	if !slices.Equal(got, want) {
		t.Fatalf("collectSubDomains() = %v, want %v", got, want)
	}
}

func TestCollectSubDomainsRespectsMaxDepth(t *testing.T) {
	wordlist := []string{"a"}
	resolver := func(host string) bool {
		return host == "a.example.com" || host == "a.a.example.com"
	}

	withoutRecursion := collectSubDomains("example.com", wordlist, true, 0, 2, resolver)
	wantWithoutRecursion := []string{"a.example.com"}
	if !slices.Equal(withoutRecursion, wantWithoutRecursion) {
		t.Fatalf("max-depth=0 result = %v, want %v", withoutRecursion, wantWithoutRecursion)
	}

	withDepthOne := collectSubDomains("example.com", wordlist, true, 1, 2, resolver)
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

	got := collectSubDomains("example.com", wordlist, false, 0, 0, resolver)
	want := []string{"api.example.com", "dev.example.com"}
	if !slices.Equal(got, want) {
		t.Fatalf("collectSubDomains() = %v, want %v", got, want)
	}
}
