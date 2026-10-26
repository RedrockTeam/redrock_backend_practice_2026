package pool

import "context"

// Job 是一个可以被 worker 执行的任务。
type Job func(context.Context) (int, error)

// Result 是某一个 Job 的执行结果。
type Result struct {
	Index int   // 这个结果属于 jobs[Index]
	Value int   // Job 的返回值
	Err   error // Job 返回的错误；任务没来得及启动时是 ctx.Err()
}

// Run 用 workers 个 worker 并发执行 jobs。
//
// 约定：
//   - 返回的切片长度等于 len(jobs)，结果按原下标归位：results[i].Index == i；
//   - 同时运行的任务数不超过 workers；
//   - ctx 取消后停止领取新任务，已经开始的任其自然结束；
//   - 没启动的任务也要占一个结果，Err 设为 ctx.Err()；
//   - workers <= 0 或 jobs 为空时返回空结果，不执行任务；
//   - 某个任务返回错误不会自动取消其他任务，每个任务有自己的 Err。
//
// TODO: 现在它只返回了一排零值。请用 channel、WaitGroup、select 和 context 完成它。
func Run(ctx context.Context, workers int, jobs []Job) []Result {
	results := make([]Result, len(jobs))
	if workers <= 0 || len(jobs) == 0 {
		return nil
	}

	return results
}
