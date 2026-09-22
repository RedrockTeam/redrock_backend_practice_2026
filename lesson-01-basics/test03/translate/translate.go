package main

import "fmt"

// TODO（清单第 9 步）：把 reference.c 里的
//
//	int square(int n) { return n * n; }
//
// 改写成 Go 的函数。Go 用 func 开头，参数写成「名字 类型」，
// 返回类型跟在参数列表后面。

func main() {
	// TODO（清单第 3~7 步）：把 reference.c 里 main 的内容改写到这里。
	//   - 用 := 或 var 声明累加变量，Go 的类型写在名字后面
	//   - for 和 if 的条件不写小括号，但 for 的三段之间仍然用分号
	//   - 左大括号必须和 for / if / else 在同一行
	//   - printf(...) 换成 fmt.Printf(...)，格式串和 C 一样
	//   - 行尾的分号全部删掉
	fmt.Println()
}
