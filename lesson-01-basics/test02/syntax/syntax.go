package syntax

import "fmt"

// Add 返回 a 与 b 的和。
// 它是导出的（首字母大写），其他包可以调用。
func Add(a, b int) int {
	return add(a, b)
}

// add 是只在 syntax 包内使用的加法函数（首字母小写）。
func add(a, b int) int {
	// TODO: 用 Go 的函数返回语法完成加法。
	// C：int add(int a, int b) { return a + b; }
	// Go：返回类型写在参数列表后面。
	return 0
}

// ClassifyNumber 在 n 大于 0 时返回 "positive"，
// 等于 0 时返回 "zero"，小于 0 时返回 "negative"。
func ClassifyNumber(n int) string {
	// TODO: 用 if / else if / else 完成三个分支。
	// Go 的 if 条件不写小括号，左大括号必须和 if 同一行。
	return ""
}

// SumTo 返回 1 + 2 + ... + n；n 小于等于 0 时返回 0。
func SumTo(n int) int {
	// TODO: 声明一个累加变量，再用 for 循环累加。
	// Go 的 for 条件不写小括号，但三段之间仍然用分号。
	return 0
}

// Describe 返回 "小明 今年 18 岁" 这样的一行描述，末尾没有换行。
// 名字和年龄之间、年龄和 "岁" 之间各有一个空格。
func Describe(name string, age int) string {
	// TODO: 用 fmt.Sprintf 拼出这句话。
	// Sprintf 和 Printf 的格式化规则一样，但它返回字符串而不是打印。
	// 字符串用 %s，整数用 %d。
	return fmt.Sprintf("")
}

// MaxScore 是导出的常量，其他包可以用 syntax.MaxScore 读到它。
const MaxScore = 100

// minScore 是未导出的常量，只有 syntax 包内部看得见。
const minScore = 0

// Clamp 把 score 收进 [minScore, MaxScore] 区间：
// 低于下限返回下限，高于上限返回上限，否则原样返回。
func Clamp(score int) int {
	// TODO: 用 if 判断两个边界，直接用上面两个常量，不要写死 0 和 100。
	return 0
}

// InitialCount 返回一个还没开始计数的计数器的初始值。
//
// 要求：用 var 声明一个 int 变量，不要给它赋值，直接返回它。
// Go 的变量未初始化时有零值：int 是 0，string 是 ""，bool 是 false。
func InitialCount() int {
	// TODO: 声明一个 var 变量后直接 return，别写 return 0。
	return -1
}
