package missing

import (
	"sort"
	"testing"
)

func TestCollect(t *testing.T) {
	got := Collect([]int{1, 2, 3, 4})
	sort.Ints(got)
	want := []int{1, 4, 9, 16}
	if len(got) != len(want) {
		t.Fatalf("Collect returned %d results, want %d; who waits for goroutines?", len(got), len(want))
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("Collect() = %#v, want %#v", got, want)
		}
	}
}
