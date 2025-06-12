package queue

import (
	"sync"

	"github.com/Jh123x/go-collections/optional"
)

var (
	_ Queue[string] = (*LockQueue[string])(nil)
)

type LockQueue[T any] struct {
	buffer  []optional.Optional[T]
	mux     *sync.Mutex
	start   int64
	end     int64
	maxSize int64
}

func NewLockQueue[T any](len int64) *LockQueue[T] {
	return &LockQueue[T]{
		buffer:  make([]optional.Optional[T], len+1),
		mux:     &sync.Mutex{},
		start:   0,
		end:     0,
		maxSize: len,
	}
}

func (q *LockQueue[T]) Len() int64 {
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

	optional := optional.NewOptional(&val)
	q.mux.Lock()
	q.buffer[q.end] = optional
	q.end++
	if q.end >= q.maxSize+1 {
		q.end = 0
	}
	q.mux.Unlock()

	return true
}

func (q *LockQueue[T]) Dequeue() optional.Optional[T] {
	if q.Len() == 0 {
		return optional.NewOptional[T](nil)
	}

	q.mux.Lock()
	v := q.buffer[q.start]
	q.start++
	if q.start >= q.maxSize+1 {
		q.start = 0
	}
	q.mux.Unlock()

	return v
}
