package diff

import (
	"sort"

	"github.com/user/envdiff/internal/parser"
)

// MergeResult represents the merged view of a key across all environments.
type MergeResult struct {
	Key    string
	Values map[string]string // env label -> value (empty string if missing)
}

// Merge takes a map of env label -> parsed env map and returns a unified
// slice of MergeResult entries, one per unique key across all environments.
func Merge(envs map[string]parser.EnvMap) []MergeResult {
	keySet := make(map[string]struct{})
	for _, env := range envs {
		for k := range env {
			keySet[k] = struct{}{}
		}
	}

	keys := make([]string, 0, len(keySet))
	for k := range keySet {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	results := make([]MergeResult, 0, len(keys))
	for _, key := range keys {
		values := make(map[string]string, len(envs))
		for label, env := range envs {
			values[label] = env[key] // empty string if key is absent
		}
		results = append(results, MergeResult{
			Key:    key,
			Values: values,
		})
	}
	return results
}

// MergeHasKey reports whether the given key is present (non-empty) in at
// least one environment within the MergeResult.
func MergeHasKey(mr MergeResult, label string) bool {
	_, ok := mr.Values[label]
	return ok
}
