package ignore_test

import (
	"testing"

	"github.com/user/envdiff/internal/diff"
	"github.com/user/envdiff/internal/ignore"
)

func makeRuleSet(keys ...string) *ignore.RuleSet {
	rs := ignore.New()
	for _, k := range keys {
		rs.Add(k)
	}
	return rs
}

func TestFilterResults_RemovesIgnoredKeys(t *testing.T) {
	results := []diff.Result{
		{Key: "APP_NAME", Status: diff.StatusMatch},
		{Key: "SECRET_KEY", Status: diff.StatusMissing},
		{Key: "DB_HOST", Status: diff.StatusMismatch},
	}
	rs := makeRuleSet("SECRET_KEY")
	got := ignore.FilterResults(results, rs)

	if len(got) != 2 {
		t.Fatalf("expected 2 results, got %d", len(got))
	}
	for _, r := range got {
		if r.Key == "SECRET_KEY" {
			t.Error("SECRET_KEY should have been filtered out")
		}
	}
}

func TestFilterResults_NilRuleSet(t *testing.T) {
	results := []diff.Result{
		{Key: "APP_NAME", Status: diff.StatusMatch},
	}
	got := ignore.FilterResults(results, nil)
	if len(got) != 1 {
		t.Errorf("expected 1 result with nil RuleSet, got %d", len(got))
	}
}

func TestFilterEnvMap_RemovesIgnoredKeys(t *testing.T) {
	env := map[string]string{
		"APP_NAME":   "myapp",
		"SECRET_KEY": "supersecret",
		"DB_HOST":    "localhost",
	}
	rs := makeRuleSet("SECRET_KEY", "DB_HOST")
	got := ignore.FilterEnvMap(env, rs)

	if len(got) != 1 {
		t.Fatalf("expected 1 key, got %d", len(got))
	}
	if _, ok := got["APP_NAME"]; !ok {
		t.Error("APP_NAME should remain")
	}
}

func TestFilterEnvMap_EmptyRuleSet(t *testing.T) {
	env := map[string]string{"A": "1", "B": "2"}
	rs := makeRuleSet()
	got := ignore.FilterEnvMap(env, rs)
	if len(got) != 2 {
		t.Errorf("expected 2 keys with empty RuleSet, got %d", len(got))
	}
}
