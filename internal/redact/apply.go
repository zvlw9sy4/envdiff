package redact

import "github.com/yourorg/envdiff/internal/diff"

// ApplyToResults returns a copy of the diff results with sensitive values
// replaced by the redacted placeholder in every environment entry.
func (r *Redactor) ApplyToResults(results []diff.Result) []diff.Result {
	out := make([]diff.Result, len(results))
	for i, res := range results {
		copy := res
		if r.IsSensitive(res.Key) {
			redacted := make(map[string]string, len(res.Values))
			for env, val := range res.Values {
				if val != "" {
					redacted[env] = redactedValue
				} else {
					redacted[env] = val
				}
			}
			copy.Values = redacted
		}
		out[i] = copy
	}
	return out
}
