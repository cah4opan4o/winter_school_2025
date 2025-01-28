package algorithms

import "container/heap"

type Item struct {
	vertex int
	dist   int
	index  int // индекс элемента в куче
}

type PriorityQueue struct {
	items []*Item
}

func (pq *PriorityQueue) Len() int {
	return len(pq.items)
}

func (pq *PriorityQueue) Less(i, j int) bool {
	return pq.items[i].dist < pq.items[j].dist
}

func (pq *PriorityQueue) Swap(i, j int) {
	pq.items[i], pq.items[j] = pq.items[j], pq.items[i]
	pq.items[i].index = i
	pq.items[j].index = j
}

func (pq *PriorityQueue) Push(x interface{}) {
	item := x.(*Item)
	pq.items = append(pq.items, item)
	item.index = len(pq.items) - 1
}

func (pq *PriorityQueue) Pop() interface{} {
	old := pq.items
	n := len(old)
	item := old[n-1]
	pq.items = old[0 : n-1]
	return item
}

func (pq *PriorityQueue) update(item *Item, dist int) {
	item.dist = dist
	heap.Fix(pq, item.index)
}

// Реализация алгоритма Дейкстры
func Dijkstra(g *Graph, start int) ([]int, []int) {
	// Инициализация массивов
	dist := make([]int, len(g.adj))
	parent := make([]int, len(g.adj))

	for i := range dist {
		dist[i] = int(^uint(0) >> 1) // ∞
		parent[i] = -1
	}

	dist[start] = 0

	// Приоритетная очередь (min heap)
	pq := &PriorityQueue{}
	heap.Push(pq, &Item{vertex: start, dist: 0})

	for pq.Len() > 0 {
		item := heap.Pop(pq).(*Item)
		u := item.vertex

		// Обрабатываем все соседей u
		for v, weight := range g.adj[u] {
			// Если найден более короткий путь
			if dist[u]+weight < dist[v] {
				dist[v] = dist[u] + weight
				parent[v] = u
				heap.Push(pq, &Item{vertex: v, dist: dist[v]})
			}
		}
	}

	return dist, parent
}
