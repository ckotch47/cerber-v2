package command

import "testing"

func TestRootCommandDoesNotRequirePositionalArgs(t *testing.T) {
	if rootCmd.Args != nil {
		t.Fatalf("expected root command args validator to be nil, got non-nil")
	}
	if !rootCmd.SilenceUsage {
		t.Fatalf("expected root command to silence usage on errors")
	}
	if !rootCmd.SilenceErrors {
		t.Fatalf("expected root command to silence cobra's default error output")
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

	reqTimeout := findPathCmd.Flags().Lookup("request-timeout")
	if reqTimeout == nil {
		t.Fatalf("expected --request-timeout flag to be present")
	}
	if reqTimeout.Shorthand != "" {
		t.Fatalf("expected --request-timeout to have no shorthand, got %q", reqTimeout.Shorthand)
	}
}

func TestCleanDomain(t *testing.T) {
	got, err := cleanDomain("https://www.example.com/")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "example.com" {
		t.Fatalf("cleanDomain() = %q, want %q", got, "example.com")
	}
}

func TestCleanDomainReturnsErrorOnEmptyInput(t *testing.T) {
	_, err := cleanDomain("   ")
	if err == nil {
		t.Fatalf("expected error for empty domain input")
	}
}
