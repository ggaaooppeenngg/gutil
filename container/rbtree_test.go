package container

import (
	"strconv"
	"strings"
	"testing"

	"github.com/ggaaooppeenngg/cpic"
)

//reb black tree asset
func rbTreeAsset(errors []string) func(root *rbNode) int {
	var f func(root *rbNode) int
	f = func(root *rbNode) int {
		if root == nil {
			return 0
		} else {
			l := root.child[left]
			r := root.child[right]
			if isRed(root) {
				if isRed(l) || isRed(r) {
					errors = append(errors, "Red violation")
					return 0
				}
			}
			lh := f(l)
			rh := f(r)
			if (l != nil && root.ele.Less(l.ele)) ||
				(r != nil && !root.ele.Less(root.ele)) {
				errors = append(errors, "Binary tree violation")
				return 0
			}
			if lh != 0 && rh != 0 && lh != rh {
				errors = append(errors, "Black violation")
				return 0
			}
			if isBlack(root) {
				return lh + 1
			} else {
				return lh
			}
		}
		return 0
	}
	return f
}

type Int int

func (i Int) Less(then BinaryNode) bool {
	return i < then.(Int)
}
func (i Int) String() string {
	return strconv.Itoa(int(i))
}

func TestTree(t *testing.T) {
	tree := NewRBtree()
	tree.Insert(Int(1))
	tree.Insert(Int(2))
	tree.Insert(Int(3))
	tree.Insert(Int(4))
	tree.Insert(Int(5))
	tree.Insert(Int(6))
	tree.Insert(Int(7))
	tree.Insert(Int(8))
	var errors []string
	rbTreeAsset(errors)(tree.root)
	if len(errors) > 0 {
		t.Fatal(strings.Join(errors, "\n"))
	}
	tree.Walk()
	o, err := cpic.Gen(output)
	if err != nil {
		t.Fatal(err)
	} else {
		println(o)
	}

}
