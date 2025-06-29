package heap

import (
	"golang.org/x/exp/constraints"
)

var (
	_ PriorityQueue[int] = (*MinHeap[int])(nil)
)

type MinHeap[T constraints.Ordered] struct {
	buffer []T
}

func NewMinHeap[T constraints.Ordered]() *MinHeap[T] {
	return &MinHeap[T]{buffer: make([]T, 0)}
}

func (p *MinHeap[T]) Insert(val T) {
	p.buffer = append(p.buffer, val)

	// Bubble up sequence
	for i := len(p.buffer) - 1; i > 0; {
		parentIdx := p.parent(i)
		if p.buffer[i] > p.buffer[parentIdx] {
			break
		}
		p.buffer[parentIdx], p.buffer[i] = p.buffer[i], p.buffer[parentIdx]
		i = parentIdx
	}
}

func (p *MinHeap[T]) Pop() (T, bool) {
	len := p.Len()
	if len <= 0 {
		var empty T
		return empty, false
	}

	result := p.buffer[0]
	if len == 1 {
		p.buffer = p.buffer[:0]
		return result, true
	}

	// Maintain heap property.
	p.buffer[0], p.buffer[len-1] = p.buffer[len-1], p.buffer[0]
	len -= 1

	for start := 0; start < len; {
		startVal := p.buffer[start]

		lIdx := p.left(start)
		rIdx := p.right(start)

		if lIdx >= len && rIdx >= len {
			break
		}

		if rIdx >= len {
			if p.buffer[lIdx] < startVal {
				p.buffer[lIdx], p.buffer[start] = p.buffer[start], p.buffer[lIdx]
			}
			break
		}

		lVal := p.buffer[lIdx]
		rVal := p.buffer[rIdx]

		if startVal < lVal && startVal < rVal {
			break
		}

		if lVal < rVal {
			p.buffer[lIdx], p.buffer[start] = p.buffer[start], p.buffer[lIdx]
			start = lIdx
			continue
		}

		p.buffer[rIdx], p.buffer[start] = p.buffer[start], p.buffer[rIdx]
		start = rIdx
	}

	p.buffer = p.buffer[:len]
	return result, true
}

func (p *MinHeap[T]) Len() int {
	return len(p.buffer)
}

func (p *MinHeap[T]) left(idx int) int {
	return idx*2 + 1
}

func (p *MinHeap[T]) right(idx int) int {
	return idx*2 + 2
}

func (p *MinHeap[T]) parent(idx int) int {
	return (idx - 1) / 2
}
