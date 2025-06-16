package stack

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

const defaultStackSize = 5

type testFn[T any] func(*testing.T, LenStack[T])

func slWrap[T any, Q LenStack[T]](fn func(int) Q, size int) func() LenStack[T] {
	return func() LenStack[T] { return fn(size) }
}

func TestStackCorrectness(t *testing.T) {
	stacks := map[string]func() LenStack[int]{
		"LockStack": slWrap(NewLockStack[int], defaultStackSize),
	}

	tests := map[string]testFn[int]{
		"FILO Property":      testFiloCond,
		"Length Correctness": testLenCond,
	}

	for name, fn := range stacks {
		t.Run(name, func(t *testing.T) {
			for name, tc := range tests {
				t.Run(name, func(t *testing.T) {
					tc(t, fn())
				})
			}
		})
	}
}

func testFiloCond(t *testing.T, s LenStack[int]) {
	for idx := range defaultStackSize {
		assert.True(t, s.Push(idx))
	}

	assert.False(t, s.Push(defaultStackSize+1))
	for idx := range defaultStackSize {
		item, ok := s.Pop()
		assert.True(t, ok)
		assert.Equal(t, defaultStackSize-idx-1, item)
	}

	_, ok := s.Pop()
	assert.False(t, ok)
}

func testLenCond(t *testing.T, s LenStack[int]) {
	assert.Equal(t, 0, s.Len())
	for idx := range defaultStackSize {
		assert.True(t, s.Push(idx))
		assert.Equal(t, idx+1, s.Len())
	}

	assert.False(t, s.Push(defaultStackSize+1))
	assert.Equal(t, defaultStackSize, s.Len())
	for idx := range defaultStackSize {
		item, ok := s.Pop()
		assert.True(t, ok)
		assert.Equal(t, defaultStackSize-idx-1, item)
		assert.Equal(t, defaultStackSize-idx-1, s.Len())
	}

	assert.Equal(t, 0, s.Len())
	_, ok := s.Pop()
	assert.False(t, ok)
	assert.Equal(t, 0, s.Len())
}

func testRaceCond(t *testing.T, s LenStack[int]) {
	assert.Equal(t, s.Len(), 0)
	wg := sync.WaitGroup{}
	wg.Add(defaultStackSize)
	for idx := range defaultStackSize {
		go func() {
			defer wg.Done()
			assert.True(t, s.Push(idx))

			_, ok := s.Pop()
			assert.True(t, ok)
		}()
	}
	wg.Wait()
	assert.Equal(t, s.Len(), 0)
}
func BenchmarkSequentialStack(b *testing.B) {
	queues := map[string]func() LenStack[int]{
		"LockQueue": slWrap(NewLockStack[int], defaultStackSize),
	}

	for name, fn := range queues {
		b.Run(name, func(b *testing.B) {
			b.ReportAllocs()
			q := fn()
			for idx := range b.N {
				q.Push(idx)
				if val, ok := q.Pop(); !ok || val != idx {
					b.Fail()
				}
			}
		})
	}
}

func TestStackParallel(t *testing.T) {
	queues := map[string]func() LenStack[int]{
		"LockStack": slWrap(NewLockStack[int], defaultStackSize),
	}

	tests := map[string]testFn[int]{
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
