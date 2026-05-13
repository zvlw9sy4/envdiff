package audit

import (
	"fmt"
	"io"
	"time"

	"github.com/user/envdiff/internal/diff"
	"github.com/user/envdiff/internal/summary"
)

// Record represents a single audit log entry for an envdiff run.
type Record struct {
	Timestamp  time.Time        `json:"timestamp"`
	Files      []string         `json:"files"`
	Summary    summary.Summary  `json:"summary"`
	Results    []diff.Result    `json:"results"`
}

// New creates a new audit Record stamped with the current time.
func New(files []string, results []diff.Result) Record {
	s := summary.Compute(results)
	return Record{
		Timestamp: time.Now().UTC(),
		Files:     files,
		Summary:   s,
		Results:   results,
	}
}

// Write serialises the audit record as a human-readable log line to w.
func Write(w io.Writer, r Record) error {
	_, err := fmt.Fprintf(
		w,
		"[%s] files=%d keys=%d missing=%d mismatched=%d matched=%d\n",
		r.Timestamp.Format(time.RFC3339),
		len(r.Files),
		r.Summary.TotalKeys,
		r.Summary.Missing,
		r.Summary.Mismatched,
		r.Summary.Matched,
	)
	return err
}
