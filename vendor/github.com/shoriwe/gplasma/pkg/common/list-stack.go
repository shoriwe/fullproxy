package common

type (
	stackNode struct {
		Value any
		Next  *stackNode
	}
	ListStack[T any] struct {
		Top *stackNode
	}
)

func (s *ListStack[T]) Push(value T) {
	s.Top = &stackNode{
		Value: value,
		Next:  s.Top,
	}
}

func (s *ListStack[T]) Peek() T {
	return s.Top.Value.(T)
}

func (s *ListStack[T]) Pop() T {
	value := s.Top.Value.(T)
	s.Top = s.Top.Next
	return value
}

func (s *ListStack[T]) HasNext() bool {
	return s.Top != nil
}
