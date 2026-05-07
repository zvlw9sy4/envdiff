package diff

import (
	"testing"
)

func TestCompare_AllMatch(t *testing.T) {
	envs := map[string]map[string]string{
		"dev":  {"APP_PORT": "8080", "DB_HOST": "localhost"},
		"prod": {"APP_PORT": "8080", "DB_HOST": "localhost"},
	}
	results := Compare(envs)
	for _, r := range results {
		if r.Status != StatusMatch {
			t.Errorf("expected key %q to match, got status %d", r.Key, r.Status)
		}
	}
}

func TestCompare_MissingKey(t *testing.T) {
	envs := map[string]map[string]string{
		"dev":  {"APP_PORT": "8080", "SECRET_KEY": "abc"},
		"prod": {"APP_PORT": "8080"},
	}
	results := Compare(envs)
	resultMap := indexByKey(results)

	if resultMap["SECRET_KEY"].Status != StatusMissing {
		t.Errorf("expected SECRET_KEY to be missing in prod")
	}
	if resultMap["APP_PORT"].Status != StatusMatch {
		t.Errorf("expected APP_PORT to match")
	}
}

func TestCompare_MismatchedValue(t *testing.T) {
	envs := map[string]map[string]string{
		"dev":  {"DB_HOST": "localhost"},
		"prod": {"DB_HOST": "db.prod.internal"},
	}
	results := Compare(envs)
	resultMap := indexByKey(results)

	if resultMap["DB_HOST"].Status != StatusMismatch {
		t.Errorf("expected DB_HOST to be a mismatch")
	}
}

func TestCompare_SortedOutput(t *testing.T) {
	envs := map[string]map[string]string{
		"dev": {"ZEBRA": "1", "ALPHA": "2", "MIDDLE": "3"},
	}
	results := Compare(envs)
	keys := []string{results[0].Key, results[1].Key, results[2].Key}
	expected := []string{"ALPHA", "MIDDLE", "ZEBRA"}
	for i, k := range keys {
		if k != expected[i] {
			t.Errorf("expected key[%d]=%q, got %q", i, expected[i], k)
		}
	}
}

func TestCompare_ThreeEnvs(t *testing.T) {
	envs := map[string]map[string]string{
		"dev":     {"A": "1", "B": "2"},
		"staging": {"A": "1", "B": "2"},
		"prod":    {"A": "1"},
	}
	results := Compare(envs)
	resultMap := indexByKey(results)

	if resultMap["B"].Status != StatusMissing {
		t.Errorf("expected B to be missing in prod")
	}
	if _, ok := resultMap["B"].Values["prod"]; !ok {
		t.Errorf("expected prod entry in values map for B")
	}
}

func indexByKey(results []Result) map[string]Result {
	m := make(map[string]Result, len(results))
	for _, r := range results {
		m[r.Key] = r
	}
	return m
}
