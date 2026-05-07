package diff

// KeyStatus represents the comparison status of a key across environments.
type KeyStatus int

const (
	StatusMatch    KeyStatus = iota // key exists in both with same value
	StatusMismatch                  // key exists in both but values differ
	StatusMissing                   // key missing in one or more envs
)

// Result holds the diff result for a single key.
type Result struct {
	Key    string
	Status KeyStatus
	Values map[string]string // env name -> value (empty string if missing)
}

// Compare takes a map of environment name -> parsed key/value pairs and
// returns a slice of Result describing differences between them.
func Compare(envs map[string]map[string]string) []Result {
	// Collect all unique keys across all envs
	allKeys := make(map[string]struct{})
	for _, kv := range envs {
		for k := range kv {
			allKeys[k] = struct{}{}
		}
	}

	envNames := make([]string, 0, len(envs))
	for name := range envs {
		envNames = append(envNames, name)
	}

	var results []Result
	for key := range allKeys {
		values := make(map[string]string, len(envs))
		missingInAny := false
		var firstVal string
		firstSet := false
		mismatch := false

		for _, name := range envNames {
			val, ok := envs[name][key]
			if !ok {
				missingInAny = true
				values[name] = ""
			} else {
				values[name] = val
				if !firstSet {
					firstVal = val
					firstSet = true
				} else if val != firstVal {
					mismatch = true
				}
			}
		}

		status := StatusMatch
		if missingInAny {
			status = StatusMissing
		} else if mismatch {
			status = StatusMismatch
		}

		results = append(results, Result{
			Key:    key,
			Status: status,
			Values: values,
		})
	}

	sortResults(results)
	return results
}

// sortResults sorts results alphabetically by key for deterministic output.
func sortResults(results []Result) {
	for i := 1; i < len(results); i++ {
		for j := i; j > 0 && results[j].Key < results[j-1].Key; j-- {
			results[j], results[j-1] = results[j-1], results[j]
		}
	}
}
