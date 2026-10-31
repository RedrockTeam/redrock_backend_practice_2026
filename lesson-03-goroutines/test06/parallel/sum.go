package parallel

// Sum 返回 values 中所有元素的和。
//
// 要求：把 values 拆成至少两段，交给 goroutine 并行求和，
// 再把各段结果加起来。第三课先用 time.Sleep 等它们；
// 第四课会换成 WaitGroup。
//
// TODO: 现在它是顺序版本。先让它通过功能测试，再用 go 关键字改写。
func Sum(values []int) int {
	total := 0
	for _, value := range values {
		total += value
	}
	return total
}
