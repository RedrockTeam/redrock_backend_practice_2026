package asyncprint_test

import (
	"bytes"
	"go/ast"
	"go/parser"
	"go/token"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"testing"

	"redrock/teaching-exercises/lesson-03-goroutines/test05/asyncprint"
)

// lockedBuffer 让多个 goroutine 同时写也不会破坏 bytes.Buffer。
// 它只是测试用来收集输出的容器，不是本题要你实现的同步。
type lockedBuffer struct {
	mu sync.Mutex
	b  bytes.Buffer
}

func (w *lockedBuffer) Write(p []byte) (int, error) {
	w.mu.Lock()
	defer w.mu.Unlock()
	return w.b.Write(p)
}

func (w *lockedBuffer) String() string {
	w.mu.Lock()
	defer w.mu.Unlock()
	return w.b.String()
}

func TestRunStartsEveryWorker(t *testing.T) {
	var out lockedBuffer
	asyncprint.Run(&out, 5)

	lines := strings.Split(strings.TrimSpace(out.String()), "\n")
	if len(lines) == 1 && lines[0] == "" {
		t.Fatal("Run 没有输出任何 worker；goroutine 是否还没跑，函数就返回了？")
	}
	sort.Strings(lines)

	want := []string{"worker 0", "worker 1", "worker 2", "worker 3", "worker 4"}
	if len(lines) != len(want) {
		t.Fatalf("Run 输出了 %d 行 %v, want %d 行；每个 worker 恰好一行", len(lines), lines, len(want))
	}
	for i := range want {
		if lines[i] != want[i] {
			t.Fatalf("Run 的输出 = %v, want %v", lines, want)
		}
	}
}

func TestRunHandlesZeroWorkers(t *testing.T) {
	var out lockedBuffer
	asyncprint.Run(&out, 0)
	if got := strings.TrimSpace(out.String()); got != "" {
		t.Fatalf("workers = 0 时输出 %q, want 空", got)
	}
}

// TestRunUsesGoroutines 是源码契约：这一题练的是 go 关键字，
// 顺序 for 循环打印虽然输出正确，但没有建立任何 goroutine。
func TestRunUsesGoroutines(t *testing.T) {
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
		t.Fatal("没有找到可解析的实现文件；请保留 run.go")
	}
	if !found {
		t.Fatal("run.go 里没有 go 语句。这一题要把每个 worker 放进独立 goroutine，而不是顺序打印")
	}
}
