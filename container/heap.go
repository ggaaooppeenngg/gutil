package container

import ()

// 保证子节点比父结点小,用完全二叉树实现,只要一个数组.
// 默认是最小堆.
type Heap struct {
	heap []*edge
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
func swim(h Heap, j int) {
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
func sink(h Heap, i int) {
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
