package optional

type Optional[T any] struct{ v *T }

func NewOptional[T any](value *T) Optional[T] {
	return Optional[T]{v: value}
}

func (o Optional[T]) IsEmpty() bool {
	return o.v == nil
}

func (o Optional[T]) Unwrap() T {
	return *o.v
}
