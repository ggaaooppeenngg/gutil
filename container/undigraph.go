package container

import (
	"sync"
)

//UndiGraph is undirected graph,implemented in an array of adjacency lists.
//1.  adjacency matrix can not represent parallel edges,and also space wasting.
//2.  array of edges(with 2 vertecies) will be examined when getting the adjant vertecies.
//Put all into considerations,an array of adjacency lists is the best option.
type UndiGraph struct {
	V        int                 //number of vertices
	E        int                 //number of edges
	sync     sync.RWMutex        // map is not cunrrency safe
	Vertices []*Vertex           //adjacency-lists
	Edges    map[*Vertex][]*edge //edge list
}

//NewUndiGraph returns an empty undirected graph
func NewUndiGraph() *UndiGraph {
	g := &UndiGraph{
		Edges: make(map[*Vertex][]*edge),
	}
	return g
}

//AddEdge adds edge v-w
func (g *UndiGraph) AddEdge(v *Vertex, w *Vertex) {
	var edges []*edge
	var e *edge

	g.sync.Lock()
	//v add edge v-w
	edges = g.Edges[v]
	e = &edge{w, new(float64)}
	edges = append(edges, e)
	g.Edges[v] = edges
	//w add edge w-v
	edges = g.Edges[w]
	e = &edge{v, new(float64)}
	edges = append(edges, e)
	g.Edges[w] = edges
	g.sync.Unlock()

	g.E++
	return
}

//Adj returns adjancent vertices of v.
func (g *UndiGraph) Adj(v *Vertex) (edges []*edge) {
	g.sync.RLock()
	defer g.sync.RUnlock()
	edges = g.Edges[v]
	return
}

//DFS is depth first search,it returns vertices that counts.
func (g *UndiGraph) DFS(s *Vertex) int {
	return dfs(g, s)
}

//dfs implements recursive DFS search method.
func dfs(g *UndiGraph, s *Vertex) int {

	var (
		marked = make(map[*Vertex]bool)
		count  int
	)

	var _dfs func(g *UndiGraph, s *Vertex) int
	_dfs = func(g *UndiGraph, s *Vertex) int {
		marked[s] = true
		count++
		edges := g.Adj(s)
		for _, edge := range edges {
			if !marked[edge.vtx] {
				_dfs(g, edge.vtx)
			}
		}
		return count
	}
	_dfs(g, s)
	return count
}

//Path returns the path from s to d in order.
func (g *UndiGraph) Path(s *Vertex, d *Vertex) []*edge {

	return nil
}

//HasPath checks if there is a path from s to d.
func (g *UndiGraph) HasPath(s *Vertex, d *Vertex) bool {
	return false
}
