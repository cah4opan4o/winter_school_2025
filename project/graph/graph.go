package graph

type Edge struct{
	u, v, w int
}

type Graph struct {
	// Список смежности
	adj map[int][]int // Ключ - номер пользователя, значение - с кем рёбра
	edges []Edge
}

func NewGraph() *Graph {
	return &Graph{adj: make(map[int][]int),edges: make([]Edge, 0)}
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

func (g *Graph) Edges() []Edge{
	return g.edges
}

func HasEdge(g *Graph, u, v int) bool {
	for _, neighbor := range g.adj[u] {
		if neighbor == v {
			return true
		}
	}
	return false
}

func (g *Graph) GetAllEdges() []Edge {
    var edges []Edge
	for _, value := range g.Edges(){
		if value.u < value.v{
			edges = append(edges, Edge{value.u, value.v, value.w})
		}
	} 
    return edges
}

func ConnectedComponents(g *Graph) (count int, comp map[int]int) {
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
