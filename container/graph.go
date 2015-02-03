package container // import "github.com/ggaaooppeenngg/util/container"

//inspired by gyuho's goraph [https://github.com/gyuho/goraph/blob/refactor/graph/graph.go]

import (
	"fmt"
	"log"
	"sync"

	"github.com/ggaaooppeenngg/util"
)

// 当函数要修改数据的时候才用指针作为接收者.
// graph contains graph datas,implemented by adjacency list.
type Graph struct {
	Vertices []*Vertex

	sync.Mutex // map is not concurrency safe.

	edgeFrom  map[*Vertex][]*edge // outBound maps vertex to outgoing edges
	edgeTo    map[*Vertex][]*edge // inBound maps vertex to outgoing edges
	vertexIDs map[string]bool     // map[id]bool,prevent duplicating IDs
}

// newGraph inits a new graph.
func NewGraph() *Graph {
	return &Graph{
		Vertices:  []*Vertex{},
		edgeFrom:  make(map[*Vertex][]*edge),
		edgeTo:    make(map[*Vertex][]*edge),
		vertexIDs: make(map[string]bool),
	}
}

//vertex is a vertex(node) in graph
type Vertex struct {
	Id  string //for debug
	Ele interface{}
}

//newVertex inits a new vertext.
func newVertex(id string) *Vertex {
	return &Vertex{
		Id: id,
	}
}

//edge from or to vtx.
type edge struct {
	vtx    *Vertex  //source or destination
	weight *float64 //default nil, because default folat64!=0
}

//weight gets dereference a float pointer
func (e *edge) Weight() (w float64, has bool) {
	if e.weight == nil {
		return 0.0, false
	} else {
		return *e.weight, true
	}
}

//Set weight
func (e *edge) SetWeight(w float64) {
	if e.weight == nil {
		e.weight = new(float64)
	}
	*e.weight = w
}

//newEdge returns a new edge with vertex and weight.
func newEdge(v *Vertex, w float64) *edge {
	wp := new(float64)
	*wp = w
	return &edge{
		v,
		wp,
	}
}

func (e edge) String() string {
	w, _ := e.Weight()
	return fmt.Sprintf("%s:%f", e.vtx.Id, w)
}

//for debug
func (g Graph) String() string {
	var output string
	for _, vtx := range g.Vertices {
		output += fmt.Sprintf("[%s]:\n", vtx.Id)
		for _, edge := range g.edgeFrom[vtx] {
			output += fmt.Sprintf("-> %f %s\n", *edge.weight, edge.vtx.Id)
		}
	}
	return output
}

func (g *Graph) AddVertex(vtx *Vertex) error {
	g.Mutex.Lock()
	if _, ok := g.vertexIDs[vtx.Id]; ok {
		g.Mutex.Unlock()
		return fmt.Errorf("'%s' already exists", vtx.Id)
	}
	//map 不是并发安全的
	g.vertexIDs[vtx.Id] = true
	g.Mutex.Unlock()
	g.Vertices = append(g.Vertices, vtx)
	return nil
}

func (g *Graph) DeleteVertext(vtx *Vertex) {
	// delete vertex from Vertices
	for idx, vertex := range g.Vertices {
		if vertex == vtx {
			util.DeleteIndexOf(idx, &g.Vertices)
			break
		}
	}
	defer g.Unlock()
	g.Lock()
	// delete coresponding edge.
	for _, fromEdge1 := range g.edgeTo[vtx] {
		// 删掉入度对应的边.
		// (from) -> x
		for idx, fromEdge2 := range g.edgeFrom[fromEdge1.vtx] {
			if fromEdge2.vtx == vtx {
				edges := g.edgeFrom[fromEdge2.vtx]
				util.DeleteIndexOf(idx, &edges)
			}
		}
	}

	for _, toEdge1 := range g.edgeFrom[vtx] {
		// 删掉出度对应的边
		// (x) -> to
		for idx, toEdge2 := range g.edgeTo[toEdge1.vtx] {
			if toEdge2.vtx == vtx {
				edges := g.edgeTo[toEdge1.vtx]
				util.DeleteIndexOf(idx, &edges)
			}
		}
	}

	//删掉自己村边的两个edges
	delete(g.edgeTo, vtx)
	delete(g.edgeFrom, vtx)
	delete(g.vertexIDs, vtx.Id)
}

func (g Graph) FindVertexById(id string) *Vertex {
	for _, vtx := range g.Vertices {
		if vtx.Id == id {
			return vtx
		}
	}
	return nil
}

