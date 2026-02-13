package set

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

type setInitFn[T comparable] func() Set[T]

var (
	AllSetFns = map[string]setInitFn[string]{
		"HashSet": testWrapper(NewHashSet[string]),
	}
)

func testWrapper[StoreType comparable, SetType Set[StoreType]](initFn func() SetType) func() Set[StoreType] {
	return func() Set[StoreType] {
		return initFn()
	}
}

func TestSetInit(t *testing.T) {
	for name, fn := range AllSetFns {
		t.Run(name, func(t *testing.T) {
			defer func() { assert.Nil(t, recover()) }()
			setVal := fn()
			assert.NotNil(t, setVal)
		})
	}
}

func TestSet_Len(t *testing.T) {
	tests := map[string]struct {
		addVals     []string
		expectedLen int
	}{
		"duplicate vals count as one": {
			addVals:     []string{"test", "test", "test2"},
			expectedLen: 2,
		},
		"non duplicates should count all": {
			addVals:     []string{"test", "test2", "test3"},
			expectedLen: 3,
		},
		"all duplicates should count as 1": {
			addVals:     []string{"test", "test", "test", "test", "test"},
			expectedLen: 1,
		},
	}

	for containerName, fn := range AllSetFns {
		for name, tc := range tests {
			t.Run(fmt.Sprintf("%s-%s", containerName, name), func(t *testing.T) {
				set := fn()
				set.Add(tc.addVals...)
				assert.Equal(t, tc.expectedLen, set.Len(), "set val: %v", set.ToSlice())
			})
		}
	}
}

func TestSet_Add_Has(t *testing.T) {
	tests := map[string]struct {
		addVal         []string
		checkMapResult map[string]bool
	}{
		"Added values should exist": {
			addVal: []string{"test"},
			checkMapResult: map[string]bool{
				"test":  true,
				"test2": false,
			},
		},
		"Added multiple times should exist": {
			addVal: []string{"test", "test"},
			checkMapResult: map[string]bool{
				"test":  true,
				"test2": false,
			},
		},
	}

	for containerName, fn := range AllSetFns {
		for name, tc := range tests {
			t.Run(fmt.Sprintf("%s-%s", containerName, name), func(t *testing.T) {
				defer func() { assert.Nil(t, recover()) }()
				set := fn()
				assert.NotNil(t, set)

				set.Add(tc.addVal...)
				for val, expectedRes := range tc.checkMapResult {
					assert.Equal(t, expectedRes, set.Has(val))
				}
			})
		}
	}
}

func TestSet_Intersection(t *testing.T) {
	tests := map[string]struct {
		setA            []string
		setB            []string
		expectedResults []string
	}{
		"no intersection": {
			setA:            []string{"test", "test1", "test2"},
			setB:            []string{"test3", "test4", "test5"},
			expectedResults: []string{},
		},
		"has single intersection": {
			setA:            []string{"test", "test1"},
			setB:            []string{"test", "test2"},
			expectedResults: []string{"test"},
		},
		"same sets": {
			setA:            []string{"s", "t", "u", "v"},
			setB:            []string{"s", "t", "u", "v"},
			expectedResults: []string{"s", "t", "u", "v"},
		},
		"intersection with empty": {
			setA:            []string{},
			setB:            []string{"test"},
			expectedResults: []string{},
		},
	}

	for setName, fn := range AllSetFns {
		for name, tc := range tests {
			t.Run(
				fmt.Sprintf("%s-%s", setName, name),
				func(t *testing.T) {
					setA := fn()
					setA.Add(tc.setA...)
					setB := fn()
					setB.Add(tc.setB...)

					result := setA.Intersect(setB)
					assert.ElementsMatch(t, tc.expectedResults, result.ToSlice())
				},
			)
		}
	}
}

func TestSet_Union(t *testing.T) {
	tests := map[string]struct {
		setA          []string
		setB          []string
		expectResults []string
	}{
		"empty sets": {
			setA:          []string{},
			setB:          []string{},
			expectResults: []string{},
		},
		"1 empty set": {
			setA:          []string{},
			setB:          []string{"test", "test2"},
			expectResults: []string{"test", "test2"},
		},
		"1 empty set inverted": {
			setA:          []string{"test", "test2"},
			setB:          []string{},
			expectResults: []string{"test", "test2"},
		},
		"non empty sets": {
			setA:          []string{"test1", "test2"},
			setB:          []string{"test3", "test4"},
			expectResults: []string{"test1", "test2", "test3", "test4"},
		},
	}

	for setName, fn := range AllSetFns {
		for name, tc := range tests {
			t.Run(
				fmt.Sprintf("%s-%s", setName, name),
				func(t *testing.T) {
					setA := fn()
					setA.Add(tc.setA...)
					setB := fn()
					setB.Add(tc.setB...)

					results := setA.Union(setB)
					assert.ElementsMatch(t, tc.expectResults, results.ToSlice())
				},
			)
		}
	}
}
