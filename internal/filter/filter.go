package filter

import "github.com/user/envdiff/internal/diff"

// StatusFilter defines which result statuses to include.
type StatusFilter struct {
	IncludeMissing   bool
	IncludeMismatch  bool
	IncludeMatch     bool
}

// DefaultFilter returns a filter that shows missing and mismatched keys only.
func DefaultFilter() StatusFilter {
	return StatusFilter{
		IncludeMissing:  true,
		IncludeMismatch: true,
		IncludeMatch:    false,
	}
}

// AllFilter returns a filter that includes every result status.
func AllFilter() StatusFilter {
	return StatusFilter{
		IncludeMissing:  true,
		IncludeMismatch: true,
		IncludeMatch:    true,
	}
}

// Apply returns only the results that match the filter criteria.
func Apply(results []diff.Result, f StatusFilter) []diff.Result {
	filtered := make([]diff.Result, 0, len(results))
	for _, r := range results {
		switch r.Status {
		case diff.StatusMissing:
			if f.IncludeMissing {
				filtered = append(filtered, r)
			}
		case diff.StatusMismatch:
			if f.IncludeMismatch {
				filtered = append(filtered, r)
			}
		case diff.StatusMatch:
			if f.IncludeMatch {
				filtered = append(filtered, r)
			}
		}
	}
	return filtered
}

// ByKey returns only results whose key matches one of the provided keys.
func ByKey(results []diff.Result, keys []string) []diff.Result {
	if len(keys) == 0 {
		return results
	}
	keySet := make(map[string]struct{}, len(keys))
	for _, k := range keys {
		keySet[k] = struct{}{}
	}
	filtered := make([]diff.Result, 0)
	for _, r := range results {
		if _, ok := keySet[r.Key]; ok {
			filtered = append(filtered, r)
		}
	}
	return filtered
}
