package queue

import (
	"sync"
)

var (
	_ Queue[string] = (*LockQueue[string])(nil)
)

type LockQueue[T any] struct {
	buffer  []T
	mux     *sync.Mutex
	start   int
	end     int
	maxSize int
}

func NewLockQueue[T any](len int) *LockQueue[T] {
	return &LockQueue[T]{
		buffer:  make([]T, len+1),
		mux:     &sync.Mutex{},
		start:   0,
		end:     0,
		maxSize: len,
	}
}

func (q *LockQueue[T]) Len() int {
	q.mux.Lock()
	diff := q.end - q.start
	q.mux.Unlock()
	if diff < 0 {
		return diff + q.maxSize + 1
	}

	return diff
}

func (q *LockQueue[T]) Enqueue(val T) bool {
	if q.Len() == q.maxSize {
		return false
	}

	q.mux.Lock()
	q.buffer[q.end] = val
	q.end++
	if q.end >= q.maxSize+1 {
		q.end = 0
	}
	q.mux.Unlock()

	return true
}

func (q *LockQueue[T]) Dequeue() (T, bool) {
	if q.Len() == 0 {
		var empty T
		return empty, false
	}

	q.mux.Lock()
	v := q.buffer[q.start]
	q.start++
	if q.start >= q.maxSize+1 {
		q.start = 0
	}
	q.mux.Unlock()

	return v, true
}
