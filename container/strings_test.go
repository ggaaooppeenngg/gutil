package container

import (
	"testing"
)

func stringSliceEqual(s1, s2 []string) bool {
	if len(s1) != len(s2) {
		return false
	}
	for i := 0; i < len(s1); i++ {
		if s1[i] != s2[i] {
			return false
		}
	}
	return true
}

func TestLSD(t *testing.T) {
	result := LSD([]string{
		"4PGC938",
		"2IYE230",
		"3CIO720",
		"1ICK750",
		"1OHV845",
		"4JZY524",
		"1ICK750",
		"3CIO720",
		"1OHV845",
		"1OHV845",
		"2RLA629",
		"2RLA629",
		"3ATW723",
	}, 7)
	stringSliceEqual(result, []string{
		"1ICK750",
		"1ICK750",
		"1OHV845",
		"1OHV845",
		"1OHV845",
		"2IYE230",
		"2RLA629",
		"2RLA629",
		"3ATW723",
		"3CIO720",
		"3CIO720",
		"4JZY524",
		"4PGC938",
	})
}
