package set


var (
	_ Set[string] = HashSet[string]{}
)

type empty struct{}

type HashSet[T comparable] map[T]empty

func NewHashSet[T comparable]() HashSet[T] {
	return make(HashSet[T])
}

func (h HashSet[T]) Add(values ...T) {
	for _, value := range values {
		if _, ok := h[value]; ok {
			continue
		}

		h[value] = empty{}
	}
}

func (h HashSet[T]) Has(value T) bool {
	_, ok := h[value]
	return ok
}

func (h HashSet[T]) Intersect(other Set[T]) Set[T] {
	results := NewHashSet[T]()

	for _, val := range other.ToSlice() {
		if _, ok := h[val]; !ok {
			continue
		}
		results.Add(val)
	}

	return results
}

func (h HashSet[T]) Union(other Set[T]) Set[T] {
	results := NewHashSet[T]()
	for _, val := range other.ToSlice() {
		results.Add(val)
	}

	for k := range h {
		results.Add(k)
	}

	return results
}

func (h HashSet[T]) ToSlice() []T {
	results := make([]T, 0, len(h))

	for k := range h {
		results = append(results, k)
	}

	return results
}

func (h HashSet[T]) Len() int {
	return len(h)
}
