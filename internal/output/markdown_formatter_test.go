package output_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/user/envdiff/internal/diff"
	"github.com/user/envdiff/internal/output"
)

func TestMarkdownFormatter_EmptyResults(t *testing.T) {
	f, err := output.NewFormatter("markdown")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var buf bytes.Buffer
	if err := f.Write(&buf, []diff.Result{}); err != nil {
		t.Fatalf("Write failed: %v", err)
	}

	got := buf.String()
	if !strings.Contains(got, "No differences found") {
		t.Errorf("expected empty message, got: %q", got)
	}
}

func TestMarkdownFormatter_Write(t *testing.T) {
	f, err := output.NewFormatter("markdown")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	results := []diff.Result{
		{
			Key:    "DB_HOST",
			Status: diff.StatusMatch,
			Values: map[string]string{"prod": "localhost", "dev": "localhost"},
		},
		{
			Key:    "API_KEY",
			Status: diff.StatusMismatch,
			Values: map[string]string{"prod": "abc123", "dev": "xyz789"},
		},
		{
			Key:    "SECRET",
			Status: diff.StatusMissing,
			Values: map[string]string{"prod": "topsecret"},
		},
	}

	var buf bytes.Buffer
	if err := f.Write(&buf, results); err != nil {
		t.Fatalf("Write failed: %v", err)
	}

	got := buf.String()

	if !strings.Contains(got, "| Key |") {
		t.Errorf("expected header row, got:\n%s", got)
	}
	if !strings.Contains(got, "DB_HOST") {
		t.Errorf("expected DB_HOST in output, got:\n%s", got)
	}
	if !strings.Contains(got, "API_KEY") {
		t.Errorf("expected API_KEY in output, got:\n%s", got)
	}
	if !strings.Contains(got, "_missing_") {
		t.Errorf("expected _missing_ placeholder, got:\n%s", got)
	}
	if !strings.Contains(got, "---") {
		t.Errorf("expected separator row, got:\n%s", got)
	}

	lines := strings.Split(strings.TrimSpace(got), "\n")
	// header + separator + 3 data rows = 5 lines
	if len(lines) != 5 {
		t.Errorf("expected 5 lines, got %d:\n%s", len(lines), got)
	}
}

func TestNewFormatter_MarkdownFormat(t *testing.T) {
	_, err := output.NewFormatter("markdown")
	if err != nil {
		t.Errorf("expected no error for 'markdown' format, got: %v", err)
	}
}
