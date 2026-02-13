package multifilter

import (
	"fmt"

	"github.com/Jh123x/go-collections/set"
)

var (
	_ Multifilter[string, string] = (*HashFilter[string, string])(nil)
)

type HashFilter[T comparable, R any] struct {
	data    []R
	filters []map[T]set.Set[int]
}

func NewHashFilter[T comparable, R any](results []FilterHelper[T, R], dimensions int) (*HashFilter[T, R], error) {
	data := make([]R, 0, len(results))
	filters := make([]map[T]set.Set[int], dimensions)

	for idx := range filters {
		filters[idx] = make(map[T]set.Set[int])
	}

	for resultIdx, val := range results {
		data = append(data, val.Result)

		if len(val.FilterDimensions) > dimensions {
			return nil, fmt.Errorf("too many dimensions found")
		}

		for filterIdx, dimensionValues := range val.FilterDimensions {
			for _, dimVal := range dimensionValues {
				if _, ok := filters[filterIdx][dimVal]; !ok {
					filters[filterIdx][dimVal] = set.NewHashSet[int]()
				}
				filters[filterIdx][dimVal].Add(resultIdx)
			}
		}
	}
	return &HashFilter[T, R]{
		data:    data,
		filters: filters,
	}, nil
}

func (h *HashFilter[T, R]) Search(dimensions ...[]T) []R {
	var results set.Set[int]
	hasFirstSet := false
	for dimIdx, dimension := range dimensions {
		if dimension == nil {
			continue
		}

		tmpAcc := set.NewHashSet[int]()
		for _, d := range dimension {
			tmp := h.filters[dimIdx][d]
			tmpAcc.Add(tmp.ToSlice()...)
		}

		if !hasFirstSet {
			results = tmpAcc
			continue
		}

		results = results.Intersect(tmpAcc)
	}

	if !hasFirstSet {
		return h.data
	}

	finalReturn := make([]R, 0, results.Len())
	for _, resultIdx := range results.ToSlice() {
		finalReturn = append(finalReturn, h.data[resultIdx])
	}

	return finalReturn
}
