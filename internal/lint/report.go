package lint

import (
	"fmt"
	"io"
	"sort"
)

// Summary holds aggregated counts from a lint run.
type Summary struct {
	Total  int
	Errors int
	Warns  int
}

// Summarize computes a Summary from a slice of Issues.
func Summarize(issues []Issue) Summary {
	s := Summary{Total: len(issues)}
	for _, i := range issues {
		switch i.Severity {
		case SeverityError:
			s.Errors++
		case SeverityWarn:
			s.Warns++
		}
	}
	return s
}

// PrintReport writes a human-readable lint report to w.
func PrintReport(w io.Writer, issues []Issue) {
	if len(issues) == 0 {
		fmt.Fprintln(w, "✔  No lint issues found.")
		return
	}

	sorted := make([]Issue, len(issues))
	copy(sorted, issues)
	sort.Slice(sorted, func(i, j int) bool {
		if sorted[i].Key != sorted[j].Key {
			return sorted[i].Key < sorted[j].Key
		}
		return sorted[i].Message < sorted[j].Message
	})

	for _, issue := range sorted {
		var prefix string
		switch issue.Severity {
		case SeverityError:
			prefix = "✘ ERROR"
		case SeverityWarn:
			prefix = "⚠ WARN "
		default:
			prefix = "  INFO "
		}
		fmt.Fprintf(w, "  %s  %-30s  %s\n", prefix, issue.Key, issue.Message)
	}

	s := Summarize(issues)
	fmt.Fprintf(w, "\n%d issue(s): %d error(s), %d warning(s)\n", s.Total, s.Errors, s.Warns)
}
