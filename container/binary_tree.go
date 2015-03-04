package container

type BinaryNode interface {
	Less(n BinaryNode) bool
}

type bNode struct {
	left  *bNode
	rigth *bNode
	ele   BinaryNode
}
type BinaryTree struct {
	root *bNode
}

func NewBinaryTree() *BinaryTree {
	return &BinaryTree{
		root: new(bNode),
	}
}

func (t *BinaryTree) Insert(n BinaryNode) {
	node := t.root
	for node.ele != nil {
		if n.Less(node.ele) {
			node = node.left
		} else {
			node = node.rigth
		}
	}
	node.ele = n
	return
}
