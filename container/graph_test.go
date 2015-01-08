package container

import (
	"testing"
)

func TestDeleteEdge(t *testing.T) {
	g := newGraph()
	a := &vertex{id: "A"}
	b := &vertex{id: "B"}
	c := &vertex{id: "C"}
	d := &vertex{id: "D"}
	g.addVertex(a)
	g.addVertex(b)
	g.addVertex(c)
	g.addVertex(d)
	//
	//A:->B,->C,->D
	//B:->A,->C,->D
	//C:->A,->B,->D
	//D:->A,->B,->C
	g.connect(a, b, 1)
	g.connect(a, c, 1)
	g.connect(a, d, 1)
}
