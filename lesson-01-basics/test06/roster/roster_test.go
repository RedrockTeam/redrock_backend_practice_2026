package roster_test

import (
	"fmt"
	"strings"
	"testing"

	"redrock/teaching-exercises/lesson-01-basics/test06/roster"
)

// sample 造一份有同名学生的小名册。
func sample() *roster.Roster {
	r := roster.New()
	r.Add(roster.Student{ID: "S001", Name: " Lin ", Score: 90})
	r.Add(roster.Student{ID: "S002", Name: "Wang", Score: 80})
	r.Add(roster.Student{ID: "S003", Name: "Lin", Score: 70})
	return r
}

// large 造一份 size 条记录的名册，学号是 S001、S002 ... 依次递增。
func large(size int) *roster.Roster {
	r := roster.New()
	for index := 1; index <= size; index++ {
		r.Add(roster.Student{
			ID:    fmt.Sprintf("S%03d", index),
			Name:  fmt.Sprintf("student-%d", index),
			Score: index % 101,
		})
	}
	return r
}

func TestAddAndLen(t *testing.T) {
	r := roster.New()
	if got := r.Len(); got != 0 {
		t.Fatalf("New().Len() = %d, want 0", got)
	}

	r.Add(roster.Student{ID: "S001", Name: "Lin", Score: 90})
	r.Add(roster.Student{ID: "S002", Name: "Wang", Score: 80})
	if got := r.Len(); got != 2 {
		t.Fatalf("Len() after two Add() = %d, want 2", got)
	}
}

func TestAddTrimsName(t *testing.T) {
	student, ok := sample().FindByID("S001")
	if !ok {
		t.Fatal(`FindByID("S001") 没找到记录`)
	}
	if student.Name != "Lin" {
		t.Fatalf("Add() 存下的 Name = %q, want %q（两端空白要清掉）", student.Name, "Lin")
	}
}

func TestFindByID(t *testing.T) {
	r := sample()

	student, ok := r.FindByID("S002")
	if !ok {
		t.Fatal(`FindByID("S002") = _, false, want _, true`)
	}
	if student.Name != "Wang" || student.Score != 80 {
		t.Fatalf("FindByID(\"S002\") = %+v, want {S002 Wang 80}", student)
	}

	missing, ok := r.FindByID("S999")
	if ok {
		t.Fatalf(`FindByID("S999") = %+v, true, want zero value, false`, missing)
	}
	if missing != (roster.Student{}) {
		t.Fatalf(`FindByID("S999") 第一个返回值 = %+v, want 零值`, missing)
	}
}

func TestTotalScoreByName(t *testing.T) {
	r := sample()

	// 名册里有两个 Lin：90 + 70。
	if got := r.TotalScoreByName("Lin"); got != 160 {
		t.Fatalf(`TotalScoreByName("Lin") = %d, want 160`, got)
	}
	if got := r.TotalScoreByName("Wang"); got != 80 {
		t.Fatalf(`TotalScoreByName("Wang") = %d, want 80`, got)
	}
	if got := r.TotalScoreByName("None"); got != 0 {
		t.Fatalf(`TotalScoreByName("None") = %d, want 0`, got)
	}
}

func TestFormatText(t *testing.T) {
	want := "S001 Lin 90\nS002 Wang 80\nS003 Lin 70\n"
	if got := roster.FormatText(sample()); got != want {
		t.Fatalf("FormatText() =\n%q\nwant\n%q", got, want)
	}
}

func TestFormatCSV(t *testing.T) {
	want := "id,name,score\nS001,Lin,90\nS002,Wang,80\nS003,Lin,70\n"
	if got := roster.FormatCSV(sample()); got != want {
		t.Fatalf("FormatCSV() =\n%q\nwant\n%q", got, want)
	}
}

func TestFormatEmptyRoster(t *testing.T) {
	if got := roster.FormatText(roster.New()); got != "" {
		t.Fatalf("空名册的 FormatText() = %q, want %q", got, "")
	}
	if got, want := roster.FormatCSV(roster.New()), "id,name,score\n"; got != want {
		t.Fatalf("空名册的 FormatCSV() = %q, want %q（表头总是要写）", got, want)
	}
}

