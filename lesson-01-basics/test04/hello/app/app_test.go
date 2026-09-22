package app_test

import (
	"go/parser"
	"go/token"
	"strconv"
	"testing"

	"redrock/teaching-exercises/lesson-01-basics/test04/hello/app"
)

const calcImportPath = "redrock/teaching-exercises/lesson-01-basics/test04/hello/calc"

func TestSum(t *testing.T) {
	if got := app.Sum(1, 2); got != 3 {
		t.Fatalf("Sum(1, 2) = %d, want 3", got)
	}
	if got := app.Sum(-5, 5); got != 0 {
		t.Fatalf("Sum(-5, 5) = %d, want 0", got)
	}
}

func TestDescribe(t *testing.T) {
	if got, want := app.Describe(1, 2), "1 + 2 = 3"; got != want {
		t.Fatalf("Describe(1, 2) = %q, want %q", got, want)
	}
	if got, want := app.Describe(-1, 1), "-1 + 1 = 0"; got != want {
		t.Fatalf("Describe(-1, 1) = %q, want %q", got, want)
	}
}

// 这一题的重点：导入路径是「module 路径 + 子目录」，不是相对路径，
// 也不是包名。写错了整个包根本编译不过，所以这里顺便把正确答案的形状讲清楚。
func TestAppImportsCalcByItsFullPath(t *testing.T) {
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, "app.go", nil, parser.ImportsOnly)
	if err != nil {
		t.Fatalf("解析 app.go 失败：%v", err)
	}

	for _, spec := range file.Imports {
		path, err := strconv.Unquote(spec.Path.Value)
		if err != nil {
			continue
		}
		if path == calcImportPath {
			return
		}
	}
	t.Fatalf("app.go 没有导入 calc 包。正确的导入路径是 %q：\n"+
		"module 路径 %q 来自仓库根目录的 go.mod，后面接上 calc 包所在的子目录。",
		calcImportPath, "redrock/teaching-exercises")
}
