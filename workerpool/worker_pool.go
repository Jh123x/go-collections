package workerpool

import (
	"sync"
)

// WorkerFn defines the function signature for a worker function that processes tasks of type T and returns results of type R.
type WorkerFn[T, R any] func(T) R

// WorkerPool represents a pool of workers that can process tasks concurrently.
type WorkerPool[T, R any] struct {
	workerChan chan T
	resultChan chan R
	wg         *sync.WaitGroup
}

// NewWorker creates a new worker pool with the specified number of workers and a worker function.
// The worker function takes an input of type T and returns a result of type R.
func NewWorker[T, R any](workerCount, qSize int, workerFn WorkerFn[T, R]) *WorkerPool[T, R] {
	wp := &WorkerPool[T, R]{
		workerChan: make(chan T, qSize),
		resultChan: make(chan R, qSize),
		wg:         &sync.WaitGroup{},
	}

	// Spawn workers
	for i := 0; i < workerCount; i++ {
		go func() {
			for task := range wp.workerChan {
				result := workerFn(task)
				wp.resultChan <- result
				wp.wg.Done()
			}
		}()
	}

	return wp
}

// Pull returns a result from the worker in the pool.
func (wp *WorkerPool[T, R]) Pull() R {
	return <-wp.resultChan
}

// Push adds a task to the worker pool.
func (wp *WorkerPool[T, R]) Push(task T) {
	wp.wg.Add(1)
	wp.workerChan <- task
}

// Done signals that no more tasks will be added to the worker pool.
func (wp *WorkerPool[T, R]) Done() {
	close(wp.workerChan)
}

// Join waits for all tasks to be processed
func (wp *WorkerPool[T, R]) Join() {
	wp.wg.Wait()
	close(wp.resultChan)
}
