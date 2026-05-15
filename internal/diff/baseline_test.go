package diff

import (
	"os"
	"path/filepath"
	"testing"
)

func TestNewBaseline_StoresKeys(t *testing.T) {
	env := map[string]string{"HOST": "localhost", "PORT": "5432"}
	b := NewBaseline("prod", env)
	if b.Label != "prod" {
		t.Errorf("expected label prod, got %s", b.Label)
	}
	if b.Keys["HOST"] != "localhost" {
		t.Errorf("expected HOST=localhost")
	}
	if b.Keys["PORT"] != "5432" {
		t.Errorf("expected PORT=5432")
	}
	if b.CreatedAt.IsZero() {
		t.Error("expected non-zero CreatedAt")
	}
}

func TestNewBaseline_DoesNotMutateInput(t *testing.T) {
	env := map[string]string{"KEY": "val"}
	b := NewBaseline("dev", env)
	env["KEY"] = "mutated"
	if b.Keys["KEY"] != "val" {
		t.Error("baseline should not reflect mutations to original map")
	}
}

func TestSaveAndLoadBaseline_RoundTrip(t *testing.T) {
	env := map[string]string{"DB_HOST": "db.prod", "DEBUG": "false"}
	b := NewBaseline("staging", env)

	tmp := filepath.Join(t.TempDir(), "baseline.json")
	if err := SaveBaseline(b, tmp); err != nil {
		t.Fatalf("SaveBaseline: %v", err)
	}

	loaded, err := LoadBaseline(tmp)
	if err != nil {
		t.Fatalf("LoadBaseline: %v", err)
	}
	if loaded.Label != "staging" {
		t.Errorf("expected label staging, got %s", loaded.Label)
	}
	if loaded.Keys["DB_HOST"] != "db.prod" {
		t.Errorf("expected DB_HOST=db.prod")
	}
}

func TestLoadBaseline_InvalidPath(t *testing.T) {
	_, err := LoadBaseline("/nonexistent/path/baseline.json")
	if err == nil {
		t.Error("expected error for missing file")
	}
}

func TestLoadBaseline_InvalidJSON(t *testing.T) {
	tmp := filepath.Join(t.TempDir(), "bad.json")
	os.WriteFile(tmp, []byte("not json{"), 0644)
	_, err := LoadBaseline(tmp)
	if err == nil {
		t.Error("expected error for invalid JSON")
	}
}

func TestDetectDrift_Added(t *testing.T) {
	b := NewBaseline("prod", map[string]string{"A": "1"})
	current := map[string]string{"A": "1", "B": "2"}
	entries := DetectDrift(b, current)
	if len(entries) != 1 || entries[0].Key != "B" || entries[0].Status != "added" {
		t.Errorf("expected one added entry for B, got %+v", entries)
	}
}

func TestDetectDrift_Removed(t *testing.T) {
	b := NewBaseline("prod", map[string]string{"A": "1", "B": "2"})
	current := map[string]string{"A": "1"}
	entries := DetectDrift(b, current)
	if len(entries) != 1 || entries[0].Key != "B" || entries[0].Status != "removed" {
		t.Errorf("expected one removed entry for B, got %+v", entries)
	}
}

func TestDetectDrift_Changed(t *testing.T) {
	b := NewBaseline("prod", map[string]string{"HOST": "old.host"})
	current := map[string]string{"HOST": "new.host"}
	entries := DetectDrift(b, current)
	if len(entries) != 1 || entries[0].Status != "changed" {
		t.Errorf("expected one changed entry, got %+v", entries)
	}
	if entries[0].OldValue != "old.host" || entries[0].NewValue != "new.host" {
		t.Errorf("unexpected values: %+v", entries[0])
	}
}

func TestDetectDrift_NoChanges(t *testing.T) {
	env := map[string]string{"X": "1", "Y": "2"}
	b := NewBaseline("prod", env)
	entries := DetectDrift(b, env)
	if len(entries) != 0 {
		t.Errorf("expected no drift, got %+v", entries)
	}
}

func TestDetectDrift_SortedOutput(t *testing.T) {
	b := NewBaseline("prod", map[string]string{"Z": "1", "A": "1", "M": "1"})
	current := map[string]string{}
	entries := DetectDrift(b, current)
	for i := 1; i < len(entries); i++ {
		if entries[i].Key < entries[i-1].Key {
			t.Errorf("entries not sorted: %s before %s", entries[i-1].Key, entries[i].Key)
		}
	}
}
