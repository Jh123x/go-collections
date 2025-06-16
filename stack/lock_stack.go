package stack

import "sync"

var (
	_ Stack[int] = (*LockStack[int])(nil)
)

type LockStack[T any] struct {
	buffer  []T
	mux     *sync.Mutex
	maxSize int
}

func NewLockStack[T any](len int) *LockStack[T] {
	return &LockStack[T]{
		buffer:  make([]T, 0, len),
		mux:     &sync.Mutex{},
		maxSize: len,
	}
}

func (s *LockStack[T]) Len() int {
	s.mux.Lock()
	defer s.mux.Unlock()
	return len(s.buffer)
}

func (s *LockStack[T]) Push(val T) bool {
	if s.Len() == s.maxSize {
		return false
	}

	s.mux.Lock()
	s.buffer = append(s.buffer, val)
	s.mux.Unlock()

	return true
}

func (s *LockStack[T]) Pop() (T, bool) {
	len := s.Len()
	if len == 0 {
		var empty T
		return empty, false
	}

	s.mux.Lock()
	last := len - 1
	val := s.buffer[last]
	s.buffer = s.buffer[:last]
	s.mux.Unlock()

	return val, true
}
