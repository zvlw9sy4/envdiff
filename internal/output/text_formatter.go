package output

import (
	"fmt"
	"io"
	"sort"
	"strings"

	"github.com/user/envdiff/internal/diff"
)

type textFormatter struct {
	w io.Writer
}

func (f *textFormatter) Write(results []diff.Result) error {
	if len(results) == 0 {
		_, err := fmt.Fprintln(f.w, "No differences found.")
		return err
	}

	for _, r := range results {
		statusLabel := strings.ToUpper(string(r.Status))
		_, err := fmt.Fprintf(f.w, "[%s] %s\n", statusLabel, r.Key)
		if err != nil {
			return err
		}

		envs := make([]string, 0, len(r.Values))
		for env := range r.Values {
			envs = append(envs, env)
		}
		sort.Strings(envs)

		for _, env := range envs {
			val := r.Values[env]
			if val == "" {
				val = "<missing>"
			}
			_, err := fmt.Fprintf(f.w, "  %-20s %s\n", env+":", val)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
