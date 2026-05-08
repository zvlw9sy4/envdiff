package output

import (
	"bytes"
	"strings"
	"testing"

	"github.com/user/envdiff/internal/diff"
)

func TestTableFormatter_EmptyResults(t *testing.T) {
	f := &tableFormatter{}
	var buf bytes.Buffer
	err := f.Write(&buf, []diff.Result{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "No differences found.") {
		t.Errorf("expected empty message, got: %s", buf.String())
	}
}

func TestTableFormatter_Write(t *testing.T) {
	results := []diff.Result{
		{
			Key:    "DB_HOST",
			Status: "mismatch",
			Values: map[string]string{
				"dev":  "localhost",
				"prod": "db.example.com",
			},
		},
		{
			Key:    "API_KEY",
			Status: "missing",
			Values: map[string]string{
				"dev":  "abc123",
				"prod": "",
			},
		},
	}

	f := &tableFormatter{}
	var buf bytes.Buffer
	err := f.Write(&buf, results)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "KEY") {
		t.Error("expected header to contain KEY")
	}
	if !strings.Contains(output, "STATUS") {
		t.Error("expected header to contain STATUS")
	}
	if !strings.Contains(output, "DB_HOST") {
		t.Error("expected output to contain DB_HOST")
	}
	if !strings.Contains(output, "mismatch") {
		t.Error("expected output to contain mismatch status")
	}
	if !strings.Contains(output, "API_KEY") {
		t.Error("expected output to contain API_KEY")
	}
}

func TestNewFormatter_TableFormat(t *testing.T) {
	f, err := NewFormatter("table")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := f.(*tableFormatter); !ok {
		t.Errorf("expected *tableFormatter, got %T", f)
	}
}
