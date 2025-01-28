package graph

func Merge(a []int, b []int) []int {
    final := []int{}
    i := 0
    j := 0
    for i < len(a) && j < len(b) {
        if a[i] < b[j] {
            final = append(final, a[i])
            i++
        } else {
            final = append(final, b[j])
            j++
        }
    }
    for ; i < len(a); i++ {
        final = append(final, a[i])
    }
    for ; j < len(b); j++ {
        final = append(final, b[j])
    }
    return final
}

func MergeSort(items []int) []int{
	first := MergeSort(items[:len(items)/2])
	second := MergeSort(items[len(items)/2:])
	return Merge(first,second)
}

func MST(n int, edges []Edge) (mst []Edge, totalWeight int){
	MergeSort(edges)
	ds := NewDisjoinSet(n)
	return 
}