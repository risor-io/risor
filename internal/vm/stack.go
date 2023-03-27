package vm

type Stack[T any] struct {
	items []T
	count int
}

func (s *Stack[T]) Push(item T) {
	// s.items = append(s.items, item)
	s.items[s.count] = item
	s.count++
}

func (s *Stack[T]) Pop() T {
	// size := len(s.items)
	// if size == 0 {
	// 	var empty T
	// 	return empty
	// }
	// top := s.items[size-1]
	// s.items = s.items[:size-1]
	// return top
	item := s.items[s.count-1]
	s.count--
	return item
}

func (s *Stack[T]) Top() (T, bool) {
	// size := len(s.items)
	// if size == 0 {
	// 	var empty T
	// 	return empty, false
	// }
	// return s.items[size-1], true
	if s.count == 0 {
		var empty T
		return empty, false
	}
	return s.items[s.count-1], true
}

func (s *Stack[T]) Get(back int) (T, bool) {
	// size := len(s.items)
	// if size == 0 {
	// 	var empty T
	// 	return empty, false
	// }
	// if back >= size {
	// 	var empty T
	// 	return empty, false
	// }
	// return s.items[size-1-back], true
	if back >= s.count {
		var empty T
		return empty, false
	}
	return s.items[s.count-1-back], true
}

func (s *Stack[T]) Size() int {
	// return len(s.items)
	return s.count
}

func NewStack[T any](size int) *Stack[T] {
	return &Stack[T]{
		items: make([]T, size),
	}
}
