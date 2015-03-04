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

//test graph:
//A-4.000--------+
// -3.000+       |
//       |       |
//       v       |
//       C-3.000----------------+
//        -3.000+|              |
//       ^      ||              |
//       |      vv              |
//       |      B -3.000-+      |
//       |               |      |
//       |               v      |
//       +--------2.000-D -2.000--------+
//                      ^^      |       |
//                      ||      v       |
//                      |+2.000-E-1.000+|
//                      |              ||
//                      |              vv
//                      |              F
//                      |
//                      |
//                      +----------------2.000-G
//result:
//A->B->D
//A->C->E

func TestDSP(t *testing.T) {
	A, B, C, D, E, F, G := newVertex("A"), newVertex("B"), newVertex("C"), newVertex("D"), newVertex("E"), newVertex("F"), newVertex("G")
	g := NewDiGraph()
	g.AddEdge(A, B, 4)
	g.AddEdge(A, C, 3)

	g.AddEdge(B, D, 3)

	g.AddEdge(C, E, 3)
	g.AddEdge(C, B, 3)

	g.AddEdge(D, C, 2)
	g.AddEdge(D, F, 2)

	g.AddEdge(E, F, 1)
	g.AddEdge(E, D, 2)

	g.AddEdge(G, F, 2)
	p := g.DSP(A)
	if p.String() != `A -> B C
B -> D
C -> E
D ->
E -> F
F ->
` {
		t.Fatal("get wrong path")

	}
	t.Log(p)
}
