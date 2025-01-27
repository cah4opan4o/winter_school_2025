package graph

type Stack struct {
	data []int
}

func NewStack() *Stack {
    return &Stack{data: make([]int, 0)}
}

func (s *Stack) Push(x int) {
    s.data = append(s.data, x)
}
func (s *Stack) Pop() (int, bool) {
    if len(s.data) == 0 {
        return 0, false
    }
    topIndex := len(s.data) - 1
    val := s.data[topIndex]
    s.data = s.data[:topIndex]
    return val, true
}
func (s *Stack) IsEmpty() bool {
    return len(s.data) == 0
}