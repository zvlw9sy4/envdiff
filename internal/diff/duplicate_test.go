package diff

import (
	"testing"
)

func TestFindDuplicates_NoDuplicates(t *testing.T) {
	envLines := map[string][]string{
		"production": {"HOST=localhost", "PORT=8080", "DEBUG=false"},
	}
	results := FindDuplicates(envLines)
	if len(results) != 0 {
		t.Errorf("expected no duplicates, got %d", len(results))
	}
}

func TestFindDuplicates_SingleDuplicate(t *testing.T) {
	envLines := map[string][]string{
		"staging": {"HOST=localhost", "HOST=remotehost", "PORT=9090"},
	}
	results := FindDuplicates(envLines)
	if len(results) != 1 {
		t.Fatalf("expected 1 duplicate, got %d", len(results))
	}
	if results[0].Key != "HOST" {
		t.Errorf("expected key HOST, got %s", results[0].Key)
	}
	if results[0].Label != "staging" {
		t.Errorf("expected label staging, got %s", results[0].Label)
	}
	if results[0].Count != 2 {
		t.Errorf("expected count 2, got %d", results[0].Count)
	}
	if len(results[0].Values) != 2 {
		t.Errorf("expected 2 values, got %d", len(results[0].Values))
	}
}

func TestFindDuplicates_MultipleEnvs(t *testing.T) {
	envLines := map[string][]string{
		"dev":  {"KEY=a", "KEY=b"},
		"prod": {"KEY=x"},
	}
	results := FindDuplicates(envLines)
	if len(results) != 1 {
		t.Fatalf("expected 1 duplicate across envs, got %d", len(results))
	}
	if results[0].Label != "dev" {
		t.Errorf("expected duplicate in dev, got %s", results[0].Label)
	}
}

func TestFindDuplicates_SkipsComments(t *testing.T) {
	envLines := map[string][]string{
		"local": {"# HOST=commented", "HOST=real", "HOST=duplicate"},
	}
	results := FindDuplicates(envLines)
	if len(results) != 1 {
		t.Fatalf("expected 1 duplicate, got %d", len(results))
	}
	if results[0].Count != 2 {
		t.Errorf("expected count 2 (comment excluded), got %d", results[0].Count)
	}
}

func TestFindDuplicates_SortedOutput(t *testing.T) {
	envLines := map[string][]string{
		"env": {"Z_KEY=1", "Z_KEY=2", "A_KEY=x", "A_KEY=y"},
	}
	results := FindDuplicates(envLines)
	if len(results) != 2 {
		t.Fatalf("expected 2 duplicates, got %d", len(results))
	}
	if results[0].Key != "A_KEY" {
		t.Errorf("expected A_KEY first, got %s", results[0].Key)
	}
	if results[1].Key != "Z_KEY" {
		t.Errorf("expected Z_KEY second, got %s", results[1].Key)
	}
}

func TestFindDuplicates_EmptyInput(t *testing.T) {
	results := FindDuplicates(map[string][]string{})
	if len(results) != 0 {
		t.Errorf("expected no duplicates for empty input, got %d", len(results))
	}
}

func TestFindDuplicates_EmptyLines(t *testing.T) {
	envLines := map[string][]string{
		"env": {"", "", "HOST=localhost"},
	}
	results := FindDuplicates(envLines)
	if len(results) != 0 {
		t.Errorf("expected no duplicates when empty lines present, got %d", len(results))
	}
}

func TestSplitLine_ValidLine(t *testing.T) {
	key, value := splitLine("DATABASE_URL=postgres://localhost/db")
	if key != "DATABASE_URL" {
		t.Errorf("expected DATABASE_URL, got %s", key)
	}
	if value != "postgres://localhost/db" {
		t.Errorf("expected postgres://localhost/db, got %s", value)
	}
}

func TestSplitLine_ValueWithEquals(t *testing.T) {
	key, value := splitLine("EQUATION=a=b=c")
	if key != "EQUATION" {
		t.Errorf("expected EQUATION, got %s", key)
	}
	if value != "a=b=c" {
		t.Errorf("expected a=b=c, got %s", value)
	}
}

func TestSplitLine_NoEquals(t *testing.T) {
	key, value := splitLine("NOEQUALS")
	if key != "" {
		t.Errorf("expected empty key for line without '=', got %s", key)
	}
	if value != "" {
		t.Errorf("expected empty value for line without '=', got %s", value)
	}
}
