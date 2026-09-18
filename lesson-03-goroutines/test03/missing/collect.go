package missing

// Collect 故意没有等待 goroutine，结果可能为空。第四课会用同步原语修复。
func Collect(values []int) []int {
	results := make([]int, 0, len(values))
	for _, value := range values {
		go func() { results = append(results, value*value) }()
	}
	return results
}
