package output

import (
	"fmt"
	"io"

	"github.com/user/envdiff/internal/diff"
)

// Formatter is the interface all output formatters implement.
type Formatter interface {
	Write(w io.Writer, results []diff.Result) error
}

// NewFormatter returns a Formatter for the given format string.
// Supported formats: "text", "json", "csv", "markdown", "table".
func NewFormatter(format string) (Formatter, error) {
	switch format {
	case "text", "":
		return &textFormatter{}, nil
	case "json":
		return &jsonFormatter{}, nil
	case "csv":
		return &csvFormatter{}, nil
	case "markdown":
		return &markdownFormatter{}, nil
	case "table":
		return &tableFormatter{}, nil
	default:
		return nil, fmt.Errorf("unknown format: %q (supported: text, json, csv, markdown, table)", format)
	}
}
