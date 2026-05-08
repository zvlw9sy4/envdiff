package ignore

import (
	"bufio"
	"os"
	"strings"
)

// RuleSet holds a set of key patterns to ignore during comparison.
type RuleSet struct {
	patterns map[string]struct{}
}

// New returns an empty RuleSet.
func New() *RuleSet {
	return &RuleSet{patterns: make(map[string]struct{})}
}

// Add inserts a key pattern into the rule set.
func (r *RuleSet) Add(key string) {
	key = strings.TrimSpace(key)
	if key != "" && !strings.HasPrefix(key, "#") {
		r.patterns[key] = struct{}{}
	}
}

// Contains reports whether the given key matches any rule.
func (r *RuleSet) Contains(key string) bool {
	_, ok := r.patterns[key]
	return ok
}

// Keys returns all patterns in the rule set.
func (r *RuleSet) Keys() []string {
	keys := make([]string, 0, len(r.patterns))
	for k := range r.patterns {
		keys = append(keys, k)
	}
	return keys
}

// LoadFile reads an ignore file (one key per line, # for comments)
// and returns a populated RuleSet.
func LoadFile(path string) (*RuleSet, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	rs := New()
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		rs.Add(line)
	}
	return rs, scanner.Err()
}
