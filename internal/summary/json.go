package summary

import (
	"encoding/json"
	"io"
)

// jsonStats is the JSON-serialisable representation of Stats.
type jsonStats struct {
	Environments int `json:"environments"`
	TotalKeys    int `json:"total_keys"`
	Matched      int `json:"matched"`
	Missing      int `json:"missing"`
	Mismatched   int `json:"mismatched"`
	HasIssues    bool `json:"has_issues"`
}

// PrintJSON writes the summary as a JSON object to w.
func PrintJSON(w io.Writer, s Stats) error {
	js := jsonStats{
		Environments: s.EnvCount,
		TotalKeys:    s.TotalKeys,
		Matched:      s.Matched,
		Missing:      s.Missing,
		Mismatched:   s.Mismatched,
		HasIssues:    s.HasIssues(),
	}
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(js)
}
