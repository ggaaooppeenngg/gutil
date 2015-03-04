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
	//println(g.String())
}
func TestConnectedComponents(t *testing.T) {
	g := NewUndiGraph()
	a, b, c, d := &Vertex{Id: "A"}, &Vertex{Id: "B"}, &Vertex{Id: "C"}, &Vertex{Id: "D"}
	g.AddEdge(a, b)
	g.AddEdge(c, d)
	cc := g.CC()
	if !cc.Connected(a, b) || cc.Connected(a, d) {
		t.Fatal("CC failed")
	}
}

func TestCycle(t *testing.T) {
	g := NewUndiGraph()
	a, b, c, d := &Vertex{Id: "A"}, &Vertex{Id: "B"}, &Vertex{Id: "C"}, &Vertex{Id: "D"}
	g.AddEdge(a, b)
	g.AddEdge(b, c)
	g.AddEdge(c, d)
	g.AddEdge(d, a)
	if !g.HasCycle() {
		t.Fatal("cycle failed")
	}
}
