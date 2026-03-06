package command

import "testing"

func TestRootCommandDoesNotRequirePositionalArgs(t *testing.T) {
	if rootCmd.Args != nil {
		t.Fatalf("expected root command args validator to be nil, got non-nil")
	}
}

func TestFindCommandWordlistFlags(t *testing.T) {
	wordlist := findCmd.Flags().Lookup("wordlist")
	if wordlist == nil {
		t.Fatalf("expected --wordlist flag to be present")
	}
	if wordlist.Shorthand != "w" {
		t.Fatalf("expected --wordlist shorthand to be -w, got %q", wordlist.Shorthand)
	}

	legacy := findCmd.Flags().Lookup("worldlis")
	if legacy == nil {
		t.Fatalf("expected deprecated --worldlis flag to be present")
	}
	if legacy.Deprecated == "" {
		t.Fatalf("expected --worldlis to be marked deprecated")
	}
}

func TestFindPathCommandWordlistFlags(t *testing.T) {
	wordlist := findPathCmd.Flags().Lookup("wordlist")
	if wordlist == nil {
		t.Fatalf("expected --wordlist flag to be present")
	}
	if wordlist.Shorthand != "w" {
		t.Fatalf("expected --wordlist shorthand to be -w, got %q", wordlist.Shorthand)
	}

	legacy := findPathCmd.Flags().Lookup("worldlis")
	if legacy == nil {
		t.Fatalf("expected deprecated --worldlis flag to be present")
	}
	if legacy.Deprecated == "" {
		t.Fatalf("expected --worldlis to be marked deprecated")
	}
}
