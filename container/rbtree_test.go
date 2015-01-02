package container

import (
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"testing"

	"github.com/ggaaooppeenngg/cpic"
)

//reb black tree asset
func rbTreeAsset(errors *[]string) func(root *rbNode) int {
	var f func(root *rbNode) int
	f = func(root *rbNode) int {
		if root == nil {
			return 0
		} else {
			l := root.child[left]
			r := root.child[right]
			//red violation check
			lh := f(l)
			rh := f(r)
			//if children has errors,stop asset
			if len(*errors) > 0 {
				return 0
			}
			if isRed(root) {
				if isRed(l) || isRed(r) {
					*errors = append(*errors, "Red violation")
					return 0
				}
			}
			if (l != nil && !l.ele.Less(root.ele)) ||
				(r != nil && !root.ele.Less(r.ele)) {
				var dir string
				if l != nil && !l.ele.Less(root.ele) {
					dir = fmt.Sprintf("Left node %v root node %v", l.ele, root.ele)
				}
				if r != nil && !root.ele.Less(r.ele) {
					dir = fmt.Sprintf("Right node %v root node %v", r.ele, root.ele)
				}
				*errors = append(*errors, dir+" binary tree violation")
				return 0
			}
			//check black violation
			if lh != rh {
				*errors = append(*errors, fmt.Sprintf("Black violation,Left hight is %d,Right hight is %d", lh, rh))
				return 0
			}
			if isRed(root) {
				return lh
			} else {
				return lh + 1
			}
		}
		return 0
	}
	return f
}
func asst(n *rbNode) (errors []string) {
	rbTreeAsset(&errors)(n)
	if isRed(n) {
		errors = append(errors, "root is red")
	}
	return errors
}

type Int int

func (i Int) Less(then BinaryNode) bool {
	return i < then.(Int)
}
func (i Int) String() string {
	return strconv.Itoa(int(i))
}

//100 random tests
func TestRandom(t *testing.T) {
	var limit = 100
	for i := 0; i < limit; i++ {
		var length = rand.Intn(100) + 1 // 0< length <=100
		var tree = NewRBtree()
		var input = make([]Int, length)
		for i := 0; i < length; i++ {
			int := Int(rand.Intn(1000))
			tree.Insert(int)
			input = append(input, int)
		}
		var errors []string
		errors = asst(tree.root)
		tree.Walk()
		_, err := cpic.Gen(output)
		if err != nil {
			t.Fatal(err)
		}
		if len(errors) > 0 {
			t.Fatal(strings.Join(errors, "\n"))
		}
	}
}

/*
func TestMannual(t *testing.T) {
	root := new(rbNode)
	root.child[left] = new(rbNode)
	root.ele = Int(6)
	root.child[left].ele = Int(12)
	fmt.Println(strings.Join(asst(root), "\n"))
}
*/
