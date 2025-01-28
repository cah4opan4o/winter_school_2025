package graph

type DisjoinSet struct{
	parent []int
	rank []int
}

func NewDisjoinSet(n int) *DisjoinSet {
	d := &DisjoinSet{parent : make([]int,n), rank : make([]int,n)}
	for i:= 0; i < n; i++{
		d.parent[i] = i
		d.rank[i] = 0
	}
	return d 
}

func (d *DisjoinSet) Find(x int) int{
	if d.parent[x] != x{
		d.parent[x] = d.Find(d.parent[x])
	}
	return d.parent[x]
}

func (d* DisjoinSet) Union(x,y int) bool{
	rx := d.Find(x)
	ry := d.Find(y)

	if rx == ry{
		return false
	} else{
		if d.rank[rx] == d.rank[ry]{
			// how do that??
			// подвесить дерево меньше ранга к большему. Если ранги равны, rank увеличиваем у нового корня.
		}
	}
	return true
}