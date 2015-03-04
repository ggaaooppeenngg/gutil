package container

import (
	"bufio"
	"io"
	"os"
	"strconv"
	"strings"
	"testing"
)

func TestTopSort(t *testing.T) {

}

func wrapErr(t *testing.T) func(error) {
	return func(e error) {
		if e != nil {
			t.Fatal(e)
		}
	}
}

func initTinyGraph(graph *EdgeWeightedDigraph, t *testing.T) {
	catchErr := wrapErr(t)
	f, e := os.Open("./test/graph/tinyEWDAG.txt")
	catchErr(e)
	defer f.Close()
	buf := bufio.NewReader(f)
	count := 0
	for {
		line, _, e := buf.ReadLine()
		if e != io.EOF {
			catchErr(e)
		} else {
			break
		}
		if count > 1 {
			args := strings.Split(string(line), " ")
			w, e := strconv.ParseFloat(args[2], 64)
			catchErr(e)
			graph.AddEdge(&Vertex{Id: args[0]}, &Vertex{Id: args[1]}, w)
		}
		count++
	}

}

func TestASP(t *testing.T) {
	g := newEdgeWeightedDigraph()
	initTinyGraph(g, t)
	path := g.ASP(&Vertex{Id: "5"})
	if path.PathString() != `To 0 = 5->4:0.35 4->0:0.38
To 1 = 5->1:0.32
To 2 = 5->7:0.28 7->2:0.34
To 3 = 5->1:0.32 1->3:0.29
To 4 = 5->4:0.35
To 5 = 
To 6 = 5->1:0.32 1->3:0.29 3->6:0.52
To 7 = 5->7:0.28` {
		t.Fail()
	}

}

func TestALP(t *testing.T) {
	g := newEdgeWeightedDigraph()
	initTinyGraph(g, t)
	path := g.ALP(&Vertex{Id: "5"})
	if path.PathString() != `To 0 = 5->1:0.32 1->3:0.29 3->6:0.52 6->4:0.93 4->0:0.38
To 1 = 5->1:0.32
To 2 = 5->1:0.32 1->3:0.29 3->6:0.52 6->4:0.93 4->7:0.37 7->2:0.34
To 3 = 5->1:0.32 1->3:0.29
To 4 = 5->1:0.32 1->3:0.29 3->6:0.52 6->4:0.93
To 5 = 
To 6 = 5->1:0.32 1->3:0.29 3->6:0.52
To 7 = 5->1:0.32 1->3:0.29 3->6:0.52 6->4:0.93 4->7:0.37` {
		t.Log(path.PathString())
		t.Fail()
	}
}
