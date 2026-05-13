package redact

import (
	"regexp"
	"strings"
)

// DefaultPatterns are key name patterns whose values should be redacted.
var DefaultPatterns = []string{
	"(?i)password",
	"(?i)secret",
	"(?i)token",
	"(?i)api_key",
	"(?i)private_key",
	"(?i)auth",
	"(?i)credential",
}

const redactedValue = "[REDACTED]"

// Redactor holds compiled patterns used to identify sensitive keys.
type Redactor struct {
	patterns []*regexp.Regexp
}

// New creates a Redactor from the given key-name patterns.
// If patterns is nil, DefaultPatterns are used.
func New(patterns []string) (*Redactor, error) {
	if patterns == nil {
		patterns = DefaultPatterns
	}
	compiled := make([]*regexp.Regexp, 0, len(patterns))
	for _, p := range patterns {
		re, err := regexp.Compile(p)
		if err != nil {
			return nil, err
		}
		compiled = append(compiled, re)
	}
	return &Redactor{patterns: compiled}, nil
}

// IsSensitive returns true when the key matches any redaction pattern.
func (r *Redactor) IsSensitive(key string) bool {
	for _, re := range r.patterns {
		if re.MatchString(key) {
			return true
		}
	}
	return false
}

// RedactValue returns the redacted placeholder if the key is sensitive,
// otherwise it returns the original value unchanged.
func (r *Redactor) RedactValue(key, value string) string {
	if r.IsSensitive(key) {
		return redactedValue
	}
	return value
}

// ApplyToMap returns a copy of the env map with sensitive values replaced.
func (r *Redactor) ApplyToMap(env map[string]string) map[string]string {
	out := make(map[string]string, len(env))
	for k, v := range env {
		out[k] = r.RedactValue(k, v)
	}
	return out
}

// RedactedValue is the placeholder string used for sensitive values.
func RedactedValue() string { return strings.Clone(redactedValue) }
