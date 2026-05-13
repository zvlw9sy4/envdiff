package diff

import "github.com/yourusername/envdiff/internal/diff"

// EnvStats holds per-environment key statistics derived from diff results.
type EnvStats struct {
	Env        string
	Total      int
	Present    int
	Missing    int
	Mismatched int
}

// Stats computes per-environment statistics from a slice of Result.
func Stats(results []Result, envNames []string) []EnvStats {
	if len(results) == 0 || len(envNames) == 0 {
		return nil
	}

	// Build a map from env name to stats index for quick lookup.
	indexOf := make(map[string]int, len(envNames))
	stats := make([]EnvStats, len(envNames))
	for i, name := range envNames {
		stats[i] = EnvStats{Env: name}
		indexOf[name] = i
	}

	for _, r := range results {
		for _, name := range envNames {
			idx := indexOf[name]
			stats[idx].Total++

			val, exists := r.Values[name]
			switch {
			case !exists || val == "":
				if r.Status == StatusMissing {
					stats[idx].Missing++
				} else {
					stats[idx].Present++
				}
			case r.Status == StatusMismatch:
				stats[idx].Present++
				stats[idx].Mismatched++
			default:
				stats[idx].Present++
			}
		}
	}

	return stats
}
