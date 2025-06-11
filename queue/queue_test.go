package queue

import (
	"strconv"
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
		"ChanQueue": qWrap(NewChanQueue[string], defaultQSize),
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
		"Peek Correctness": func(t *testing.T, q Queue[string]) {
			for idx := range defaultQSize {
				assert.True(t, q.Enqueue("test"+strconv.FormatInt(int64(idx), 10)))

				pItem := q.Peek()
				if !assert.False(t, pItem.IsEmpty(), "Iteration %d", idx+1) {
					t.FailNow()
				}
				assert.Equal(t, pItem.Unwrap(), "test0")
			}

			assert.False(t, q.Enqueue("test"+strconv.FormatInt(int64(defaultQSize), 10)))

			for idx := range defaultQSize {
				expectedItem := "test" + strconv.FormatInt(int64(idx), 10)
				for _ = range 5 {
					pItem := q.Peek()
					assert.False(t, pItem.IsEmpty())
					assert.Equal(t, pItem.Unwrap(), expectedItem)
				}

				item := q.Dequeue()
				assert.False(t, item.IsEmpty())
				assert.Equal(t, item.Unwrap(), expectedItem)
			}

			item := q.Dequeue()
			assert.True(t, item.IsEmpty())
		},
		"Length Correctness": func(t *testing.T, q Queue[string]) {
			for idx := range defaultQSize {
				assert.True(t, q.Enqueue("test"+strconv.FormatInt(int64(idx), 10)))

				pItem := q.Peek()
				if !assert.False(t, pItem.IsEmpty(), "Iteration %d", idx+1) {
					t.FailNow()
				}
				assert.Equal(t, pItem.Unwrap(), "test0")
				assert.Equal(t, q.Len(), int64(idx+1))
			}

			assert.False(t, q.Enqueue("test"+strconv.FormatInt(int64(defaultQSize), 10)))
			assert.Equal(t, q.Len(), int64(defaultQSize))

			for idx := range defaultQSize {
				expectedItem := "test" + strconv.FormatInt(int64(idx), 10)
				for _ = range 5 {
					pItem := q.Peek()
					assert.False(t, pItem.IsEmpty())
					assert.Equal(t, pItem.Unwrap(), expectedItem)
					assert.Equal(t, q.Len(), int64(defaultQSize-idx))
				}

				item := q.Dequeue()
				assert.False(t, item.IsEmpty())
				assert.Equal(t, item.Unwrap(), expectedItem)
				assert.Equal(t, q.Len(), int64(defaultQSize-idx-1))
			}

			item := q.Dequeue()
			assert.True(t, item.IsEmpty())
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
