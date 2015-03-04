package container

import (
	"fmt"
)

// Edge represents an directed edge.
type DirectedEdge struct {
	From   *Vertex
	To     *Vertex
	Weight float64
}

// for debug
func (de *DirectedEdge) String() string {
	return fmt.Sprintf("%v->%v:%.2f", de.From, de.To, de.Weight)
}

// SimplexEdge only stores one vertex.
type SimplexEdge struct {
	Vtx    *Vertex
	Weight float64
}

// for debug
func (se *SimplexEdge) String() string {
	return fmt.Sprintf("-%v:%2.f", se.Vtx, se.Weight)
}
