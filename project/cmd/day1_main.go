package main

import (
	"fmt"
	"winter_school_2025/project/graph"
)

func main2() {
	g := graph.NewGraph()

	g.AddEdge(0, 1, true)
	g.AddEdge(1, 2, true)
	// g.AddEdge(2, 3, true)
	g.AddEdge(3, 0, true)
	g.AddEdge(4, 5, true)
	g.AddEdge(5, 6, true)
	g.AddEdge(6, 4, true)

	fmt.Println("Graph_list")
	fmt.Println(g.Adj())
	// for node, neighbors := range g.Adj() {
	// 	fmt.Printf("%d: %v\n", node, neighbors)
	// }

	fmt.Println("user 1 and 4, friends?", graph.HasEdge(g, 1, 4))
	fmt.Println("user 1 and 4, friends?", graph.HasEdge(g, 1, 2))
	fmt.Println("BFS: ", graph.BFS(g, 0))
	fmt.Println("DFS: ", graph.DFS(g, 0))

	count, comp := graph.ConnectedComponents(g)
	fmt.Printf("count: %d\n", count)
	fmt.Print(comp)

}
