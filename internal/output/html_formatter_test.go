package output_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/user/envdiff/internal/diff"
	"github.com/user/envdiff/internal/output"
)

func TestHTMLFormatter_EmptyResults(t *testing.T) {
	f, err := output.NewFormatter("html")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var buf bytes.Buffer
	if err := f.Write(&buf, []diff.Result{}); err != nil {
		t.Fatalf("Write error: %v", err)
	}

	out := buf.String()
	if !strings.Contains(out, "No results") {
		t.Errorf("expected 'No results' message, got:\n%s", out)
	}
	if !strings.Contains(out, "</html>") {
		t.Errorf("expected closing html tag, got:\n%s", out)
	}
}

func TestHTMLFormatter_Write(t *testing.T) {
	f, err := output.NewFormatter("html")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	results := []diff.Result{
		{
			Key:    "DB_HOST",
			Status: diff.StatusMatch,
			Values: map[string]string{"dev": "localhost", "prod": "localhost"},
		},
		{
			Key:    "API_KEY",
			Status: diff.StatusMismatch,
			Values: map[string]string{"dev": "abc", "prod": "xyz"},
		},
		{
			Key:    "SECRET",
			Status: diff.StatusMissing,
			Values: map[string]string{"dev": "s3cr3t", "prod": ""},
		},
	}

	var buf bytes.Buffer
	if err := f.Write(&buf, results); err != nil {
		t.Fatalf("Write error: %v", err)
	}

	out := buf.String()

	for _, want := range []string{
		"<!DOCTYPE html>",
		"envdiff Report",
		"DB_HOST",
		"API_KEY",
		"SECRET",
		"match",
		"mismatch",
		"missing",
		"</table>",
		"</html>",
	} {
		if !strings.Contains(out, want) {
			t.Errorf("expected output to contain %q", want)
		}
	}
}

func TestNewFormatter_HTMLFormat(t *testing.T) {
	f, err := output.NewFormatter("html")
	if err != nil {
		t.Fatalf("expected no error for 'html' format, got: %v", err)
	}
	if f == nil {
		t.Fatal("expected non-nil formatter")
	}
}
