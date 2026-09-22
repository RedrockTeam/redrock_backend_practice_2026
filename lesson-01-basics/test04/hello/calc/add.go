// Package calc 就是课件第三部分那个 example.com/hello/calc 的真身。
// 它所在的目录是 lesson-01-basics/test04/hello/calc，
// 所以它的完整导入路径 = go.mod 里的 module 路径 + 这串子目录。
package calc

// Add 返回 a 与 b 的和。首字母大写，别的包可以调用它。
func Add(a, b int) int {
	// TODO: 返回两数之和。
	return 0
}

// AddDoubled 把两个数各翻一倍再相加。
// 要求用下面的 helper 完成，不要自己写 *2：
// helper 首字母小写，别的包看不见它，但同一个包里随便用。
func AddDoubled(a, b int) int {
	// TODO: 调用 helper 完成。
	return 0
}

// helper 返回 x 的两倍。首字母小写，只在 calc 包内可见。
func helper(x int) int {
	// TODO: 返回 x 的两倍。
	return 0
}
