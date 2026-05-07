package output_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/user/envdiff/internal/diff"
	"github.com/user/envdiff/internal/output"
)

var sampleResults = []diff.Result{
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

func TestNewFormatter_UnknownFormat(t *testing.T) {
	_, err := output.NewFormatter("xml")
	if err == nil {
		t.Fatal("expected error for unknown format, got nil")
	}
}

func TestTextFormatter_Write(t *testing.T) {
	f, err := output.NewFormatter(output.FormatText)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var buf bytes.Buffer
	if err := f.Write(&buf, sampleResults); err != nil {
		t.Fatalf("Write error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "DB_HOST") {
		t.Error("expected DB_HOST in text output")
	}
	if !strings.Contains(out, "[mismatch]") {
		t.Error("expected mismatch status in text output")
	}
}

func TestJSONFormatter_Write(t *testing.T) {
	f, err := output.NewFormatter(output.FormatJSON)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var buf bytes.Buffer
	if err := f.Write(&buf, sampleResults); err != nil {
		t.Fatalf("Write error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, `"key"`) {
		t.Error("expected JSON key field")
	}
	if !strings.Contains(out, `"status"`) {
		t.Error("expected JSON status field")
	}
	if !strings.Contains(out, "API_KEY") {
		t.Error("expected API_KEY in JSON output")
	}
}

func TestCSVFormatter_Write(t *testing.T) {
	f, err := output.NewFormatter(output.FormatCSV)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var buf bytes.Buffer
	if err := f.Write(&buf, sampleResults); err != nil {
		t.Fatalf("Write error: %v", err)
	}
	lines := strings.Split(strings.TrimSpace(buf.String()), "\n")
	// header + 3 data rows
	if len(lines) != 4 {
		t.Errorf("expected 4 lines, got %d", len(lines))
	}
	if !strings.HasPrefix(lines[0], "key,status") {
		t.Errorf("expected CSV header to start with 'key,status', got %q", lines[0])
	}
}

func TestCSVFormatter_EmptyResults(t *testing.T) {
	f, _ := output.NewFormatter(output.FormatCSV)
	var buf bytes.Buffer
	if err := f.Write(&buf, []diff.Result{}); err != nil {
		t.Fatalf("unexpected error on empty results: %v", err)
	}
	if buf.Len() != 0 {
		t.Error("expected empty output for empty results")
	}
}
