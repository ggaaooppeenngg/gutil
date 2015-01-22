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
