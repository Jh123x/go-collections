package stack

type Stack[T any] interface {
	// Push pushes the element to the top of the stack.
	// Return true if the push is successful
	Push(val T) bool

	// Pop returns the element at the top and removes it.
	// Returns false if the stack is empty
	Pop() (val T, ok bool)
}

type LenStack[T any] interface {
	// Len returns the size of the stack.
	Len() int

	// Should also implement the other stack methods.
	Stack[T]
}
