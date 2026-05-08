package output

import (
	"fmt"
	"io"

	"github.com/user/envdiff/internal/diff"
)

// Formatter writes diff results to an io.Writer.
type Formatter interface {
	Write(w io.Writer, results []diff.Result) error
}

// NewFormatter returns a Formatter for the given format name.
// Supported formats: text, json, csv, markdown, table, yaml, html.
func NewFormatter(format string) (Formatter, error) {
	switch format {
	case "text", "":
		return &TextFormatter{}, nil
	case "json":
		return &JSONFormatter{}, nil
	case "csv":
		return &CSVFormatter{}, nil
	case "markdown", "md":
		return &MarkdownFormatter{}, nil
	case "table":
		return &TableFormatter{}, nil
	case "yaml":
		return &YAMLFormatter{}, nil
	case "html":
		return &HTMLFormatter{}, nil
	default:
		return nil, fmt.Errorf("unknown output format: %q", format)
	}
}
