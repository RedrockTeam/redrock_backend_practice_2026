package missing

// Collect 返回 values 中每个元素的平方，顺序与输入一致。
//
// 当前实现启动了 goroutine，却在它们写完之前就返回了。
// 这是第三课 test07 的观察对象：这一题不改它，也不要用 time.Sleep 掩盖。
// 第四课 test01 会用 WaitGroup 把这条边界补上。
func Collect(values []int) []int {
	results := make([]int, len(values))
	for i, value := range values {
		i, value := i, value
		go func() {
			results[i] = value * value
		}()
	}
	return results
}
