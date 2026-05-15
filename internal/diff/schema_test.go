package diff

import (
	"bytes"
	"strings"
	"testing"
)

func schemaEnvs() map[string]map[string]string {
	return map[string]map[string]string{
		"dev": {
			"PORT":     "8080",
			"DEBUG":    "true",
			"API_URL":  "https://dev.example.com",
			"TIMEOUT":  "30",
		},
		"prod": {
			"PORT":     "443",
			"DEBUG":    "false",
			"API_URL":  "https://prod.example.com",
			// TIMEOUT missing in prod
		},
	}
}

func TestBuildSchema_KeyCount(t *testing.T) {
	envs := schemaEnvs()
	entries := BuildSchema(envs)
	if len(entries) != 4 {
		t.Errorf("expected 4 schema entries, got %d", len(entries))
	}
}

func TestBuildSchema_SortedKeys(t *testing.T) {
	entries := BuildSchema(schemaEnvs())
	for i := 1; i < len(entries); i++ {
		if entries[i].Key < entries[i-1].Key {
			t.Errorf("entries not sorted: %s before %s", entries[i-1].Key, entries[i].Key)
		}
	}
}

func TestBuildSchema_AllSame_WhenTypesMatch(t *testing.T) {
	entries := BuildSchema(schemaEnvs())
	for _, e := range entries {
		if e.Key == "PORT" && !e.AllSame {
			t.Errorf("PORT should be AllSame=true (both integer)")
		}
		if e.Key == "DEBUG" && !e.AllSame {
			t.Errorf("DEBUG should be AllSame=true (both bool)")
		}
	}
}

func TestBuildSchema_NotAllSame_WhenMissing(t *testing.T) {
	entries := BuildSchema(schemaEnvs())
	for _, e := range entries {
		if e.Key == "TIMEOUT" {
			if e.AllSame {
				t.Errorf("TIMEOUT should be AllSame=false (missing in prod)")
			}
			if e.EnvTypes["prod"] != "missing" {
				t.Errorf("expected prod TIMEOUT type=missing, got %s", e.EnvTypes["prod"])
			}
		}
	}
}

func TestPrintSchemaReport_ContainsKeys(t *testing.T) {
	var buf bytes.Buffer
	entries := BuildSchema(schemaEnvs())
	PrintSchemaReport(&buf, entries)
	out := buf.String()
	for _, key := range []string{"PORT", "DEBUG", "API_URL", "TIMEOUT"} {
		if !strings.Contains(out, key) {
			t.Errorf("expected output to contain key %q", key)
		}
	}
}

func TestPrintSchemaReport_EmptyEntries(t *testing.T) {
	var buf bytes.Buffer
	PrintSchemaReport(&buf, []SchemaEntry{})
	if !strings.Contains(buf.String(), "No schema entries") {
		t.Errorf("expected empty message, got: %s", buf.String())
	}
}

func TestBuildSchema_EmptyInput(t *testing.T) {
	entries := BuildSchema(map[string]map[string]string{})
	if len(entries) != 0 {
		t.Errorf("expected 0 entries for empty input, got %d", len(entries))
	}
}
