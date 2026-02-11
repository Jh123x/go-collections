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
		"non duplicates should count all":{
			addVals: []string{"test","test2","test3"},
			expectedLen: 3,
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
