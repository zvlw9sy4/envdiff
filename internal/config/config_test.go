package config

import (
	"testing"
)

func baseConfig() *Config {
	return &Config{
		Files:  []string{"a.env", "b.env"},
		Output: "text",
	}
}

func TestValidate_TooFewFiles(t *testing.T) {
	c := &Config{Files: []string{"a.env"}}
	if err := c.Validate(); err == nil {
		t.Fatal("expected error for fewer than 2 files")
	}
}

func TestValidate_LabelMismatch(t *testing.T) {
	c := baseConfig()
	c.Labels = []string{"prod"}
	if err := c.Validate(); err == nil {
		t.Fatal("expected error for mismatched labels")
	}
}

func TestValidate_InvalidFormat(t *testing.T) {
	c := baseConfig()
	c.Output = "xml"
	if err := c.Validate(); err == nil {
		t.Fatal("expected error for unsupported format")
	}
}

func TestValidate_MutuallyExclusiveFlags(t *testing.T) {
	c := baseConfig()
	c.OnlyMissing = true
	c.OnlyMismatch = true
	if err := c.Validate(); err == nil {
		t.Fatal("expected error for mutually exclusive flags")
	}
}

func TestValidate_DefaultsEmptyOutput(t *testing.T) {
	c := baseConfig()
	c.Output = ""
	if err := c.Validate(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if c.Output != "text" {
		t.Errorf("expected default output 'text', got %q", c.Output)
	}
}

func TestValidate_ValidConfig(t *testing.T) {
	c := baseConfig()
	c.Labels = []string{"dev", "prod"}
	if err := c.Validate(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestFilterKeySet(t *testing.T) {
	c := baseConfig()
	c.FilterKeys = []string{"DB_HOST", " API_KEY ", "PORT"}
	set := c.FilterKeySet()
	for _, k := range []string{"DB_HOST", "API_KEY", "PORT"} {
		if !set[k] {
			t.Errorf("expected key %q in set", k)
		}
	}
	if set[" API_KEY "] {
		t.Error("expected trimmed key, not raw padded key")
	}
}
