package container

import (
	"container/heap"
	"container/list"
	"errors"
	"fmt"
	"math"
	"sync"
)

// UndiGraph is undirected graph,implemented in an slice of adjacency lists.

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
	return &UndiGraph{
		Edges:      make(map[*Vertex][]*edge),
		verticesId: make(map[string]bool),
	}
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

func (g *UndiGraph) String() string {
	var output string
	for _, v := range g.Vertices {
		edges := g.Adj(v)
		output += v.Id + " ->"
		for _, edge := range edges {
			output += " " + edge.vtx.Id
		}
		output += "\n"
	}
	return output
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

//AddEdge adds edge v-w,parallel edges is allowed,and self-loop is allowed.
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

	//w add edge w-v
	e = &edge{v, weight}
	if _, ok := g.Edges[w]; !ok {
		g.Edges[w] = []*edge{e}
	} else {
		edges = g.Edges[w]
		edges = append(edges, e)
		g.Edges[w] = edges
	}
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

//HasCycle detects cycle,returns false if g is acylic,else true.
//assuming no self-loops or parral edges.
func (g *UndiGraph) HasCycle() bool {
	var (
		marked   = make(map[*Vertex]bool)
		prev     *Vertex
		hasCycle bool
	)
	if len(g.Vertices) == 0 {
		return true
	}
	s := g.Vertices[0]
	dfs(g, s, nil, marked, func(v *Vertex) {
		edges := g.Adj(v)
		for _, edge := range edges {
			//如果指向了一个marked果的结点,并且这个结点不是v的前驱,说明v连接到了之前的点上.
			//也就是说成环
			if marked[edge.vtx] {
				if edge.vtx != prev {
					hasCycle = true
				}
			}
		}
		prev = v
	})
	return hasCycle
}

//IsBipartite returns true if g is bipartite else returns false.
func (g *UndiGraph) IsBipartite() bool {
	var (
		marked = make(map[*Vertex]bool)
		color  = make(map[*Vertex]bool)
		ok     bool
	)
	if len(g.Vertices) == 0 {
		return true
	}
	s := g.Vertices[0]
	dfs(g, s, nil, marked, func(v *Vertex) {
		for _, edge := range g.Adj(v) {
			//遍历的时候间隔标志颜色
			color[edge.vtx] = !color[v]
			//如果之前的点中颜色一样说明不能分成两个等价类.
			if marked[edge.vtx] {
				if color[edge.vtx] == color[v] {
					ok = false
				}
			}
		}
	})
	return ok
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
//1 This Path is the shortest path if the weights are all the same.
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

//proof:
//因为队列当中的分布是这样的 [k|k+1],队列里面前段是离s为k的结点,后段是离s为k+1的结点,k大于0(k指的是边数).
//那么对于一个结点v来说,当它进入这个队列的时候,之前都没有扫描到它,也就是说,它没有出现在路径小于等于k的地方.
//那么v出现在了k+1的地方,在这之后v不会出现在k+2,k+3这样的距离上了.
//通俗的讲,当v入队的时候,之前没有遍历到v,也就是说之前的k层距离遍历的点内都没有v,那么现在有v了,说明这个距离就是最短距离,
//另外在之后更长的距离上也不会看到v了.

//dfs implements recursive DFS search method.
//Time complex is O(V+E).
//一般是和这个图的度成比例.
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

//CC returns Connected Components,if g is empty,returns nil.
func (g *UndiGraph) CC() (cc *CC) {
	return g.divide()
}

//CC is Connected Componentes used to divide vertices into connected components(equivalent classes).
type CC struct {
	id map[*Vertex]int //equivalent classes id.
}

//Connected returns true if v and w are connected else false.
func (cc *CC) Connected(v, w *Vertex) bool {
	return cc.id[v] == cc.id[w]
}

//Divide divides vertices into connected components.
func (g *UndiGraph) divide() *CC {
	cc := &CC{make(map[*Vertex]int)}
	if len(g.Vertices) == 0 {
		return nil
	}
	count := 0
	marked := make(map[*Vertex]bool)
	walk := func(v *Vertex) {
		cc.id[v] = count
	}
	for _, v := range g.Vertices {
		if !marked[v] {
			dfs(g, v, nil, marked, walk)
			count++
		}
	}
	return cc
}

//minimum spanning tree.

//proof:图的a cut 是把定点分成两个非空结点集合,而一个crossing edge 的 a cut就是连接一个集合到另外一个集合的一条边.
//对于任何的cut,最小的那个crossing edge就正在MST里面,可以用反证法,如果有一个e是最小的,不在MST里面,那么假设f连接两个集合并且满足最小生成树.
//这个时候把e加到里面,就会有一个环(树连接任意两点就会成环),从一个集合到f再到另一个并且通过e返回,这个时候把f去掉,那么权重比原来小,所以除了e不可能有别的结点.
//连接两个集合的crossing edge不一定只有一个,比如g和h是一个集合里面的,剩下的是另外一个的,连接g和h就有两个.
//a-+
//  |
//  v
//  b-+
//    |
//    v
//    c---+
//     -+ |
//      | |
//      v |
//      d---+
//        | |
//        v |
//        g |
//          |
//          v
//          e---+
//           -+ |
//            | |
//            v |
//            f |
//              |
//              v
//              h

//解决方法: Prim's algorithm,贪婪算法 Greedy MST algorithm.
//把定点分成两部分,每次都找到最小的crossing edge然后把对应结点加到另外一个集合里面.
//PrimMST returns a new graph containing minimun spanning tree.
func PMST(g *UndiGraph) *UndiGraph {
	//记录这颗数
	edgeTo := make(map[*Vertex]*edge)
	//记录距离
	distTo := make(map[*Vertex]float64) //distTo[w]=edgeTo[w].weight
	//marked,用于dfs
	marked := make(map[*Vertex]bool)
	//优先队列
	pq := make(PQ, g.V)
	heap.Init(&pq)
	if len(g.Vertices) == 0 {
		return nil
	}
	s := g.Vertices[0]
	for _, v := range g.Vertices {
		distTo[v] = math.MaxFloat64
	}
	distTo[s] = 0.0
	pq.Push(&edge{s, new(float64)})
	for pq.Len() > 0 {
		edgeV := pq.Pop().(*edge)
		marked[edgeV.vtx] = true
		v := edgeV.vtx
		edges := g.Adj(v)
		for _, e := range edges {
			if marked[e.vtx] {
				continue
			}
			var w float64
			if e.weight == nil {
				w = 0
			} else {
				w = *(e.weight)
			}
			if w < distTo[e.vtx] {
				edgeTo[e.vtx] = &edge{v, e.weight}
				distTo[e.vtx] = w
				if index := pq.contains(e.vtx); index != -1 {
					item := pq[index]
					pq.update(item, e)
				} else {
					pq.Push(e)
				}
			}
		}
	}
	retG := NewUndiGraph()
	for prev, e := range edgeTo {
		if e.weight != nil {
			retG.AddEdge(prev, e.vtx, *e.weight)
		} else {
			retG.AddEdge(prev, e.vtx)
		}
	}
	return retG
}

//Kruskal’s algorithm
func kMST(g *UndiGraph) *UndiGraph {
	//按照顺序组装最短的不成环的边直到组装了V-1次.
	//TODO:因为要用到边的两头,存储的边没有这个值,很难找.
	//改变存储模式,或者先转换再运行算法.
	return nil
}
