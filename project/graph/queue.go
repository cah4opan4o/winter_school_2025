package graph

// type Queue struct {
//     data []int
//     head int
//     tail int
// }

type SimpleQueue []int

func (q *SimpleQueue) Enqueue(x int) {
    *q = append(*q, x)
}
func (q *SimpleQueue) Dequeue() (int, bool) {
    if len(*q) == 0 {
        return 0, false
    }
    val := (*q)[0]
    *q = (*q)[1:]
    return val, true
}