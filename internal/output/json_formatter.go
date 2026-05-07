package output

import (
	"encoding/json"
	"io"

	"github.com/user/envdiff/internal/diff"
)

// jsonResult is the serialisable representation of a diff.Result.
type jsonResult struct {
	Key    string            `json:"key"`
	Status string            `json:"status"`
	Values map[string]string `json:"values"`
}

// JSONFormatter writes results as a JSON array.
type JSONFormatter struct{}

func (f *JSONFormatter) Write(w io.Writer, results []diff.Result) error {
	output := make([]jsonResult, 0, len(results))
	for _, r := range results {
		output = append(output, jsonResult{
			Key:    r.Key,
			Status: string(r.Status),
			Values: r.Values,
		})
	}

	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(output)
}
