package queue

import (
	"sync/atomic"

	"github.com/Jh123x/go-collections/optional"
)

var (
	_ Queue[string] = (*ChanQueue[string])(nil)
)

type ChanQueue[T any] struct {
	buffer    chan T
	currVal   atomic.Pointer[T]
	currSize  int64
	totalSize int64
}

func NewChanQueue[T any](len int64) *ChanQueue[T] {
	return &ChanQueue[T]{
		buffer:    make(chan T, len),
		currVal:   atomic.Pointer[T]{},
		currSize:  0,
		totalSize: len,
	}
}

func (q *ChanQueue[T]) Peek() optional.Optional[T] {
	if q.currSize <= 0 {
		return optional.NewOptional[T](nil)
	}

	if val := q.currVal.Load(); val != nil {
		return optional.NewOptional(val)
	}

	front := <-q.buffer
	q.currVal.Store(&front)
	return optional.NewOptional(&front)
}

func (q *ChanQueue[T]) Len() int64 {
	return q.currSize
}

func (q *ChanQueue[T]) Enqueue(val T) bool {
	if q.currSize >= q.totalSize {
		return false
	}

	q.buffer <- val
	atomic.AddInt64(&q.currSize, 1)
	return true
}

func (q *ChanQueue[T]) Dequeue() optional.Optional[T] {
	v := q.Peek()
	q.currVal.Store(nil)
	if v.IsEmpty() {
		return v
	}

	atomic.AddInt64(&q.currSize, -1)
	return v
}
