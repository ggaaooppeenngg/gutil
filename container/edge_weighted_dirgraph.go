// 最短路径算法的比较:
// 算法名称          限制                    一般情况　 最差情况　内存 sweet spot
// Dijkstra(eager)   positive edge weights   ElogV      ElogV     V    worst-case guarantee
// topological 		 sort edge-weighted DAGs E+V        E+V       V    optimal for acyclic
// Bellman-Ford      (queue-based) no        E+V        VE        V    widely applicable
//                   negative cycles
package container

import (
	"fmt"
	"sort"
	"sync"

	"github.com/ggaaooppeenngg/util"
)

type EdgeWeightedDigraph struct {
	V             int // number of vertices
	E             int // number of edges
	Vertices      []*Vertex
	AdjacencyList map[*Vertex][]*DirectedEdge
	verticesId    map[string]bool
	sync          sync.RWMutex
}

// newEdgeWeightedDigraph returns new Edge Weighted Digraph.
func newEdgeWeightedDigraph() *EdgeWeightedDigraph {
	return &EdgeWeightedDigraph{
		AdjacencyList: make(map[*Vertex][]*DirectedEdge),
		verticesId:    make(map[string]bool),
	}
}
func (g EdgeWeightedDigraph) String() string {

	var output string
	vertices := vtxSlice(g.Vertices)
	sort.Sort(vertices)
	for _, v := range vertices {
		edges := g.Adj(v)
		output += v.Id + " ->"
		for _, edge := range edges {
			output += " " + edge.To.Id
		}
		output += "\n"
	}
	return output

}

//
func (g *EdgeWeightedDigraph) AddVertex(v *Vertex) error {
	g.sync.Lock()
	defer g.sync.Unlock()
	if _, ok := g.verticesId[v.Id]; ok {
		return fmt.Errorf("'%s' already exists", v.Id)
	}
	g.verticesId[v.Id] = true
	g.Vertices = append(g.Vertices, v)
	g.V++
	return nil
}

//
func (g *EdgeWeightedDigraph) GetVertexById(id string) *Vertex {
	for _, vtx := range g.Vertices {
		if vtx.Id == id {
			return vtx
		}
	}
	return nil
}

// AddEdge adds new edege to graph.
func (g *EdgeWeightedDigraph) AddEdge(from *Vertex, to *Vertex, weight float64) {
	var (
		e     *DirectedEdge
		edges []*DirectedEdge
		err   error
	)
	err = g.AddVertex(from)
	if err != nil {
		from = g.GetVertexById(from.Id)

	}
	err = g.AddVertex(to)
	if err != nil {
		to = g.GetVertexById(to.Id)
	}
	g.sync.Lock()
	defer g.sync.Unlock()
	e = &DirectedEdge{from, to, weight}
	if _, ok := g.AdjacencyList[from]; !ok {
		g.AdjacencyList[from] = []*DirectedEdge{e}
	} else {
		edges = g.AdjacencyList[from]
		edges = append(edges, e)
		g.AdjacencyList[from] = edges
	}
	g.E++
}

//
func (g *EdgeWeightedDigraph) Adj(v *Vertex) []*DirectedEdge {
	g.sync.RLock()
	defer g.sync.RUnlock()
	return g.AdjacencyList[v]
}

// 通过这条边的点.距离能够缩短就relax.
func (g *EdgeWeightedDigraph) relaxEdge(path *DirectedPath, e *DirectedEdge) {
	if path.DistTo[e.To] > path.DistTo[e.From]+e.Weight {
		path.DistTo[e.To] = path.DistTo[e.From] + e.Weight
		path.EdgeTo[e.To] = e
	}
}

// 通过这个点距离能够缩短就relax.
func (g *EdgeWeightedDigraph) relaxVertex(path *DirectedPath, v *Vertex) {
	edges := g.Adj(v)
	for _, e := range edges {
		w := e.To
		if path.DistTo[w] > path.DistTo[v]+e.Weight {
			path.DistTo[w] = path.DistTo[v] + e.Weight
			path.EdgeTo[w] = e
		}
	}
}

// 扩大通过这个点的距离
func (g *EdgeWeightedDigraph) tightVertex(path *DirectedPath, v *Vertex) {
	edges := g.Adj(v)
	for _, e := range edges {
		w := e.To
		if path.DistTo[w] < path.DistTo[v]+e.Weight {
			path.DistTo[w] = path.DistTo[v] + e.Weight
			path.EdgeTo[w] = e
		}
	}

}

// 拓扑排序
func (g *EdgeWeightedDigraph) TopSort() []*Vertex {
	var (
		marked      = make(map[*Vertex]bool)
		topVertices = make([]*Vertex, 0, g.V)
	)
	if g.HasCycle() {
		return nil
	}
	var vist = func(v *Vertex) {
		topVertices = append(topVertices, v)
	}
	if g.V > 0 {
		s := g.Vertices[0]
		g.dfsPost(marked, s, vist)
		util.Reverse(&topVertices)
	}
	return topVertices
}

// HasCycle探测是否成环
func (g *EdgeWeightedDigraph) HasCycle() bool {
	var (
		hasCycle = false
		onStack  = make(map[*Vertex]bool)
		marked   = make(map[*Vertex]bool)
	)
	var dfs func(g *EdgeWeightedDigraph, v *Vertex)
	dfs = func(g *EdgeWeightedDigraph, v *Vertex) {
		onStack[v] = true
		marked[v] = true
		for _, edge := range g.Adj(v) {
			if hasCycle {
				return
			} else if !marked[edge.To] {
				dfs(g, edge.To)
			} else if onStack[edge.To] {
				hasCycle = true
				return
			}
		}
		onStack[v] = false
	}
	return hasCycle
}

