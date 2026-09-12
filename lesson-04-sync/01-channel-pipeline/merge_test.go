package pipeline

import (
	"testing"
	"time"
)

func TestMerge(t *testing.T) {
	a, b := make(chan int, 2), make(chan int, 1)
	a <- 1
	a <- 2
	close(a)
	b <- 3
	close(b)
	out := Merge(a, b)
	got := map[int]bool{}
	deadline := time.After(time.Second)
	for len(got) < 3 {
		select {
		case v, ok := <-out:
			if !ok {
				t.Fatalf("output closed early: %#v", got)
			}
			got[v] = true
		case <-deadline:
			t.Fatal("Merge timed out; who closes output?")
		}
	}
	if _, ok := <-out; ok {
		t.Fatal("output should be closed")
	}
}
