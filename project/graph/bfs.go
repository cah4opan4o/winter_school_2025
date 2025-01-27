package graph

func BFS(g *Graph, start int) []int {
	visited := make(map[int]bool)
	queue := SimpleQueue{}
	queue.Enqueue(start)
	// список отмеченных вершин
	order := []int{}

	for {
		u, ok := queue.Dequeue()
		if !ok {
			break
		}
		if visited[u] {
			continue
		}
		visited[u] = true
		order = append(order, u)

		for _, neighbor := range g.adj[u] {
			if !visited[neighbor] {
				queue.Enqueue(neighbor)
			}
		}
	}

	return order
}
