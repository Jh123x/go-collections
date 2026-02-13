package set

type Set[T comparable] interface {
	// Add adds the value of type T to the set.
	// Does nothing if it already exists
	Add(values ...T)

	// Has returns true if value is contained within the set.
	Has(value T) bool

	// Intersect returns the set Intersection between the current set and the other set.
	// Returns the Set[T] of the current concrete type.
	Intersect(other Set[T]) Set[T]

	// Union returns the Union between the current set and the other set.
	// Returns the Set[T]  of the current concrete type.
	Union(other Set[T]) Set[T]

	// Returns the list of contained value.
	ToSlice() []T

	// Returns the size of the container.
	Len() int
}
