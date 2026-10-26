package scoreboard_test

import (
	"reflect"
	"testing"

	"redrock/teaching-exercises/lesson-02-collections/test02/scoreboard"
)

func TestTotalsSumsAndDoesNotMutateInput(t *testing.T) {
	records := map[string][]int{
		"Lin":  {80, 90},
		"Wang": {100},
		"Zhao": {},
	}
	before := map[string][]int{
		"Lin":  {80, 90},
		"Wang": {100},
		"Zhao": {},
	}

	got := scoreboard.Totals(records)
	want := map[string]int{"Lin": 170, "Wang": 100, "Zhao": 0}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("Totals() = %#v, want %#v", got, want)
	}
	if !reflect.DeepEqual(records, before) {
		t.Fatalf("Totals 修改了输入 records：%#v -> %#v", before, records)
	}
}

func TestTotalsEmptyAndNil(t *testing.T) {
	if got := scoreboard.Totals(nil); len(got) != 0 {
		t.Fatalf("Totals(nil) = %#v, want 空 map", got)
	}
	got := scoreboard.Totals(map[string][]int{"Lin": nil})
	if len(got) != 1 || got["Lin"] != 0 {
		t.Fatalf("Totals({Lin: nil}) = %#v, want {Lin: 0}", got)
	}
}

func TestTopSortsByScoreThenName(t *testing.T) {
	totals := map[string]int{"Wang": 100, "Lin": 170, "Zhao": 170, "Chen": 100}

	got := scoreboard.Top(totals, 3)
	want := []string{"Lin", "Zhao", "Chen"}
	if len(got) != len(want) {
		t.Fatalf("Top(..., 3) = %#v, want 3 条", got)
	}
	for i, name := range want {
		if got[i].Name != name {
			t.Fatalf("Top(..., 3)[%d].Name = %q, want %q；同分时按姓名升序", i, got[i].Name, name)
		}
	}
	if got[0].Score != 170 || got[2].Score != 100 {
		t.Fatalf("Top(..., 3) = %#v, 分数不正确", got)
	}
}

func TestTopEdgeCases(t *testing.T) {
	totals := map[string]int{"Lin": 90, "Wang": 80}

	for _, n := range []int{0, -1} {
		if got := scoreboard.Top(totals, n); len(got) != 0 {
			t.Errorf("Top(..., %d) = %#v, want 空结果", n, got)
		}
	}
	if got := scoreboard.Top(totals, 10); len(got) != 2 {
		t.Errorf("Top(..., 10) 返回 %d 条, want 2", len(got))
	}
	if got := scoreboard.Top(nil, 3); len(got) != 0 {
		t.Errorf("Top(nil, 3) = %#v, want 空结果", got)
	}
}

func TestTopDoesNotMutateTotals(t *testing.T) {
	totals := map[string]int{"Lin": 170, "Wang": 100}
	before := map[string]int{"Lin": 170, "Wang": 100}

	scoreboard.Top(totals, 2)

	if !reflect.DeepEqual(totals, before) {
		t.Fatalf("Top 修改了 totals：%#v -> %#v", before, totals)
	}
}
