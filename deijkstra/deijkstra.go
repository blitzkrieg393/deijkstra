package deijkstra

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/valyala/fasthttp"
	"io/ioutil"
	"math"
	"sort"
	"strconv"
	"strings"
	"sync"
)

type Node struct {
	Id     int    `json:"id"`
	Name   string `json:"name"`
	Weight int
}

type Edge struct {
	From     int `json:"from"`
	To       int `json:"to"`
	Weight   int `json:"weight"`
	ToNode   *Node
	FromNode *Node
}

type Graph struct {
	Nodes []*Node `json:"nodes"`
	Edges []*Edge `json:"edges"`
}

type Ways struct{
	graph *Graph
	way []string
}

var way []byte

func NewWays(fileName string) *Ways {
	w := &Ways{}
	return &Ways{
		graph: w.createGraph(fileName),
	}
}

func (w *Ways) createGraph(fileName string) *Graph {
	var graph *Graph
	content, _ := ioutil.ReadFile("graph.json")
	_ = json.Unmarshal(content, &graph)

	return graph

}

func (w *Ways) Full(ctx *fasthttp.RequestCtx){
	w.way = []string{}
	route := ctx.QueryArgs().Peek("route")
	if route == nil {
		ctx.SetStatusCode(400)
		fmt.Fprint(ctx, "route must not be empty\n")

		return
	}

	res := strings.Split(string(route), ",")
	from, err := strconv.Atoi(res[0])
	if err != nil {
		ctx.SetStatusCode(400)
		fmt.Fprint(ctx, "error on parsing string\n")

		return
	}

	to, err := strconv.Atoi(res[1])
	if err != nil {
		ctx.SetStatusCode(400)
		fmt.Fprint(ctx, "error on parsing string\n")

		return
	}

	for _, node := range w.graph.Nodes {
		if node.Id == from {
			node.Weight = 0
			w.fullVertex(node, to, way)
		}
	}

	encoder := json.NewEncoder(ctx)
	_ = encoder.Encode(w.way)


	return


}
func (w *Ways) Short(ctx *fasthttp.RequestCtx){
	w.way = []string{}
	route := ctx.QueryArgs().Peek("route")
	if route == nil {
		ctx.SetStatusCode(400)
		fmt.Fprint(ctx, "route must not be empty\n")

		return
	}

	res := strings.Split(string(route), ",")
	from, err := strconv.Atoi(res[0])
	if err != nil {
		ctx.SetStatusCode(400)
		fmt.Fprint(ctx, "error on parsing string\n")

		return
	}

	to, err := strconv.Atoi(res[1])
	if err != nil {
		ctx.SetStatusCode(400)
		fmt.Fprint(ctx, "error on parsing string\n")

		return
	}

	for _, node := range w.graph.Nodes {
		if node.Id == from {
			node.Weight = 0
			w.shortVertex(node, to, way)
		}
	}

	encoder := json.NewEncoder(ctx)
	_ = encoder.Encode(w.way)

	return
}

func (w *Ways) shortVertex(node *Node, end int, way []byte) {
	edgesLen := len(w.graph.Edges)
	j := 100
	ch := make(chan *Edge)
	wg := sync.WaitGroup{}

	for i := 0; i < edgesLen; {
		if j > edgesLen {
			j = edgesLen
		}

		wg.Add(1)
		go func(i int, j int) {
			for _, edge := range w.graph.Edges[i:j] {
				if edge.From == node.Id && edge.To != 0 {
					ch <- edge
				}
			}

			wg.Done()
		}(i, j)

		i, j = j, j+100
	}

	var curPosEdges []*Edge
	go func() {
		for edge := range ch {
			for _, graphNode := range w.graph.Nodes {
				if graphNode.Id == edge.To {
					// todo зарефакторить
					if graphNode.Weight > 0 && graphNode.Weight < math.MaxInt64 {
						edge.ToNode = graphNode
					} else {
						graphNode.Weight = math.MaxInt64
						edge.ToNode = graphNode
					}
				}
			}
			edge.FromNode = node
			curPosEdges = append(curPosEdges, edge)
		}
		close(ch)
	}()
	wg.Wait()

	var suitableEdges []*Edge
	if len(curPosEdges) == 0 {
		return
	}

	for _, edge := range curPosEdges {
		if edge.ToNode.Weight > (edge.Weight + edge.FromNode.Weight) {
			edge.ToNode.Weight = edge.Weight + edge.FromNode.Weight
			suitableEdges = append(suitableEdges, edge)
		}
	}

	sort.Slice(suitableEdges, func(i, j int) bool {
		return suitableEdges[i].Weight < suitableEdges[j].Weight
	})

	for _, suitableEdge := range suitableEdges {

		if !bytes.Contains(way, []byte(fmt.Sprintf("%d", node.Id))) {
			way = append(way, []byte(fmt.Sprintf("%d", node.Id))...)
		}

		if suitableEdge.To == end {
			if !bytes.Contains(way, []byte(fmt.Sprintf("%d", suitableEdge.To))) {
				way = append(way, []byte(fmt.Sprintf("%d", suitableEdge.To))...)
			}

			w.way = append(w.way, string(way))
			way = []byte{}

			w.shortVertex(suitableEdge.ToNode, end, way)
		}

		w.shortVertex(suitableEdge.ToNode, end, way)
	}
}

func (w *Ways) fullVertex(node *Node, end int, way []byte) {
	if node.Id == 7 {
		fmt.Println(node)
	}
	edgesLen := len(w.graph.Edges)
	j := 100
	ch := make(chan *Edge)
	wg := sync.WaitGroup{}
	for i := 0; i < edgesLen; {
		if j > edgesLen {
			j = edgesLen
		}

		wg.Add(1)
		go func(i int, j int) {
			for _, edge := range w.graph.Edges[i:j] {
				if edge.From == node.Id && edge.To != 0 {
					ch <- edge
				}
			}

			wg.Done()
		}(i, j)

		i, j = j, j+100
	}

	var curPosEdges []*Edge
	go func() {
		for edge := range ch {
			for _, graphNode := range w.graph.Nodes {
				if graphNode.Id == edge.To {
					// todo зарефакторить
					if graphNode.Weight > 0 && graphNode.Weight < math.MaxInt64 {
						edge.ToNode = graphNode
					} else {
						graphNode.Weight = math.MaxInt64
						edge.ToNode = graphNode
					}
				}
			}
			edge.FromNode = node
			curPosEdges = append(curPosEdges, edge)
		}
		close(ch)
	}()
	wg.Wait()

	var suitableEdges []*Edge
	if len(curPosEdges) == 0 {
		return
	}

	for _, edge := range curPosEdges {
		edge.ToNode.Weight = edge.Weight + edge.FromNode.Weight
		suitableEdges = append(suitableEdges, edge)
	}

	sort.Slice(suitableEdges, func(i, j int) bool {
		return suitableEdges[i].Weight < suitableEdges[j].Weight
	})

	for _, suitableEdge := range suitableEdges {

		if !bytes.Contains(way, []byte(fmt.Sprintf("%d", node.Id))) {
			way = append(way, []byte(fmt.Sprintf("%d", node.Id))...)
		}

		if suitableEdge.To == end {
			if !bytes.Contains(way, []byte(fmt.Sprintf("%d", suitableEdge.To))) {
				way = append(way, []byte(fmt.Sprintf("%d", suitableEdge.To))...)
			}

			w.way = append(w.way, string(way))
			way = []byte{}

			w.fullVertex(suitableEdge.ToNode, end, way)
		}

		w.fullVertex(suitableEdge.ToNode, end, way)
	}
}
