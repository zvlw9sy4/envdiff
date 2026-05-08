package output

import (
	"fmt"
	"html"
	"io"
	"strings"

	"github.com/user/envdiff/internal/diff"
)

type HTMLFormatter struct{}

func (f *HTMLFormatter) Write(w io.Writer, results []diff.Result) error {
	fmt.Fprintln(w, `<!DOCTYPE html>`)
	fmt.Fprintln(w, `<html lang="en">`)
	fmt.Fprintln(w, `<head><meta charset="UTF-8"><title>envdiff Report</title>`)
	fmt.Fprintln(w, `<style>`)
	fmt.Fprintln(w, `body{font-family:sans-serif;padding:1rem;}table{border-collapse:collapse;width:100%;}th,td{border:1px solid #ccc;padding:0.5rem 1rem;text-align:left;}th{background:#f4f4f4;}.match{color:#2d7a2d;}.missing{color:#b05c00;}.mismatch{color:#a00;}`)
	fmt.Fprintln(w, `</style></head>`)
	fmt.Fprintln(w, `<body><h1>envdiff Report</h1>`)

	if len(results) == 0 {
		fmt.Fprintln(w, `<p>No results to display.</p></body></html>`)
		return nil
	}

	// Collect env labels from first result
	envs := make([]string, 0, len(results[0].Values))
	for _, r := range results {
		for env := range r.Values {
			envs = append(envs, env)
			break
		}
		break
	}
	// Use stable order from first result
	for _, r := range results[:1] {
		envs = envs[:0]
		for k := range r.Values {
			envs = append(envs, k)
		}
	}

	fmt.Fprintln(w, `<table>`)
	fmt.Fprintf(w, "<tr><th>Key</th><th>Status</th>")
	for _, e := range envs {
		fmt.Fprintf(w, "<th>%s</th>", html.EscapeString(e))
	}
	fmt.Fprintln(w, "</tr>")

	for _, r := range results {
		statusClass := strings.ToLower(string(r.Status))
		fmt.Fprintf(w, "<tr><td>%s</td><td class=%q>%s</td>",
			html.EscapeString(r.Key),
			statusClass,
			html.EscapeString(string(r.Status)),
		)
		for _, e := range envs {
			v, ok := r.Values[e]
			if !ok {
				v = "—"
			}
			fmt.Fprintf(w, "<td>%s</td>", html.EscapeString(v))
		}
		fmt.Fprintln(w, "</tr>")
	}

	fmt.Fprintln(w, `</table></body></html>`)
	return nil
}
