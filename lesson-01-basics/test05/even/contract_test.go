package even

import (
	"go/ast"
	"go/parser"
	"go/token"
	"path/filepath"
	"runtime"
	"strconv"
	"testing"
)

// requiredInputs 是 TestIsEven 至少要覆盖的输入。
// 0 和 -2 正是当前实现漏掉的那一类。
var requiredInputs = []int{2, 3, 0, -2}

// failureReporters 是 testing 包里用来报告失败的方法。
var failureReporters = map[string]bool{
	"Error": true, "Errorf": true, "Fatal": true, "Fatalf": true,
}

// 这个契约测试检查你是不是真的写了测试，
// 而不是只把 student_test.go 里的 t.Fatal 删掉。
func TestStudentWroteIsEvenCases(t *testing.T) {
	body := parseTestIsEvenBody(t)

	covered := make(map[int]bool, len(requiredInputs))
	reportsFailure := false

	ast.Inspect(body, func(node ast.Node) bool {
		call, ok := node.(*ast.CallExpr)
		if !ok {
			return true
		}

		// 记录 IsEven(<整数字面量>) 覆盖了哪些输入。
		if callee, ok := call.Fun.(*ast.Ident); ok && callee.Name == "IsEven" && len(call.Args) == 1 {
			if value, ok := integerLiteral(call.Args[0]); ok {
				covered[value] = true
			}
		}

		// 记录是否调用了 t.Error / t.Errorf / t.Fatal / t.Fatalf。
		if selector, ok := call.Fun.(*ast.SelectorExpr); ok {
			if receiver, ok := selector.X.(*ast.Ident); ok && receiver.Name == "t" {
				if failureReporters[selector.Sel.Name] {
					reportsFailure = true
				}
			}
		}
		return true
	})

	for _, input := range requiredInputs {
		if !covered[input] {
			t.Errorf("TestIsEven 还需要检查 IsEven(%d)", input)
		}
	}
	if !reportsFailure {
		t.Error("TestIsEven 必须在结果不符合预期时调用 t.Error / t.Errorf / t.Fatal / t.Fatalf")
	}
}

// parseTestIsEvenBody 从 student_test.go 中取出 TestIsEven 的函数体。
func parseTestIsEvenBody(t *testing.T) *ast.BlockStmt {
	t.Helper()

	_, currentFile, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("无法定位测试文件")
	}
	studentTest := filepath.Join(filepath.Dir(currentFile), "student_test.go")

	parsed, err := parser.ParseFile(token.NewFileSet(), studentTest, nil, 0)
	if err != nil {
		t.Fatalf("解析 student_test.go 失败：%v", err)
	}

	for _, declaration := range parsed.Decls {
		function, ok := declaration.(*ast.FuncDecl)
		if ok && function.Name.Name == "TestIsEven" && function.Body != nil {
			return function.Body
		}
	}
	t.Fatal("student_test.go 中必须保留函数 TestIsEven")
	return nil
}

// integerLiteral 识别 42 和 -42 这两种写法。
func integerLiteral(expression ast.Expr) (int, bool) {
	switch value := expression.(type) {
	case *ast.BasicLit:
		if value.Kind != token.INT {
			return 0, false
		}
		parsed, err := strconv.Atoi(value.Value)
		if err != nil {
			return 0, false
		}
		return parsed, true
	case *ast.UnaryExpr:
		if value.Op != token.SUB {
			return 0, false
		}
		parsed, ok := integerLiteral(value.X)
		if !ok {
			return 0, false
		}
		return -parsed, true
	}
	return 0, false
}
