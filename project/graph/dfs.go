package graph

func DFS(g *Graph, start int) []int {
	visited := make(map[int]bool)
	s := Stack{}
	s.Push(start)
	order := []int{}
	for {
		u, ok := s.Pop()
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
				s.Push(neighbor)
			}
		}
	}
	return order
}
