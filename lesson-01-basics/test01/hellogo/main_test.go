package main

import (
	"io"
	"os"
	"testing"
)

func TestHelloWorld(t *testing.T) {
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

	if got, want := string(output), "hello, go\n"; got != want {
		t.Fatalf("main() 输出 %q，期望 %q", got, want)
	}
}
