package container

import (
	"fmt"
	"testing"
)

func TestDiGraph(t *testing.T) {
	g := NewDiGraph()
	a, b, c, d := &Vertex{Id: "a"}, &Vertex{Id: "b"}, &Vertex{Id: "c"}, &Vertex{Id: "d"}
	g.AddEdge(a, b)
	g.AddEdge(b, c)
	g.AddEdge(c, d)
	fmt.Println(g)
}

func TestHasDiCycle(t *testing.T) {
	g := NewDiGraph()
	a, b, c, d := &Vertex{Id: "a"}, &Vertex{Id: "b"}, &Vertex{Id: "c"}, &Vertex{Id: "d"}
	g.AddEdge(a, b)
	g.AddEdge(b, c)
	g.AddEdge(c, b)
	g.AddVertex(d)
	/*
		cycle := g.HasDirectedCycle()
		if len(cycle) == 0 {
			t.Fail()
		}
	*/
}
