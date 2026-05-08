package config

import (
	"errors"
	"strings"
)

// Config holds the parsed CLI configuration for an envdiff run.
type Config struct {
	Files   []string
	Labels  []string
	Output  string
	ShowAll bool
	OnlyMissing  bool
	OnlyMismatch bool
	FilterKeys   []string
	NoColor bool
}

// Validate checks that the config is valid before use.
func (c *Config) Validate() error {
	if len(c.Files) < 2 {
		return errors.New("at least two .env files are required")
	}
	if len(c.Labels) > 0 && len(c.Labels) != len(c.Files) {
		return errors.New("number of labels must match number of files")
	}
	validOutputs := map[string]bool{
		"text":     true,
		"json":     true,
		"csv":      true,
		"markdown": true,
		"table":    true,
		"yaml":     true,
	}
	format := strings.ToLower(c.Output)
	if format == "" {
		c.Output = "text"
	} else if !validOutputs[format] {
		return errors.New("unsupported output format: " + c.Output)
	}
	if c.OnlyMissing && c.OnlyMismatch {
		return errors.New("--only-missing and --only-mismatch are mutually exclusive")
	}
	return nil
}

// FilterKeySet returns the FilterKeys as a set for O(1) lookup.
func (c *Config) FilterKeySet() map[string]bool {
	set := make(map[string]bool, len(c.FilterKeys))
	for _, k := range c.FilterKeys {
		set[strings.TrimSpace(k)] = true
	}
	return set
}
