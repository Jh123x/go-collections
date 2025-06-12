package queue

import "container/list"

var (
	_ Queue[string] = (*StdQueue[string])(nil)
)

type StdQueue[T any] struct {
	list    *list.List
	maxSize int
}

func NewStdQueue[T any](maxSize int) *StdQueue[T] {
	return &StdQueue[T]{list: list.New(), maxSize: maxSize}
}

func (s *StdQueue[T]) Len() int {
	return s.list.Len()
}

func (s *StdQueue[T]) Enqueue(val T) bool {
	if s.Len() >= s.maxSize {
		return false
	}

	s.list.PushBack(val)
	return true
}

func (s *StdQueue[T]) Dequeue() (T, bool) {
	v := s.list.Front()
	if v == nil {
		var empty T
		return empty, false
	}

	return s.list.Remove(v).(T), true
}
