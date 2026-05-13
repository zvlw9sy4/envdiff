package redact_test

import (
	"testing"

	"github.com/yourorg/envdiff/internal/diff"
	"github.com/yourorg/envdiff/internal/redact"
)

func makeResults() []diff.Result {
	return []diff.Result{
		{
			Key:    "APP_ENV",
			Status: diff.StatusMatch,
			Values: map[string]string{"dev": "production", "staging": "production"},
		},
		{
			Key:    "DB_PASSWORD",
			Status: diff.StatusMismatch,
			Values: map[string]string{"dev": "devpass", "staging": "stagingpass"},
		},
		{
			Key:    "API_TOKEN",
			Status: diff.StatusMissing,
			Values: map[string]string{"dev": "tok123", "staging": ""},
		},
	}
}

func TestApplyToResults_RedactsSensitiveValues(t *testing.T) {
	r, _ := redact.New(nil)
	results := makeResults()
	out := r.ApplyToResults(results)

	for _, res := range out {
		if res.Key == "APP_ENV" {
			for env, val := range res.Values {
				if val != "production" {
					t.Errorf("APP_ENV[%s] should not be redacted", env)
				}
			}
		}
		if res.Key == "DB_PASSWORD" {
			for env, val := range res.Values {
				if val != redact.RedactedValue() {
					t.Errorf("DB_PASSWORD[%s] should be redacted, got %q", env, val)
				}
			}
		}
	}
}

func TestApplyToResults_PreservesEmptyMissingValues(t *testing.T) {
	r, _ := redact.New(nil)
	results := makeResults()
	out := r.ApplyToResults(results)

	for _, res := range out {
		if res.Key == "API_TOKEN" {
			if res.Values["staging"] != "" {
				t.Error("empty (missing) value should remain empty, not redacted")
			}
			if res.Values["dev"] != redact.RedactedValue() {
				t.Errorf("non-empty token value should be redacted")
			}
		}
	}
}

func TestApplyToResults_DoesNotMutateInput(t *testing.T) {
	r, _ := redact.New(nil)
	results := makeResults()
	_ = r.ApplyToResults(results)
	if results[1].Values["dev"] != "devpass" {
		t.Error("original results should not be mutated")
	}
}

func TestApplyToResults_StatusUnchanged(t *testing.T) {
	r, _ := redact.New(nil)
	results := makeResults()
	out := r.ApplyToResults(results)
	for i, res := range out {
		if res.Status != results[i].Status {
			t.Errorf("status changed for key %s", res.Key)
		}
	}
}
