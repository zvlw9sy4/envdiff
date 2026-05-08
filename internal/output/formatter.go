package output

import (
	"fmt"
	"io"

	"github.com/user/envdiff/internal/diff"
)

// Formatter writes diff results to an io.Writer in a specific format.
type Formatter interface {
	Write(w io.Writer, results []diff.Result) error
}

// NewFormatter returns a Formatter for the given format name.
// Supported formats: text, json, csv, markdown, table, yaml.
func NewFormatter(format string) (Formatter, error) {
	switch format {
	case "text", "":
		return &textFormatter{}, nil
	case "json":
		return &jsonFormatter{}, nil
	case "csv":
		return &csvFormatter{}, nil
	case "markdown", "md":
		return &markdownFormatter{}, nil
	case "table":
		return &tableFormatter{}, nil
	case "yaml", "yml":
		return &yamlFormatter{}, nil
	default:
		return nil, fmt.Errorf("unknown format %q: supported formats are text, json, csv, markdown, table, yaml", format)
	}
}
