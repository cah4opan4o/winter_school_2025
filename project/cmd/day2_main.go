package main

import (
	"container/heap"
	"fmt"
	"math"
	"sort"
)

// Структура DisjointSet
type DisjointSet struct {
	parent []int
	rank   []int
}

// Создание множества
func NewDisjointSet(n int) *DisjointSet {
	ds := &DisjointSet{
		parent: make([]int, n),
		rank:   make([]int, n),
	}
	for i := 0; i < n; i++ {
		ds.parent[i] = i // Каждый сам себе корень
		ds.rank[i] = 0
	}
	return ds
}

// Поиск сжатия пути
func (ds *DisjointSet) Find(x int) int {
	if ds.parent[x] != x {
		ds.parent[x] = ds.Find(ds.parent[x]) // Сжатие пути
	}
	return ds.parent[x]
}

// Объединение множеств
func (ds *DisjointSet) Union(x, y int) bool {
	rootX := ds.Find(x)
	rootY := ds.Find(y)

	if rootX == rootY {
		return false
	}

	// Ранг - балансировка
	if ds.rank[rootX] > ds.rank[rootY] {
		ds.parent[rootY] = rootX
	} else if ds.rank[rootX] < ds.rank[rootY] {
		ds.parent[rootX] = rootY
	} else {
		ds.parent[rootY] = rootX
		ds.rank[rootX]++
	}

	return true
}

type Edge struct {
	u, v, w int
}

type Graph struct {
	edges []Edge
}

// Добавление ребра
func (g *Graph) AddEdge(u, v, w int) {
	g.edges = append(g.edges, Edge{u, v, w})
}

// Получение всех рёбер
func (g *Graph) GetAllEdges() []Edge {
	return g.edges
}

func MST(n int, edges []Edge) ([]Edge, int) {
	sort.Slice(edges, func(i, j int) bool {
		return edges[i].w < edges[j].w
	})

	ds := NewDisjointSet(n)
	mst := []Edge{}
	totalWeight := 0

	for _, edge := range edges {
		if ds.Find(edge.u) != ds.Find(edge.v) {
			ds.Union(edge.u, edge.v)
			mst = append(mst, edge)
			totalWeight += edge.w
		}
		if len(mst) == n-1 {
			break
		}
	}

	return mst, totalWeight
}

type Item struct {
	vertex, dist int
}

type PriorityQueue []Item

func (pq *PriorityQueue) Len() int { return len(*pq) }
func (pq *PriorityQueue) Less(i, j int) bool {
	return (*pq)[i].dist < (*pq)[j].dist
}
func (pq *PriorityQueue) Swap(i, j int) { (*pq)[i], (*pq)[j] = (*pq)[j], (*pq)[i] }
func (pq *PriorityQueue) Push(x interface{}) {
	*pq = append(*pq, x.(Item))
}
func (pq *PriorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	*pq = old[:n-1]
	return item
}

func Dijkstra(g *Graph, start int, n int) ([]int, []int) {
	dist := make([]int, n)
	parent := make([]int, n)

	for i := range dist {
		dist[i] = math.MaxInt32
		parent[i] = -1
	}

	dist[start] = 0
	pq := &PriorityQueue{}
	heap.Init(pq)
	heap.Push(pq, Item{start, 0})

	for pq.Len() > 0 {
		u := heap.Pop(pq).(Item).vertex

		for _, edge := range g.edges {
			if edge.u == u {
				v, w := edge.v, edge.w
				if dist[u]+w < dist[v] {
					dist[v] = dist[u] + w
					parent[v] = u
					heap.Push(pq, Item{v, dist[v]})
				}
			}
		}
	}

	return dist, parent
}

func BellmanFord(g *Graph, start int, n int) ([]int, []int, bool) {
	dist := make([]int, n)
	parent := make([]int, n)

	for i := range dist {
		dist[i] = math.MaxInt32
		parent[i] = -1
	}

	dist[start] = 0

	for i := 0; i < n-1; i++ {
		for _, edge := range g.edges {
			if dist[edge.u] != math.MaxInt32 && dist[edge.u]+edge.w < dist[edge.v] {
				dist[edge.v] = dist[edge.u] + edge.w
				parent[edge.v] = edge.u
			}
		}
	}

	// Проверка отрицательных циклов
	for _, edge := range g.edges {
		if dist[edge.u] != math.MaxInt32 && dist[edge.u]+edge.w < dist[edge.v] {
			return dist, parent, true
		}
	}

	return dist, parent, false
}

func main() {
	fmt.Println("=== Union-Find ===")
	ds := NewDisjointSet(6)
	ds.Union(0, 1)
	ds.Union(1, 2)
	ds.Union(3, 4)

	fmt.Println("Find(0):", ds.Find(0))
	fmt.Println("Find(1):", ds.Find(1))
	fmt.Println("Find(2):", ds.Find(2))
	fmt.Println("Find(3):", ds.Find(3))
	fmt.Println("Find(4):", ds.Find(4))
	fmt.Println("Find(5):", ds.Find(5))

	fmt.Println("\n=== MST (Kruskal) ===")
	g := &Graph{}
	g.AddEdge(0, 1, 4)
	g.AddEdge(0, 2, 4)
	g.AddEdge(1, 2, 2)
	g.AddEdge(1, 3, 5)
	g.AddEdge(2, 3, 5)
	g.AddEdge(3, 4, 3)

	mst, totalWeight := MST(5, g.GetAllEdges())
	fmt.Println("MST edges:", mst)
	fmt.Println("Total weight:", totalWeight)

	fmt.Println("\n=== Dijkstra ===")
	dist, _ := Dijkstra(g, 0, 5)
	fmt.Println("Dijkstra distances:", dist)

	fmt.Println("\n=== Bellman-Ford ===")
	dist, _, hasNegativeCycle := BellmanFord(g, 0, 5)
	fmt.Println("Bellman-Ford distances:", dist)
	fmt.Println("Has negative cycle:", hasNegativeCycle)
}
