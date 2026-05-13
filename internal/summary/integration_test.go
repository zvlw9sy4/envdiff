package summary_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/user/envdiff/internal/diff"
	"github.com/user/envdiff/internal/summary"
)

// TestSummary_FullPipeline simulates computing and printing a summary after a
// real Compare call with multiple environments.
func TestSummary_FullPipeline(t *testing.T) {
	envs := map[string]map[string]string{
		"dev": {"HOST": "localhost", "PORT": "8080", "DEBUG": "true"},
		"prod": {"HOST": "example.com", "PORT": "443"},
		"staging": {"HOST": "staging.example.com", "PORT": "8080", "SECRET": "abc"},
	}

	results := diff.Compare(envs)
	s := summary.ComputeWithEnvs(results, len(envs))

	if s.EnvCount != 3 {
		t.Errorf("expected 3 envs, got %d", s.EnvCount)
	}
	if s.TotalKeys == 0 {
		t.Error("expected non-zero total keys")
	}
	if !s.HasIssues() {
		t.Error("expected issues to be present")
	}

	var buf bytes.Buffer
	summary.Print(&buf, s)
	out := buf.String()
	if !strings.Contains(out, "3") {
		t.Error("expected summary output to mention the env count")
	}
}

func TestSummary_AllMatch(t *testing.T) {
	envs := map[string]map[string]string{
		"dev":  {"HOST": "localhost", "PORT": "8080"},
		"prod": {"HOST": "localhost", "PORT": "8080"},
	}

	results := diff.Compare(envs)
	s := summary.ComputeWithEnvs(results, len(envs))

	if s.HasIssues() {
		t.Errorf("expected no issues, got missing=%d mismatched=%d", s.Missing, s.Mismatched)
	}
	if s.Matched == 0 {
		t.Error("expected at least one matched key")
	}
}
