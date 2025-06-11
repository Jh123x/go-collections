package queue

import "github.com/Jh123x/go-collections/optional"

type Queue[T any] interface {
	// Peek returns the element at the front of the queue without removing it.
	// Returns nil if the queue is empty.
	Peek() optional.Optional[T]

	// Len returns the length of the queue.
	Len() int64

	// Enqueue adds the element to the end of the queue.
	// Returns true if enqueuing is successful.
	Enqueue(T) bool

	// Dequeue removes the element at the front of the queue and returns it.
	// Returns nil if the queue is empty
	Dequeue() optional.Optional[T]
}
