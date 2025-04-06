package stack_test

import (
	"github.com/a1sarpi/gocalc/src/stack"
	"testing"
)

func TestStack(t *testing.T) {
	s := stack.New[int]()
	s.Push(1)
	s.Push(2)

	if got := s.Pop(); got != 2 {
		t.Errorf("Pop() = %v, want 2", got)
	}

	if got := s.Top(); got != 1 {
		t.Errorf("Peek() = %v, want 1", got)
	}

	if got := s.Len(); got != 1 {
		t.Errorf("Size() = %v, want 1", got)
	}
}
