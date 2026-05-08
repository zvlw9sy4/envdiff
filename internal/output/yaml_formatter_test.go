package output

import (
	"bytes"
	"strings"
	"testing"

	"github.com/user/envdiff/internal/diff"
)

func TestYAMLFormatter_EmptyResults(t *testing.T) {
	f := &yamlFormatter{}
	var buf bytes.Buffer

	if err := f.Write(&buf, []diff.Result{}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := buf.String()
	if !strings.Contains(out, "results: []") {
		t.Errorf("expected empty results marker, got: %s", out)
	}
}

func TestYAMLFormatter_Write(t *testing.T) {
	f := &yamlFormatter{}
	var buf bytes.Buffer

	results := []diff.Result{
		{
			Key:    "DB_HOST",
			Status: diff.StatusMatch,
			Values: map[string]string{"dev": "localhost", "prod": "localhost"},
		},
		{
			Key:    "API_KEY",
			Status: diff.StatusMissing,
			Values: map[string]string{"dev": "abc123", "prod": ""},
		},
	}

	if err := f.Write(&buf, results); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := buf.String()

	if !strings.Contains(out, "results:") {
		t.Errorf("expected 'results:' header")
	}
	if !strings.Contains(out, "key: DB_HOST") {
		t.Errorf("expected DB_HOST key")
	}
	if !strings.Contains(out, "status: match") {
		t.Errorf("expected status: match")
	}
	if !strings.Contains(out, "key: API_KEY") {
		t.Errorf("expected API_KEY key")
	}
	if !strings.Contains(out, "status: missing") {
		t.Errorf("expected status: missing")
	}
	if !strings.Contains(out, "prod: ~") {
		t.Errorf("expected null marker for empty prod value")
	}
}

func TestNewFormatter_YAMLFormat(t *testing.T) {
	f, err := NewFormatter("yaml")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := f.(*yamlFormatter); !ok {
		t.Errorf("expected *yamlFormatter, got %T", f)
	}
}
