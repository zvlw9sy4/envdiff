package loader_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/user/envdiff/internal/loader"
)

func writeTempEnv(t *testing.T, content string) string {
	t.Helper()
	tmpDir := t.TempDir()
	path := filepath.Join(tmpDir, ".env")
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatalf("failed to write temp file: %v", err)
	}
	return path
}

func TestLoadFiles_BasicLoad(t *testing.T) {
	p := writeTempEnv(t, "KEY=value\nFOO=bar\n")
	files, err := loader.LoadFiles([]string{p}, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(files) != 1 {
		t.Fatalf("expected 1 file, got %d", len(files))
	}
	if files[0].Data["KEY"] != "value" {
		t.Errorf("expected KEY=value, got %s", files[0].Data["KEY"])
	}
}

func TestLoadFiles_CustomLabels(t *testing.T) {
	p := writeTempEnv(t, "A=1\n")
	files, err := loader.LoadFiles([]string{p}, []string{"production"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if files[0].Label != "production" {
		t.Errorf("expected label 'production', got %s", files[0].Label)
	}
}

func TestLoadFiles_DefaultLabelIsBasename(t *testing.T) {
	p := writeTempEnv(t, "X=1\n")
	files, err := loader.LoadFiles([]string{p}, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if files[0].Label != ".env" {
		t.Errorf("expected label '.env', got %s", files[0].Label)
	}
}

func TestLoadFiles_MismatchedLabels(t *testing.T) {
	p := writeTempEnv(t, "X=1\n")
	_, err := loader.LoadFiles([]string{p}, []string{"a", "b"})
	if err == nil {
		t.Fatal("expected error for mismatched labels, got nil")
	}
}

func TestLoadFiles_FileNotFound(t *testing.T) {
	_, err := loader.LoadFiles([]string{"/nonexistent/.env"}, nil)
	if err == nil {
		t.Fatal("expected error for missing file, got nil")
	}
}

func TestToEnvMaps(t *testing.T) {
	p1 := writeTempEnv(t, "A=1\n")
	p2 := writeTempEnv(t, "B=2\n")
	files, err := loader.LoadFiles([]string{p1, p2}, []string{"dev", "prod"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	labels, maps := loader.ToEnvMaps(files)
	if labels[0] != "dev" || labels[1] != "prod" {
		t.Errorf("unexpected labels: %v", labels)
	}
	if maps[0]["A"] != "1" || maps[1]["B"] != "2" {
		t.Errorf("unexpected map values")
	}
}
