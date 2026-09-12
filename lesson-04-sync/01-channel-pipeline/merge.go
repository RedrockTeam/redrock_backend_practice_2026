package pipeline

// Merge 将多个输入合并到一个输出，并在所有输入关闭后关闭输出。
func Merge(inputs ...<-chan int) <-chan int {
	// TODO: start one forwarding goroutine per input and coordinate close
	return make(chan int)
}
