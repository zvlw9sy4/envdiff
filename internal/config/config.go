package config

import (
	"errors"
	"fmt"
)

// Config holds all runtime configuration for envdiff.
type Config struct {
	Files       []string
	Labels      []string
	Format      string
	ShowAll     bool
	OnlyMissing bool
	OnlyMismatch bool
	FilterKeys  []string
	IgnoreFile  string
}

var validFormats = map[string]bool{
	"text":     true,
	"json":     true,
	"csv":      true,
	"markdown": true,
	"table":    true,
	"yaml":     true,
	"html":     true,
}

// Validate checks that the Config is consistent and complete.
func (c *Config) Validate() error {
	if len(c.Files) < 2 {
		return errors.New("at least two .env files are required")
	}
	if len(c.Labels) > 0 && len(c.Labels) != len(c.Files) {
		return fmt.Errorf("number of labels (%d) must match number of files (%d)",
			len(c.Labels), len(c.Files))
	}
	if c.Format != "" && !validFormats[c.Format] {
		return fmt.Errorf("unknown format %q; valid formats: text, json, csv, markdown, table, yaml, html", c.Format)
	}
	if c.OnlyMissing && c.OnlyMismatch {
		return errors.New("--only-missing and --only-mismatch are mutually exclusive")
	}
	if c.ShowAll && (c.OnlyMissing || c.OnlyMismatch) {
		return errors.New("--show-all cannot be combined with --only-missing or --only-mismatch")
	}
	return nil
}

// DefaultFormat returns the configured format or "text" as fallback.
func (c *Config) DefaultFormat() string {
	if c.Format == "" {
		return "text"
	}
	return c.Format
}
