package workerpool

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type empty struct{}

func TestWorkerPool(t *testing.T) {
	// Example worker function that squares an integer
	const count = 1000
	workerFn := func(x int) int { return x * x }

	// Create a worker pool with 3 workers
	pool := NewWorker(3, count, workerFn)
	defer pool.Join()

	// Push tasks to the worker pool
	expectedResults := make(map[int]empty, count)
	for i := 0; i < count; i++ {
		expectedResults[workerFn(i)] = empty{}
		pool.Push(i)
	}
	pool.Done()

	// Pull results from the worker pool
	for i := 0; i < count; i++ {
		result := pool.Pull()
		_, ok := expectedResults[result]
		assert.True(t, ok)
		delete(expectedResults, result)
	}

	assert.Empty(t, expectedResults)
}

func BenchmarkWorkerPool(b *testing.B) {
	tests := map[string]struct {
		workerCount int
	}{
		"1 worker":   {workerCount: 1},
		"2 workers":  {workerCount: 2},
		"4 workers":  {workerCount: 4},
		"8 workers":  {workerCount: 8},
		"16 workers": {workerCount: 16},
		"32 workers": {workerCount: 32},
		"64 workers": {workerCount: 64},
	}

	for name, tc := range tests {
		b.Run(name, func(b *testing.B) {
			workerFn := func(x int) int { return x * x }
			pool := NewWorker(tc.workerCount, b.N, workerFn)
			defer pool.Join()
			for i := 0; i < b.N; i++ {
				pool.Push(i)
			}
			pool.Done()
			for i := 0; i < b.N; i++ {
				_ = pool.Pull()
			}
		})
	}

}
