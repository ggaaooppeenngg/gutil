package container

import (
	"math"
	"sort"
	"strconv"
	"strings"

	"github.com/ggaaooppeenngg/util"
)

// Path is a representation of paths from source vertex.
// single source path with simplex edge.
type Path struct {
	Source *Vertex
	edgeTo map[*Vertex]*edge
	distTo map[*Vertex]float64
}

// DirectedPath is the representation of directed path.
type DirectedPath struct {
	Source *Vertex
	EdgeTo map[*Vertex]*DirectedEdge
	DistTo map[*Vertex]float64
}

//
func NewDirectedPath(s *Vertex) *DirectedPath {
	return &DirectedPath{
		s,
		make(map[*Vertex]*DirectedEdge),
		make(map[*Vertex]float64),
	}
}

func (p DirectedPath) PathString() string {
	var outputs []string
	for v, _ := range p.DistTo {
		var (
			final = v
		)
		var paths []string
		for {
			if edge, ok := p.EdgeTo[v]; ok {
				weight := strconv.FormatFloat(edge.Weight, 'g', 3, 64)
				paths = append(paths, edge.From.Id+"->"+edge.To.Id+":"+weight)
				v = edge.From
			} else {
				util.Reverse(&paths)
				path := strings.Join(paths, " ")
				outputs = append(outputs, "To "+final.Id+" = "+path)
				break
			}
		}
	}
	sort.Strings(outputs)
	return strings.Join(outputs, "\n")
}

func (p DirectedPath) String() string {
	g := newEdgeWeightedDigraph()
	for _, v := range p.EdgeTo {
		g.AddEdge(v.From, v.To, v.Weight)
	}
	return g.String()
}

// InitPosInf initates the DistTo to positive inifinity.
func (dp *DirectedPath) InitPosInf(vertices []*Vertex) {
	for _, v := range vertices {
		dp.DistTo[v] = math.MaxFloat64
	}
}

// InitNegInf initates the DistTo to negtive inifinity.
func (dp *DirectedPath) InitNegInf(vertices []*Vertex) {
	for _, v := range vertices {
		dp.DistTo[v] = -math.MaxFloat64
	}
}

// newPath returns a initiated Path.
func newPath(s *Vertex) *Path {
	return &Path{
		s,
		make(map[*Vertex]*edge),
		make(map[*Vertex]float64),
	}
}

func (p *Path) String() string {
	g := NewDiGraph()
	for k, v := range p.edgeTo {
		w, _ := v.Weight()
		g.AddEdge(v.vtx, k, w)
	}
	return g.String()
}

//HasPathto return true,if Path reachs d.
func (p *Path) HasPathTo(d *Vertex) bool {
	if _, ok := p.edgeTo[d]; ok {
		return true
	} else {
		return false
	}
}
