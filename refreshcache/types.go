package refreshcache

import (
	"sync/atomic"
	"time"
)

type LoaderFunc[T any] func() (T, error)

type RefreshCache[T any] struct {
	value *atomic.Pointer[T]
	err   *atomic.Pointer[error]
}

// NewRefreshCache creates a new RefreshCache that uses the provided loader function.
// The cache is refreshed at the specified refreshInterval.
func NewRefreshCache[T any](loader LoaderFunc[T], refreshInterval time.Duration) RefreshCache[T] {
	cache := RefreshCache[T]{
		value: &atomic.Pointer[T]{},
		err:   &atomic.Pointer[error]{},
	}

	// Initial Load of the cache
	v, err := loader()
	cache.value.Store(&v)
	cache.err.Store(&err)

	go func() {
		ticker := time.NewTicker(refreshInterval)
		defer ticker.Stop()

		for {
			// Wait for the next tick to refresh
			<-ticker.C
			v, err := loader()

			// Let GC Collect the result of the previous load
			_ = cache.value.Swap(&v)
			_ = cache.err.Swap(&err)
		}
	}()

	return cache
}

func (r RefreshCache[T]) Get() (T, error) {
	v := r.value.Load()
	err := r.err.Load()

	if v == nil || err == nil {
		return *new(T), nil
	}

	return *v, *err
}

func main() {
	// Example usage
	loader := func() (int, error) {
		return 42, nil
	}

	cache := NewRefreshCache(loader, time.Second*10)
	value, _ := cache.Get()
	println(value) // Should print 42

	// Shallow copy of the cache
	cache2 := cache
	value2, _ := cache2.Get()
	println(value2) // Should also print 42
}
