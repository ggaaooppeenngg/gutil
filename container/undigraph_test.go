package container

import (
	"testing"
)

func TestUndiGraph(t *testing.T) {
	g := NewUndiGraph()
	a := &Vertex{Id: "A"}
	b := &Vertex{Id: "B"}
	c := &Vertex{Id: "C"}
	d := &Vertex{Id: "D"}
	g.AddVertex(a)
	g.AddVertex(b)
	g.AddVertex(c)
	g.AddVertex(d)
	g.AddEdge(a, b)
	g.AddEdge(b, c)
	println(g.String())
}
