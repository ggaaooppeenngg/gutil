package container

import (
	"container/heap"
	"fmt"
)

//PQ is the index minimun priority queue containing *edges.
type PQ struct {
	edges []*edge
	index map[*Vertex]int // 用结点作索引,假设每个边的结点唯一.
}

func (pq PQ) String() string {
	var output string
	for index, edge := range pq.edges {
		if edge.weight != nil {
			output += edge.vtx.Id + ":" + fmt.Sprintf("%f", *edge.weight) + ","
		} else {
			output += edge.vtx.Id + ","
		}
		if index == pq.Len() {
			output = output[:len(output)-1]
		}
	}
	return output
}

func (pq PQ) Len() int { return len(pq.edges) }
func (pq PQ) Less(i, j int) bool {
	var w1, w2 float64
	if pq.edges[i].weight == nil {
		w1 = 0.0
	} else {
		w1 = *(pq.edges[i].weight)
	}
	if pq.edges[j].weight == nil {
		w2 = 0.0
	} else {
		w2 = *(pq.edges[j].weight)
	}
	return w1 < w2
}
func (pq PQ) Swap(i, j int) {
	pq.edges[i], pq.edges[j] = pq.edges[j], pq.edges[i]
}

// 不能用这个push pop,要用heap的push和pop.
func (pq *PQ) Push(x interface{}) {
	edge := x.(*edge)
	pq.index[edge.vtx] = len(pq.edges)
	pq.edges = append(pq.edges, edge)
}
func (pq *PQ) Pop() interface{} {
	old := pq.edges
	n := len(old)
	item := old[n-1]
	pq.edges = old[0 : n-1]
	delete(pq.index, item.vtx)
	return item
}

//含有某个结点.
func (pq PQ) Contains(v *Vertex) bool {
	_, ok := pq.index[v]
	return ok
}

//更改某个结点
func (pq *PQ) Change(v *Vertex, weight float64) {
	i := pq.index[v]
	pq.edges[i].SetWeight(weight)
	heap.Fix(pq, i)
}

//Insert inserts element.
func (pq *PQ) Insert(e *edge) {
	heap.Push(pq, e)
}

//Get gets the minimun element and removes it.
func (pq *PQ) Get() *edge {
	if pq.Len() > 0 {
		v := heap.Pop(pq)
		e, ok := v.(*edge)
		if ok {
			return e
		}
	}
	return nil
}

//newPQ returns a new PQ.edges is optional.
func newPQ(edges ...[]*edge) *PQ {
	var pq *PQ
	if len(edges) == 0 {
		pq = &PQ{
			index: make(map[*Vertex]int),
		}
	}
	if len(edges) == 1 {
		pq = &PQ{
			edges[0],
			make(map[*Vertex]int),
		}
		for i, edge := range edges[0] {
			pq.index[edge.vtx] = i
		}
		heap.Init(pq)
	}
	return pq
}
