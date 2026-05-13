package summary_test

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/user/envdiff/internal/diff"
	"github.com/user/envdiff/internal/summary"
)

func TestPrintJSON_ValidJSON(t *testing.T) {
	results := []diff.Result{
		{Key: "HOST", Status: diff.StatusMatch},
		{Key: "PORT", Status: diff.StatusMissing},
	}
	s := summary.ComputeWithEnvs(results, 2)

	var buf bytes.Buffer
	if err := summary.PrintJSON(&buf, s); err != nil {
		t.Fatalf("PrintJSON returned error: %v", err)
	}

	var out map[string]interface{}
	if err := json.Unmarshal(buf.Bytes(), &out); err != nil {
		t.Fatalf("output is not valid JSON: %v", err)
	}
}

func TestPrintJSON_FieldValues(t *testing.T) {
	results := []diff.Result{
		{Key: "A", Status: diff.StatusMatch},
		{Key: "B", Status: diff.StatusMismatch},
		{Key: "C", Status: diff.StatusMissing},
	}
	s := summary.ComputeWithEnvs(results, 3)

	var buf bytes.Buffer
	_ = summary.PrintJSON(&buf, s)

	var out map[string]interface{}
	_ = json.Unmarshal(buf.Bytes(), &out)

	if int(out["total_keys"].(float64)) != 3 {
		t.Errorf("expected total_keys 3, got %v", out["total_keys"])
	}
	if out["has_issues"].(bool) != true {
		t.Error("expected has_issues to be true")
	}
	if int(out["environments"].(float64)) != 3 {
		t.Errorf("expected environments 3, got %v", out["environments"])
	}
}

func TestPrintJSON_NoIssues(t *testing.T) {
	results := []diff.Result{
		{Key: "HOST", Status: diff.StatusMatch},
	}
	s := summary.ComputeWithEnvs(results, 1)

	var buf bytes.Buffer
	_ = summary.PrintJSON(&buf, s)

	var out map[string]interface{}
	_ = json.Unmarshal(buf.Bytes(), &out)

	if out["has_issues"].(bool) != false {
		t.Error("expected has_issues to be false")
	}
}
