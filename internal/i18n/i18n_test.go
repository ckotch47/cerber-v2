package i18n

import "testing"

func TestSetLangAndTranslate(t *testing.T) {
	if err := SetLang("en"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got := T("msg_not_found"); got != "Not found" {
		t.Fatalf("unexpected translation: %q", got)
	}

	if err := SetLang("ru"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got := T("msg_not_found"); got != "Не найдено" {
		t.Fatalf("unexpected translation: %q", got)
	}
}

func TestSetLangUnsupported(t *testing.T) {
	if err := SetLang("de"); err == nil {
		t.Fatalf("expected error for unsupported language")
	}
}

func TestDetectFromLocale(t *testing.T) {
	if got := detectFromLocale("ru_RU.UTF-8"); got != "ru" {
		t.Fatalf("expected ru, got %q", got)
	}
	if got := detectFromLocale("en_US.UTF-8"); got != "en" {
		t.Fatalf("expected en, got %q", got)
	}
	if got := detectFromLocale("de_DE.UTF-8"); got != "" {
		t.Fatalf("expected empty for unsupported locale, got %q", got)
	}
}
