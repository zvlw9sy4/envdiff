package diff

import (
	"testing"

	"github.com/user/envdiff/internal/parser"
)

func sampleEnvs() map[string]parser.EnvMap {
	return map[string]parser.EnvMap{
		"dev": {
			"APP_NAME": "myapp",
			"DEBUG":    "true",
			"DB_HOST":  "localhost",
		},
		"prod": {
			"APP_NAME": "myapp",
			"DB_HOST":  "prod.db.example.com",
			"SECRET":   "s3cr3t",
		},
	}
}

func TestMerge_AllKeysPresent(t *testing.T) {
	results := Merge(sampleEnvs())
	keys := make([]string, len(results))
	for i, r := range results {
		keys[i] = r.Key
	}
	want := []string{"APP_NAME", "DB_HOST", "DEBUG", "SECRET"}
	if len(keys) != len(want) {
		t.Fatalf("expected %d keys, got %d: %v", len(want), len(keys), keys)
	}
	for i, k := range want {
		if keys[i] != k {
			t.Errorf("key[%d]: want %q, got %q", i, k, keys[i])
		}
	}
}

func TestMerge_ValuesPresentForEachEnv(t *testing.T) {
	results := Merge(sampleEnvs())
	for _, r := range results {
		if _, ok := r.Values["dev"]; !ok {
			t.Errorf("key %q missing 'dev' entry in Values map", r.Key)
		}
		if _, ok := r.Values["prod"]; !ok {
			t.Errorf("key %q missing 'prod' entry in Values map", r.Key)
		}
	}
}

func TestMerge_MissingKeyIsEmptyString(t *testing.T) {
	results := Merge(sampleEnvs())
	for _, r := range results {
		if r.Key == "DEBUG" {
			if v := r.Values["prod"]; v != "" {
				t.Errorf("expected empty value for DEBUG in prod, got %q", v)
			}
			return
		}
	}
	t.Error("key DEBUG not found in merge results")
}

func TestMerge_EmptyInput(t *testing.T) {
	results := Merge(map[string]parser.EnvMap{})
	if len(results) != 0 {
		t.Errorf("expected 0 results, got %d", len(results))
	}
}

func TestMergeHasKey_Present(t *testing.T) {
	mr := MergeResult{
		Key:    "FOO",
		Values: map[string]string{"dev": "bar", "prod": ""},
	}
	if !MergeHasKey(mr, "dev") {
		t.Error("expected MergeHasKey to return true for 'dev'")
	}
}

func TestMergeHasKey_Absent(t *testing.T) {
	mr := MergeResult{
		Key:    "FOO",
		Values: map[string]string{"dev": "bar"},
	}
	if MergeHasKey(mr, "prod") {
		t.Error("expected MergeHasKey to return false for 'prod'")
	}
}
