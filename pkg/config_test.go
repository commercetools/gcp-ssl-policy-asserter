package sslpolicy

import "testing"

// TestListToContainsMap validates listToContainsMap
// creates a map for "does contain checking" from a list of
// strings correctly.
func TestListToContainsMap(t *testing.T) {
	input := []string{"cat", "dog", "potato"}
	output := listToContainsMap(input)

	for _, v := range input {
		_, ok := output[v]
		if !ok {
			t.Fatalf("output did not contain %s", v)
		}
	}
}

func TestEmptyListToContainsMap(t *testing.T) {
	input := []string{}
	output := listToContainsMap(input)

	for _, v := range []string{"x", "y", "z"} {
		_, ok := output[v]
		if ok {
			t.Fatalf("output should not contain %s.", v)
		}
	}
}
