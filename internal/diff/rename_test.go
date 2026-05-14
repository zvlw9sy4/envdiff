package diff

import (
	"testing"
)

func TestDetectRenames_BasicRename(t *testing.T) {
	ref := map[string]string{
		"DB_HOST": "localhost",
		"DB_PORT": "5432",
	}
	others := map[string]map[string]string{
		"prod": {
			"DATABASE_HOST": "localhost",
			"DB_PORT":       "5432",
		},
	}

	results := DetectRenames(ref, others)

	if len(results) != 1 {
		t.Fatalf("expected 1 rename candidate, got %d", len(results))
	}
	if results[0].OldKey != "DB_HOST" {
		t.Errorf("expected OldKey=DB_HOST, got %s", results[0].OldKey)
	}
	if results[0].NewKey != "DATABASE_HOST" {
		t.Errorf("expected NewKey=DATABASE_HOST, got %s", results[0].NewKey)
	}
	if results[0].EnvName != "prod" {
		t.Errorf("expected EnvName=prod, got %s", results[0].EnvName)
	}
}

func TestDetectRenames_NoRename(t *testing.T) {
	ref := map[string]string{"APP_KEY": "abc123"}
	others := map[string]map[string]string{
		"staging": {"APP_KEY": "abc123"},
	}

	results := DetectRenames(ref, others)
	if len(results) != 0 {
		t.Errorf("expected 0 rename candidates, got %d", len(results))
	}
}

func TestDetectRenames_EmptyValues(t *testing.T) {
	ref := map[string]string{"KEY": ""}
	others := map[string]map[string]string{
		"prod": {"OTHER_KEY": ""},
	}

	results := DetectRenames(ref, others)
	if len(results) != 0 {
		t.Errorf("expected no renames for empty values, got %d", len(results))
	}
}

func TestDetectRenames_MultipleEnvs(t *testing.T) {
	ref := map[string]string{
		"OLD_TOKEN": "secret",
	}
	others := map[string]map[string]string{
		"staging": {"NEW_TOKEN": "secret"},
		"prod":    {"API_TOKEN": "secret"},
	}

	results := DetectRenames(ref, others)
	if len(results) != 2 {
		t.Fatalf("expected 2 rename candidates, got %d", len(results))
	}
}

func TestApplyRenames_Basic(t *testing.T) {
	envMap := map[string]string{
		"OLD_KEY": "value1",
		"KEEP_KEY": "value2",
	}
	rm := RenameMap{"OLD_KEY": "NEW_KEY"}

	result := ApplyRenames(envMap, rm)

	if _, ok := result["OLD_KEY"]; ok {
		t.Error("OLD_KEY should have been renamed")
	}
	if result["NEW_KEY"] != "value1" {
		t.Errorf("expected NEW_KEY=value1, got %s", result["NEW_KEY"])
	}
	if result["KEEP_KEY"] != "value2" {
		t.Errorf("expected KEEP_KEY=value2, got %s", result["KEEP_KEY"])
	}
}

func TestApplyRenames_EmptyMap(t *testing.T) {
	envMap := map[string]string{"KEY": "val"}
	rm := RenameMap{}

	result := ApplyRenames(envMap, rm)
	if result["KEY"] != "val" {
		t.Errorf("expected KEY=val, got %s", result["KEY"])
	}
}
