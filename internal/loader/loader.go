package loader

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/user/envdiff/internal/parser"
)

// EnvFile represents a loaded environment file with its label and parsed data.
type EnvFile struct {
	Label string
	Path  string
	Data  map[string]string
}

// LoadFiles loads multiple .env files from the given paths.
// The label for each file defaults to its base filename unless overridden.
// Returns an error if the number of labels does not match the number of paths,
// if a file does not exist, or if parsing fails.
func LoadFiles(paths []string, labels []string) ([]EnvFile, error) {
	if len(labels) > 0 && len(labels) != len(paths) {
		return nil, fmt.Errorf("number of labels (%d) must match number of files (%d)", len(labels), len(paths))
	}

	result := make([]EnvFile, 0, len(paths))

	for i, p := range paths {
		if _, err := os.Stat(p); os.IsNotExist(err) {
			return nil, fmt.Errorf("file not found: %s", p)
		}

		data, err := parser.ParseFile(p)
		if err != nil {
			return nil, fmt.Errorf("failed to parse %s: %w", p, err)
		}

		label := filepath.Base(p)
		if len(labels) > 0 {
			label = labels[i]
		}

		result = append(result, EnvFile{
			Label: label,
			Path:  p,
			Data:  data,
		})
	}

	return result, nil
}

// ToEnvMaps extracts labels and data maps from a slice of EnvFiles.
func ToEnvMaps(files []EnvFile) ([]string, []map[string]string) {
	labels := make([]string, len(files))
	maps := make([]map[string]string, len(files))
	for i, f := range files {
		labels[i] = f.Label
		maps[i] = f.Data
	}
	return labels, maps
}

// KeyCount returns the total number of unique keys across all provided EnvFiles.
// A key is counted once regardless of how many files it appears in.
func KeyCount(files []EnvFile) int {
	seen := make(map[string]struct{})
	for _, f := range files {
		for k := range f.Data {
			seen[k] = struct{}{}
		}
	}
	return len(seen)
}
