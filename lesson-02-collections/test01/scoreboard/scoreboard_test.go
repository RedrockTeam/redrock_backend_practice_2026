package scoreboard

import "testing"

func TestTotals(t *testing.T) {
	in := map[string][]int{"Lin": {80, 90}, "Wang": {100}}
	got := Totals(in)
	if got["Lin"] != 170 || got["Wang"] != 100 {
		t.Fatalf("Totals() = %#v", got)
	}
}
func TestTop(t *testing.T) {
	got := Top(map[string]int{"Wang": 100, "Lin": 170, "Zhao": 170}, 2)
	if len(got) != 2 || got[0].Name != "Lin" || got[1].Name != "Zhao" {
		t.Fatalf("Top() = %#v", got)
	}
}
