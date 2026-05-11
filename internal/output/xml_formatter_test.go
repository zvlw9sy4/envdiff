package output_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/user/envdiff/internal/diff"
	"github.com/user/envdiff/internal/output"
)

func TestXMLFormatter_EmptyResults(t *testing.T) {
	f, err := output.NewFormatter("xml")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var buf bytes.Buffer
	if err := f.Write(&buf, []diff.Result{}); err != nil {
		t.Fatalf("Write returned error: %v", err)
	}

	out := buf.String()
	if !strings.Contains(out, "<envdiff>") {
		t.Errorf("expected <envdiff> tag, got: %s", out)
	}
}

func TestXMLFormatter_Write(t *testing.T) {
	f, err := output.NewFormatter("xml")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	results := []diff.Result{
		{
			Key:    "DB_HOST",
			Status: diff.StatusMismatch,
			Values: map[string]string{
				"production": "prod-db",
				"staging":    "stage-db",
			},
		},
		{
			Key:    "SECRET_KEY",
			Status: diff.StatusMissing,
			Values: map[string]string{
				"production": "abc123",
				"staging":    "",
			},
		},
	}

	var buf bytes.Buffer
	if err := f.Write(&buf, results); err != nil {
		t.Fatalf("Write returned error: %v", err)
	}

	out := buf.String()
	if !strings.Contains(out, "<?xml") {
		t.Errorf("expected XML declaration, got: %s", out)
	}
	if !strings.Contains(out, "DB_HOST") {
		t.Errorf("expected DB_HOST in output, got: %s", out)
	}
	if !strings.Contains(out, "SECRET_KEY") {
		t.Errorf("expected SECRET_KEY in output, got: %s", out)
	}
	if !strings.Contains(out, "mismatch") {
		t.Errorf("expected mismatch status in output, got: %s", out)
	}
	if !strings.Contains(out, "missing") {
		t.Errorf("expected missing status in output, got: %s", out)
	}
}

func TestNewFormatter_XMLFormat(t *testing.T) {
	_, err := output.NewFormatter("xml")
	if err != nil {
		t.Errorf("expected no error for xml format, got: %v", err)
	}
}
