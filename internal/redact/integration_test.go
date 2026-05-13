package redact_test

import (
	"testing"

	"github.com/yourorg/envdiff/internal/diff"
	"github.com/yourorg/envdiff/internal/redact"
)

// TestRedact_FullPipeline simulates loading env maps, comparing them, then
// redacting sensitive values before output.
func TestRedact_FullPipeline(t *testing.T) {
	envs := map[string]map[string]string{
		"dev": {
			"APP_NAME":    "myapp",
			"DB_PASSWORD": "devpass",
			"API_KEY":     "dev-key-abc",
		},
		"prod": {
			"APP_NAME":    "myapp",
			"DB_PASSWORD": "prodpass",
			"API_KEY":     "prod-key-xyz",
		},
	}

	results := diff.Compare(envs)

	r, err := redact.New(nil)
	if err != nil {
		t.Fatalf("failed to create redactor: %v", err)
	}

	redacted := r.ApplyToResults(results)

	for _, res := range redacted {
		switch res.Key {
		case "APP_NAME":
			for env, val := range res.Values {
				if val != "myapp" {
					t.Errorf("APP_NAME[%s] should not be redacted, got %q", env, val)
				}
			}
		case "DB_PASSWORD", "API_KEY":
			for env, val := range res.Values {
				if val != redact.RedactedValue() {
					t.Errorf("%s[%s] should be redacted, got %q", res.Key, env, val)
				}
			}
		}
	}
}

func TestRedact_CustomPatterns(t *testing.T) {
	r, err := redact.New([]string{"(?i)internal"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	env := map[string]string{
		"INTERNAL_URL": "http://internal.svc",
		"PUBLIC_URL":   "http://example.com",
	}

	out := r.ApplyToMap(env)
	if out["INTERNAL_URL"] != redact.RedactedValue() {
		t.Error("INTERNAL_URL should be redacted")
	}
	if out["PUBLIC_URL"] != "http://example.com" {
		t.Error("PUBLIC_URL should not be redacted")
	}
}
