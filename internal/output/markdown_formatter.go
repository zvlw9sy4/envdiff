package output

import (
	"fmt"
	"io"
	"strings"

	"github.com/user/envdiff/internal/diff"
)

type markdownFormatter struct{}

func (f *markdownFormatter) Write(w io.Writer, results []diff.Result) error {
	if len(results) == 0 {
		_, err := fmt.Fprintln(w, "_No differences found._")
		return err
	}

	// Collect all environment labels from the first result
	var envLabels []string
	if len(results) > 0 {
		for label := range results[0].Values {
			envLabels = append(envLabels, label)
		}
	}
	// Sort labels for deterministic output
	sortedLabels := make([]string, len(envLabels))
	copy(sortedLabels, envLabels)
	for i := 0; i < len(sortedLabels)-1; i++ {
		for j := i + 1; j < len(sortedLabels); j++ {
			if sortedLabels[i] > sortedLabels[j] {
				sortedLabels[i], sortedLabels[j] = sortedLabels[j], sortedLabels[i]
			}
		}
	}

	// Header row
	header := "| Key | Status |" + strings.Repeat(" Env |" , 0)
	sep := "|-----|--------|" + strings.Repeat("-----|" , 0)

	columns := []string{"Key", "Status"}
	columns = append(columns, sortedLabels...)
	header = "| " + strings.Join(columns, " | ") + " |"

	sepParts := make([]string, len(columns))
	for i := range columns {
		sepParts[i] = "---"
	}
	sep = "| " + strings.Join(sepParts, " | ") + " |"

	if _, err := fmt.Fprintln(w, header); err != nil {
		return err
	}
	if _, err := fmt.Fprintln(w, sep); err != nil {
		return err
	}

	for _, r := range results {
		row := []string{r.Key, string(r.Status)}
		for _, label := range sortedLabels {
			val, ok := r.Values[label]
			if !ok {
				row = append(row, "_missing_")
			} else {
				row = append(row, val)
			}
		}
		line := "| " + strings.Join(row, " | ") + " |"
		if _, err := fmt.Fprintln(w, line); err != nil {
			return err
		}
	}

	return nil
}
