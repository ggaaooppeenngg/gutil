package container

import (
	"errors"
	"fmt"
	"math"
	"sort"
	"sync"

	"github.com/ggaaooppeenngg/util"
)

type Order uint

const (
	PRE  Order = 1 << iota // preorder
	POST                   // postorder
	REV                    // reverse order
)

// DiGraph is the representation fo directed graph.
// underlying datastructure is adjancency lists of the vertex out degree.
// Adj gets the vertex out degree,use Reverse to get
// DiGraph with adjancency lists storing indegree in lists.
type DiGraph struct {
	V          int                 // number of vertices
	E          int                 // number of edges
	Vertices   []*Vertex           // adjacency-lists
	verticesId map[string]bool     // id map in order to prevent duplicate ids.
	Edges      map[*Vertex][]*edge // edge list
	sync       sync.RWMutex        // map is not cunrrency safe
}

// NewDigraph returns new directed graph.
func NewDiGraph() *DiGraph {
	return &DiGraph{
		Edges:      make(map[*Vertex][]*edge),
		verticesId: make(map[string]bool),
	}
}

// AddVertex adds a Vertex in graph,returns error if vertex exists.
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

// GetVertexById returns *Vertex if id matches.
func (g *DiGraph) GetVertexById(id string) *Vertex {
	for _, vtx := range g.Vertices {
		if vtx.Id == id {
			return vtx
		}
	}
	return nil
}

type vtxSlice []*Vertex

func (vs vtxSlice) Less(i, j int) bool { return vs[i].Id < vs[j].Id }
func (vs vtxSlice) Len() int           { return len(vs) }
func (vs vtxSlice) Swap(i, j int)      { vs[i], vs[j] = vs[j], vs[i] }

func (g *DiGraph) String() string {
	var output string
	vertices := vtxSlice(g.Vertices)
	sort.Sort(vertices)
	for _, v := range vertices {
		edges := g.Adj(v)
		output += v.Id + " ->"
		for _, edge := range edges {
			output += " " + edge.vtx.Id
		}
		output += "\n"
	}
	return output
}

// AddEdge from v to w,default weight is 0.
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
		//TODO:有必要重新赋值么?
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

// Adj returns adjancent vertices of v.
func (g *DiGraph) Adj(v *Vertex) (edges []*edge) {
	g.sync.RLock()
	defer g.sync.RUnlock()
	edges = g.Edges[v]
	return
}

// dfs mutile sources.
// 这个可以用在标记清扫的垃圾回收中,
// 每个object有一位标记位,当dfs不能到达说明,没有被引用,会被标记,
// 以便清扫.
func (g *DiGraph) DirectedDFS(source *Vertex) []*Vertex {
	var (
		marked   = make(map[*Vertex]bool)
		vertices []*Vertex
	)
	walk := func(v *Vertex) {
		vertices = append(vertices, v)
	}
	g.dfs(source, nil, marked, walk, nil)
	return vertices
}

// DFSInOrder 按照order的顺序返回DFS遍历结果.
func (g *DiGraph) DFSInOrder(order Order) []*Vertex {
	if len(g.Vertices) == 0 {
		return g.Vertices
	}
	var (
		s      = g.Vertices[0]
		result []*Vertex
		marked = make(map[*Vertex]bool)
	)
	if order&PRE == PRE {
		g.dfs(s, nil, marked, func(v *Vertex) {
			result = append(result, v)
		}, nil)
		return result
	}
	if order&POST == POST {
		g.dfs(s, nil, marked, nil, func(v *Vertex) {
			result = append(result, v)
		})
	}
	if order&REV == REV {
		util.Reverse(&result)
	}
	return result
}

// TopSort returns vertices in topological sort order(post reverse),if
// there is cycle returns nil.
func (g *DiGraph) TopSort() []*Vertex {
	if len(g.DirectedCycle()) != 0 {
		return nil
	}
	//后序遍历是个回溯的过程,dfs(v)发生在dfs(w)之前,w就要指向v.
	//然后把序列倒过来就是一个拓扑排序.
	//直觉上说,DFS的后序遍历就是拓扑排序.
	return g.DFSInOrder(POST | REV)
}

// 另外实现一个dfs
// 可以用stack表示当前扫描的路劲,用onStack标记这些点,
// 当下个点存在onStack的时候,说明成环.
// 返回的cycle数组,头尾相等,如果无环长度为0.
func (g *DiGraph) DirectedCycle() []*Vertex {
	var (
		onStack  = make(map[*Vertex]bool)
		marked   = make(map[*Vertex]bool)
		vertexTo = make(map[*Vertex]*Vertex)
		hasCycle bool
		cycle    []*Vertex
	)

	// dfs
	var dfs func(g *DiGraph, v *Vertex)
	dfs = func(g *DiGraph, v *Vertex) {
		onStack[v] = true
		marked[v] = true
		for _, edge := range g.Adj(v) {
			if hasCycle {
				return
			} else if !marked[edge.vtx] {
				vertexTo[edge.vtx] = v
				dfs(g, edge.vtx)
			} else if onStack[edge.vtx] {
				for x := v; x != edge.vtx; x = vertexTo[x] {
					cycle = append([]*Vertex{x}, cycle...)
				}
				cycle = append([]*Vertex{v, edge.vtx}, cycle...)
				hasCycle = true
				return
			}
		}
		// 从扫描路径上退栈.
		onStack[v] = false
	}

	for _, v := range g.Vertices {
		if !marked[v] {
			dfs(g, v)
		}
	}
	return cycle
}

