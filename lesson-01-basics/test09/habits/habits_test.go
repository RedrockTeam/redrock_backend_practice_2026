package habits_test

import (
	"testing"

	"redrock/teaching-exercises/lesson-01-basics/test09/habits"
)

func TestIsWeekend(t *testing.T) {
	want := map[int]bool{1: false, 2: false, 3: false, 4: false, 5: false, 6: true, 7: true}
	for day := 1; day <= 7; day++ {
		if got := habits.IsWeekend(day); got != want[day] {
			t.Errorf("IsWeekend(%d) = %v, want %v；Go 的 case 执行完就结束，不会落到下一个 case", day, got, want[day])
		}
	}
}

func TestAverageKeepsFraction(t *testing.T) {
	cases := []struct {
		scores []int
		want   float64
	}{
		{scores: []int{90, 85}, want: 87.5},
		{scores: []int{1, 2}, want: 1.5},
		{scores: []int{60}, want: 60},
		{scores: nil, want: 0},
	}
	for _, c := range cases {
		if got := habits.Average(c.scores); got != c.want {
			t.Errorf("Average(%v) = %v, want %v；想一想除法是在什么类型上做的", c.scores, got, c.want)
		}
	}
}

func TestCountCharsCountsCharactersNotBytes(t *testing.T) {
	cases := map[string]int{
		"":        0,
		"hello":   5,
		"红岩网校":    4,
		"Go语言":    4,
		"后端 2026": 7,
	}
	for input, want := range cases {
		if got := habits.CountChars(input); got != want {
			t.Errorf("CountChars(%q) = %d, want %d；len 数的是字节，一个汉字在 UTF-8 里占 3 个字节", input, got, want)
		}
	}
}

func TestReverseKeepsCharactersIntact(t *testing.T) {
	cases := map[string]string{
		"":     "",
		"abc":  "cba",
		"红岩":   "岩红",
		"Go语言": "言语oG",
	}
	for input, want := range cases {
		if got := habits.Reverse(input); got != want {
			t.Errorf("Reverse(%q) = %q, want %q；按字节倒序会把一个汉字的 3 个字节拆散", input, got, want)
		}
	}
}

func TestSumScoresAddsValidInput(t *testing.T) {
	got, err := habits.SumScores([]string{"90", "85", "-5"})
	if err != nil || got != 170 {
		t.Fatalf(`SumScores(["90" "85" "-5"]) = %d, %v, want 170, nil`, got, err)
	}
	if got, err := habits.SumScores(nil); err != nil || got != 0 {
		t.Fatalf("SumScores(nil) = %d, %v, want 0, nil", got, err)
	}
}

func TestSumScoresReportsBadInput(t *testing.T) {
	got, err := habits.SumScores([]string{"90", "八十", "70"})
	if err == nil {
		t.Fatalf(`SumScores(["90" "八十" "70"]) 的 error = nil；解析失败的错误没有被带出循环，检查 := 在哪里新建了变量`)
	}
	if got != 0 {
		t.Fatalf(`SumScores(["90" "八十" "70"]) = %d, want 0`, got)
	}
}
