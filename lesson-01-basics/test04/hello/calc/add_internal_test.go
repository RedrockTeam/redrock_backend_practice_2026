package calc

import (
	"go/ast"
	"go/parser"
	"go/token"
	"testing"
)

// 这个测试文件声明的是 package calc，和被测代码同包，
// 所以可以直接调用小写的 helper。app 包就做不到这件事。
func TestHelperIsVisibleInsideThePackage(t *testing.T) {
	if got := helper(21); got != 42 {
		t.Fatalf("helper(21) = %d, want 42", got)
	}
	if got := helper(0); got != 0 {
		t.Fatalf("helper(0) = %d, want 0", got)
	}
}

// AddDoubled 必须真的走 helper，否则「包内可见」就没被用上。
func TestAddDoubledGoesThroughHelper(t *testing.T) {
	if got := AddDoubled(3, 4); got != 14 {
		t.Fatalf("AddDoubled(3, 4) = %d, want 14", got)
	}

	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, "add.go", nil, 0)
	if err != nil {
		t.Fatalf("解析 add.go 失败：%v", err)
	}

	var body *ast.FuncDecl
	for _, decl := range file.Decls {
		if fn, ok := decl.(*ast.FuncDecl); ok && fn.Name.Name == "AddDoubled" {
			body = fn
		}
	}
	if body == nil {
		t.Fatal("add.go 里没有找到 AddDoubled")
	}

	called := false
	ast.Inspect(body, func(node ast.Node) bool {
		call, ok := node.(*ast.CallExpr)
		if !ok {
			return true
		}
		if ident, ok := call.Fun.(*ast.Ident); ok && ident.Name == "helper" {
			called = true
		}
		return true
	})
	if !called {
		t.Fatal("AddDoubled 没有调用 helper：这一题要的就是「同包的小写函数可以直接用」")
	}
}
