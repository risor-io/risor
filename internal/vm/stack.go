package vm

type Stack[T any] struct {
	items []T
}

func (s *Stack[T]) Push(item T) {
	s.items = append(s.items, item)
}

func (s *Stack[T]) Pop() (T, bool) {
	size := len(s.items)
	if size == 0 {
		var empty T
		return empty, false
	}
	top := s.items[size-1]
	s.items = s.items[:size-1]
	return top, true
}

func (s *Stack[T]) Top() (T, bool) {
	size := len(s.items)
	if size == 0 {
		var empty T
		return empty, false
	}
	return s.items[size-1], true
}

func (s *Stack[T]) Get(back int) (T, bool) {
	size := len(s.items)
	if size == 0 {
		var empty T
		return empty, false
	}
	if back >= size {
		var empty T
		return empty, false
	}
	return s.items[size-1-back], true
}

func (s *Stack[T]) Size() int {
	return len(s.items)
}

func NewStack[T any](size int) *Stack[T] {
	return &Stack[T]{
		items: make([]T, 0, size),
	}
}
