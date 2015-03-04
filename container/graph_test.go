package container

import (
	"testing"
)

func TestDeleteEdge(t *testing.T) {
	g := NewGraph()
	a := &Vertex{Id: "A"}
	b := &Vertex{Id: "B"}
	c := &Vertex{Id: "C"}
	d := &Vertex{Id: "D"}
	g.AddVertex(a)
	g.AddVertex(b)
	g.AddVertex(c)
	g.AddVertex(d)
	//
	//A:->B,->C,->D
	//B:->A,->C,->D
	//C:->A,->B,->D
	//D:->A,->B,->C
	g.Connect(a, b, 1)
	g.Connect(a, c, 1)
	g.Connect(a, d, 1)
}
