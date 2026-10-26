package runall_test

import (
	"go/ast"
	"go/parser"
	"go/token"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"redrock/teaching-exercises/lesson-04-sync/test01/runall"
)

// TestRunAllWaitsForEveryTask 用“任务完成得比函数返回慢”的方式，
// 稳定地暴露“没有等待 goroutine”这个问题。
func TestRunAllWaitsForEveryTask(t *testing.T) {
	const (
		size   = 64
		rounds = 20
	)

	for round := 0; round < rounds; round++ {
		tasks := make([]func() int, size)
		for i := 0; i < size; i++ {
			i := i
			tasks[i] = func() int {
				time.Sleep(time.Duration(i%4) * time.Millisecond)
				return i*i + 1
			}
		}

		got := runall.RunAll(tasks)
		if len(got) != size {
			t.Fatalf("round %d: len(RunAll(...)) = %d, want %d", round, len(got), size)
		}
		for i, value := range got {
			want := i*i + 1
			if value != want {
				t.Fatalf(
					"round %d: 第 %d 个结果 = %d, want %d；RunAll 是不是没等 goroutine 结束就返回了？",
					round, i, value, want,
				)
			}
		}
	}
}

// TestRunAllPreservesOrder 检查结果按下标归位，而不是按完成先后 append。
func TestRunAllPreservesOrder(t *testing.T) {
	const (
		size   = 32
		rounds = 50
	)

	for round := 0; round < rounds; round++ {
		tasks := make([]func() int, size)
		for i := range tasks {
			i := i
			tasks[i] = func() int { return 1000 + i }
		}

		got := runall.RunAll(tasks)
		if len(got) != size {
			t.Fatalf("round %d: len(RunAll(...)) = %d, want %d", round, len(got), size)
		}
		for i, value := range got {
			if value != 1000+i {
				t.Fatalf(
					"round %d: 第 %d 个结果 = %d, want %d；结果必须按输入下标写回，不能按完成顺序追加",
					round, i, value, 1000+i,
				)
			}
		}
	}
}

func TestRunAllHandlesEmpty(t *testing.T) {
	if got := runall.RunAll(nil); len(got) != 0 {
		t.Fatalf("RunAll(nil) 的长度 = %d, want 0", len(got))
	}
	if got := runall.RunAll([]func() int{}); len(got) != 0 {
		t.Fatalf("RunAll(空切片) 的长度 = %d, want 0", len(got))
	}
}

// TestRunAllDoesNotUseSleep 是一道源码契约测试：
// 这个练习要建立的是“谁等谁”的明确顺序，不是“睡一会儿大概就好了”。
func TestRunAllDoesNotUseSleep(t *testing.T) {
	files, err := filepath.Glob("*.go")
	if err != nil {
		t.Fatalf("列目录失败: %v", err)
	}

	fset := token.NewFileSet()
	foundAt := ""
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

		timeNames := map[string]bool{"time": true}
		for _, imp := range file.Imports {
			path := strings.Trim(imp.Path.Value, `"`)
			if path != "time" {
				continue
			}
			if imp.Name != nil {
				timeNames[imp.Name.Name] = true
			}
		}

		ast.Inspect(file, func(node ast.Node) bool {
			call, ok := node.(*ast.CallExpr)
			if !ok {
				return true
			}
			selector, ok := call.Fun.(*ast.SelectorExpr)
			if !ok {
				return true
			}
			ident, ok := selector.X.(*ast.Ident)
			if !ok {
				return true
			}
			if timeNames[ident.Name] && selector.Sel.Name == "Sleep" {
				foundAt = fset.Position(call.Pos()).String()
			}
			return true
		})
	}

	if checked == 0 {
		t.Fatal("没有找到可解析的实现文件；请保留 runall.go")
	}
	if foundAt != "" {
		t.Fatalf(
			"%s 调用了 time.Sleep。这里要建立明确的等待关系：用 WaitGroup 或 channel，让 RunAll 在最后一个任务结束时才返回。",
			foundAt,
		)
	}
}
