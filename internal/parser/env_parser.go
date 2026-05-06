package parser

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// EnvMap represents the key-value pairs parsed from a .env file.
type EnvMap map[string]string

// ParseFile reads a .env file and returns a map of key-value pairs.
// It skips blank lines and comments (lines starting with '#').
// It returns an error if the file cannot be opened or contains malformed lines.
func ParseFile(path string) (EnvMap, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("opening env file %q: %w", path, err)
	}
	defer f.Close()

	env := make(EnvMap)
	scanner := bufio.NewScanner(f)
	lineNum := 0

	for scanner.Scan() {
		lineNum++
		line := strings.TrimSpace(scanner.Text())

		// Skip blank lines and comments
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		// Strip inline export keyword
		line = strings.TrimPrefix(line, "export ")

		idx := strings.IndexByte(line, '=')
		if idx < 0 {
			return nil, fmt.Errorf("%s:%d: malformed line (missing '='): %q", path, lineNum, line)
		}

		key := strings.TrimSpace(line[:idx])
		value := strings.TrimSpace(line[idx+1:])

		// Strip surrounding quotes from value
		value = stripQuotes(value)

		if key == "" {
			return nil, fmt.Errorf("%s:%d: empty key", path, lineNum)
		}

		env[key] = value
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("reading env file %q: %w", path, err)
	}

	return env, nil
}

// stripQuotes removes matching surrounding single or double quotes from a string.
func stripQuotes(s string) string {
	if len(s) >= 2 {
		if (s[0] == '"' && s[len(s)-1] == '"') ||
			(s[0] == '\'' && s[len(s)-1] == '\'') {
			return s[1 : len(s)-1]
		}
	}
	return s
}
