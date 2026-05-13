package diff

import (
	"testing"
)

func statsResults() []Result {
	return []Result{
		{
			Key:    "APP_NAME",
			Status: StatusMatch,
			Values: map[string]string{"dev": "myapp", "prod": "myapp"},
		},
		{
			Key:    "DB_URL",
			Status: StatusMismatch,
			Values: map[string]string{"dev": "localhost", "prod": "db.prod.example.com"},
		},
		{
			Key:    "SECRET_KEY",
			Status: StatusMissing,
			Values: map[string]string{"dev": "abc123", "prod": ""},
		},
	}
}

func TestStats_BasicCounts(t *testing.T) {
	results := statsResults()
	envNames := []string{"dev", "prod"}

	stats := Stats(results, envNames)

	if len(stats) != 2 {
		t.Fatalf("expected 2 env stats, got %d", len(stats))
	}

	dev := stats[0]
	if dev.Env != "dev" {
		t.Errorf("expected dev, got %s", dev.Env)
	}
	if dev.Total != 3 {
		t.Errorf("dev total: expected 3, got %d", dev.Total)
	}
	if dev.Missing != 0 {
		t.Errorf("dev missing: expected 0, got %d", dev.Missing)
	}
	if dev.Present != 3 {
		t.Errorf("dev present: expected 3, got %d", dev.Present)
	}
}

func TestStats_ProdMissingKey(t *testing.T) {
	results := statsResults()
	envNames := []string{"dev", "prod"}

	stats := Stats(results, envNames)
	prod := stats[1]

	if prod.Missing != 1 {
		t.Errorf("prod missing: expected 1, got %d", prod.Missing)
	}
	if prod.Mismatched != 1 {
		t.Errorf("prod mismatched: expected 1, got %d", prod.Mismatched)
	}
}

func TestStats_EmptyResults(t *testing.T) {
	stats := Stats([]Result{}, []string{"dev", "prod"})
	if stats != nil {
		t.Errorf("expected nil stats for empty results, got %v", stats)
	}
}

func TestStats_EmptyEnvNames(t *testing.T) {
	stats := Stats(statsResults(), []string{})
	if stats != nil {
		t.Errorf("expected nil stats for empty env names, got %v", stats)
	}
}

func TestStats_EnvNamesPreserved(t *testing.T) {
	results := statsResults()
	envNames := []string{"dev", "prod"}
	stats := Stats(results, envNames)

	for i, name := range envNames {
		if stats[i].Env != name {
			t.Errorf("expected env %s at index %d, got %s", name, i, stats[i].Env)
		}
	}
}
