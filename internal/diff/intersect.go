package diff

import "sort"

// KeySet returns the union of all keys across the given env maps.
func KeySet(envs map[string]map[string]string) []string {
	seen := make(map[string]struct{})
	for _, env := range envs {
		for k := range env {
			seen[k] = struct{}{}
		}
	}
	keys := make([]string, 0, len(seen))
	for k := range seen {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

// Intersection returns keys that appear in ALL provided env maps.
func Intersection(envs map[string]map[string]string) []string {
	if len(envs) == 0 {
		return nil
	}

	// Count occurrences of each key
	counts := make(map[string]int)
	for _, env := range envs {
		for k := range env {
			counts[k]++
		}
	}

	keys := make([]string, 0)
	for k, count := range counts {
		if count == len(envs) {
			keys = append(keys, k)
		}
	}
	sort.Strings(keys)
	return keys
}

// Difference returns keys present in base but missing in at least one other env.
func Difference(base string, envs map[string]map[string]string) []string {
	baseEnv, ok := envs[base]
	if !ok {
		return nil
	}

	missing := make(map[string]struct{})
	for k := range baseEnv {
		for label, env := range envs {
			if label == base {
				continue
			}
			if _, found := env[k]; !found {
				missing[k] = struct{}{}
				break
			}
		}
	}

	keys := make([]string, 0, len(missing))
	for k := range missing {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
