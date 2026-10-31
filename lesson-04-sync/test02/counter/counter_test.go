package counter_test

import (
	"go/ast"
	"go/parser"
	"go/token"
	"path/filepath"
	"strings"
	"sync"
	"testing"

	"redrock/teaching-exercises/lesson-04-sync/test02/counter"
)

// adder 只是测试用来把两种实现放进同一段代码里的接口——第二课 interface 的用法。
type adder interface {
	Add(delta int)
	Value() int
	Reset()
}

func implementations() map[string]func() adder {
	return map[string]func() adder{
		"Counter":       func() adder { return &counter.Counter{} },
		"AtomicCounter": func() adder { return &counter.AtomicCounter{} },
	}
}

func TestZeroValueAndSequential(t *testing.T) {
	for name, build := range implementations() {
		c := build()
		if got := c.Value(); got != 0 {
			t.Errorf("零值 %s.Value() = %d, want 0", name, got)
		}
		c.Add(5)
		c.Add(-2)
		c.Add(10)
		if got := c.Value(); got != 13 {
			t.Errorf("%s: Add(5), Add(-2), Add(10) 之后 Value() = %d, want 13", name, got)
		}
		c.Reset()
		if got := c.Value(); got != 0 {
			t.Errorf("%s: Reset() 之后 Value() = %d, want 0", name, got)
		}
		c.Add(1)
		if got := c.Value(); got != 1 {
			t.Errorf("%s: Reset 之后再 Add(1)，Value() = %d, want 1", name, got)
		}
	}
}

func TestConcurrentAddDoesNotLoseUpdates(t *testing.T) {
	const (
		goroutines = 100
		perG       = 100
		rounds     = 5
	)
	for name, build := range implementations() {
		for round := 0; round < rounds; round++ {
			c := build()
			var wg sync.WaitGroup
			wg.Add(goroutines)
			for g := 0; g < goroutines; g++ {
				go func() {
					defer wg.Done()
					for i := 0; i < perG; i++ {
						c.Add(1)
					}
				}()
			}
			wg.Wait()
			if got := c.Value(); got != goroutines*perG {
				t.Fatalf("%s round %d: 100 个 goroutine 各 Add(1) 100 次，Value() = %d, want %d；"+
					"c.n += delta 是读-改-写三步，和第三课 test03 里交错丢更新是同一件事",
					name, round, got, goroutines*perG)
			}
		}
	}
}

func TestConcurrentMixedOperations(t *testing.T) {
	for name, build := range implementations() {
		c := build()
		var wg sync.WaitGroup
		// 一组每次 +3、一组每次 -1，同时还有读者在读：最终结果必须是 50×200×(3-1) = 20000。
		for g := 0; g < 50; g++ {
			wg.Add(3)
			go func() {
				defer wg.Done()
				for i := 0; i < 200; i++ {
					c.Add(3)
				}
			}()
			go func() {
				defer wg.Done()
				for i := 0; i < 200; i++ {
					c.Add(-1)
				}
			}()
			go func() {
				defer wg.Done()
				for i := 0; i < 50; i++ {
					_ = c.Value()
				}
			}()
		}
		wg.Wait()
		if got := c.Value(); got != 20000 {
			t.Fatalf("%s: 并发 +3 与 -1 各 10000 次后 Value() = %d, want 20000；读者也在同时读，Value 也要同步", name, got)
		}
	}
}

// TestCounterUsesMutexAndAtomicCounterUsesAtomic 是源码契约：
// 功能测试分不出你是用锁还是用原子操作，所以直接读 counter.go。
func TestCounterUsesMutexAndAtomicCounterUsesAtomic(t *testing.T) {
	files, err := filepath.Glob("*.go")
	if err != nil {
		t.Fatalf("列目录失败: %v", err)
	}

	fset := token.NewFileSet()
	structs := map[string]*ast.StructType{}
	syncNames := map[string]bool{}
	atomicNames := map[string]bool{}
	atomicCallIn := map[string]bool{} // 接收者类型名 -> 方法体里调用过 atomic.Xxx

	var parsed []*ast.File
	for _, name := range files {
		if strings.HasSuffix(name, "_test.go") {
			continue
		}
		file, err := parser.ParseFile(fset, name, nil, 0)
		if err != nil {
			t.Fatalf("解析 %s 失败: %v", name, err)
		}
		parsed = append(parsed, file)
	}
	if len(parsed) == 0 {
		t.Fatal("没有找到可解析的实现文件；请保留 counter.go")
	}

	for _, file := range parsed {
		for _, imp := range file.Imports {
			path := strings.Trim(imp.Path.Value, `"`)
			local := path[strings.LastIndex(path, "/")+1:]
			if imp.Name != nil {
				local = imp.Name.Name
			}
			switch path {
			case "sync":
				syncNames[local] = true
			case "sync/atomic":
				atomicNames[local] = true
			}
		}
		ast.Inspect(file, func(node ast.Node) bool {
			if spec, ok := node.(*ast.TypeSpec); ok {
				if st, ok := spec.Type.(*ast.StructType); ok {
					structs[spec.Name.Name] = st
				}
			}
			return true
		})
	}

	for _, file := range parsed {
		for _, decl := range file.Decls {
			fn, ok := decl.(*ast.FuncDecl)
			if !ok || fn.Recv == nil || fn.Body == nil {
				continue
			}
			recv := receiverName(fn.Recv.List[0].Type)
			ast.Inspect(fn.Body, func(node ast.Node) bool {
				if sel, ok := node.(*ast.SelectorExpr); ok {
					if ident, ok := sel.X.(*ast.Ident); ok && atomicNames[ident.Name] {
						atomicCallIn[recv] = true
					}
				}
				return true
			})
		}
	}

	fieldFrom := func(st *ast.StructType, pkgs map[string]bool, names ...string) bool {
		if st == nil {
			return false
		}
		for _, field := range st.Fields.List {
			sel, ok := field.Type.(*ast.SelectorExpr)
			if !ok {
				continue
			}
			ident, ok := sel.X.(*ast.Ident)
			if !ok || !pkgs[ident.Name] {
				continue
			}
			if len(names) == 0 {
				return true
			}
			for _, n := range names {
				if sel.Sel.Name == n {
					return true
				}
			}
		}
		return false
	}

	if !fieldFrom(structs["Counter"], syncNames, "Mutex", "RWMutex") {
		t.Error("Counter 结构体里没有 sync.Mutex / sync.RWMutex 字段；这一版要用锁保护计数")
	}
	if fieldFrom(structs["AtomicCounter"], syncNames, "Mutex", "RWMutex") {
		t.Error("AtomicCounter 里有锁字段；这一版要求不用锁，只用 sync/atomic")
	}
	if !fieldFrom(structs["AtomicCounter"], atomicNames) && !atomicCallIn["AtomicCounter"] {
		t.Error("AtomicCounter 既没有 atomic.Int64 这类字段，方法里也没有调用 sync/atomic；这一版要用原子操作")
	}
}

func receiverName(expr ast.Expr) string {
	if star, ok := expr.(*ast.StarExpr); ok {
		expr = star.X
	}
	if ident, ok := expr.(*ast.Ident); ok {
		return ident.Name
	}
	return ""
}
