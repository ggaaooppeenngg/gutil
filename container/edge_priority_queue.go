package container

import (
	"container/heap"
)

//Item is the container of the
type Item struct {
	*edge
	index int // The index of the item in the heap,用来fix的时候用.
}

//PQ is the priority queue containing *Vertex.

type PQ struct {
	items []*Item
	index map[*Vertex]struct{}
}

func (pq PQ) Len() int { return len(pq.items) }
func (pq PQ) Less(i, j int) bool {
	if pq.items[i].weight == nil {
		return true
	}
	if pq.items[j].weight == nil {
		return false
	}
	if *(pq.items[i].weight) < *(pq.items[j].weight) {
		return true
	}
	return false
}
func (pq PQ) Swap(i, j int) {
	pq.items[i], pq.items[j] = pq.items[j], pq.items[i]
}

func (pq *PQ) Push(x interface{}) {
	n := pq.Len()
	item := new(Item)
	item.edge = x.(*edge)
	item.index = n
	pq.items = append(pq.items, item)
}

//BUG:pop有问题
func (pq *PQ) Pop() interface{} {
	old := pq.items
	if len(old) > 0 {
		item := old[0]
		old[0], old[len(old)-1] = old[len(old)-1], old[0]
		item.index = -1 //for safety
		*pq = old[0 : len(old)-1]
		return item.edge
	} else {

	}
}

//含有某个结点.
func (pq PQ) contains(v *Vertex) bool {
	_, ok := pq.index[v]
	return ok
}
