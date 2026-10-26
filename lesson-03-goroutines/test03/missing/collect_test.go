package missing_test

import (
	"go/ast"
	"go/parser"
	"go/token"
	"path/filepath"
	"strings"
	"testing"

	"redrock/teaching-exercises/lesson-03-goroutines/test03/missing"
)

func TestCollectWaitsForEveryResult(t *testing.T) {
	const size = 200

	values := make([]int, size)
	for i := range values {
		values[i] = i + 1
	}

	// 多跑几轮，避免偶尔调度恰好把所有 goroutine 都跑完。
	for round := 0; round < 5; round++ {
		got := missing.Collect(values)
		if len(got) != size {
			t.Fatalf("round %d: len(Collect(...)) = %d, want %d", round, len(got), size)
		}
		for i, value := range got {
			want := (i + 1) * (i + 1)
			if value != want {
				t.Fatalf(
					"round %d: 第 %d 个结果 = %d, want %d；Collect 在 goroutine 写完之前就返回了",
					round, i, value, want,
				)
			}
		}
	}
}

func TestCollectHandlesEmpty(t *testing.T) {
	if got := missing.Collect(nil); len(got) != 0 {
		t.Fatalf("Collect(nil) = %#v, want 空切片", got)
	}
}

// TestCollectDoesNotUseSleep 挡住「睡一会儿大概就好了」的修法。
// 这一题要的是建立“写 → 返回”的顺序，第四课会用 WaitGroup 讲。
func TestCollectDoesNotUseSleep(t *testing.T) {
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
		t.Fatal("没有找到可解析的实现文件；请保留 collect.go")
	}
	if foundAt != "" {
		t.Fatalf("%s 使用了 time.Sleep。Sleep 只是让当前函数等一会儿，解决不了“函数已经返回、调用方看不到结果”的问题", foundAt)
	}
}
