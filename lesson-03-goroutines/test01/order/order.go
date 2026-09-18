package order

// Run 按输入顺序返回任务结果；实现可以先用顺序版本，再尝试并发优化。
func Run(tasks []func() string) []string {
	out := make([]string, 0, len(tasks))
	for _, task := range tasks {
		out = append(out, task())
	}
	return out
}
