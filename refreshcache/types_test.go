package refreshcache

import (
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestRefreshCache_Get(t *testing.T) {
	counter := 0
	loader := func() (int, error) {
		if counter == 1 {
			return 0, fmt.Errorf("error on 5")
		}

		counter++
		return counter, nil
	}

	cache := NewRefreshCache(loader, time.Microsecond)
	v, err := cache.Get()
	assert.Equal(t, 1, v, "initial value should be 1")
	assert.NoError(t, err, "initial error should be nil")

	// Wait for the cache to refresh
	time.Sleep(time.Millisecond)

	v, err = cache.Get()
	assert.Equal(t, 0, v, "value should be zero on error")
	assert.Equal(t, fmt.Errorf("error on 5"), err, "error should be propagated")
}

func TestRefreshCache_InitialLoad(t *testing.T) {
	loader := func() (int, error) {
		return 10, nil
	}

	cache := NewRefreshCache(loader, time.Minute)
	v, err := cache.Get()
	assert.Equal(t, 10, v, "initial value should be 10")
	assert.NoError(t, err, "initial error should be nil")
}

func TestRefreshCache_ParallelAccess(t *testing.T) {
	refreshCount := 0
	loader := func() (int, error) {
		refreshCount++
		return 5, nil
	}

	cache := NewRefreshCache(loader, time.Nanosecond)
	wg := &sync.WaitGroup{}
	wg.Add(1000)
	for range 1000 {
		go func() {
			v, err := cache.Get()
			assert.Equal(t, 5, v, "value should be 5")
			assert.NoError(t, err, "error should be nil")
			wg.Done()
		}()
	}

	wg.Wait()
	assert.Greater(t, refreshCount, 1, "loader should have been called multiple times")
}

func BenchmarkRefreshCache(b *testing.B) {
	loader := func() (int, error) {
		result := 0
		for i := range 10000000 {
			result += i
		}
		return result, nil
	}

	cache := NewRefreshCache(loader, time.Minute)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = cache.Get()
	}
}
