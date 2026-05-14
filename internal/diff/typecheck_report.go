package diff

import (
	"fmt"
	"io"
	"sort"
	"strings"
)

// PrintTypeMismatchReport writes a human-readable report of type mismatches to w.
func PrintTypeMismatchReport(w io.Writer, mismatches []TypeMismatch) {
	if len(mismatches) == 0 {
		fmt.Fprintln(w, "No value type mismatches detected.")
		return
	}

	fmt.Fprintf(w, "Found %d key(s) with value type mismatches:\n\n", len(mismatches))

	for _, m := range mismatches {
		fmt.Fprintf(w, "  %-30s\n", m.Key)
		labels := sortedLabels(m.Types)
		for _, label := range labels {
			fmt.Fprintf(w, "    %-20s -> %s\n", label, m.Types[label])
		}
	}
}

// TypeMismatchSummary returns a compact one-line summary string.
func TypeMismatchSummary(mismatches []TypeMismatch) string {
	if len(mismatches) == 0 {
		return "type check passed: no mismatches"
	}
	keys := make([]string, 0, len(mismatches))
	for _, m := range mismatches {
		keys = append(keys, m.Key)
	}
	sort.Strings(keys)
	return fmt.Sprintf("type mismatches in %d key(s): %s", len(mismatches), strings.Join(keys, ", "))
}

func sortedLabels(types map[string]ValueType) []string {
	labels := make([]string, 0, len(types))
	for l := range types {
		labels = append(labels, l)
	}
	sort.Strings(labels)
	return labels
}
