package stack

import "fmt"

var Debug bool

type Stack[T any] struct {
	items []T
}

func New[T any]() *Stack[T] {
	return &Stack[T]{items: make([]T, 0)}
}

func (s *Stack[T]) Push(x T) {
	if Debug {
		fmt.Printf("Stack: pushing %v\n", x)
	}
	s.items = append(s.items, x)
	if Debug {
		fmt.Printf("Stack after push: %v\n", s.items)
	}
}

func (s *Stack[T]) Pop() T {
	if len(s.items) == 0 {
		if Debug {
			fmt.Printf("Stack: attempted to pop from empty stack\n")
		}
		var zero T
		return zero
	}

	item := s.items[len(s.items)-1]
	s.items = s.items[:len(s.items)-1]
	if Debug {
		fmt.Printf("Stack: popped %v\n", item)
		fmt.Printf("Stack after pop: %v\n", s.items)
	}
	return item
}

func (s *Stack[T]) Top() T {
	if len(s.items) == 0 {
		if Debug {
			fmt.Printf("Stack: attempted to get top from empty stack\n")
		}
		var zero T
		return zero
	}

	item := s.items[len(s.items)-1]
	if Debug {
		fmt.Printf("Stack: top is %v\n", item)
	}
	return item
}

func (s *Stack[T]) IsEmpty() bool {
	empty := len(s.items) == 0
	if Debug {
		fmt.Printf("Stack: IsEmpty() = %v\n", empty)
	}
	return empty
}

func (s *Stack[T]) Len() int {
	length := len(s.items)
	if Debug {
		fmt.Printf("Stack: Len() = %d\n", length)
	}
	return length
}

func (s *Stack[T]) String() string {
	return fmt.Sprintf("%v", s.items)
}
