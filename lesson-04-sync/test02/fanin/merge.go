package fanin

// Merge 并发读取所有 inputs，把读到的值转发到返回的 channel。
//
// 约定：
//   - 返回的 channel 由 Merge 负责关闭；
//   - 只有当所有输入都关闭、所有转发 goroutine 都退出后，才能 close(out)；
//   - 输入 channel 只由它们的所有者关闭，Merge 不能关闭输入。
//
// TODO: 下面的代码启动了转发 goroutine，但没有在所有输入结束后关闭 out，
// 所以消费者会一直等下去。请补全“等待 + 关闭”的部分。
func Merge(inputs ...<-chan int) <-chan int {
	out := make(chan int)

	for _, input := range inputs {
		input := input
		go func() {
			for value := range input {
				out <- value
			}
		}()
	}

	return out
}
