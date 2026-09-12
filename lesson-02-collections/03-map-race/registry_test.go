package registry

import (
	"sync"
	"testing"
)

func TestConcurrentRegistryAccess(t *testing.T) {
	r := New()
	var wg sync.WaitGroup
	for i := 0; i < 4; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for j := 0; j < 100; j++ {
				r.Set("shared", id*100+j)
				_, _ = r.Get("shared")
			}
		}(i)
	}
	wg.Wait()
}
