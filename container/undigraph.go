package container

import (
	"container/list"
	"errors"
	"fmt"
	"sync"
)

//UndiGraph is undirected graph,implemented in an array of adjacency lists.

type UndiGraph struct {
	V          int                 // number of vertices
	E          int                 // number of edges
	sync       sync.RWMutex        // map is not cunrrency safe
	Vertices   []*Vertex           // adjacency-lists
	verticesId map[string]bool     // id map in order to prevent duplicate ids.
	Edges      map[*Vertex][]*edge // edge list
}

//Path is a representation of paths from source vertex.
type Path struct {
	Source *Vertex
	edgeTo map[*Vertex]*edge
}

//HasPathto return true,if Path reachs d.
func (p *Path) HasPathTo(d *Vertex) bool {
	if _, ok := p.edgeTo[d]; ok {
		return true
	} else {
		return false
	}
}

//NewUndiGraph returns an empty undirected graph
func NewUndiGraph() *UndiGraph {
	g := &UndiGraph{
		Edges: make(map[*Vertex][]*edge),
	}
	return g
}

//TODO:optimize search method.
//GetVertexById returns *Vertex if id matches.
func (g *UndiGraph) GetVertexById(id string) *Vertex {
	for _, vtx := range g.Vertices {
		if vtx.Id == id {
			return vtx
		}
	}
	return nil
}

//AddVertex adds v to graph if duplicate return error
func (g *UndiGraph) AddVertex(v *Vertex) error {
	g.sync.Lock()
	defer g.sync.Unlock()
	if _, ok := g.verticesId[v.Id]; ok {
		return fmt.Errorf("'%s' already exists", v.Id)
	}
	g.verticesId[v.Id] = true
	g.Vertices = append(g.Vertices, v)
	return nil
}

//AddEdge adds edge v-w,parallel edges is allowed
func (g *UndiGraph) AddEdge(v *Vertex, w *Vertex, weights ...float64) error {
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

	//v add edge v-w
	e = &edge{w, weight}
	if _, ok := g.Edges[v]; !ok {
		g.Edges[v] = []*edge{e}
	} else {
		edges = g.Edges[v]
		edges = append(edges, e)
		g.Edges[v] = edges
	}

	//w add edge w-v
	e = &edge{v, weight}
	if _, ok := g.Edges[w]; !ok {
		g.Edges[w] = []*edge{e}
	} else {
		edges = g.Edges[w]
		edges = append(edges, e)
		g.Edges[w] = edges
	}

	g.sync.Unlock()

	g.E++
	return nil
}

//GetEdgeWeight returns weight from s to d,if edges are parallal,weights are accumulated.
func (g *UndiGraph) GetEdgeWeight(s, d *Vertex) *float64 {
	g.sync.RLock()
	defer g.sync.RUnlock()
	var weight *float64
	for _, edge := range g.Edges[s] {
		if edge.vtx == d {
			if weight != nil {
				weight = new(float64)
			}
			*weight += *edge.weight
		}
	}
	return weight
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
	var (
		marked = make(map[*Vertex]bool)
		count  int //vertecies counter
	)
	dfs(g, s, &count, marked, nil)
	return count
}

//GetPathDFS returns the path from s in order in DFS.
func (g *UndiGraph) GetPathDFS(s *Vertex) *Path {
	var (
		marked       = make(map[*Vertex]bool)
		edgeTo       = make(map[*Vertex]*edge)
		tmp    *edge = nil
	)
	dfs(g, s, nil, marked, func(v *Vertex) {
		if tmp != nil {
			edgeTo[v] = tmp
		}
		tmp = &edge{v, nil}
	})
	return &Path{s, edgeTo}
}

//GetPathBFS returns the path from s in BFS.
func (g *UndiGraph) GetPathBFS(s *Vertex) *Path {
	var (
		tmp    *edge
		edgeTo = make(map[*Vertex]*edge)
	)
	bfs(g, s, func(v *Vertex) {
		if tmp != nil {
			edgeTo[v] = tmp
		}
		tmp = &edge{v, nil}
	})
	return &Path{s, edgeTo}
}

//dfs implements recursive DFS search method.
func dfs(g *UndiGraph, s *Vertex, count *int, marked map[*Vertex]bool, walk func(v *Vertex)) {
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
			dfs(g, edge.vtx, count, marked, walk)
		}
	}
}

//bfs implements iterative bfs search method.
func bfs(g *UndiGraph, s *Vertex, walk func(v *Vertex)) int {
	var (
		count  int
		marked = make(map[*Vertex]bool)
		queue  = list.New()
	)
	if s == nil {
		return 0
	}

	marked[s] = true
	if walk != nil {
		walk(s)
	}
	count++
	queue.PushBack(s)

	for queue.Len() > 0 {
		ele := queue.Front()
		v := ele.Value.(*Vertex)
		edges := g.Adj(v)
		for _, edge := range edges {
			if !marked[edge.vtx] {
				marked[edge.vtx] = true
				if walk != nil {
					walk(edge.vtx)
				}
				count++
				queue.PushBack(edge.vtx)
			}
		}
	}
	return count
}
