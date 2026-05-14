package lint

import (
	"fmt"
	"strings"
)

// Severity indicates the level of a lint issue.
type Severity string

const (
	SeverityWarn  Severity = "warn"
	SeverityError Severity = "error"
)

// Issue represents a single lint finding for a key.
type Issue struct {
	Key      string
	Message  string
	Severity Severity
}

func (i Issue) String() string {
	return fmt.Sprintf("[%s] %s: %s", i.Severity, i.Key, i.Message)
}

// Rule is a function that inspects a key-value pair and returns issues.
type Rule func(key, value string) []Issue

// EmptyValueRule warns when a key is present but has an empty value.
func EmptyValueRule(key, value string) []Issue {
	if value == "" {
		return []Issue{{Key: key, Message: "key is present but has an empty value", Severity: SeverityWarn}}
	}
	return nil
}

// UpperCaseKeyRule warns when a key contains lowercase letters.
func UpperCaseKeyRule(key, value string) []Issue {
	if key != strings.ToUpper(key) {
		return []Issue{{Key: key, Message: "key should be uppercase", Severity: SeverityWarn}}
	}
	return nil
}

// NoSpaceInKeyRule errors when a key contains spaces.
func NoSpaceInKeyRule(key, _ string) []Issue {
	if strings.Contains(key, " ") {
		return []Issue{{Key: key, Message: "key must not contain spaces", Severity: SeverityError}}
	}
	return nil
}

// Run applies all provided rules to each key-value pair in the env map
// and returns a slice of all discovered issues.
func Run(env map[string]string, rules []Rule) []Issue {
	var issues []Issue
	for k, v := range env {
		for _, rule := range rules {
			issues = append(issues, rule(k, v)...)
		}
	}
	return issues
}

// DefaultRules returns the standard set of lint rules.
func DefaultRules() []Rule {
	return []Rule{
		EmptyValueRule,
		UpperCaseKeyRule,
		NoSpaceInKeyRule,
	}
}
