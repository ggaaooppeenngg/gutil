package container

import ()

// 保证子节点比父结点小,用完全二叉树实现,只要一个数组.
// 默认是最小堆.
type Heap struct {
	heap []*edge
}

// Push adds an element to the heap.
func (h *Heap) Push(e *edge) {
	h.heap = append(h.heap, e)
	h.swim(len(h.heap) - 1)
}

// Pop returns minimun element or nil if heap is empty,and removes it.
func (h *Heap) Pop() *edge {
	if len(h.heap) > 0 {
		min := h.heap[0]
		h.Swap(0, len(h.heap)-1)
		h.heap = h.heap[:len(h.heap)-1]
		h.sink(0)
		return min
	} else {
		return nil
	}
}

// implement Swap,Less,Len of sort.Interface.
func (h *Heap) Swap(i, j int) {
	h.heap[i], h.heap[j] = h.heap[j], h.heap[i]
}
func (h *Heap) Less(i, j int) bool {
	var w1, w2 float64
	if h.heap[i].weight == nil {
		w1 = 0.0
	} else {
		w1 = *(h.heap[i].weight)
	}
	if h.heap[j].weight == nil {
		w2 = 0.0
	} else {
		w2 = *(h.heap[j].weight)
	}
	return w1 < w2
}
func (h *Heap) Len() int {
	return len(h.heap)
}

// travel up 去调整heap,不断向上替换结点.
func (h *Heap) swim(j int) {
	for {
		i := (j - 1) / 2 //parent
		if i == j || !h.Less(j, i) {
			break
		}
		h.Swap(i, j)
		j = i
	}
}

// travel down 去调整heap,如果index溢出的话就停止.
func (h *Heap) sink(i int) {
	for {
		j1 := 2*i + 1                //left child.
		if j1 >= h.Len() || j1 < 0 { //j1 < 0 after overflow.
			break
		}
		j := j1 //left child
		if j2 := j1 + 1; j2 < h.Len() && !h.Less(j1, j2) {
			j = j2 // = 2*i+2 ,right child
		}
		if !h.Less(j, i) {
			break
		}
		h.Swap(i, j)
		i = j
	}
}
