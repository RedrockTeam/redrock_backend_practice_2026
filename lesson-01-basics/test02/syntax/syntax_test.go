package syntax_test

import (
	"testing"

	"redrock/teaching-exercises/lesson-01-basics/test02/syntax"
)

func TestAdd(t *testing.T) {
	if got := syntax.Add(3, 4); got != 7 {
		t.Fatalf("Add(3, 4) = %d, want 7", got)
	}
	if got := syntax.Add(-2, 2); got != 0 {
		t.Fatalf("Add(-2, 2) = %d, want 0", got)
	}
}

func TestClassifyNumber(t *testing.T) {
	if got := syntax.ClassifyNumber(8); got != "positive" {
		t.Fatalf("ClassifyNumber(8) = %q, want %q", got, "positive")
	}
	if got := syntax.ClassifyNumber(0); got != "zero" {
		t.Fatalf("ClassifyNumber(0) = %q, want %q", got, "zero")
	}
	if got := syntax.ClassifyNumber(-3); got != "negative" {
		t.Fatalf("ClassifyNumber(-3) = %q, want %q", got, "negative")
	}
}

func TestSumTo(t *testing.T) {
	if got := syntax.SumTo(4); got != 10 {
		t.Fatalf("SumTo(4) = %d, want 10", got)
	}
	if got := syntax.SumTo(1); got != 1 {
		t.Fatalf("SumTo(1) = %d, want 1", got)
	}
	if got := syntax.SumTo(0); got != 0 {
		t.Fatalf("SumTo(0) = %d, want 0", got)
	}
	if got := syntax.SumTo(-2); got != 0 {
		t.Fatalf("SumTo(-2) = %d, want 0", got)
	}
}

func TestDescribe(t *testing.T) {
	if got, want := syntax.Describe("小明", 18), "小明 今年 18 岁"; got != want {
		t.Fatalf("Describe(%q, %d) = %q, want %q", "小明", 18, got, want)
	}
	if got, want := syntax.Describe("Lin", 0), "Lin 今年 0 岁"; got != want {
		t.Fatalf("Describe(%q, %d) = %q, want %q", "Lin", 0, got, want)
	}
}
