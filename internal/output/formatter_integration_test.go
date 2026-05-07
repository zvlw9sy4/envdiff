package output_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/user/envdiff/internal/diff"
	"github.com/user/envdiff/internal/output"
)

var integrationResults = []diff.Result{
	{
		Key:    "DB_HOST",
		Status: diff.StatusMismatch,
		Values: map[string]string{"prod": "db.prod.example.com", "staging": "db.staging.example.com"},
	},
	{
		Key:    "SECRET_KEY",
		Status: diff.StatusMissing,
		Values: map[string]string{"prod": "supersecret", "staging": ""},
	},
}

func TestTextFormatter_Integration(t *testing.T) {
	var buf bytes.Buffer
	f, err := output.NewFormatter("text", &buf)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if err := f.Write(integrationResults); err != nil {
		t.Fatalf("Write failed: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "[MISMATCH]") {
		t.Errorf("expected MISMATCH label in text output, got:\n%s", out)
	}
	if !strings.Contains(out, "DB_HOST") {
		t.Errorf("expected DB_HOST in text output, got:\n%s", out)
	}
	if !strings.Contains(out, "<missing>") {
		t.Errorf("expected <missing> placeholder in text output, got:\n%s", out)
	}
}

func TestCSVFormatter_Integration(t *testing.T) {
	var buf bytes.Buffer
	f, err := output.NewFormatter("csv", &buf)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if err := f.Write(integrationResults); err != nil {
		t.Fatalf("Write failed: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "key,status,env,value") {
		t.Errorf("expected CSV header in output, got:\n%s", out)
	}
	if !strings.Contains(out, "SECRET_KEY") {
		t.Errorf("expected SECRET_KEY in CSV output, got:\n%s", out)
	}
}
