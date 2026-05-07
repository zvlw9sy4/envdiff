package filter_test

import (
	"testing"

	"github.com/user/envdiff/internal/diff"
	"github.com/user/envdiff/internal/filter"
)

func sampleResults() []diff.Result {
	return []diff.Result{
		{Key: "APP_ENV", Status: diff.StatusMatch, Values: map[string]string{"dev": "production", "prod": "production"}},
		{Key: "DB_HOST", Status: diff.StatusMismatch, Values: map[string]string{"dev": "localhost", "prod": "db.prod.example.com"}},
		{Key: "SECRET_KEY", Status: diff.StatusMissing, Values: map[string]string{"dev": "abc123", "prod": ""}},
	}
}

func TestApply_DefaultFilter(t *testing.T) {
	results := sampleResults()
	f := filter.DefaultFilter()
	out := filter.Apply(results, f)
	if len(out) != 2 {
		t.Fatalf("expected 2 results, got %d", len(out))
	}
	for _, r := range out {
		if r.Status == diff.StatusMatch {
			t.Errorf("default filter should exclude match status, got key %s", r.Key)
		}
	}
}

func TestApply_AllFilter(t *testing.T) {
	results := sampleResults()
	f := filter.AllFilter()
	out := filter.Apply(results, f)
	if len(out) != 3 {
		t.Fatalf("expected 3 results, got %d", len(out))
	}
}

func TestApply_OnlyMismatch(t *testing.T) {
	results := sampleResults()
	f := filter.StatusFilter{IncludeMismatch: true}
	out := filter.Apply(results, f)
	if len(out) != 1 || out[0].Key != "DB_HOST" {
		t.Fatalf("expected only DB_HOST, got %v", out)
	}
}

func TestByKey_FilterKeys(t *testing.T) {
	results := sampleResults()
	out := filter.ByKey(results, []string{"APP_ENV", "SECRET_KEY"})
	if len(out) != 2 {
		t.Fatalf("expected 2 results, got %d", len(out))
	}
}

func TestByKey_EmptyKeys(t *testing.T) {
	results := sampleResults()
	out := filter.ByKey(results, nil)
	if len(out) != len(results) {
		t.Fatalf("expected all results when keys is empty, got %d", len(out))
	}
}

func TestByKey_NoMatch(t *testing.T) {
	results := sampleResults()
	out := filter.ByKey(results, []string{"NONEXISTENT"})
	if len(out) != 0 {
		t.Fatalf("expected 0 results, got %d", len(out))
	}
}