func (g *Graph) Connect(src, dst *Vertex, weights ...float64) {
	if len(weights) > 1 {
		log.Printf("only one weight allowed")
		return
	}
	var weight *float64
	if len(weights) == 1 {
		weight = new(float64)
		*weight = weights[0]
	}
	err := g.AddVertex(src)
	if err != nil {
		//log.Printf("'%s' was previously added to graph\n", src.Id)
		src = g.FindVertexById(src.Id)
	} else {
		//log.Printf("'%s' is added to graph\n", src.Id)
	}
	err = g.AddVertex(dst)
	if err != nil {
		//log.Printf("'%s' was previously added to graph\n", dst.Id)
		dst = g.FindVertexById(dst.Id)
	} else {
		//log.Printf("'%s' is added to graph\n", dst.Id)
	}
	edgeSrc := &edge{
		vtx:    src,
		weight: weight,
	}
	edgeDst := &edge{
		vtx:    dst,
		weight: weight,
	}
	defer g.Unlock()
	g.Mutex.Lock()
	if _, ok := g.edgeFrom[src]; !ok {
		g.edgeFrom[src] = []*edge{edgeDst}
	} else {
		isDuplicate := false
		for _, edge := range g.edgeFrom[src] {
			if edge.vtx == dst {
				log.Println("Duplicate(Parallel) Edge Found. Overwriting the Weight value.")
				log.Printf("%v --> %v + %v\n", edge.weight, edge.weight, weight)
				if edge.weight == nil {
					edge.weight = new(float64)
				}
				*edge.weight += *weight
				isDuplicate = true
				break
			}
		}
		if !isDuplicate {
			g.edgeFrom[src] = append(g.edgeFrom[src], edgeDst)
		}
	}
	if _, ok := g.edgeTo[dst]; !ok {
		g.edgeTo[dst] = []*edge{edgeSrc}
	} else {
		isDuplicate := false
		for _, edge := range g.edgeTo[dst] {
			if edge.vtx == src {
				log.Println("Duplicate(Parallel) Edge Found. Overwriting the Weight value.")
				log.Printf("%v --> %v + %v\n", edge.weight, edge.weight, weight)
				if edge.weight == nil {
					edge.weight = new(float64)
				}
				*edge.weight += *weight
				isDuplicate = true
				break
			}
		}
		if !isDuplicate {
			g.edgeTo[dst] = append(g.edgeTo[dst], edgeSrc)
		}
	}
}

//delete a edge
//same pointer,not the same id.
func (g *Graph) DeleteEdge(src, dst *Vertex) {
	defer g.Unlock()
	g.Lock()

	//delete edge from src
	for idx, edge := range g.edgeFrom[src] {
		if edge.vtx == dst {
			edges := g.edgeFrom[src]
			util.DeleteIndexOf(idx, &edges)
			break
		}
	}

	//delete edge to dst
	for idx, edge := range g.edgeTo[dst] {
		if edge.vtx == src {
			edges := g.edgeTo[dst]
			util.DeleteIndexOf(idx, &edges)
			break
		}
	}

}

//vertextSize returns number of vertex.
func (g Graph) VertexSize() int {
	return len(g.Vertices)
}

//EdgeExist returns true if edge from src to dst exists,else returns false
func (g Graph) EdgeExists(src, dst *Vertex) bool {
	defer g.Unlock()
	g.Lock()
	for _, edge := range g.edgeFrom[src] {
		if edge.vtx == dst {
			return true
		}
	}
	return false
}

//getEdgeWeight returns weight value of an edge from src to dst.
func (g Graph) GetEdgeWeight(src, dst *Vertex) *float64 {
	defer g.Unlock()
	g.Lock()
	for _, edge := range g.edgeFrom[src] {
		if edge.vtx == dst {
			return edge.weight
		}
	}
	return nil
}

//OutBound returns the outbound verticies for vtx
func (g Graph) OutBound(vtx *Vertex) []*Vertex {
	var outbound []*Vertex
	for _, e := range g.edgeFrom[vtx] {
		outbound = append(outbound, e.vtx)
	}
	return outbound
}

//InBound returns the inboud verticies for vtx
func (g Graph) InBound(vtx *Vertex) []*Vertex {
	var inbound []*Vertex
	for _, e := range g.edgeTo[vtx] {
		inbound = append(inbound, e.vtx)
	}
	return inbound
}