//
func (g *EdgeWeightedDigraph) dfsPost(marked map[*Vertex]bool, s *Vertex, post func(v *Vertex)) {
	edges := g.Adj(s)
	marked[s] = true
	for _, edge := range edges {
		if !marked[edge.To] {
			g.dfsPost(marked, edge.To, post)
		}
	}
	post(s)
}

// 拓扑排序的最短路径,因为寻找拓扑序的时候,是不断寻找入度为0的点,所以每次结点被访问的时候,之前的结点一定已经遍历过了.
// 时间复杂度是 O(E+V) 的.其实最短路径在拓扑排序里面,一旦一个结点被访问到,那么所有指向这个结点或者间接指向这个结点的路径都被访问过了.
// 那么之前的路径是最短路径的话,这条路径就是最短路径了,这个道理是一样的,拓扑排序可以减少不必要的遍历.
// 拓扑排序的时间负责度是 O(E+V),遍历的时候每条边都只relax一次,所以时间复杂度是O(V+E).
// ASP 是Acyclic Shortest Path 的缩写,也是一个单源最短路径算法.
func (g *EdgeWeightedDigraph) ASP(s *Vertex) *DirectedPath {
	s = g.GetVertexById(s.Id)
	path := NewDirectedPath(s)
	vs := g.TopSort()
	if vs != nil {
		path.InitPosInf(vs)
		path.DistTo[s] = 0.0
		for _, v := range vs {
			g.relaxVertex(path, v)
		}
	}
	return path
}

// 最长路径也是一样的,根据拓扑排序的规则,访问这个结点的时候,之前可能达到的结点都已经访问到了,那么最长路径只可能从这里面来,
// 类似于优化了的穷举法,这是典型的动态规划问题.
// ALP is Acyclic Longest Path.
func (g *EdgeWeightedDigraph) ALP(s *Vertex) *DirectedPath {
	s = g.GetVertexById(s.Id)
	path := NewDirectedPath(s)
	vs := g.TopSort()
	if vs != nil {
		path.InitNegInf(vs)
		path.DistTo[s] = 0.0
		for _, v := range vs {
			g.tightVertex(path, v)
		}
	}
	return path
}

//
// Bellman-Ford Algorithm 传统方式就是这样的,是书上的java代码,当伪代码看.
// for (int pass = 0; pass < G.V(); pass++)
// 	for (v = 0; v < G.V(); v++)
// 	 for (DirectedEdge e : G.adj(v))
// 		relax(e);
// 有三层循环,每次所有的边都会被遍历一次,遍历的次数等于点的个数.(内层两个循环是所有的边,最外面是点的次数)
// 时间复杂度是O(EV)(每个点都要relaxE条边),空间复杂度是O(V).
// proof:
// V0->V1->...->Vk,是一条最短路径,如果没有负环,这条路径是存在的.
// 假设 经过i次pass,这个算法可以算出V0到Vi的最短路径.
// 基本是 (i=0) V0 是一条最短路径.
// 推理是 如果V0->V1->...->Vi成立,然后Vi被relax之后,Vi+1的距离就不会超过Vi的距离加上边Vi->Vi+1的距离.
// 当然也不可能小于,因为它是最短路径.
// 整个算法代价很高,需要优化,原本的论文里面也提到了这个优化,现在讨论这个优化.
// 其实并不是说每次都要遍历所有的边,很容易发现,需要遍历的其实是那些之前导致自身distTo距离改变的结点.
// 这些结点的出度才是有效的,其它的还会保持原样,所以优化的过程就是基于这一点,依赖队列的结构实现.
// BFSP is Bellman-Ford alogrithm (queue based)
func (g *EdgeWeightedDigraph) BFSP(s *Vertex) *DirectedPath {
	path := NewDirectedPath(s)
	onQ := make(map[*Vertex]bool)
	queue := new(VertexQueue)
	path.InitPosInf(g.Vertices)
	path.DistTo[s] = 0.0
	onQ[s] = true
	cost := 0
	// TODO:负环检查不能统一.
	for !queue.Empty() && g.hasCycle() {
		v := queue.Front()
		queue.Pop()
		onQ[v] = false
		g.relaxOnQ(path, v, &cost, onQ, queue)
	}
	return nil
}

// relax on Queue is relax function for the BFSP
func (g *EdgeWeightedDigraph) relaxOnQ(path *DirectedPath, v *Vertex, cost *int, onQ map[*Vertex]bool, queue *VertexQueue) {
	for _, edge := range g.Adj(v) {
		if path.DistTo[edge.To] > path.DistTo[edge.From]+edge.Weight {
			path.DistTo[edge.To] = path.DistTo[edge.From] + edge.Weight
			path.EdgeTo[edge.To] = edge
			if !onQ[edge.To] {
				queue.Push(edge.To)
				onQ[edge.To] = true
			}
		}
		*cost++
		if *cost%g.V == 0 {
			g.findNegativeCycle()
		}
	}
}

// 检查和查找负环的相关函数.
func (g *EdgeWeightedDigraph) findNegativeCycle() {

}
func (g *EdgeWeightedDigraph) hasCycle() bool {
	return false
}
func (g *EdgeWeightedDigraph) NegativeCycle() []*Vertex {
	return nil
}
