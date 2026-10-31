package pipeline

import "context"

// 这是课件第七部分第 1 小节的流水线：
//
//	Generate ──► Square ──► Filter ──► 消费者
//
// 每一级都是「输入 channel → 输出 channel」，由本级的 goroutine 负责
// defer close(out)。正常跑完时它已经是对的。
//
// TODO: 它不认识 ctx。消费者只想要前两个结果、然后 cancel() 时，
// 上游各级会永远卡在 out <- v 上，goroutine 泄漏；上游如果一直不发也不关，
// 本级会永远卡在接收上。请让每一级的「发送」和「接收」都能被 ctx 取消。

// Generate 依次发送 nums，发完关闭输出。
func Generate(ctx context.Context, nums ...int) <-chan int {
	out := make(chan int)
	go func() {
		defer close(out)
		for _, n := range nums {
			out <- n
		}
	}()
	return out
}

// Square 把每个输入值平方后发送，输入关闭后关闭输出。
func Square(ctx context.Context, in <-chan int) <-chan int {
	out := make(chan int)
	go func() {
		defer close(out)
		for n := range in {
			out <- n * n
		}
	}()
	return out
}

// Filter 只转发 keep(v) 为 true 的值，输入关闭后关闭输出。
func Filter(ctx context.Context, in <-chan int, keep func(int) bool) <-chan int {
	out := make(chan int)
	go func() {
		defer close(out)
		for n := range in {
			if keep(n) {
				out <- n
			}
		}
	}()
	return out
}
