package heap

import "golang.org/x/exp/constraints"

type PriorityQueue[T constraints.Ordered] interface {
	Insert(v T)
	Pop() (T, bool)
	Len() int
}
