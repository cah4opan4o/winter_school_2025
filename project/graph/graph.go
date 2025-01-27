package graph

type Graph struct {
    // Список смежности
    adj map[int][]int
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

func HasEdge(g *Graph, u, v int) bool{
// вернуть true, если v есть в g.adj[u].
	for _, n := range g.adj[u]{
		if n == v{
			return true
		}
	}
	return false
}