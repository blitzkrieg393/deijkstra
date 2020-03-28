package main

import (
	"battle/deijkstra"
	"fmt"
	"github.com/buaazp/fasthttprouter"
	"github.com/valyala/fasthttp"
)

func main() {
	paths := deijkstra.NewWays("graph.json")
	httpRouter := fasthttprouter.New()
	httpRouter.GET("/paths/", paths.Full)
	httpRouter.GET("/short-paths/", paths.Short)

	fmt.Println("Http server started")
	err := fasthttp.ListenAndServe("localhost:8080", httpRouter.Handler)
	if err != nil {
		fmt.Printf("Can't start http server: %s \n", err.Error())
	}
}

