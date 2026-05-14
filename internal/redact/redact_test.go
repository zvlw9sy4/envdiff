package redact_test

import (
	"testing"

	"github.com/yourorg/envdiff/internal/redact"
)

func TestNew_DefaultPatterns(t *testing.T) {
	r, err := redact.New(nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if r == nil {
		t.Fatal("expected non-nil Redactor")
	}
}

func TestNew_InvalidPattern(t *testing.T) {
	_, err := redact.New([]string{"[invalid"})
	if err == nil {
		t.Fatal("expected error for invalid pattern")
	}
}

func TestIsSensitive_MatchesPassword(t *testing.T) {
	r, _ := redact.New(nil)
	if !r.IsSensitive("DB_PASSWORD") {
		t.Error("expected DB_PASSWORD to be sensitive")
	}
}

func TestIsSensitive_MatchesToken(t *testing.T) {
	r, _ := redact.New(nil)
	if !r.IsSensitive("GITHUB_TOKEN") {
		t.Error("expected GITHUB_TOKEN to be sensitive")
	}
}

func TestIsSensitive_NonSensitiveKey(t *testing.T) {
	r, _ := redact.New(nil)
	if r.IsSensitive("APP_PORT") {
		t.Error("expected APP_PORT to NOT be sensitive")
	}
}

func TestRedactValue_SensitiveKey(t *testing.T) {
	r, _ := redact.New(nil)
	got := r.RedactValue("API_KEY", "super-secret-123")
	if got != redact.RedactedValue() {
		t.Errorf("expected redacted placeholder, got %q", got)
	}
}

func TestRedactValue_SafeKey(t *testing.T) {
	r, _ := redact.New(nil)
	got := r.RedactValue("APP_ENV", "production")
	if got != "production" {
		t.Errorf("expected original value, got %q", got)
	}
}

func TestApplyToMap_RedactsSensitiveKeys(t *testing.T) {
	r, _ := redact.New(nil)
	env := map[string]string{
		"APP_ENV":     "production",
		"DB_PASSWORD": "hunter2",
		"SECRET_KEY":  "abc123",
	}
	out := r.ApplyToMap(env)
	if out["APP_ENV"] != "production" {
		t.Errorf("APP_ENV should not be redacted")
	}
	if out["DB_PASSWORD"] != redact.RedactedValue() {
		t.Errorf("DB_PASSWORD should be redacted")
	}
	if out["SECRET_KEY"] != redact.RedactedValue() {
		t.Errorf("SECRET_KEY should be redacted")
	}
}

func TestApplyToMap_DoesNotMutateOriginal(t *testing.T) {
	r, _ := redact.New(nil)
	env := map[string]string{"DB_PASSWORD": "hunter2"}
	_ = r.ApplyToMap(env)
	if env["DB_PASSWORD"] != "hunter2" {
		t.Error("original map should not be mutated")
	}
}

func TestIsSensitive_CaseInsensitiveMatch(t *testing.T) {
	tests := []struct {
		key      string
		want     bool
	}{
		{"db_password", true},
		{"Db_Password", true},
		{"DB_PASSWORD", true},
		{"app_port", false},
	}
	r, _ := redact.New(nil)
	for _, tt := range tests {
		got := r.IsSensitive(tt.key)
		if got != tt.want {
			t.Errorf("IsSensitive(%q) = %v, want %v", tt.key, got, tt.want)
		}
	}
}
