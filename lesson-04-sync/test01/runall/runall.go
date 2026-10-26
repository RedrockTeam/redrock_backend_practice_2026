package runall

// RunAll 并发执行 tasks，并按输入顺序返回每个任务的结果。
//
// TODO: 现在它启动了 goroutine 后立刻返回，也没有等待它们结束。
// 请让 RunAll 在所有任务都结束后再返回。
//
// 要求：
//   - 每个任务恰好执行一次；
//   - 结果写回它原来的下标，顺序不能乱；
//   - 不能用 time.Sleep 猜“大概跑完了”。
func RunAll(tasks []func() int) []int {
	result := make([]int, len(tasks))
	for i, task := range tasks {
		i, task := i, task
		go func() {
			result[i] = task()
		}()
	}
	return result
}
