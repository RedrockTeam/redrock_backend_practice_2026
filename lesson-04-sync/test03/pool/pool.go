package pool

import "context"

type Job func(context.Context) (int, error)
type Result struct {
	Value int
	Err   error
}

// Run 使用 workers 个 worker 执行 jobs，保证结果数量和输入顺序一致。
func Run(ctx context.Context, workers int, jobs []Job) []Result {
	// TODO: use channels and WaitGroup; handle workers <= 0 and cancellation
	return nil
}
