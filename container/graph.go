package container // import "github.com/ggaaooppeenngg/util/container"

//inspired by gyuho's goraph [https://github.com/gyuho/goraph/blob/refactor/graph/graph.go]

import (
	"fmt"
	"log"
	"sync"
)

// 当函数要修改数据的时候才用指针作为接收者.
// graph contains graph datas,implemented by adjacency list.
type graph struct {
	vertices []*vertex

	sync.Mutex // map is not concurrency safe.

	edgeFrom  map[*vertex][]*edge // outBound maps vertex to outgoing edges
	edgeTo    map[*vertex][]*edge // inBound maps vertex to outgoing edges
	vertexIDs map[string]bool     // map[id]bool,prevent duplicating IDs
}

// newGraph inits a new graph.
func newGraph() *graph {
	return &graph{
		vertices:  []*vertex{},
		edgeFrom:  make(map[*vertex][]*edge),
		edgeTo:    make(map[*vertex][]*edge),
		vertexIDs: make(map[string]bool),
	}
}

//vertex is a vertex(node) in graph
type vertex struct {
	id  string //for debug
	ele interface{}
}

//newVertex inits a new vertext.
func newVertext(id string) *vertex {
	return &vertex{
		id: id,
	}
}

//edge from or to vtx.
type edge struct {
	vtx    *vertex //source or destination
	weight float64 //default 0
}

//for debug
func (g graph) String() string {
	var output string
	for _, vtx := range g.vertices {
		output += fmt.Sprintf("[%s]:\n", vtx.id)
		for _, edge := range g.edgeFrom[vtx] {
			output += fmt.Sprintf("-> %f %s\n", edge.weight, edge.vtx.id)
		}
	}
	return output
}

func (g *graph) addVertex(vtx *vertex) error {
	g.Mutex.Lock()
	if _, ok := g.vertexIDs[vtx.id]; ok {
		g.Mutex.Unlock()
		return fmt.Errorf("'%s' already exists", vtx.id)
	}
	//map 不是并发安全的
	g.vertexIDs[vtx.id] = true
	g.Mutex.Unlock()
	g.vertices = append(g.vertices, vtx)
	return nil
}
func (g *graph) deleteVertext(vtx *vertex) {
	for _, vertex := range g.vertices {
		if vertex == vtx {

		}
	}
}

func (g graph) findVertexById(id string) *vertex {
	for _, vtx := range g.vertices {
		if vtx.id == id {
			return vtx
		}
	}
	return nil
}

func (g *graph) connect(src, dst *vertex, weight float64) {
	err := g.addVertex(src)
	if err != nil {
		log.Printf("'%s' was previously added to graph\n", src.id)
		src = g.findVertexById(src.id)
	} else {
		log.Printf("'%s' is added to graph\n", src.id)
	}
	err = g.addVertex(dst)
	if err != nil {
		log.Printf("'%s' was previously added to graph\n", dst.id)
		dst = g.findVertexById(dst.id)
	} else {
		log.Printf("'%s' is added to graph\n", dst.id)
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
				edge.weight += weight
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
				edge.weight += weight
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
func (g *graph) deleteEdge(src, dst *vertex) {
	defer g.Unlock()
	g.Lock()

	//delete edge from src
	for idx, edge := range g.edgeFrom[src] {
		if edge.vtx == dst {
			fmt.Println("delete vtx")
			if idx == len(g.edgeFrom[src])-1 {
				g.edgeFrom[src] = g.edgeFrom[src][:idx]
			} else {
				g.edgeFrom[src] = append(g.edgeFrom[src][:idx], g.edgeFrom[src][idx+1:]...)
			}
			break
		}
	}

	//delete edge to dst
	for idx, edge := range g.edgeTo[dst] {
		if edge.vtx == src {
			if idx == len(g.edgeTo[dst])-1 {
				g.edgeTo[dst] = g.edgeTo[dst][:idx]
			} else {
				g.edgeTo[dst] = append(g.edgeTo[dst][:idx], g.edgeTo[dst][idx+1:]...)
			}
			break
		}
	}

}

//vertextSize returns number of vertex.
func (g graph) vertexSize() int {
	return len(g.vertices)
}

//getEdgeWeight returns weight value of an edge from src to dst.
func (g graph) getEdgeWeight(src, dst *vertex) float64 {
	defer g.Unlock()
	g.Lock()
	for _, edge := range g.edgeFrom[src] {
		if edge.vtx == dst {
			return edge.weight
		}
	}
	return 0.0
}
