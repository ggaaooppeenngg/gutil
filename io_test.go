package util

import (
	"testing"
)

func deleteIndexOfInt(i int, slice []int) []int {
	if i >= len(slice)-1 {
		return slice[:len(slice)-1]
	} else {
		return append(slice[:i], slice[i+1:]...)
	}
}
func addIndexOfInt(i int, slice []int, ele int) []int {
	if i < len(slice)-1 {
		return append(slice[:i], append([]int{ele}, slice[i:]...)...)
	} else {
		slice[i] = ele
		return slice
	}
}

func sliceIntEqual(s1, s2 []int) bool {
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
func TestRunOutput(t *testing.T) {
	out, err := RunOutput("echo","Hello world")
	if string(out) != "Hello world\n" || err != nil {
		t.Log(err)
		t.Log(string(out))
		t.Fail()
	}
}
func TestReverse(t *testing.T) {
	var slice = []int{0, 1, 2}
	var reverse = []int{2, 1, 0}
	Reverse(&slice)
	if !sliceIntEqual(slice, reverse) {
		t.Log(slice)
		t.Log(reverse)
		t.Fail()
	}
}

func TestIndexOf(t *testing.T) {
	var slice = []int{0, 1, 2}
	if IndexOf(0, slice) != 0 {
		t.Fatal("IndexOf failed,get index", IndexOf(0, slice))
	}
	type Int int
	var sliceI = []Int{Int(0), Int(1), Int(2)}
	if IndexOf(Int(0), sliceI) != 0 {
		t.Fatal("IndexOf failed,get index", IndexOf(Int(0), sliceI))
	}
}

func TestDeleteIndexOf(t *testing.T) {
	var slice = []int{0, 1, 2, 3}
	DeleteIndexOf(0, &slice)
	if !sliceIntEqual(slice, []int{1, 2, 3}) {
		t.Fatal("delete index %d at begining,failed.", 0)
	}
	slice = addIndexOfInt(0, slice, 0)
	DeleteIndexOf(1, &slice)
	if !sliceIntEqual(slice, []int{0, 2, 3}) {
		t.Fatal("delete index %d at middle ,failed.", 1)
	}
	slice = addIndexOfInt(1, slice, 1)
	DeleteIndexOf(3, &slice)
	if !sliceIntEqual(slice, []int{0, 1, 2}) {
		t.Fatal("delete index %d at last,failed.", 3)
	}
}
