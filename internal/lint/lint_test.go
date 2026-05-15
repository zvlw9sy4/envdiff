package lint_test

import (
	"strings"
	"testing"

	"github.com/user/envdiff/internal/lint"
)

func TestEmptyValueRule_Warn(t *testing.T) {
	issues := lint.EmptyValueRule("DB_HOST", "")
	if len(issues) != 1 {
		t.Fatalf("expected 1 issue, got %d", len(issues))
	}
	if issues[0].Severity != lint.SeverityWarn {
		t.Errorf("expected warn severity")
	}
}

func TestEmptyValueRule_NoIssue(t *testing.T) {
	issues := lint.EmptyValueRule("DB_HOST", "localhost")
	if len(issues) != 0 {
		t.Errorf("expected no issues, got %d", len(issues))
	}
}

func TestUpperCaseKeyRule_Warn(t *testing.T) {
	issues := lint.UpperCaseKeyRule("db_host", "localhost")
	if len(issues) != 1 {
		t.Fatalf("expected 1 issue, got %d", len(issues))
	}
	if issues[0].Severity != lint.SeverityWarn {
		t.Errorf("expected warn severity")
	}
}

func TestUpperCaseKeyRule_NoIssue(t *testing.T) {
	issues := lint.UpperCaseKeyRule("DB_HOST", "localhost")
	if len(issues) != 0 {
		t.Errorf("expected no issues for uppercase key")
	}
}

func TestNoSpaceInKeyRule_Error(t *testing.T) {
	issues := lint.NoSpaceInKeyRule("DB HOST", "localhost")
	if len(issues) != 1 {
		t.Fatalf("expected 1 issue, got %d", len(issues))
	}
	if issues[0].Severity != lint.SeverityError {
		t.Errorf("expected error severity")
	}
}

func TestRun_MultipleRules(t *testing.T) {
	env := map[string]string{
		"db_host": "",
		"PORT":    "8080",
	}
	issues := lint.Run(env, lint.DefaultRules())
	// db_host triggers empty + uppercase; PORT is clean
	if len(issues) < 2 {
		t.Errorf("expected at least 2 issues, got %d", len(issues))
	}
}

func TestRun_EmptyEnv(t *testing.T) {
	issues := lint.Run(map[string]string{}, lint.DefaultRules())
	if len(issues) != 0 {
		t.Errorf("expected no issues for empty env, got %d", len(issues))
	}
}

func TestIssue_String(t *testing.T) {
	i := lint.Issue{Key: "FOO", Message: "test message", Severity: lint.SeverityWarn}
	s := i.String()
	if s == "" {
		t.Error("expected non-empty string representation")
	}
}

func TestIssue_String_ContainsKeyAndMessage(t *testing.T) {
	i := lint.Issue{Key: "MY_VAR", Message: "something is wrong", Severity: lint.SeverityError}
	s := i.String()
	if !strings.Contains(s, "MY_VAR") {
		t.Errorf("expected string to contain key %q, got %q", "MY_VAR", s)
	}
	if !strings.Contains(s, "something is wrong") {
		t.Errorf("expected string to contain message %q, got %q", "something is wrong", s)
	}
}
