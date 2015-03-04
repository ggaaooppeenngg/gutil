package container

import (
	"testing"
)

func TestQueue(t *testing.T) {
	e := newEdge(newVertex("d"), 4.0)
	var es = []*edge{newEdge(newVertex("a"), 1.0),
		e,
		newEdge(newVertex("b"), 2.0),
		newEdge(newVertex("c"), 3.0),
	}
	pq := newPQ(es)
	if !pq.Contains(e.vtx) {
		t.Fatalf("do not get e %v", e)
	}
	for i := 0; i < 4; i++ {
		if i+1 != int(*(pq.Get().weight)) {
			t.Fatal("Is not priority queue")
		}
	}
}
