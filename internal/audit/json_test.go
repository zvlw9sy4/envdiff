package audit_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/user/envdiff/internal/audit"
)

func TestWriteJSON_ValidJSON(t *testing.T) {
	r := audit.New([]string{"dev.env", "prod.env"}, sampleResults())
	var buf bytes.Buffer

	if err := audit.WriteJSON(&buf, r); err != nil {
		t.Fatalf("WriteJSON returned error: %v", err)
	}

	if !strings.Contains(buf.String(), "\"timestamp\"") {
		t.Error("expected 'timestamp' field in JSON output")
	}
	if !strings.Contains(buf.String(), "\"files\"") {
		t.Error("expected 'files' field in JSON output")
	}
	if !strings.Contains(buf.String(), "\"summary\"") {
		t.Error("expected 'summary' field in JSON output")
	}
}

func TestWriteJSON_RoundTrip(t *testing.T) {
	orig := audit.New([]string{"a.env", "b.env"}, sampleResults())
	var buf bytes.Buffer

	if err := audit.WriteJSON(&buf, orig); err != nil {
		t.Fatalf("WriteJSON error: %v", err)
	}

	decoded, err := audit.ReadJSON(&buf)
	if err != nil {
		t.Fatalf("ReadJSON error: %v", err)
	}

	if len(decoded.Files) != len(orig.Files) {
		t.Errorf("files mismatch: got %d want %d", len(decoded.Files), len(orig.Files))
	}
	if decoded.Summary.TotalKeys != orig.Summary.TotalKeys {
		t.Errorf("TotalKeys mismatch: got %d want %d", decoded.Summary.TotalKeys, orig.Summary.TotalKeys)
	}
	if decoded.Summary.Missing != orig.Summary.Missing {
		t.Errorf("Missing mismatch: got %d want %d", decoded.Summary.Missing, orig.Summary.Missing)
	}
}

func TestReadJSON_InvalidInput(t *testing.T) {
	buf := bytes.NewBufferString("{not valid json")
	_, err := audit.ReadJSON(buf)
	if err == nil {
		t.Error("expected error for invalid JSON, got nil")
	}
}
