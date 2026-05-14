package diff

import (
	"sort"
)

// DuplicateKey represents a key that appears more than once within a single env file.
type DuplicateKey struct {
	Key    string
	Label  string
	Count  int
	Values []string
}

// FindDuplicates scans a raw line-by-line representation of env files for duplicate keys.
// envLines maps a label to a slice of "KEY=VALUE" raw entries (pre-dedup).
func FindDuplicates(envLines map[string][]string) []DuplicateKey {
	var results []DuplicateKey

	for label, lines := range envLines {
		counts := make(map[string][]string)
		for _, line := range lines {
			key, value := splitLine(line)
			if key == "" {
				continue
			}
			counts[key] = append(counts[key], value)
		}

		for key, values := range counts {
			if len(values) > 1 {
				results = append(results, DuplicateKey{
					Key:    key,
					Label:  label,
					Count:  len(values),
					Values: values,
				})
			}
		}
	}

	sort.Slice(results, func(i, j int) bool {
		if results[i].Label != results[j].Label {
			return results[i].Label < results[j].Label
		}
		return results[i].Key < results[j].Key
	})

	return results
}

// splitLine splits a raw "KEY=VALUE" string into key and value.
// Returns empty strings if the line is a comment or blank.
func splitLine(line string) (string, string) {
	if len(line) == 0 || line[0] == '#' {
		return "", ""
	}
	for i := 0; i < len(line); i++ {
		if line[i] == '=' {
			return line[:i], line[i+1:]
		}
	}
	return "", ""
}
