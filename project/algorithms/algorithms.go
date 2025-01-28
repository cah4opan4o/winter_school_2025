package algorithms

import "math"

func BellmanFord(g *Graph, start int) ([]int, []int, bool) {
	// Инициализация расстояний и родителей
	dist := make([]int, len(g.adj))
	parent := make([]int, len(g.adj))

	for i := range dist {
		dist[i] = math.MaxInt
		parent[i] = -1
	}
	dist[start] = 0

	// Основной цикл алгоритма (V-1 проходов)
	for i := 1; i < len(g.adj); i++ {
		for u := range g.adj {
			for v, w := range g.adj[u] {
				if dist[u] != math.MaxInt && dist[u]+w < dist[v] {
					dist[v] = dist[u] + w
					parent[v] = u
				}
			}
		}
	}

	// Проверка на наличие отрицательного цикла
	negativeCycle := false
	for u := range g.adj {
		for v, w := range g.adj[u] {
			if dist[u] != math.MaxInt && dist[u]+w < dist[v] {
				negativeCycle = true
				break
			}
		}
		if negativeCycle {
			break
		}
	}

	return dist, parent, negativeCycle
}
