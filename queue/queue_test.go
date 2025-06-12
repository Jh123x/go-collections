package queue

import (
	"strconv"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

const defaultQSize = 5

type testFn[T any] func(*testing.T, Queue[T])

func qWrap[T any, Q Queue[T]](fn func(int) Q, size int) func() Queue[T] {
	return func() Queue[T] { return fn(size) }
}

func TestQueueCorrectness(t *testing.T) {
	queues := map[string]func() Queue[string]{
		"LockQueue": qWrap(NewLockQueue[string], defaultQSize),
		"StdQueue":  qWrap(NewStdQueue[string], defaultQSize),
	}

	tests := map[string]testFn[string]{
		"FIFO Property": func(t *testing.T, q Queue[string]) {
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
		},
		"Length Correctness": func(t *testing.T, q Queue[string]) {
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
		},
		"Queue Dequeue correctness": func(t *testing.T, q Queue[string]) {
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
		},
		"Check Race Condition": func(t *testing.T, q Queue[string]) {
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
		},
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

func BenchmarkQueues(b *testing.B) {
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
