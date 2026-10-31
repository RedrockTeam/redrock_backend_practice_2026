package parallel_test

import (
	"go/ast"
	"go/parser"
	"go/token"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"redrock/teaching-exercises/lesson-03-goroutines/test06/parallel"
)

func TestSum(t *testing.T) {
	large := make([]int, 1000)
	for i := range large {
		large[i] = i + 1
	}

	cases := []struct {
		values []int
		want   int
	}{
		{values: []int{1, 2, 3, 4}, want: 10},
		{values: []int{-5, 5, 10}, want: 10},
		{values: []int{7}, want: 7},
		{values: nil, want: 0},
		{values: []int{}, want: 0},
		{values: large, want: 500500},
	}

	for _, c := range cases {
		if got := parallel.Sum(c.values); got != c.want {
			t.Errorf("Sum(%d 个元素) = %d, want %d", len(c.values), got, c.want)
		}
	}
}

func TestSumDoesNotModifyInput(t *testing.T) {
	values := []int{1, 2, 3, 4}
	before := append([]int(nil), values...)

	parallel.Sum(values)

	if !reflect.DeepEqual(values, before) {
		t.Fatalf("Sum 修改了输入：%v -> %v", before, values)
	}
}

// TestSumUsesGoroutines 是源码契约：顺序求和虽然结果正确，
// 但这一题要练的是把一段工作拆给多个 goroutine。
func TestSumUsesGoroutines(t *testing.T) {
	files, err := filepath.Glob("*.go")
	if err != nil {
		t.Fatalf("列目录失败: %v", err)
	}

	fset := token.NewFileSet()
	found := false
	checked := 0

	for _, name := range files {
		if strings.HasSuffix(name, "_test.go") {
			continue
		}
		checked++
		file, err := parser.ParseFile(fset, name, nil, 0)
		if err != nil {
			t.Fatalf("解析 %s 失败: %v", name, err)
		}
		ast.Inspect(file, func(node ast.Node) bool {
			if _, ok := node.(*ast.GoStmt); ok {
				found = true
			}
			return true
		})
	}

	if checked == 0 {
		t.Fatal("没有找到可解析的实现文件；请保留 sum.go")
	}
	if !found {
		t.Fatal("sum.go 里没有 go 语句。这一题要把 values 拆给 goroutine，而不是一个 for 循环从头加到尾")
	}
}
