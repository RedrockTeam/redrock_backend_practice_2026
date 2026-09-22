package main

import (
	"go/ast"
	"go/parser"
	"go/scanner"
	"go/token"
	"io"
	"os"
	"testing"
)

const sourceFile = "translate.go"

const wantOutput = `1 odd 1
2 even 4
3 odd 9
4 even 16
5 odd 25
total 55
`

func readSource(t *testing.T) []byte {
	t.Helper()
	source, err := os.ReadFile(sourceFile)
	if err != nil {
		t.Fatalf("读取 %s 失败：%v", sourceFile, err)
	}
	return source
}

func parseSource(t *testing.T, fset *token.FileSet, source []byte) *ast.File {
	t.Helper()
	file, err := parser.ParseFile(fset, sourceFile, source, 0)
	if err != nil {
		t.Fatalf("解析 %s 失败：%v", sourceFile, err)
	}
	return file
}

func findFunc(file *ast.File, name string) *ast.FuncDecl {
	for _, decl := range file.Decls {
		if fn, ok := decl.(*ast.FuncDecl); ok && fn.Recv == nil && fn.Name.Name == name {
			return fn
		}
	}
	return nil
}

// TestTranslationOutput 检查改写后的程序输出和 C 程序逐字一致。
func TestTranslationOutput(t *testing.T) {
	previousStdout := os.Stdout
	reader, writer, err := os.Pipe()
	if err != nil {
		t.Fatalf("创建输出管道失败：%v", err)
	}
	os.Stdout = writer
	defer func() {
		os.Stdout = previousStdout
	}()

	main()

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

	if got := string(output); got != wantOutput {
		t.Fatalf("程序输出与 reference.c 不一致。\n得到：\n%s\n期望：\n%s", got, wantOutput)
	}
}

// TestTranslationKeepsSquare 检查 square 被真的翻译成了 Go 函数并被 main 调用，
// 而不是把结果直接算在 main 里。对应清单第 9 步。
func TestTranslationKeepsSquare(t *testing.T) {
	fset := token.NewFileSet()
	file := parseSource(t, fset, readSource(t))

	square := findFunc(file, "square")
	if square == nil {
		t.Fatal("没有找到函数 square：C 里的 int square(int n) 也要翻译过来，不要把 n*n 直接写进 main")
	}
	if got := square.Type.Params.NumFields(); got != 1 {
		t.Fatalf("square 有 %d 个参数，期望 1 个（C 里是 int square(int n)）", got)
	}
	if got := square.Type.Results.NumFields(); got != 1 {
		t.Fatalf("square 有 %d 个返回值，期望 1 个；Go 的返回类型写在参数列表后面", got)
	}

	mainFunc := findFunc(file, "main")
	if mainFunc == nil {
		t.Fatal("没有找到 func main：Go 程序的入口没有参数也没有返回值（清单第 3、4 步）")
	}
	called := false
	ast.Inspect(mainFunc, func(node ast.Node) bool {
		call, ok := node.(*ast.CallExpr)
		if !ok {
			return true
		}
		if ident, ok := call.Fun.(*ast.Ident); ok && ident.Name == "square" {
			called = true
		}
		return true
	})
	if !called {
		t.Fatal("main 没有调用 square：翻译要保持原来的结构，不要把平方展开成字面量")
	}
}

// TestTranslationKeepsControlFlow 检查 for 和 if/else 都被保留下来了，
// 而不是把六行输出硬编码成六条 Println。对应清单第 7 步与第二部分 2.5。
func TestTranslationKeepsControlFlow(t *testing.T) {
	fset := token.NewFileSet()
	file := parseSource(t, fset, readSource(t))

	mainFunc := findFunc(file, "main")
	if mainFunc == nil {
		t.Fatal("没有找到 func main")
	}

	var loops, elseBranches int
	ast.Inspect(mainFunc, func(node ast.Node) bool {
		switch stmt := node.(type) {
		case *ast.ForStmt, *ast.RangeStmt:
			loops++
		case *ast.IfStmt:
			if stmt.Else != nil {
				elseBranches++
			}
		}
		return true
	})

	if loops == 0 {
		t.Error("main 里没有循环：C 的 for 要翻译成 Go 的 for，不能把六行输出一条条写死")
	}
	if elseBranches == 0 {
		t.Error("main 里没有带 else 的 if：C 的 if/else 要一起翻译过来")
	}
}

// TestTranslationDropsSemicolons 检查 C 风格的行尾分号都删掉了。
// for 三段式头部里的分号是 Go 语法要求的，不算。对应清单第 6 步与第二部分 2.2。
func TestTranslationDropsSemicolons(t *testing.T) {
	source := readSource(t)

	parseFset := token.NewFileSet()
	file := parseSource(t, parseFset, source)

	type span struct{ from, to int }
	var headers []span
	ast.Inspect(file, func(node ast.Node) bool {
		if loop, ok := node.(*ast.ForStmt); ok && loop.Body != nil {
			headers = append(headers, span{
				from: parseFset.Position(loop.For).Offset,
				to:   parseFset.Position(loop.Body.Lbrace).Offset,
			})
		}
		return true
	})

	scanFset := token.NewFileSet()
	target := scanFset.AddFile(sourceFile, scanFset.Base(), len(source))
	var s scanner.Scanner
	s.Init(target, source, nil, 0)

	for {
		pos, tok, lit := s.Scan()
		if tok == token.EOF {
			break
		}
		// 编译器自动插入的分号，字面量是 "\n"；手写的才是 ";"。
		if tok != token.SEMICOLON || lit != ";" {
			continue
		}
		offset := scanFset.Position(pos).Offset
		inHeader := false
		for _, header := range headers {
			if offset > header.from && offset < header.to {
				inHeader = true
				break
			}
		}
		if !inHeader {
			t.Errorf("%s:%d 还留着 C 风格的行尾分号。Go 会自动插入分号，手写的要删掉（清单第 6 步）",
				sourceFile, scanFset.Position(pos).Line)
		}
	}
}
