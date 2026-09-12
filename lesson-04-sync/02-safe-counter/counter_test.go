package counter

import (
	"sync"
	"testing"
)

func TestCounterIsSafe(t *testing.T) {
	var c Counter
	var wg sync.WaitGroup
	for i := 0; i < 8; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < 1000; j++ {
				c.Add(1)
			}
		}()
	}
	wg.Wait()
	if c.Value() != 8000 {
		t.Fatalf("Value() = %d, want 8000", c.Value())
	}
}
