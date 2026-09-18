package even

// IsEven 判断 n 是否为偶数。
//
// 当前实现是错的：它漏掉了一整类输入。
// 本题的顺序是「先写测试暴露问题，再修实现」，
// 所以请先在 student_test.go 里把问题测出来，再回到这里修。
func IsEven(n int) bool {
	return n > 0 && n%2 == 0
}
