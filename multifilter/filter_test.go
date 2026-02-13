package multifilter

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func filterWrapper[DimType comparable, R any, FilterType Multifilter[DimType, R]](
	fn func([]FilterHelper[DimType, R], int) (FilterType, error),
) func([]FilterHelper[DimType, R], int) (Multifilter[DimType, R], error) {
	return func(f []FilterHelper[DimType, R], dim int) (Multifilter[DimType, R], error) {
		return fn(f, dim)
	}
}

type filterFnType func([]FilterHelper[string, string], int) (Multifilter[string, string], error)

var AllFilterTypes = map[string]filterFnType{
	"Hashfilter": filterWrapper(
		NewHashFilter[string, string],
	),
}

type searchCase struct {
	dimensions      [][]string
	expectedResults []string
}

func TestFilter_Search(t *testing.T) {
	tests := map[string]struct {
		initVal         []FilterHelper[string, string]
		dimensions      int
		searchCases     []searchCase
		expectedResults []string
		expectedErr     error
	}{
		"success case with 1 dimension": {
			initVal: []FilterHelper[string, string]{
				{
					FilterDimensions: [][]string{
						{"sedan"},
						{"hybrid"},
						{"4 door"},
					},
					Result: "Car 1",
				},
			},
			dimensions: 3,
			searchCases: []searchCase{
				{
					dimensions:      [][]string{{"sedan"}},
					expectedResults: []string{"Car 1"},
				},
			},
			expectedErr: nil,
		},
	}

	for filterName, fn := range AllFilterTypes {
		for name, tc := range tests {
			t.Run(fmt.Sprintf("%s-%s", filterName, name), func(t *testing.T) {
				multiFilter, err := fn(tc.initVal, tc.dimensions)
				assert.Equal(t, tc.expectedErr, err)

				for _, sc := range tc.searchCases {
					results := multiFilter.Search(sc.dimensions...)
					assert.ElementsMatch(t, sc.expectedResults, results)
				}
			})
		}
	}
}
