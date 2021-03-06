// Generated by: main
// TypeWriter: queue
// Directive: +gen on Vertex

package container

// Generated by gen:http://clipperhouse.github.io/gen,
// with package github.com/ggaaooppeenngg/queue.
// Do not modify it.
// Synchronization is not guaranteed.

type VertexQueue struct {
	nr   int
	head *VertexNode
	tail *VertexNode
}

type VertexNode struct {
	data *Vertex
	next *VertexNode
}

func (iq VertexQueue) Empty() bool {
	return iq.nr == 0
}
func (iq VertexQueue) Size() int {
	return iq.nr
}
func (iq VertexQueue) Front() *Vertex {
	if iq.head == nil {
		panic("Invalid memory reference")
	}
	return iq.head.data
}

func (iq *VertexQueue) Push(data *Vertex) {
	newNode := &VertexNode{
		data: data,
	}
	if iq.head == nil {
		iq.head = newNode
	} else {
		iq.tail.next = newNode
	}
	iq.tail = newNode
	iq.nr++
}

func (iq *VertexQueue) Pop() {
	if iq.head == nil {
		panic("Invalid memory reference")
	}
	iq.head = iq.head.next
	iq.nr--
}
