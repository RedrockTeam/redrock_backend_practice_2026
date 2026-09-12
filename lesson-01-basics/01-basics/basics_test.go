package basics

import "testing"

func TestGreeting(t *testing.T) {
	if got := Greeting("Lin"); got != "Hello, Lin!" {
		t.Fatalf("Greeting() = %q", got)
	}
	if got := Greeting(""); got != "Hello, learner!" {
		t.Fatalf("empty Greeting() = %q", got)
	}
}

func TestAverage(t *testing.T) {
	if got := Average([]int{80, 90, 100}); got != 90 {
		t.Fatalf("Average() = %v", got)
	}
	if got := Average(nil); got != 0 {
		t.Fatalf("empty Average() = %v", got)
	}
}
