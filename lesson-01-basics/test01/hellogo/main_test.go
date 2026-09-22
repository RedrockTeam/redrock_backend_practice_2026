package main

import (
	"go/ast"
	"go/parser"
	"go/token"
	"io"
	"os"
	"testing"
)

const wantOutput = "hello, go\nhello, 红岩网校\n"

// captureStdout 把 main 打印的内容接到管道里读回来。
// 现在不用看懂它，只要知道它拿到的就是程序输出的那几行。
func captureStdout(t *testing.T, run func()) string {
	t.Helper()

	previousStdout := os.Stdout
	reader, writer, err := os.Pipe()
	if err != nil {
		t.Fatalf("创建输出管道失败：%v", err)
	}
	os.Stdout = writer
	defer func() {
		os.Stdout = previousStdout
	}()

	run()

	if err := writer.Close(); err != nil {
		t.Fatalf("关闭输出管道失败：%v", err)
	}
	output, err := io.ReadAll(reader)
	if err != nil {
		t.Fatalf("读取程序输出失败：%v", err)
	}
	if err := reader.Close(); err != nil {
		t.Fatalf("关闭输入管道失败：%v", err)
	}
	return string(output)
}

func TestHelloWorld(t *testing.T) {
	if got := captureStdout(t, main); got != wantOutput {
		t.Fatalf("main() 输出 %q，期望 %q", got, wantOutput)
	}
}

// Println 自动换行、Printf 不自动换行，是这一题真正要记住的东西。
// 光看输出分不出用了哪个，所以这里直接读源码。
func TestUsesBothPrintlnAndPrintf(t *testing.T) {
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, "main.go", nil, 0)
	if err != nil {
		t.Fatalf("解析 main.go 失败：%v", err)
	}

	used := map[string]bool{}
	ast.Inspect(file, func(node ast.Node) bool {
		call, ok := node.(*ast.CallExpr)
		if !ok {
			return true
		}
		selector, ok := call.Fun.(*ast.SelectorExpr)
		if !ok {
			return true
		}
		if pkg, ok := selector.X.(*ast.Ident); ok && pkg.Name == "fmt" {
			used[selector.Sel.Name] = true
		}
		return true
	})

	if !used["Println"] {
		t.Error("第一行要用 fmt.Println：它会自动帮你换行")
	}
	if !used["Printf"] {
		t.Error("第二行要用 fmt.Printf：它不会自动换行，所以格式串末尾得自己写 \\n")
	}
}
