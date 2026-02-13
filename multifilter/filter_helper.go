package multifilter

type FilterHelper[T comparable, R any] struct {
	FilterDimensions [][]T
	Result           R
}

// NewFilterHelper returns a typed filter helper based on the dimensions given.
// This is the way to setup for other multifilter data structures.
func NewFilterHelper[T comparable, R any](resultValue R, filters ...[]T) *FilterHelper[T, R] {
	return &FilterHelper[T, R]{
		FilterDimensions: filters,
		Result:           resultValue,
	}
}
