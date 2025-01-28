package graph

func MergeSort(edges []Edge) []Edge {
	if len(edges) <= 1 {
		return edges
	}

	mid := len(edges) / 2
	left := MergeSort(edges[:mid])
	right := MergeSort(edges[mid:])

	return merge(left, right)
}

func merge(left, right []Edge) []Edge {
	result := []Edge{}
	i, j := 0, 0

	for i < len(left) && j < len(right) {
		if left[i].w < right[j].w {
			result = append(result, left[i])
			i++
		} else {
			result = append(result, right[j])
			j++
		}
	}

	// Добавляем оставшиеся элементы
	result = append(result, left[i:]...)
	result = append(result, right[j:]...)

	return result
}

func MST(n int, edges []Edge) (mst []Edge, totalWeight int) {
	MergeSort(edges)
	ds := NewDisjoinSet(n)
	for _, edge := range edges {
		u, v, w := edge.u, edge.v, edge.w
		if ds.Find(u) != ds.Find(v) {
			ds.Union(u, v)
			mst = append(mst, edge)
			totalWeight += w

			if len(mst) == n-1 {
				break
			}
		}
	}
	return mst, totalWeight
}
