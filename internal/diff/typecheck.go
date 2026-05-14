package diff

import (
	"regexp"
	"strconv"
	"strings"
)

// ValueType represents the inferred type of an env value.
type ValueType string

const (
	TypeString  ValueType = "string"
	TypeInteger ValueType = "integer"
	TypeFloat   ValueType = "float"
	TypeBool    ValueType = "bool"
	TypeURL     ValueType = "url"
	TypeEmpty   ValueType = "empty"
)

var urlPattern = regexp.MustCompile(`^https?://`)

// InferType returns the inferred ValueType of a given env value string.
func InferType(value string) ValueType {
	if value == "" {
		return TypeEmpty
	}
	lower := strings.ToLower(value)
	if lower == "true" || lower == "false" {
		return TypeBool
	}
	if urlPattern.MatchString(value) {
		return TypeURL
	}
	if _, err := strconv.ParseInt(value, 10, 64); err == nil {
		return TypeInteger
	}
	if _, err := strconv.ParseFloat(value, 64); err == nil {
		return TypeFloat
	}
	return TypeString
}

// TypeMismatch describes a key whose value type differs across environments.
type TypeMismatch struct {
	Key   string
	Types map[string]ValueType // env label -> inferred type
}

// DetectTypeMismatches scans env maps and returns keys where the inferred
// value type is not consistent across all environments.
func DetectTypeMismatches(envs map[string]map[string]string) []TypeMismatch {
	keys := KeySet(envs)
	var mismatches []TypeMismatch

	for _, key := range keys {
		types := make(map[string]ValueType)
		for label, env := range envs {
			if val, ok := env[key]; ok {
				types[label] = InferType(val)
			}
		}
		if !allSameType(types) {
			mismatches = append(mismatches, TypeMismatch{Key: key, Types: types})
		}
	}
	return mismatches
}

func allSameType(types map[string]ValueType) bool {
	var first ValueType
	for _, t := range types {
		if first == "" {
			first = t
			continue
		}
		if t != first {
			return false
		}
	}
	return true
}
