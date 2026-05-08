package ignore

import "github.com/user/envdiff/internal/diff"

// FilterResults removes diff results whose Key is present in the RuleSet.
func FilterResults(results []diff.Result, rs *RuleSet) []diff.Result {
	if rs == nil || len(rs.Keys()) == 0 {
		return results
	}
	filtered := make([]diff.Result, 0, len(results))
	for _, r := range results {
		if !rs.Contains(r.Key) {
			filtered = append(filtered, r)
		}
	}
	return filtered
}

// FilterEnvMap removes keys present in the RuleSet from a parsed env map.
func FilterEnvMap(env map[string]string, rs *RuleSet) map[string]string {
	if rs == nil || len(rs.Keys()) == 0 {
		return env
	}
	out := make(map[string]string, len(env))
	for k, v := range env {
		if !rs.Contains(k) {
			out[k] = v
		}
	}
	return out
}
