// Package app 扮演课件第三部分里的 main.go：它是「使用别人」的那一方。
//
// 课件的例子用的是 main 包，这里改成普通的 app 包，
// 因为第四部分说过 main 包不能被其他包导入——包括外部测试文件。
package app

import (
	"fmt"
	// TODO: 在这里补上 calc 包的导入路径。
	//
	// 导入路径 = module 路径 + 包所在的子目录。
	// module 路径写在仓库根目录 go.mod 的第一行，
	// calc 包在 lesson-01-basics/test04/hello/calc 目录下。
	//
	// 补好之后，下面就能用 calc.Add(...) 了。
)

// Sum 借助 calc 包完成加法。
func Sum(a, b int) int {
	// TODO: 调用 calc.Add。
	return 0
}

// Describe 返回 "1 + 2 = 3" 这样的一行，末尾没有换行。
func Describe(a, b int) string {
	// TODO: 用 fmt.Sprintf 把算式和 calc.Add 的结果拼起来。
	return fmt.Sprintf("")
}
