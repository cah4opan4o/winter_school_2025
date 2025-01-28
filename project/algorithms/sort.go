package algorithms

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
