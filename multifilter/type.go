package multifilter

type Multifilter[T comparable, R any] interface {
	// Search returns data based on dimensions T.
	// To exclude a dimension from search, pass
	// in nil as its search value.
	Search(dimensions ...[]T) []R
}
