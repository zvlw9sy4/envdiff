package lint_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/user/envdiff/internal/lint"
)

func TestPrintReport_NoIssues(t *testing.T) {
	var buf bytes.Buffer
	lint.PrintReport(&buf, nil)
	if !strings.Contains(buf.String(), "No lint issues") {
		t.Errorf("expected no-issues message, got: %s", buf.String())
	}
}

func TestPrintReport_WithIssues(t *testing.T) {
	issues := []lint.Issue{
		{Key: "db_host", Message: "key should be uppercase", Severity: lint.SeverityWarn},
		{Key: "BAD KEY", Message: "key must not contain spaces", Severity: lint.SeverityError},
	}
	var buf bytes.Buffer
	lint.PrintReport(&buf, issues)
	out := buf.String()

	if !strings.Contains(out, "ERROR") {
		t.Errorf("expected ERROR in output")
	}
	if !strings.Contains(out, "WARN") {
		t.Errorf("expected WARN in output")
	}
	if !strings.Contains(out, "2 issue(s)") {
		t.Errorf("expected issue count summary, got: %s", out)
	}
}

func TestSummarize_Counts(t *testing.T) {
	issues := []lint.Issue{
		{Severity: lint.SeverityError},
		{Severity: lint.SeverityWarn},
		{Severity: lint.SeverityWarn},
	}
	s := lint.Summarize(issues)
	if s.Total != 3 {
		t.Errorf("expected total 3, got %d", s.Total)
	}
	if s.Errors != 1 {
		t.Errorf("expected 1 error, got %d", s.Errors)
	}
	if s.Warns != 2 {
		t.Errorf("expected 2 warns, got %d", s.Warns)
	}
}

func TestSummarize_Empty(t *testing.T) {
	s := lint.Summarize(nil)
	if s.Total != 0 || s.Errors != 0 || s.Warns != 0 {
		t.Errorf("expected all zeros for empty input")
	}
}

func TestPrintReport_SortedOutput(t *testing.T) {
	issues := []lint.Issue{
		{Key: "Z_KEY", Message: "z issue", Severity: lint.SeverityWarn},
		{Key: "A_KEY", Message: "a issue", Severity: lint.SeverityWarn},
	}
	var buf bytes.Buffer
	lint.PrintReport(&buf, issues)
	out := buf.String()
	aIdx := strings.Index(out, "A_KEY")
	zIdx := strings.Index(out, "Z_KEY")
	if aIdx > zIdx {
		t.Errorf("expected A_KEY before Z_KEY in sorted output")
	}
}
