package container

import (
	"testing"
)

func TestGraph(t *testing.T) {
	g := newGraph()
	a := &vertex{id: "A"}
	g.addVertex(a)
	g.connect(&vertex{id: "A"}, &vertex{id: "B"}, 1)
	g.connect(&vertex{id: "B"}, &vertex{id: "A"}, 1)
	c := &vertex{id: "C"}
	g.connect(&vertex{id: "A"}, c, 1)
	g.deleteEdge(a, c)
	t.Log(g)
}