// ---------------------------------------------------------------------------
// 下面三个 TestBridge 开头的测试，是第一课到第二课的接口。
// 它们检查的都是正确行为，但同时把「只用切片」的代价量出来，
// 第二课会用 map 和 interface 把这些数字和重复代码消掉。
// ---------------------------------------------------------------------------

// TestBridgeLookupCostGrowsWithRosterSize 量出线性查找的代价。
func TestBridgeLookupCostGrowsWithRosterSize(t *testing.T) {
	const size = 200
	r := large(size)

	measure := func(id string) int {
		before := r.Comparisons()
		r.FindByID(id)
		return r.Comparisons() - before
	}

	first := measure("S001")
	last := measure(fmt.Sprintf("S%03d", size))
	missing := measure("NOT-A-REAL-ID")

	if first != 1 {
		t.Errorf("查第 1 条记录用了 %d 次比较, want 1", first)
	}
	if last != size {
		t.Errorf("查第 %d 条记录用了 %d 次比较, want %d", size, last, size)
	}
	if missing != size {
		t.Errorf("查一个不存在的学号用了 %d 次比较, want %d（必须把整份名册看完）", missing, size)
	}

	t.Logf("名册 %d 条：命中首条 %d 次比较，命中末条 %d 次，查不到 %d 次。", size, first, last, missing)
	t.Log("这就是线性查找：最坏情况和名册长度成正比。第二课的 map 会把它变成一次哈希。")
}

// TestBridgeAggregationScansWholeRoster 量出「按名字统计」的代价。
func TestBridgeAggregationScansWholeRoster(t *testing.T) {
	const size = 200
	r := large(size)

	before := r.Comparisons()
	r.TotalScoreByName("student-7")
	first := r.Comparisons() - before

	before = r.Comparisons()
	r.TotalScoreByName("student-7")
	second := r.Comparisons() - before

	if first != size {
		t.Errorf("统计一个名字用了 %d 次比较, want %d", first, size)
	}
	if second != first {
		t.Errorf("同一个名字统计两次的代价 %d != %d；切片实现没有任何记忆", second, first)
	}

	t.Logf("每问一次「这个名字总分多少」就要重扫 %d 条，问两次就是 %d 次比较。", size, first+second)
	t.Log("第二课会先用 map[string]int 把统计结果攒起来，一次遍历回答所有名字。")
}

// TestBridgeTwoFormattersShareOneSignature 说明两个格式化函数已经可以被统一对待。
func TestBridgeTwoFormattersShareOneSignature(t *testing.T) {
	// 之所以能把它们塞进同一个切片，是因为两者签名完全一样：func(*roster.Roster) string。
	formatters := []struct {
		name       string
		format     func(*roster.Roster) string
		wantLines  int
		wantHeader bool
	}{
		{name: "text", format: roster.FormatText, wantLines: 3, wantHeader: false},
		{name: "csv", format: roster.FormatCSV, wantLines: 4, wantHeader: true},
	}

	for _, formatter := range formatters {
		output := formatter.format(sample())

		if output == "" {
			t.Errorf("%s: 输出为空", formatter.name)
			continue
		}
		if !strings.HasSuffix(output, "\n") {
			t.Errorf("%s: 输出 %q 的最后一行没有以换行结尾", formatter.name, output)
		}

		lines := strings.Split(strings.TrimSuffix(output, "\n"), "\n")
		if len(lines) != formatter.wantLines {
			t.Errorf("%s: 输出 %d 行, want %d 行", formatter.name, len(lines), formatter.wantLines)
		}
		if formatter.wantHeader && lines[0] != "id,name,score" {
			t.Errorf("%s: 第一行 = %q, want 表头", formatter.name, lines[0])
		}
	}

	t.Log("两个函数各自写了一遍「遍历名册、拼接每行」的循环，只有每行长什么样不同。")
	t.Log("第二课的 interface 会把「怎么遍历」留在一处，把「每行长什么样」交给不同实现。")
}
