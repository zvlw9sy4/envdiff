package audit_test

import (
	"bytes"
	"strings"
	"testing"
	"time"

	"github.com/user/envdiff/internal/audit"
	"github.com/user/envdiff/internal/diff"
)

func sampleResults() []diff.Result {
	return []diff.Result{
		{Key: "APP_NAME", Status: diff.StatusMatch, Values: map[string]string{"dev": "app", "prod": "app"}},
		{Key: "DB_URL", Status: diff.StatusMismatch, Values: map[string]string{"dev": "localhost", "prod": "db.prod"}},
		{Key: "SECRET", Status: diff.StatusMissing, Values: map[string]string{"dev": "s3cr3t"}},
	}
}

func TestNew_SetsTimestamp(t *testing.T) {
	before := time.Now().UTC()
	r := audit.New([]string{"dev.env", "prod.env"}, sampleResults())
	after := time.Now().UTC()

	if r.Timestamp.Before(before) || r.Timestamp.After(after) {
		t.Errorf("timestamp %v outside expected range [%v, %v]", r.Timestamp, before, after)
	}
}

func TestNew_SummaryPopulated(t *testing.T) {
	r := audit.New([]string{"a.env", "b.env"}, sampleResults())

	if r.Summary.TotalKeys != 3 {
		t.Errorf("expected TotalKeys=3, got %d", r.Summary.TotalKeys)
	}
	if r.Summary.Missing != 1 {
		t.Errorf("expected Missing=1, got %d", r.Summary.Missing)
	}
	if r.Summary.Mismatched != 1 {
		t.Errorf("expected Mismatched=1, got %d", r.Summary.Mismatched)
	}
}

func TestNew_FilesStored(t *testing.T) {
	files := []string{"dev.env", "staging.env", "prod.env"}
	r := audit.New(files, nil)

	if len(r.Files) != 3 {
		t.Errorf("expected 3 files, got %d", len(r.Files))
	}
}

func TestWrite_ContainsExpectedFields(t *testing.T) {
	r := audit.New([]string{"dev.env", "prod.env"}, sampleResults())
	var buf bytes.Buffer

	if err := audit.Write(&buf, r); err != nil {
		t.Fatalf("Write returned error: %v", err)
	}

	line := buf.String()
	for _, want := range []string{"files=2", "keys=3", "missing=1", "mismatched=1", "matched=1"} {
		if !strings.Contains(line, want) {
			t.Errorf("expected %q in output: %s", want, line)
		}
	}
}

func TestWrite_EmptyResults(t *testing.T) {
	r := audit.New([]string{"a.env"}, nil)
	var buf bytes.Buffer

	if err := audit.Write(&buf, r); err != nil {
		t.Fatalf("Write returned error: %v", err)
	}

	if !strings.Contains(buf.String(), "keys=0") {
		t.Errorf("expected keys=0 in output: %s", buf.String())
	}
}
