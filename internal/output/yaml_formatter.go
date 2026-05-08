package output

import (
	"fmt"
	"io"
	"strings"

	"github.com/user/envdiff/internal/diff"
)

type yamlFormatter struct{}

func (f *yamlFormatter) Write(w io.Writer, results []diff.Result) error {
	if len(results) == 0 {
		_, err := fmt.Fprintln(w, "results: []")
		return err
	}

	_, err := fmt.Fprintln(w, "results:")
	if err != nil {
		return err
	}

	for _, r := range results {
		_, err = fmt.Fprintf(w, "  - key: %s\n", r.Key)
		if err != nil {
			return err
		}

		_, err = fmt.Fprintf(w, "    status: %s\n", strings.ToLower(string(r.Status)))
		if err != nil {
			return err
		}

		_, err = fmt.Fprintf(w, "    values:\n")
		if err != nil {
			return err
		}

		for env, val := range r.Values {
			safeVal := val
			if val == "" {
				safeVal = "~"
			}
			_, err = fmt.Fprintf(w, "      %s: %s\n", env, safeVal)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
