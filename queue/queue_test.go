package queue

import (
	"strconv"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

const defaultQSize = 5

type testFn[T any] func(*testing.T, Queue[T])

func qWrap[T any, Q Queue[T]](fn func(int64) Q, size int64) func() Queue[T] {
	return func() Queue[T] { return fn(size) }
}

func TestQueueCorrectness(t *testing.T) {
	queues := map[string]func() Queue[string]{
		"LockQueue": qWrap(NewLockQueue[string], defaultQSize),
	}

	tests := map[string]testFn[string]{
		"FIFO Property": func(t *testing.T, q Queue[string]) {
			for idx := range defaultQSize {
				assert.True(t, q.Enqueue("test"+strconv.FormatInt(int64(idx), 10)))
			}

			assert.False(t, q.Enqueue("test"+strconv.FormatInt(int64(defaultQSize), 10)))

			for idx := range defaultQSize {
				item := q.Dequeue()
				assert.False(t, item.IsEmpty())
				assert.Equal(t, item.Unwrap(), "test"+strconv.FormatInt(int64(idx), 10))
			}

			item := q.Dequeue()
			assert.True(t, item.IsEmpty())
		},
		"Length Correctness": func(t *testing.T, q Queue[string]) {
			for idx := range defaultQSize {
				assert.True(t, q.Enqueue("test"+strconv.FormatInt(int64(idx), 10)))
				assert.Equal(t, q.Len(), int64(idx+1))
			}

			assert.False(t, q.Enqueue("test"+strconv.FormatInt(int64(defaultQSize), 10)))
			assert.Equal(t, q.Len(), int64(defaultQSize))

			for idx := range defaultQSize {
				expectedItem := "test" + strconv.FormatInt(int64(idx), 10)
				assert.Equal(t, q.Len(), int64(defaultQSize-idx))

				item := q.Dequeue()
				assert.False(t, item.IsEmpty())
				assert.Equal(t, item.Unwrap(), expectedItem)
				assert.Equal(t, q.Len(), int64(defaultQSize-idx-1))
			}

			item := q.Dequeue()
			assert.True(t, item.IsEmpty())
			assert.Equal(t, q.Len(), int64(0))
		},
		"Queue Dequeue correctness": func(t *testing.T, q Queue[string]) {
			assert.Equal(t, q.Len(), int64(0))
			for idx := range 1000 {
				addVal := "test" + strconv.FormatInt(int64(idx), 10)
				assert.True(t, q.Enqueue(addVal))
				assert.Equal(t, q.Len(), int64(1))

				item := q.Dequeue()
				assert.False(t, item.IsEmpty())
				assert.Equal(t, item.Unwrap(), addVal)
				assert.Equal(t, q.Len(), int64(0))
			}
			assert.Equal(t, q.Len(), int64(0))
		},
		"Check Race Condition": func(t *testing.T, q Queue[string]) {
			assert.Equal(t, q.Len(), int64(0))
			wg := sync.WaitGroup{}
			wg.Add(defaultQSize)
			for idx := range defaultQSize {
				go func() {
					defer wg.Done()
					addVal := "test" + strconv.FormatInt(int64(idx), 10)
					assert.True(t, q.Enqueue(addVal))

					item := q.Dequeue()
					assert.False(t, item.IsEmpty())
				}()
			}
			wg.Wait()
			assert.Equal(t, q.Len(), int64(0))
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
