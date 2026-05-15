package diff

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"time"
)

// Baseline represents a saved snapshot of env keys and their values
// from a specific environment, used to detect drift over time.
type Baseline struct {
	Label     string            `json:"label"`
	CreatedAt time.Time         `json:"created_at"`
	Keys      map[string]string `json:"keys"`
}

// NewBaseline creates a Baseline from a labelled env map.
func NewBaseline(label string, env map[string]string) *Baseline {
	keys := make(map[string]string, len(env))
	for k, v := range env {
		keys[k] = v
	}
	return &Baseline{
		Label:     label,
		CreatedAt: time.Now().UTC(),
		Keys:      keys,
	}
}

// SaveBaseline writes a Baseline to a JSON file at path.
func SaveBaseline(b *Baseline, path string) error {
	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("baseline: create %q: %w", path, err)
	}
	defer f.Close()
	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")
	return enc.Encode(b)
}

// LoadBaseline reads a Baseline from a JSON file at path.
func LoadBaseline(path string) (*Baseline, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("baseline: open %q: %w", path, err)
	}
	defer f.Close()
	var b Baseline
	if err := json.NewDecoder(f).Decode(&b); err != nil {
		return nil, fmt.Errorf("baseline: decode %q: %w", path, err)
	}
	return &b, nil
}

// DriftEntry describes a single key that has changed relative to the baseline.
type DriftEntry struct {
	Key      string
	Status   string // "added", "removed", "changed"
	OldValue string
	NewValue string
}

// DetectDrift compares a current env map against a saved Baseline and
// returns a sorted slice of DriftEntry describing what changed.
func DetectDrift(b *Baseline, current map[string]string) []DriftEntry {
	var entries []DriftEntry

	for k, oldVal := range b.Keys {
		newVal, ok := current[k]
		if !ok {
			entries = append(entries, DriftEntry{Key: k, Status: "removed", OldValue: oldVal})
		} else if newVal != oldVal {
			entries = append(entries, DriftEntry{Key: k, Status: "changed", OldValue: oldVal, NewValue: newVal})
		}
	}

	for k, newVal := range current {
		if _, exists := b.Keys[k]; !exists {
			entries = append(entries, DriftEntry{Key: k, Status: "added", NewValue: newVal})
		}
	}

	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Key < entries[j].Key
	})
	return entries
}
