package output

import (
	"fmt"
	"io"
	"strings"

	"github.com/user/envdiff/internal/diff"
)

// Format defines the output format type.
type Format string

const (
	FormatText Format = "text"
	FormatJSON Format = "json"
	FormatCSV  Format = "csv"
)

// Formatter writes diff results in a specific format.
type Formatter interface {
	Write(w io.Writer, results []diff.Result) error
}

// NewFormatter returns a Formatter for the given format string.
func NewFormatter(format Format) (Formatter, error) {
	switch format {
	case FormatText:
		return &TextFormatter{}, nil
	case FormatJSON:
		return &JSONFormatter{}, nil
	case FormatCSV:
		return &CSVFormatter{}, nil
	default:
		return nil, fmt.Errorf("unknown format %q: must be one of text, json, csv", format)
	}
}

// TextFormatter writes results as human-readable text.
type TextFormatter struct{}

func (f *TextFormatter) Write(w io.Writer, results []diff.Result) error {
	for _, r := range results {
		values := make([]string, 0, len(r.Values))
		for env, val := range r.Values {
			values = append(values, fmt.Sprintf("%s=%q", env, val))
		}
		_, err := fmt.Fprintf(w, "[%s] %s: %s\n", r.Status, r.Key, strings.Join(values, "  "))
		if err != nil {
			return err
		}
	}
	return nil
}

// CSVFormatter writes results as CSV rows.
type CSVFormatter struct{}

func (f *CSVFormatter) Write(w io.Writer, results []diff.Result) error {
	if len(results) == 0 {
		return nil
	}
	// Collect all env names in stable order from first result.
	envNames := make([]string, 0)
	for env := range results[0].Values {
		envNames = append(envNames, env)
	}

	header := append([]string{"key", "status"}, envNames...)
	_, err := fmt.Fprintln(w, strings.Join(header, ","))
	if err != nil {
		return err
	}

	for _, r := range results {
		row := []string{r.Key, string(r.Status)}
		for _, env := range envNames {
			row = append(row, r.Values[env])
		}
		_, err := fmt.Fprintln(w, strings.Join(row, ","))
		if err != nil {
			return err
		}
	}
	return nil
}
