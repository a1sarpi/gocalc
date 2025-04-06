package stack

import "fmt"

type Stack[T any] struct {
	items []T
}

func New[T any]() *Stack[T] {
	return &Stack[T]{items: make([]T, 0)}
}

func (s *Stack[T]) Push(x T) {
	s.items = append(s.items, x)
}

func (s *Stack[T]) Pop() T {
	if len(s.items) == 0 {
		var zero T
		return zero
	}

	item := s.items[len(s.items)-1]
	s.items = s.items[:len(s.items)-1]
	return item
}

func (s *Stack[T]) Top() T {
	if len(s.items) == 0 {
		var zero T
		return zero
	}

	return s.items[len(s.items)-1]
}

func (s *Stack[T]) IsEmpty() bool {
	return len(s.items) == 0
}

func (s *Stack[T]) Len() int {
	return len(s.items)
}

func (s *Stack[T]) String() string {
	return fmt.Sprintf("%v", s.items)
}
