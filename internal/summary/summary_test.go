package summary_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/user/envdiff/internal/diff"
	"github.com/user/envdiff/internal/summary"
)

func sampleResults() []diff.Result {
	return []diff.Result{
		{Key: "HOST", Status: diff.StatusMatch},
		{Key: "PORT", Status: diff.StatusMismatch},
		{Key: "SECRET", Status: diff.StatusMissing},
		{Key: "DEBUG", Status: diff.StatusMissing},
	}
}

func TestCompute_BasicCounts(t *testing.T) {
	s := summary.Compute(sampleResults())
	if s.TotalKeys != 4 {
		t.Errorf("expected 4 total keys, got %d", s.TotalKeys)
	}
	if s.Matched != 1 {
		t.Errorf("expected 1 matched, got %d", s.Matched)
	}
	if s.Mismatched != 1 {
		t.Errorf("expected 1 mismatched, got %d", s.Mismatched)
	}
	if s.Missing != 2 {
		t.Errorf("expected 2 missing, got %d", s.Missing)
	}
}

func TestComputeWithEnvs_SetsEnvCount(t *testing.T) {
	s := summary.ComputeWithEnvs(sampleResults(), 3)
	if s.EnvCount != 3 {
		t.Errorf("expected EnvCount 3, got %d", s.EnvCount)
	}
}

func TestHasIssues_True(t *testing.T) {
	s := summary.Compute(sampleResults())
	if !s.HasIssues() {
		t.Error("expected HasIssues to be true")
	}
}

func TestHasIssues_False(t *testing.T) {
	results := []diff.Result{
		{Key: "HOST", Status: diff.StatusMatch},
	}
	s := summary.Compute(results)
	if s.HasIssues() {
		t.Error("expected HasIssues to be false")
	}
}

func TestPrint_ContainsExpectedFields(t *testing.T) {
	s := summary.ComputeWithEnvs(sampleResults(), 2)
	var buf bytes.Buffer
	summary.Print(&buf, s)
	out := buf.String()
	for _, want := range []string{"Environments", "Total keys", "Matched", "Missing", "Mismatched"} {
		if !strings.Contains(out, want) {
			t.Errorf("expected output to contain %q", want)
		}
	}
}

func TestCompute_EmptyResults(t *testing.T) {
	s := summary.Compute(nil)
	if s.TotalKeys != 0 || s.Matched != 0 || s.Missing != 0 || s.Mismatched != 0 {
		t.Error("expected all zeros for empty results")
	}
}
