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

type PQ []*Item

func (pq PQ) Len() int { return len(pq) }
func (pq PQ) Less(i, j int) bool {
	if pq[i].weight == nil {
		return true
	}
	if pq[j].weight == nil {
		return false
	}
	if *(pq[i].weight) < *(pq[j].weight) {
		return true
	}
	return false
}
func (pq PQ) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].index = i
	pq[j].index = j
}
func (pq *PQ) Push(x interface{}) {
	n := len(*pq)
	item := x.(*Item)
	item.index = n
	*pq = append(*pq, item)
}
func (pq *PQ) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	item.index = -1 //for safety
	*pq = old[0 : n-1]
	return item
}

//更改item.index的内容并且调整.
func (pq *PQ) update(item *Item, e *edge) {
	item.edge = e
	heap.Fix(pq, item.index)
}

//含有某个结点.
func (pq PQ) contains(v *Vertex) int {
	index := -1
	for _, item := range pq {
		if item.vtx == v {
			return item.index
		}
	}
	return index
}
