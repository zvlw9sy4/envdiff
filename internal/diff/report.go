package diff

import (
	"fmt"
	"io"
	"sort"
	"strings"
)

const (
	colorReset  = "\033[0m"
	colorRed    = "\033[31m"
	colorYellow = "\033[33m"
	colorGreen  = "\033[32m"
)

// PrintReport writes a human-readable diff report to w.
// envOrder controls the column order; if nil, alphabetical order is used.
func PrintReport(w io.Writer, results []Result, envOrder []string) {
	if len(results) == 0 {
		fmt.Fprintln(w, "No differences found.")
		return
	}

	// Determine env column order
	envCols := envOrder
	if len(envCols) == 0 {
		envSet := make(map[string]struct{})
		for _, r := range results {
			for env := range r.Values {
				envSet[env] = struct{}{}
			}
		}
		for env := range envSet {
			envCols = append(envCols, env)
		}
		sort.Strings(envCols)
	}

	// Header
	header := fmt.Sprintf("%-30s  %-10s  ", "KEY", "STATUS")
	for _, env := range envCols {
		header += fmt.Sprintf("%-20s  ", env)
	}
	fmt.Fprintln(w, header)
	fmt.Fprintln(w, strings.Repeat("-", len(header)))

	for _, r := range results {
		statusLabel, color := statusInfo(r.Status)
		line := fmt.Sprintf("%s%-30s  %-10s%s  ", color, r.Key, statusLabel, colorReset)
		for _, env := range envCols {
			val, ok := r.Values[env]
			if !ok {
				val = "<missing>"
			} else if val == "" {
				val = "<empty>"
			}
			line += fmt.Sprintf("%-20s  ", val)
		}
		fmt.Fprintln(w, line)
	}
}

// statusInfo returns the display label and ANSI color code for a given Status.
func statusInfo(s Status) (label, color string) {
	switch s {
	case StatusMatch:
		return "OK", colorGreen
	case StatusMismatch:
		return "MISMATCH", colorYellow
	case StatusMissing:
		return "MISSING", colorRed
	default:
		return "UNKNOWN", colorReset
	}
}
