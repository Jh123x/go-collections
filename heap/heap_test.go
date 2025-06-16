package heap

import (
	"fmt"
	"math/rand/v2"
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/exp/constraints"
)

func wrap[T constraints.Ordered, heap PriorityQueue[T]](fn func() heap) func() PriorityQueue[T] {
	return func() PriorityQueue[T] { return fn() }
}

func TestHeap(t *testing.T) {
	heaps := map[string]func() PriorityQueue[int]{
		"Heap": wrap(NewMinHeap[int]),
	}
	tests := map[string]func(t *testing.T, heap PriorityQueue[int]){
		"Push Pop consecutive": testPushPop,
		"Random Sample":        testRandomSample,
		"Push Pop interleave":  testPushPopInterleaf,
	}

	for name, fn := range heaps {
		t.Run(name, func(t *testing.T) {
			for name, tc := range tests {
				t.Run(name, func(t *testing.T) {
					tc(t, fn())
				})
			}
		})
	}
}

func testPushPopInterleaf(t *testing.T, pq PriorityQueue[int]) {
	pq.Insert(45)
	assert.Equal(t, 1, pq.Len())
	pq.Insert(20)
	assert.Equal(t, 2, pq.Len())
	pq.Insert(14)
	assert.Equal(t, 3, pq.Len())
	pq.Insert(12)
	assert.Equal(t, 4, pq.Len())
	pq.Insert(31)
	assert.Equal(t, 5, pq.Len())
	pq.Insert(7)
	assert.Equal(t, 6, pq.Len())
	pq.Insert(11)
	assert.Equal(t, 7, pq.Len())
	pq.Insert(13)
	assert.Equal(t, 8, pq.Len())
	pq.Insert(7)
	assert.Equal(t, 9, pq.Len())
	v, ok := pq.Pop()
	assert.True(t, ok)
	assert.Equal(t, 7, v)
	assert.Equal(t, 8, pq.Len())
	pq.Insert(0)
	assert.Equal(t, 9, pq.Len())

	v, ok = pq.Pop()
	assert.True(t, ok)
	assert.Equal(t, 0, v)
	assert.Equal(t, 8, pq.Len())

	v, ok = pq.Pop()
	assert.True(t, ok)
	assert.Equal(t, 7, v)
	assert.Equal(t, 7, pq.Len())

	v, ok = pq.Pop()
	assert.True(t, ok)
	assert.Equal(t, 11, v)
	assert.Equal(t, 6, pq.Len())

	v, ok = pq.Pop()
	assert.True(t, ok)
	assert.Equal(t, 12, v)
	assert.Equal(t, 5, pq.Len())

	v, ok = pq.Pop()
	assert.True(t, ok)
	assert.Equal(t, 13, v)
	assert.Equal(t, 4, pq.Len())
}

func testPushPop(t *testing.T, pq PriorityQueue[int]) {
	for idx := range 10 {
		pq.Insert(10 - idx)
		assert.Equal(t, idx+1, pq.Len())
	}

	for idx := range 10 {
		v, ok := pq.Pop()
		assert.True(t, ok)
		assert.Equal(t, idx+1, v)
		assert.Equal(t, 10-idx-1, pq.Len())
	}
}

func testRandomSample(t *testing.T, pq PriorityQueue[int]) {
	insertSize := 10
	for idx := range insertSize {
		no := rand.IntN(5000)
		pq.Insert(no)
		fmt.Println(no)
		assert.Equal(t, idx+1, pq.Len())
	}
	assert.Equal(t, insertSize, pq.Len())

	prev, ok := pq.Pop()
	assert.True(t, ok)
	assert.Equal(t, insertSize-1, pq.Len())

	for i := range insertSize - 1 {
		v, ok := pq.Pop()
		assert.True(t, ok)
		assert.True(t, prev <= v, "%d <= %d, seq: %v", prev, v, pq)
		assert.Equal(t, insertSize-i-2, pq.Len())
		prev = v
	}
}
