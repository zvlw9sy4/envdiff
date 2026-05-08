package ignore_test

import (
	"os"
	"testing"

	"github.com/user/envdiff/internal/ignore"
)

func writeTempIgnore(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp("", "*.envignore")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := f.WriteString(content); err != nil {
		t.Fatal(err)
	}
	f.Close()
	t.Cleanup(func() { os.Remove(f.Name()) })
	return f.Name()
}

func TestRuleSet_AddAndContains(t *testing.T) {
	rs := ignore.New()
	rs.Add("SECRET_KEY")
	rs.Add("  DB_PASSWORD  ")

	if !rs.Contains("SECRET_KEY") {
		t.Error("expected SECRET_KEY to be in rule set")
	}
	if !rs.Contains("DB_PASSWORD") {
		t.Error("expected DB_PASSWORD (trimmed) to be in rule set")
	}
	if rs.Contains("OTHER_KEY") {
		t.Error("expected OTHER_KEY to not be in rule set")
	}
}

func TestRuleSet_SkipsComments(t *testing.T) {
	rs := ignore.New()
	rs.Add("# this is a comment")
	rs.Add("VALID_KEY")

	if rs.Contains("# this is a comment") {
		t.Error("comment line should not be added")
	}
	if !rs.Contains("VALID_KEY") {
		t.Error("VALID_KEY should be present")
	}
}

func TestLoadFile_Basic(t *testing.T) {
	path := writeTempIgnore(t, "# ignore these\nSECRET_KEY\nDB_PASSWORD\n\nAPI_TOKEN\n")
	rs, err := ignore.LoadFile(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for _, key := range []string{"SECRET_KEY", "DB_PASSWORD", "API_TOKEN"} {
		if !rs.Contains(key) {
			t.Errorf("expected %q in rule set", key)
		}
	}
}

func TestLoadFile_NotFound(t *testing.T) {
	_, err := ignore.LoadFile("/nonexistent/.envignore")
	if err == nil {
		t.Error("expected error for missing file")
	}
}

func TestRuleSet_Keys(t *testing.T) {
	rs := ignore.New()
	rs.Add("A")
	rs.Add("B")
	if len(rs.Keys()) != 2 {
		t.Errorf("expected 2 keys, got %d", len(rs.Keys()))
	}
}
