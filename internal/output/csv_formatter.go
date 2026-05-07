package output

import (
	"encoding/csv"
	"fmt"
	"io"

	"github.com/user/envdiff/internal/diff"
)

type csvFormatter struct {
	w io.Writer
}

func (f *csvFormatter) Write(results []diff.Result) error {
	cw := csv.NewWriter(f.w)

	header := []string{"key", "status", "env", "value"}
	if err := cw.Write(header); err != nil {
		return fmt.Errorf("csv: writing header: %w", err)
	}

	for _, r := range results {
		for env, val := range r.Values {
			row := []string{
				r.Key,
				string(r.Status),
				env,
				val,
			}
			if err := cw.Write(row); err != nil {
				return fmt.Errorf("csv: writing row for key %q: %w", r.Key, err)
			}
		}
	}

	cw.Flush()
	return cw.Error()
}
