package diff

import (
	"sort"

	"github.com/user/envdiff/internal/diff"
)

// RenameMap maps old key names to new key names across environments.
type RenameMap map[string]string

// RenameResult holds the outcome of a rename detection check.
type RenameResult struct {
	OldKey  string
	NewKey  string
	EnvName string
	Value   string
}

// DetectRenames attempts to find keys that may have been renamed between
// a reference environment and others. A rename candidate is a key missing
// in one env but whose value matches a key present only in that env.
func DetectRenames(reference map[string]string, others map[string]map[string]string) []RenameResult {
	var results []RenameResult

	// Build a value->key reverse map for the reference env
	valToKey := make(map[string]string, len(reference))
	for k, v := range reference {
		if v != "" {
			valToKey[v] = k
		}
	}

	for envName, envMap := range others {
		for k, v := range envMap {
			if v == "" {
				continue
			}
			// Key exists in other but not in reference
			if _, inRef := reference[k]; !inRef {
				// Check if the value matches a key in reference
				if oldKey, found := valToKey[v]; found {
					if _, inOther := envMap[oldKey]; !inOther {
						results = append(results, RenameResult{
							OldKey:  oldKey,
							NewKey:  k,
							EnvName: envName,
							Value:   v,
						})
					}
				}
			}
		}
	}

	sort.Slice(results, func(i, j int) bool {
		if results[i].EnvName != results[j].EnvName {
			return results[i].EnvName < results[j].EnvName
		}
		return results[i].OldKey < results[j].OldKey
	})

	return results
}

// ApplyRenames returns a copy of envMap with keys renamed according to rm.
func ApplyRenames(envMap map[string]string, rm RenameMap) map[string]string {
	out := make(map[string]string, len(envMap))
	for k, v := range envMap {
		if newKey, ok := rm[k]; ok {
			out[newKey] = v
		} else {
			out[k] = v
		}
	}
	return out
}
