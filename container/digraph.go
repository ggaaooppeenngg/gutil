package container

import (
	"errors"
	"fmt"
	"sync"
)

//DiGraph is the representation fo directed graph.
//underlying datastructure is adjancency lists of the vertex out degree.
//Adj gets the vertex out degree,use Reverse to get
//DiGraph with adjancency lists storing indegree in lists.
type DiGraph struct {
	V          int                 // number of vertices
	E          int                 // number of edges
	Vertices   []*Vertex           // adjacency-lists
	verticesId map[string]bool     // id map in order to prevent duplicate ids.
	Edges      map[*Vertex][]*edge // edge list
	sync       sync.RWMutex        // map is not cunrrency safe
}

//NewDigraph returns new directed graph.
func NewDiGraph() *DiGraph {
	return &DiGraph{
		Edges:      make(map[*Vertex][]*edge),
		verticesId: make(map[string]bool),
	}
}

//AddVertex adds a Vertex in graph,returns error if vertex exists.
func (g *DiGraph) AddVertex(v *Vertex) error {
	g.sync.Lock()
	defer g.sync.Unlock()
	if _, ok := g.verticesId[v.Id]; ok {
		return fmt.Errorf("'%s' already exists", v.Id)
	}
	g.verticesId[v.Id] = true
	g.Vertices = append(g.Vertices, v)
	return nil
}

//GetVertexById returns *Vertex if id matches.
func (g *DiGraph) GetVertexById(id string) *Vertex {
	for _, vtx := range g.Vertices {
		if vtx.Id == id {
			return vtx
		}
	}
	return nil
}

func (g *DiGraph) String() string {
	var output string
	for _, v := range g.Vertices {
		edges := g.Adj(v)
		if len(edges) == 0 {
			break
		}
		output += v.Id + " ->"
		for _, edge := range edges {
			output += " " + edge.vtx.Id
		}
		output += "\n"
	}
	return output
}

//AddEdge from v to w,default weight is 0.
func (g *DiGraph) AddEdge(v *Vertex, w *Vertex, weights ...float64) error {
	if len(weights) > 1 {
		return errors.New("more than one weight is not allowed")
	}
	var weight *float64
	if len(weights) == 1 {
		weight = new(float64)
		*weight = weights[0]
	}

	var (
		edges []*edge
		e     *edge
		err   error
	)
	//if vertex is duplicate,use the original vertex instead.
	err = g.AddVertex(v)
	if err != nil {
		v = g.GetVertexById(v.Id)
	}
	err = g.AddVertex(w)
	if err != nil {
		w = g.GetVertexById(w.Id)
	}

	g.sync.Lock()
	defer g.sync.Unlock()
	//v add edge v-w
	e = &edge{w, weight}
	if _, ok := g.Edges[v]; !ok {
		g.Edges[v] = []*edge{e}
	} else {
		edges = g.Edges[v]
		edges = append(edges, e)
		g.Edges[v] = edges
	}
	g.E++
	return nil
}

//Adj returns adjancent vertices of v.
func (g *DiGraph) Adj(v *Vertex) (edges []*edge) {
	g.sync.RLock()
	defer g.sync.RUnlock()
	edges = g.Edges[v]
	return
}

//dfs mutile sources.
//这个可以用在标记清扫的垃圾回收中,
//每个object有一位标记位,当dfs不能到达说明,没有被引用,会被标记,
//以便清扫.

func (g *DiGraph) DirectedDFS(sources []*Vertex) {
	var (
		marked = make(map[*Vertex]bool)
	)
	for _, s := range sources {
		g.dfs(s, nil, marked, nil)
	}
}

//另外实现一个dfs
func (g *DiGraph) HasDirectedCycle() bool {
	var hasCycle = false
	for _, v := range g.Vertices {
		if hasCycle {
			return hasCycle
		}
		var (
			onStack = make(map[*Vertex]bool)
			prev    = new(Vertex)
		)
		g.dfs(v, nil, onStack, func(v *Vertex) {

		})
	}
	return hasCycle
}

//walks the graph by dfs.
func (g *DiGraph) dfs(s *Vertex, count *int, marked map[*Vertex]bool, walk func(v *Vertex)) {
	marked[s] = true
	if walk != nil {
		walk(s)
	}
	if count != nil {
		*count++
	}
	edges := g.Adj(s)
	for _, edge := range edges {
		if !marked[edge.vtx] {
			g.dfs(edge.vtx, count, marked, walk)
		}
	}
}

//Reverse returns copy of the graph but with adjancency lists of in degree.
func (g *DiGraph) Reverse() *DiGraph {
	rg := NewDiGraph()
	for _, v := range g.Vertices {
		edges := g.Adj(v)
		for _, edge := range edges {
			rg.AddEdge(edge.vtx, v)
		}
	}
	return rg
}
