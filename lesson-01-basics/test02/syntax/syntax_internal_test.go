package syntax

import (
	"go/ast"
	"go/parser"
	"go/token"
	"testing"
)

// 这个测试文件声明的是 package syntax，和被测代码同属一个包，
// 因此可以直接调用小写的 add。外部测试文件做不到这一点。
func TestPackagePrivateAdd(t *testing.T) {
	if got := add(5, 6); got != 11 {
		t.Fatalf("add(5, 6) = %d, want 11", got)
	}
}

// minScore 是小写常量，只有同包的测试文件读得到。
func TestPackagePrivateMinScore(t *testing.T) {
	if minScore != 0 {
		t.Fatalf("minScore = %d, want 0", minScore)
	}
}

// 这一题考的是「零值」，所以不能写 return 0 蒙混过去：
// 必须真的声明一个没有赋值的变量，让 Go 给它零值。
func TestInitialCountReliesOnTheZeroValue(t *testing.T) {
	if got := InitialCount(); got != 0 {
		t.Fatalf("InitialCount() = %d, want 0", got)
	}

	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, "syntax.go", nil, 0)
	if err != nil {
		t.Fatalf("解析 syntax.go 失败：%v", err)
	}

	var fn *ast.FuncDecl
	for _, decl := range file.Decls {
		if candidate, ok := decl.(*ast.FuncDecl); ok && candidate.Name.Name == "InitialCount" {
			fn = candidate
		}
	}
	if fn == nil {
		t.Fatal("syntax.go 里没有找到 InitialCount")
	}

	declaredWithoutValue := false
	ast.Inspect(fn, func(node ast.Node) bool {
		spec, ok := node.(*ast.ValueSpec)
		if ok && len(spec.Values) == 0 && spec.Type != nil {
			declaredWithoutValue = true
		}
		return true
	})
	if !declaredWithoutValue {
		t.Fatal("InitialCount 里要写 `var count int` 这样不赋值的声明，再直接返回它——这一题考的就是零值")
	}
}
