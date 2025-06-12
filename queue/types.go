package queue

type Queue[T any] interface {

	// Enqueue adds the element to the end of the queue.
	// Returns true if enqueuing is successful.
	Enqueue(T) bool

	// Dequeue removes the element at the front of the queue and returns it.
	// Returns nil if the queue is empty
	Dequeue() (T, bool)
}

type LenQueue[T any] interface {
	// Len returns the length of the queue.
	Len() int

	// Should also implement the interface for a queue.
	Queue[T]
}
