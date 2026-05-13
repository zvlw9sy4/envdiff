package diff

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func envMaps() map[string]map[string]string {
	return map[string]map[string]string{
		"dev": {"DB_HOST": "localhost", "API_KEY": "abc", "DEBUG": "true"},
		"staging": {"DB_HOST": "staging-db", "API_KEY": "xyz"},
		"prod": {"DB_HOST": "prod-db", "API_KEY": "secret", "LOG_LEVEL": "warn"},
	}
}

func TestKeySet_ReturnsAllUniqueKeys(t *testing.T) {
	keys := KeySet(envMaps())
	assert.ElementsMatch(t, []string{"API_KEY", "DB_HOST", "DEBUG", "LOG_LEVEL"}, keys)
}

func TestKeySet_EmptyInput(t *testing.T) {
	keys := KeySet(map[string]map[string]string{})
	assert.Empty(t, keys)
}

func TestKeySet_IsSorted(t *testing.T) {
	keys := KeySet(envMaps())
	for i := 1; i < len(keys); i++ {
		assert.LessOrEqual(t, keys[i-1], keys[i])
	}
}

func TestIntersection_CommonKeys(t *testing.T) {
	keys := Intersection(envMaps())
	assert.Equal(t, []string{"API_KEY", "DB_HOST"}, keys)
}

func TestIntersection_EmptyInput(t *testing.T) {
	keys := Intersection(map[string]map[string]string{})
	assert.Nil(t, keys)
}

func TestIntersection_SingleEnv(t *testing.T) {
	envs := map[string]map[string]string{
		"dev": {"A": "1", "B": "2"},
	}
	keys := Intersection(envs)
	assert.Equal(t, []string{"A", "B"}, keys)
}

func TestDifference_KeysMissingInOthers(t *testing.T) {
	keys := Difference("dev", envMaps())
	// DEBUG is in dev but not in staging or prod
	assert.Contains(t, keys, "DEBUG")
}

func TestDifference_NoMissingKeys(t *testing.T) {
	envs := map[string]map[string]string{
		"a": {"X": "1"},
		"b": {"X": "2"},
	}
	keys := Difference("a", envs)
	assert.Empty(t, keys)
}

func TestDifference_UnknownBase(t *testing.T) {
	keys := Difference("nonexistent", envMaps())
	assert.Nil(t, keys)
}
