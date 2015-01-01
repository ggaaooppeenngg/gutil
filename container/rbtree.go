package container //import "github.com/ggaaooppeenngg/util/container"

//insipred by this article http://www.eternallyconfuzzled.com/tuts/datastructures/jsw_tut_rbtree.aspx

import (
	"fmt"
	"log"
	"os"
)

type RBtree struct {
	root *rbNode
}

const (
	left = iota
	right
)

type rbNode struct {
	ele   BinaryNode
	child [2]*rbNode
	depth int
	red   bool
}

func NewRBtree() *RBtree {
	return &RBtree{}
}

var logger *log.Logger

func init() {
	logger = log.New(os.Stdout, "^_^:", log.LstdFlags)
}

//单旋转
func roat(root *rbNode, dir int) *rbNode {
	dir_ := (dir + 1) % 2
	x := root.child[dir_]
	root.child[dir_] = x.child[dir]
	x.child[dir] = root
	root.red = true
	x.red = false
	return x
}

//双旋转
func roatDouble(root *rbNode, dir int) *rbNode {
	dir_ := (dir + 1) % 2
	root.child[dir_] = roat(root.child[dir_], dir_)
	return roat(root, dir)
}

//递归方法
func insert(root *rbNode, n BinaryNode) *rbNode {
	if root == nil {
		logger.Println("insert a new node", n)
		root = new(rbNode)
		root.ele = n
		root.red = true
		//新结点本身不调整.
	} else {
		var dir int
		if root.ele.Less(n) {
			logger.Println("search right")
			dir = right
		} else {
			logger.Println("search left")
			dir = left
		}
		root.child[dir] = insert(root.child[dir], n)
		//begin to fix up
		//pic is assuming dir is left
		var dir_ = (dir + 1) % 2
		if isRed(root.child[dir]) {
			if isRed(root.child[dir_]) {
				fmt.Println("change color")
				//case 1:
				//if 2 children are red,
				//1 and 3 can not be both red,so 2 must be black.
				//and red violation will propogate up.
				//
				//TREE                           TREE
				//B_2                            R_2
				//|  \             ->            |  \
				//R_1 R_3                        B_1 B_3
				//|                              |
				//R_0                            R_0

				root.red = true
				root.child[dir].red = false
				root.child[dir_].red = false
			} else {
				if isRed(root.child[dir].child[dir]) {
					fmt.Println("single")
					//case 2:
					//if left child is red and left.left child is red too.
					//
					//TREE   	             TREE
					//B_2                    B_1
					//|  \                   |  \
					//R_1 B_3   ->           R_0 R_2
					//|                          |
					//R_0                        B_3
					root = roat(root, dir_)
					//旋转以后root指向B_1,root结点就是指向根节点的指针,这样就省去了父亲节点.
				} else if isRed(root.child[dir].child[dir_]) {
					fmt.Println("double")
					root = roatDouble(root, dir_)
				}
			}
		}
	}
	return root
}

func isRed(n *rbNode) bool {
	if n != nil && n.red {
		return true
	}
	return false
}

func isBlack(n *rbNode) bool {
	if n == nil || !n.red {
		return true
	}
	return false
}

func (t *RBtree) Insert(n BinaryNode) {
	logger.Println("insert a node")
	t.root = insert(t.root, n)
	t.root.red = false
	//fix up
	//如果子节点是红色并且子节点的任意一个结点也是红色,就是发生冲突.
	//只有其中一个是,不然在这之前就已经存在红色冲突了
	return
}

var output string

func (t *RBtree) Walk() {
	output = ""
	var depth = 0
	fmt.Println("tree:")
	output += fmt.Sprintln("tree:")
	walk(t.root, depth)
}

func walk(root *rbNode, depth int) {
	var f = func(n *rbNode, depth int) {
		for i := 0; i < depth; i++ {
			fmt.Print("\t")
			output += fmt.Sprint("\t")
		}
		if n.red {
			fmt.Print("->R_")
			output += fmt.Sprint("->R_")
		} else {
			fmt.Print("->B_")
			output += fmt.Sprint("->B_")
		}
		fmt.Println(n.ele)
		output += fmt.Sprintln(n.ele)
	}
	if root != nil {
		depth++
		f(root, depth)
		if root.child[left] != nil {
			walk(root.child[left], depth)
		}
		if root.child[right] != nil {
			walk(root.child[right], depth)
		}
	}
}

func (t *RBtree) Print() {
	var depth = 0
	print(t.root, depth)
}

func print(n *rbNode, depth int) {
	if n != nil {
		depth++
		if isRed(n) {
			for i := 0; i < depth; i++ {
				fmt.Print("*")
			}
			fmt.Println("red_node")
		} else {
			for i := 0; i < depth; i++ {
				fmt.Print("*")
			}
			fmt.Println("black_node")
		}
	}
}