// walks the graph by dfs.
func (g *DiGraph) dfs(s *Vertex, count *int, marked map[*Vertex]bool, prewalk func(v *Vertex), postwalk func(v *Vertex)) {
	marked[s] = true
	//先序
	if prewalk != nil {
		prewalk(s)
	}
	if count != nil {
		*count++
	}
	edges := g.Adj(s)
	for _, edge := range edges {
		if !marked[edge.vtx] {
			g.dfs(edge.vtx, count, marked, prewalk, postwalk)
		}
	}
	//后序
	if postwalk != nil {
		postwalk(s)
	}
}

// Reverse returns copy of the graph but with adjancency lists of in degree.
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

//-------------------------------path algorithm--------------------------
// 最长路径
func (g *DiGraph) LP(s *Vertex) *Path {
	//先拓扑排序,然后遍历,用动态规划.
	return nil
}

// 权重为负的最短路径.
// 必须存在一条路径,并且该路径不能含有负环,不然能一直减小,一个点在负环上,可以当做负无穷大.
// 所以算法要负责探测负环,并计算最短路径.
// Queue-based Bellman-Ford 算法.
func (g *DiGraph) BSP() *Path {
	// 在队列上只能出现vertex的一份拷贝
	// 每次拷贝的distTo和edgeTo被更新,需要在下一次pass里面被处理.
	return nil
}

//TODO:Dijkstra's algorithm,等路劲算法.
func (g *DiGraph) DSP(s *Vertex) *Path {
	path := newPath(s)
	//TODO:封装成一个函数
	for _, v := range g.Vertices {
		path.distTo[v] = math.MaxFloat64
	}
	path.distTo[s] = 0.0
	pq := newPQ()
	pq.Insert(newEdge(s, 0.0))
	for pq.Len() > 0 {
		//然后每次都把一个距离最小的结点加入到.
		edge := pq.Get()
		//每次更新周围的点,如果可以松就松.
		//TODO:如果已经加入最短路劲的点不用处理.
		relax(path, g, edge.vtx, pq)
	}
	return path
}

// 放松点v周围的crossing edge.
func relax(path *Path, g *DiGraph, v *Vertex, pq *PQ) {
	edges := g.Adj(v)
	for _, edge := range edges {
		w, _ := edge.Weight()
		vtx := edge.vtx
		if path.distTo[vtx] > path.distTo[v]+w {
			path.distTo[vtx] = path.distTo[v] + w
			path.edgeTo[vtx] = newEdge(v, w)
			if pq.Contains(vtx) {
				pq.Change(vtx, w)
			} else {
				pq.Insert(newEdge(vtx, w))
			}
		}
	}
}

// SCC return strong connected components,if g is empty returns nil.
func (g *DiGraph) SCC() (scc *SCC) {
	return g.divide()
}

// SCC is Strong Connected Components.
type SCC struct {
	id map[*Vertex]int
}

// divide returns SCC,Kosaraju algorithm,
// it's easy to code but difficult to understand.
func (g *DiGraph) divide() *SCC {
	//反置矩阵的后序遍历.
	scc := &SCC{make(map[*Vertex]int)}
	if len(g.Vertices) == 0 {
		return nil
	}
	count := 0
	marked := make(map[*Vertex]bool)
	walk := func(v *Vertex) {
		scc.id[v] = count
	}
	vertices := g.Reverse().DFSInOrder(POST | REV)
	for _, v := range vertices {
		if !marked[v] {
			g.dfs(v, nil, marked, walk, nil)
			count++
		}
	}
	return scc
}

//proof:
//证明充分条件:如果v和s连通,假设dfs(s)不能到达v,那么说明v之前已经被marked了,
//但是因为有路径v可达s,s也肯定被marked过,所以不可能调用dfs(s),推出矛盾.
//所以v-s连通,是dfs(s)可达v的充分条件.
//证明必要条件:假设dfs(s)有一个可达结点v,那么说明s可以到达v,并且在reverse graph里面有一条路径v可达s,
//剩下只差,证明v到达s,也就是证明s可达v在reverse graph里面存在.
//证明的关键是:reverse post order保证reverse graph dfs时,对于s可达v,肯定有dfs(v)在dfs(s)之前调用(因为有一个v->s),要么是在s结束之前开始,在s结束之前结束.
//要么是在s开始之前开始(在s之前结束)
//又因为有一条v到s的可达路径,所以第一种情况不可能,只有第二种情况,所以得证.

//直觉上来讲,reverse graph 的 reverse post order 产生的是一串后面指向前面的序列,特征是一段一段连续的,分隔的地方是有一个点不能直接接到前面的结点.
//如果从s能到v,那么v就能指向s,(v之前没有被遍历到过)

// TC is TransitiveClosure.
type TC struct {
	all map[*Vertex]map[*Vertex]bool
}

func (tc *TC) Reachable(v, w *Vertex) bool {
	return tc.all[v][w]
}

// 传递闭包是平方的,把这个降低还是一个没有解决的问题.
func (g *DiGraph) TransitiveClosure() *TC {
	tc := &TC{make(map[*Vertex]map[*Vertex]bool)}
	for _, v := range g.Vertices {
		for _, w := range g.DirectedDFS(v) {
			tc.all[v][w] = true
		}
	}
	return nil
}
