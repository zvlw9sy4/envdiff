package diff

import (
	"fmt"
	"io"
	"os"
	"text/tabwriter"
)

// PrintRenameReport writes a human-readable rename detection report to stdout.
func PrintRenameReport(results []RenameResult) {
	WriteRenameReport(os.Stdout, results)
}

// WriteRenameReport writes a rename detection report to the given writer.
func WriteRenameReport(w io.Writer, results []RenameResult) {
	if len(results) == 0 {
		fmt.Fprintln(w, "No rename candidates detected.")
		return
	}

	fmt.Fprintf(w, "Rename candidates detected: %d\n\n", len(results))

	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)
	fmt.Fprintln(tw, "ENV\tOLD KEY\tNEW KEY\tVALUE")
	fmt.Fprintln(tw, "---\t-------\t-------\t-----")

	for _, r := range results {
		val := r.Value
		if len(val) > 20 {
			val = val[:17] + "..."
		}
		fmt.Fprintf(tw, "%s\t%s\t%s\t%s\n", r.EnvName, r.OldKey, r.NewKey, val)
	}

	tw.Flush()
}

// RenameReportSummary returns a short summary string of rename results.
func RenameReportSummary(results []RenameResult) string {
	if len(results) == 0 {
		return "no rename candidates found"
	}
	return fmt.Sprintf("%d rename candidate(s) found across %d environment(s)",
		len(results), countUniqueEnvs(results))
}

func countUniqueEnvs(results []RenameResult) int {
	seen := make(map[string]struct{})
	for _, r := range results {
		seen[r.EnvName] = struct{}{}
	}
	return len(seen)
}
