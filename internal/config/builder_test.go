package config

import (
	"testing"
)

func TestBuildFromArgs_BasicFiles(t *testing.T) {
	cfg, err := BuildFromArgs([]string{"dev.env", "prod.env"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(cfg.Files) != 2 {
		t.Errorf("expected 2 files, got %d", len(cfg.Files))
	}
	if cfg.Output != "text" {
		t.Errorf("expected default output 'text', got %q", cfg.Output)
	}
}

func TestBuildFromArgs_OutputFlag(t *testing.T) {
	cfg, err := BuildFromArgs([]string{"--output", "json", "a.env", "b.env"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Output != "json" {
		t.Errorf("expected 'json', got %q", cfg.Output)
	}
}

func TestBuildFromArgs_Labels(t *testing.T) {
	cfg, err := BuildFromArgs([]string{"--label", "dev", "--label", "prod", "a.env", "b.env"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(cfg.Labels) != 2 || cfg.Labels[0] != "dev" || cfg.Labels[1] != "prod" {
		t.Errorf("unexpected labels: %v", cfg.Labels)
	}
}

func TestBuildFromArgs_FilterKeys(t *testing.T) {
	cfg, err := BuildFromArgs([]string{"--keys", "DB_HOST, API_KEY ,PORT", "a.env", "b.env"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(cfg.FilterKeys) != 3 {
		t.Fatalf("expected 3 filter keys, got %d", len(cfg.FilterKeys))
	}
	if cfg.FilterKeys[1] != "API_KEY" {
		t.Errorf("expected trimmed key 'API_KEY', got %q", cfg.FilterKeys[1])
	}
}

func TestBuildFromArgs_BoolFlags(t *testing.T) {
	cfg, err := BuildFromArgs([]string{"--all", "--no-color", "a.env", "b.env"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !cfg.ShowAll {
		t.Error("expected ShowAll=true")
	}
	if !cfg.NoColor {
		t.Error("expected NoColor=true")
	}
}

func TestBuildFromArgs_InvalidFlag(t *testing.T) {
	_, err := BuildFromArgs([]string{"--unknown-flag", "a.env"})
	if err == nil {
		t.Fatal("expected error for unknown flag")
	}
}
