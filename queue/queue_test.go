package queue

import (
	"strconv"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

const defaultQSize = 5

type testFn[T any] func(*testing.T, LenQueue[T])

func qlWrap[T any, Q LenQueue[T]](fn func(int) Q, size int) func() LenQueue[T] {
	return func() LenQueue[T] { return fn(size) }
}

func qWrap[T any, Q Queue[T]](fn func(int) Q, size int) func() Queue[T] {
	return func() Queue[T] { return fn(size) }
}

func TestQueueCorrectness(t *testing.T) {
	queues := map[string]func() LenQueue[string]{
		"LockQueue": qlWrap(NewLockQueue[string], defaultQSize),
		"StdQueue":  qlWrap(NewStdQueue[string], defaultQSize),
	}

	tests := map[string]testFn[string]{
		"FIFO Property":             testFifoCond,
		"Length Correctness":        testLenCond,
		"Queue Dequeue correctness": testDequeueCond,
	}

	for name, fn := range queues {
		t.Run(name, func(t *testing.T) {
			for name, tc := range tests {
				t.Run(name, func(t *testing.T) {
					tc(t, fn())
				})
			}
		})
	}
}

func testFifoCond(t *testing.T, q LenQueue[string]) {
	for idx := range defaultQSize {
		assert.True(t, q.Enqueue("test"+strconv.FormatInt(int64(idx), 10)))
	}

	assert.False(t, q.Enqueue("test"+strconv.FormatInt(int64(defaultQSize), 10)))

	for idx := range defaultQSize {
		item, ok := q.Dequeue()
		assert.True(t, ok)
		assert.Equal(t, item, "test"+strconv.FormatInt(int64(idx), 10))
	}

	_, ok := q.Dequeue()
	assert.False(t, ok)
}

func testLenCond(t *testing.T, q LenQueue[string]) {
	for idx := range defaultQSize {
		assert.True(t, q.Enqueue("test"+strconv.FormatInt(int64(idx), 10)))
		assert.Equal(t, q.Len(), idx+1)
	}

	assert.False(t, q.Enqueue("test"+strconv.FormatInt(int64(defaultQSize), 10)))
	assert.Equal(t, q.Len(), defaultQSize)

	for idx := range defaultQSize {
		expectedItem := "test" + strconv.FormatInt(int64(idx), 10)
		assert.Equal(t, q.Len(), defaultQSize-idx)

		item, ok := q.Dequeue()
		assert.True(t, ok)
		assert.Equal(t, item, expectedItem)
		assert.Equal(t, q.Len(), defaultQSize-idx-1)
	}

	_, ok := q.Dequeue()
	assert.False(t, ok)
	assert.Equal(t, q.Len(), 0)
}

func testDequeueCond(t *testing.T, q LenQueue[string]) {
	assert.Equal(t, 0, q.Len())
	for idx := range 1000 {
		addVal := "test" + strconv.FormatInt(int64(idx), 10)
		assert.True(t, q.Enqueue(addVal))
		assert.Equal(t, q.Len(), 1)

		item, ok := q.Dequeue()
		assert.True(t, ok)
		assert.Equal(t, item, addVal)
		assert.Equal(t, q.Len(), 0)
	}
	assert.Equal(t, q.Len(), 0)
}

func testRaceCond(t *testing.T, q LenQueue[string]) {
	assert.Equal(t, q.Len(), 0)
	wg := sync.WaitGroup{}
	wg.Add(defaultQSize)
	for idx := range defaultQSize {
		go func() {
			defer wg.Done()
			addVal := "test" + strconv.FormatInt(int64(idx), 10)
			assert.True(t, q.Enqueue(addVal))

			_, ok := q.Dequeue()
			assert.True(t, ok)
		}()
	}
	wg.Wait()
	assert.Equal(t, q.Len(), 0)
}

func BenchmarkSequentialQueues(b *testing.B) {
	queues := map[string]func() Queue[int]{
		"LockQueue": qWrap(NewLockQueue[int], defaultQSize),
		"StdQueue":  qWrap(NewStdQueue[int], defaultQSize),
	}

	for name, fn := range queues {
		b.Run(name, func(b *testing.B) {
			b.ReportAllocs()
			q := fn()
			for idx := range b.N {
				q.Enqueue(idx)
				if val, ok := q.Dequeue(); !ok || val != idx {
					b.Fail()
				}
			}
		})
	}
}

func TestQueueParallel(t *testing.T) {
	queues := map[string]func() LenQueue[string]{
		"LockQueue": qlWrap(NewLockQueue[string], defaultQSize),
	}

	tests := map[string]testFn[string]{
		"Check Race Condition": testRaceCond,
	}

	for name, fn := range queues {
		t.Run(name, func(t *testing.T) {
			for name, tc := range tests {
				t.Run(name, func(t *testing.T) {
					tc(t, fn())
				})
			}
		})
	}
}

func BenchmarkParallelQueues(b *testing.B) {
	queues := map[string]func() Queue[int]{
		"LockQueue": qWrap(NewLockQueue[int], defaultQSize),
	}

	for name, fn := range queues {
		b.Run(name, func(b *testing.B) {
			b.ReportAllocs()
			q := fn()
			wg := sync.WaitGroup{}
			wg.Add(b.N)
			for idx := range b.N {
				wg.Done()
				q.Enqueue(idx)
				if _, ok := q.Dequeue(); !ok {
					b.Fail()
				}
			}
		})
	}
}
