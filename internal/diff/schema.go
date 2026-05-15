package diff

import (
	"fmt"
	"io"
	"sort"
	"strings"
)

// SchemaEntry describes a key's presence and inferred type across all environments.
type SchemaEntry struct {
	Key      string
	EnvTypes map[string]string // label -> inferred type
	AllSame  bool
}

// BuildSchema produces a schema summary: for each key found across all envs,
// record what type each environment's value is inferred to be.
func BuildSchema(envs map[string]map[string]string) []SchemaEntry {
	keys := KeySet(envs)
	entries := make([]SchemaEntry, 0, len(keys))

	for _, key := range keys {
		types := make(map[string]string, len(envs))
		for label, env := range envs {
			if val, ok := env[key]; ok {
				types[label] = InferType(val)
			} else {
				types[label] = "missing"
			}
		}
		entries = append(entries, SchemaEntry{
			Key:      key,
			EnvTypes: types,
			AllSame:  allTypesEqual(types),
		})
	}

	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Key < entries[j].Key
	})
	return entries
}

func allTypesEqual(types map[string]string) bool {
	var first string
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

// PrintSchemaReport writes a human-readable schema comparison to w.
func PrintSchemaReport(w io.Writer, entries []SchemaEntry) {
	if len(entries) == 0 {
		fmt.Fprintln(w, "No schema entries to display.")
		return
	}

	// Collect sorted label list from first entry
	var labels []string
	for l := range entries[0].EnvTypes {
		labels = append(labels, l)
	}
	sort.Strings(labels)

	header := fmt.Sprintf("%-30s  %-10s  %s", "KEY", "CONSISTENT", strings.Join(labels, "  "))
	fmt.Fprintln(w, header)
	fmt.Fprintln(w, strings.Repeat("-", len(header)+10))

	for _, e := range entries {
		consistent := "yes"
		if !e.AllSame {
			consistent = "NO"
		}
		typeParts := make([]string, 0, len(labels))
		for _, l := range labels {
			typeParts = append(typeParts, fmt.Sprintf("%-10s", e.EnvTypes[l]))
		}
		fmt.Fprintf(w, "%-30s  %-10s  %s\n", e.Key, consistent, strings.Join(typeParts, "  "))
	}
}
