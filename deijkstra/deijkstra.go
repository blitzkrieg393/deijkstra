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
	"time"
)

type Node struct {
	Id     int    `json:"id"`
	Name   string `json:"name"`
	Weight int
	Edges  []*Edge
}

type Edge struct {
	From     int `json:"from"`
	To       int `json:"to"`
	Weight   int `json:"weight"`
}

type Graph struct {
	Nodes []*Node `json:"nodes"`
	Edges []*Edge `json:"edges"`
}

type Ways struct {
	graph *Graph
	way   map[string][]string
}

var way []byte

func New(fileName string) *Ways {
	w := &Ways{
		way:   make(map[string][]string),
		graph: createGraph(fileName),
	}

	w.fillNodes()

	return w
}
func createGraph(fileName string) *Graph {
	var graph *Graph
	content, _ := ioutil.ReadFile(fileName)
	_ = json.Unmarshal(content, &graph)

	fmt.Println("Graph created")

	return graph

}
func (w *Ways) clearNodesWeight() {
	for _, node := range w.graph.Nodes {
		node.Weight = math.MaxInt64
	}
	fmt.Println("Nodes weight cleared")
}
func (w *Ways) fillNodes() {
	for _, edge := range w.graph.Edges {
		w.graph.Nodes[edge.From].Edges = append(w.graph.Nodes[edge.From].Edges, edge)
	}
	for _, node := range w.graph.Nodes {
		sort.Slice(node.Edges, func(i, j int) bool {
			return node.Edges[i].Weight < node.Edges[j].Weight
		})
	}

	fmt.Println("Nodes filled")
}

func (w *Ways) Full(ctx *fasthttp.RequestCtx) {
	now := time.Now()
	if len(w.way) != 0 {
		w.way = make(map[string][]string)
	}
	routes := ctx.QueryArgs().PeekMulti("route")

	if len(routes) == 0 {
		ctx.SetStatusCode(400)
		fmt.Fprint(ctx, "route must not be empty\n")

		return
	}

	for _, route := range routes {
		res := strings.Split(string(route), ",")
		from, err := strconv.Atoi(res[0])
		if err != nil {
			ctx.SetStatusCode(400)
			fmt.Fprint(ctx, "error on parsing vertex from\n")

			return
		}

		to, err := strconv.Atoi(res[1])
		if err != nil {
			ctx.SetStatusCode(400)
			fmt.Fprint(ctx, "error on parsing vertex to\n")

			return
		}

		w.fullVertex(w.graph.Nodes[from], to, way, string(route))
		fmt.Println("finished route ", string(route))
	}

	fmt.Println(time.Since(now))
	encoder := json.NewEncoder(ctx)
	_ = encoder.Encode(w.way)

	return
}
func (w *Ways) Short(ctx *fasthttp.RequestCtx) {
	now := time.Now()
	if len(w.way) != 0 {
		w.way = make(map[string][]string)
	}

	routes := ctx.QueryArgs().PeekMulti("route")
	if len(routes) == 0 {
		ctx.SetStatusCode(400)
		fmt.Fprint(ctx, "route must not be empty\n")

		return
	}

	for _, route := range routes {

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

		w.clearNodesWeight()
		w.graph.Nodes[from].Weight = 0
		w.shortVertex(w.graph.Nodes[from], to, way, string(route))
		fmt.Println("finished route ", string(route))
	}

	fmt.Println(time.Since(now))

	encoder := json.NewEncoder(ctx)
	_ = encoder.Encode(w.way)

	return
}

func (w *Ways) fullVertex(node *Node, end int, way []byte, key string) {
	if node.Id == end {
		nodeId := strconv.FormatInt(int64(node.Id), 10)
		way = append(way, []byte(nodeId)...)
		w.way[key] = append(w.way[key], string(way))
		way = []byte{}

		return
	}

	if len(node.Edges) == 0 {
		return
	}

	for _, nodeEdge := range node.Edges {
		if nodeEdge.To == 0 {
			return
		}

		nodeTo := strconv.FormatInt(int64(nodeEdge.To), 10)
		if bytes.Contains(way, []byte(nodeTo+",")) {
			continue
		}

		nodeId := strconv.FormatInt(int64(node.Id), 10)
		if !bytes.Contains(way, []byte(nodeId+",")) {
			way = append(way, []byte(nodeId+",")...)
		}

		w.fullVertex(w.graph.Nodes[nodeEdge.To], end, way, key)
	}
}
func (w *Ways) shortVertex(node *Node, end int, way []byte, key string) {

	if node.Id == end {
		way = append(way, []byte(fmt.Sprintf("%d", node.Id))...)
		w.way[key] = append(w.way[key], string(way))
		way = []byte{}

		return
	}

	if len(node.Edges) == 0 {
		return
	}

	for _, nodeEdge := range node.Edges {
		if nodeEdge.To == 0 {
			return
		}

		if bytes.Contains(way, []byte(fmt.Sprintf(",%d,", nodeEdge.To))) {
			continue
		}

		newWeight := nodeEdge.Weight + node.Id
		if w.graph.Nodes[nodeEdge.To].Weight > newWeight {
			w.graph.Nodes[nodeEdge.To].Weight = newWeight
			if !bytes.Contains(way, []byte(fmt.Sprintf(",%d,", node.Id))) {
				way = append(way, []byte(fmt.Sprintf("%d,", node.Id))...)
			}

			w.shortVertex(w.graph.Nodes[nodeEdge.To], end, way, key)
		}
	}
}
