package deijkstra

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
)

type Generator struct{
	Nodes []*GenNode
	Edges []*GenEdge
}

type GenNode struct {
	Id     int
	Name   string
	Weight int
}

type GenEdge struct {
	From     int
	To       int
	Weight   int
}

func NewGenerator() *Generator {
	return &Generator{}
}

func(g *Generator) Generate() string {
	g.generateNodes()
	g.generateEdges()
	fileName, _ := g.writeToFile()

	fmt.Println("File with graph created")

	return fileName
}

func (g *Generator) generateNodes() {
	for i:=0; i < 10000; i++ {
		v := &GenNode{
			Id:     i,
			Name:   fmt.Sprintf("node_name%d", i),
			Weight: 0,
		}

		g.Nodes = append(g.Nodes, v)
	}
}

func (g *Generator) generateEdges() {
	for i:=0; i < 100000; i++ {
		var v = &GenEdge{
			From:   rand.Intn(9999),
			To:     rand.Intn(9999),
			Weight: rand.Intn(9999),
		}
		g.Edges = append(g.Edges, v)
	}
}

func (g *Generator) writeToFile() (string, error) {
	file, err := ioutil.TempFile("/home/ost/goProjects/test/", "graph_")

	defer file.Close()

	if err != nil {
		return "", err
	}
	marshaled, err := json.Marshal(g)
	if err != nil {
		return "", err
	}
	_, err = file.Write(marshaled)
	if err != nil {
		return "", err
	}
	fileName := file.Name()
	return fileName, nil
}