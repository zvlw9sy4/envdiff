package diff

import (
	"testing"
)

func TestInferType_Empty(t *testing.T) {
	if got := InferType(""); got != TypeEmpty {
		t.Errorf("expected TypeEmpty, got %s", got)
	}
}

func TestInferType_Bool(t *testing.T) {
	for _, v := range []string{"true", "false", "True", "FALSE"} {
		if got := InferType(v); got != TypeBool {
			t.Errorf("InferType(%q) = %s, want TypeBool", v, got)
		}
	}
}

func TestInferType_Integer(t *testing.T) {
	for _, v := range []string{"0", "42", "-7", "1000"} {
		if got := InferType(v); got != TypeInteger {
			t.Errorf("InferType(%q) = %s, want TypeInteger", v, got)
		}
	}
}

func TestInferType_Float(t *testing.T) {
	for _, v := range []string{"3.14", "-0.5", "1.0e10"} {
		if got := InferType(v); got != TypeFloat {
			t.Errorf("InferType(%q) = %s, want TypeFloat", v, got)
		}
	}
}

func TestInferType_URL(t *testing.T) {
	for _, v := range []string{"http://example.com", "https://api.example.com/v1"} {
		if got := InferType(v); got != TypeURL {
			t.Errorf("InferType(%q) = %s, want TypeURL", v, got)
		}
	}
}

func TestInferType_String(t *testing.T) {
	for _, v := range []string{"hello", "my-value", "some_key"} {
		if got := InferType(v); got != TypeString {
			t.Errorf("InferType(%q) = %s, want TypeString", v, got)
		}
	}
}

func TestDetectTypeMismatches_NoMismatch(t *testing.T) {
	envs := map[string]map[string]string{
		"dev":  {"PORT": "8080", "DEBUG": "true"},
		"prod": {"PORT": "443", "DEBUG": "false"},
	}
	result := DetectTypeMismatches(envs)
	if len(result) != 0 {
		t.Errorf("expected no mismatches, got %d", len(result))
	}
}

func TestDetectTypeMismatches_WithMismatch(t *testing.T) {
	envs := map[string]map[string]string{
		"dev":  {"PORT": "8080", "TIMEOUT": "30"},
		"prod": {"PORT": "https://proxy", "TIMEOUT": "30"},
	}
	result := DetectTypeMismatches(envs)
	if len(result) != 1 {
		t.Fatalf("expected 1 mismatch, got %d", len(result))
	}
	if result[0].Key != "PORT" {
		t.Errorf("expected mismatch on PORT, got %s", result[0].Key)
	}
	if result[0].Types["dev"] != TypeInteger {
		t.Errorf("expected dev PORT to be TypeInteger, got %s", result[0].Types["dev"])
	}
	if result[0].Types["prod"] != TypeURL {
		t.Errorf("expected prod PORT to be TypeURL, got %s", result[0].Types["prod"])
	}
}

func TestDetectTypeMismatches_EmptyEnvs(t *testing.T) {
	result := DetectTypeMismatches(map[string]map[string]string{})
	if len(result) != 0 {
		t.Errorf("expected no mismatches for empty input, got %d", len(result))
	}
}

func TestDetectTypeMismatches_MissingKeySkipped(t *testing.T) {
	// KEY_A only in dev — should not produce a type mismatch (it's missing, not mismatched)
	envs := map[string]map[string]string{
		"dev":  {"KEY_A": "123"},
		"prod": {},
	}
	result := DetectTypeMismatches(envs)
	if len(result) != 0 {
		t.Errorf("expected no type mismatches for missing key, got %d", len(result))
	}
}
