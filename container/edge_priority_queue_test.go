package container

import (
	"testing"
)

func TestQueue(t *testing.T) {
	e := newEdge(newVertext("d"), 4.0)
	var es = []*edge{newEdge(newVertext("a"), 1.0),
		e,
		newEdge(newVertext("b"), 2.0),
		newEdge(newVertext("c"), 3.0),
	}
	pq := newPQ(es)
	if !pq.Contains(e) {
		t.Fatalf("do not get e %v", e)
	}
	for i := 0; i < 4; i++ {
		if i+1 != int(*(pq.Get().weight)) {
			t.Fatal("Is not priority queue")
		}
	}
}
