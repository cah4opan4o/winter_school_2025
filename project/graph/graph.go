package graph

type Graph struct {
	// Список смежности
	adj map[int][]int // Ключ - номер пользователя, значение - с кем рёбра
	// При необходимости хранить веса:
	// edges []Edge
}

func NewGraph() *Graph {
	return &Graph{adj: make(map[int][]int)}
}

func (g *Graph) AddEdge(u, v int, undirected bool) {
	g.adj[u] = append(g.adj[u], v)
	if undirected {
		g.adj[v] = append(g.adj[v], u)
	}
}

func (g *Graph) Adj() map[int][]int {
	return g.adj
}

func HasEdge(g *Graph, u, v int) bool {
	for _, neighbor := range g.adj[u] {
		if neighbor == v {
			return true
		}
	}
	return false
}

func ConnectedComponents(g *Graph) (count int, comp map[int]int) {
	// comp[v] = номер компоненты (1..count)
	visited := make(map[int]bool)
	comp = make(map[int]int)
	count = 0

	for v, _ := range g.Adj() {
		if !visited[v] {
			count++
			order := []int{}
			order = DFS(g, v)
			for _, u := range order {
				visited[u] = true
				comp[u] = count
			}
		}
	}
	return count, comp
}
