package node

type Node[T any] struct {
	Val  T
	Next *Node[T]
}
