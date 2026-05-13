package output

import (
	"fmt"
	"io"
	"sort"
	"strings"

	"github.com/user/envdiff/internal/diff"
)

type tableFormatter struct{}

func (f *tableFormatter) Write(w io.Writer, results []diff.Result) error {
	if len(results) == 0 {
		fmt.Fprintln(w, "No differences found.")
		return nil
	}

	// Collect all environment labels from first result and sort for consistent output
	envLabels := make([]string, 0)
	if len(results) > 0 {
		for label := range results[0].Values {
			envLabels = append(envLabels, label)
		}
		sort.Strings(envLabels)
	}

	// Calculate column widths
	keyWidth := len("KEY")
	statusWidth := len("STATUS")
	valWidths := make(map[string]int)
	for _, label := range envLabels {
		valWidths[label] = len(label)
	}

	for _, r := range results {
		if len(r.Key) > keyWidth {
			keyWidth = len(r.Key)
		}
		if len(r.Status) > statusWidth {
			statusWidth = len(r.Status)
		}
		for label, val := range r.Values {
			if len(val) > valWidths[label] {
				valWidths[label] = len(val)
			}
		}
	}

	// Build header
	header := fmt.Sprintf("| %-*s | %-*s", keyWidth, "KEY", statusWidth, "STATUS")
	for _, label := range envLabels {
		header += fmt.Sprintf(" | %-*s", valWidths[label], label)
	}
	header += " |"

	// Build separator
	sep := "|" + strings.Repeat("-", keyWidth+2) + "|" + strings.Repeat("-", statusWidth+2)
	for _, label := range envLabels {
		sep += "|" + strings.Repeat("-", valWidths[label]+2)
	}
	sep += "|"

	fmt.Fprintln(w, header)
	fmt.Fprintln(w, sep)

	for _, r := range results {
		row := fmt.Sprintf("| %-*s | %-*s", keyWidth, r.Key, statusWidth, r.Status)
		for _, label := range envLabels {
			val := r.Values[label]
			row += fmt.Sprintf(" | %-*s", valWidths[label], val)
		}
		row += " |"
		fmt.Fprintln(w, row)
	}

	return nil
}
