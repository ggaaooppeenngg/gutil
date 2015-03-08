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
	inputs := []string{
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
	}

	LSD(inputs, 7)

	if !stringSliceEqual(inputs, []string{
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
	}) {
		t.Fail()
	}
}

func TestMSD(t *testing.T) {
	inputs := []string{
		"she",
		"sells",
		"seashells",
		"by",
		"the",
		"sea",
		"shore",
		"the",
		"shells",
		"she",
		"sells",
		"are",
		"surely",
		"seashells",
	}

	MSD(inputs)

	if !stringSliceEqual(inputs, []string{
		"are",
		"by",
		"sea",
		"seashells",
		"seashells",
		"sells",
		"sells",
		"she",
		"she",
		"shells",
		"shore",
		"surely",
		"the",
		"the",
	}) {
		t.Fail()
	}
}
