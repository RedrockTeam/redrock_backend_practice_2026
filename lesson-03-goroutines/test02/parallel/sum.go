package parallel

// Sum 计算整数和。先保证行为正确，再用基准测试比较并行版本是否值得。
func Sum(values []int) int {
	total := 0
	for _, value := range values {
		total += value
	}
	return total
}
