package summary

import (
	"fmt"
	"io"

	"github.com/user/envdiff/internal/diff"
)

// Stats holds aggregated statistics from a diff comparison.
type Stats struct {
	TotalKeys  int
	Matched    int
	Missing    int
	Mismatched int
	EnvCount   int
}

// Compute calculates statistics from a slice of diff results.
func Compute(results []diff.Result) Stats {
	s := Stats{}
	seen := make(map[string]bool)

	for _, r := range results {
		if !seen[r.Key] {
			seen[r.Key] = true
			s.TotalKeys++
		}
		switch r.Status {
		case diff.StatusMatch:
			s.Matched++
		case diff.StatusMissing:
			s.Missing++
		case diff.StatusMismatch:
			s.Mismatched++
		}
	}
	return s
}

// ComputeWithEnvs sets EnvCount alongside the other stats.
func ComputeWithEnvs(results []diff.Result, envCount int) Stats {
	s := Compute(results)
	s.EnvCount = envCount
	return s
}

// Print writes a human-readable summary to w.
func Print(w io.Writer, s Stats) {
	fmt.Fprintf(w, "Summary\n")
	fmt.Fprintf(w, "-------\n")
	fmt.Fprintf(w, "Environments : %d\n", s.EnvCount)
	fmt.Fprintf(w, "Total keys   : %d\n", s.TotalKeys)
	fmt.Fprintf(w, "Matched      : %d\n", s.Matched)
	fmt.Fprintf(w, "Missing      : %d\n", s.Missing)
	fmt.Fprintf(w, "Mismatched   : %d\n", s.Mismatched)
}

// HasIssues returns true when there are any missing or mismatched keys.
func (s Stats) HasIssues() bool {
	return s.Missing > 0 || s.Mismatched > 0
}
