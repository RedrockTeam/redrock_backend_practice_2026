package pool_test

import (
	"context"
	"errors"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"redrock/teaching-exercises/lesson-04-sync/test07/pool"
)

func TestRunReturnsOrderedResults(t *testing.T) {
	const jobCount = 24

	jobs := make([]pool.Job, jobCount)
	for i := 0; i < jobCount; i++ {
		i := i
		jobs[i] = func(context.Context) (int, error) {
			// 让完成顺序和输入顺序不同，才能看出结果有没有按下标归位。
			time.Sleep(time.Duration((jobCount-i)%4) * time.Millisecond)
			return i * 10, nil
		}
	}

	got := pool.Run(context.Background(), 4, jobs)
	if len(got) != jobCount {
		t.Fatalf("len(Run(...)) = %d, want %d", len(got), jobCount)
	}
	for i, result := range got {
		if result.Index != i {
			t.Fatalf("results[%d].Index = %d；结果必须按输入下标归位", i, result.Index)
		}
		if result.Value != i*10 {
			t.Fatalf("results[%d].Value = %d, want %d", i, result.Value, i*10)
		}
		if result.Err != nil {
			t.Fatalf("results[%d].Err = %v, want nil", i, result.Err)
		}
	}
}

func TestRunRespectsWorkerLimit(t *testing.T) {
	const (
		jobCount = 12
		workers  = 3
	)

	started := make(chan struct{}, jobCount)
	release := make(chan struct{})
	var releaseOnce sync.Once
	releaseAll := func() { releaseOnce.Do(func() { close(release) }) }
	defer releaseAll()

	var running atomic.Int32
	var peak atomic.Int32

	jobs := make([]pool.Job, jobCount)
	for i := 0; i < jobCount; i++ {
		i := i
		jobs[i] = func(context.Context) (int, error) {
			current := running.Add(1)
			for {
				old := peak.Load()
				if current <= old || peak.CompareAndSwap(old, current) {
					break
				}
			}

			started <- struct{}{}
			<-release
			running.Add(-1)
			return i, nil
		}
	}

	done := make(chan []pool.Result, 1)
	go func() {
		done <- pool.Run(context.Background(), workers, jobs)
	}()

	for i := 0; i < workers; i++ {
		select {
		case <-started:
		case <-time.After(2 * time.Second):
			t.Fatalf("只启动了 %d 个任务，want 至少 %d 个 worker 同时在跑", i, workers)
		}
	}

	select {
	case <-started:
		t.Fatalf("同时启动了超过 %d 个任务；worker 数量没有被限制住", workers)
	case <-time.After(150 * time.Millisecond):
	}

	releaseAll()

	var results []pool.Result
	select {
	case results = <-done:
	case <-time.After(2 * time.Second):
		t.Fatal("释放任务后 Run 没有返回；可能有 worker 永远卡住")
	}

	if len(results) != jobCount {
		t.Fatalf("len(Run(...)) = %d, want %d", len(results), jobCount)
	}
	if got := peak.Load(); got != workers {
		t.Fatalf("观察到的峰值并发 = %d, want %d", got, workers)
	}
}

func TestRunKeepsPerJobErrors(t *testing.T) {
	sentinel := errors.New("boom")
	jobs := []pool.Job{
		func(context.Context) (int, error) { return 7, nil },
		func(context.Context) (int, error) { return 0, sentinel },
		func(context.Context) (int, error) { return 9, nil },
	}

	got := pool.Run(context.Background(), 2, jobs)
	if len(got) != 3 {
		t.Fatalf("len(Run(...)) = %d, want 3", len(got))
	}
	if got[0].Value != 7 || got[0].Err != nil {
		t.Fatalf("results[0] = %+v, want Value 7, Err nil", got[0])
	}
	if !errors.Is(got[1].Err, sentinel) {
		t.Fatalf("results[1].Err = %v, want %v；一个任务失败不应吞掉自己的错误", got[1].Err, sentinel)
	}
	if got[2].Value != 9 || got[2].Err != nil {
		t.Fatalf("results[2] = %+v, want Value 9, Err nil；一个任务失败不应取消其他任务", got[2])
	}
}

func TestRunCancellationDoesNotHang(t *testing.T) {
	const jobCount = 50

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	jobs := make([]pool.Job, jobCount)
	for i := range jobs {
		jobs[i] = func(ctx context.Context) (int, error) {
			<-ctx.Done()
			return 0, ctx.Err()
		}
	}

	done := make(chan []pool.Result, 1)
	go func() {
		done <- pool.Run(ctx, 4, jobs)
	}()

	time.Sleep(30 * time.Millisecond)
	cancel()

	select {
	case got := <-done:
		if len(got) != jobCount {
			t.Fatalf("取消后 len(Run(...)) = %d, want %d；没启动的任务也要留下结果", len(got), jobCount)
		}
		for i, result := range got {
			if result.Index != i {
				t.Fatalf("results[%d].Index = %d, want %d", i, result.Index, i)
			}
			if result.Err == nil {
				t.Fatalf("results[%d].Err = nil；取消后没执行的任务要带 ctx.Err()", i)
			}
		}
	case <-time.After(2 * time.Second):
		t.Fatal("ctx 取消后 2 秒内 Run 仍未返回；发送任务或等 worker 的地方没有监听 ctx.Done()")
	}
}

func TestRunAlreadyCanceledAndEdgeCases(t *testing.T) {
	var called atomic.Int32
	job := func(context.Context) (int, error) {
		called.Add(1)
		return 1, nil
	}

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	got := pool.Run(ctx, 2, []pool.Job{job})
	if len(got) != 1 {
		t.Fatalf("已取消时 len(Run(...)) = %d, want 1", len(got))
	}
	if got[0].Index != 0 {
		t.Fatalf("已取消时 results[0].Index = %d, want 0", got[0].Index)
	}
	if got[0].Err == nil {
		t.Fatal("已取消时 results[0].Err = nil；没执行的任务要带 ctx.Err()")
	}
	if called.Load() != 0 {
		t.Fatalf("ctx 已取消仍执行了 %d 个任务；取消后不应再启动新任务", called.Load())
	}

	if got := pool.Run(context.Background(), 0, []pool.Job{job}); len(got) != 0 {
		t.Fatalf("workers = 0 时 len(Run(...)) = %d, want 0", len(got))
	}
	if got := pool.Run(context.Background(), 2, nil); len(got) != 0 {
		t.Fatalf("jobs 为 nil 时 len(Run(...)) = %d, want 0", len(got))
	}
	if called.Load() != 0 {
		t.Fatalf("workers = 0 或 jobs 为空时执行了任务")
	}
}

func TestRunManyJobs(t *testing.T) {
	const jobCount = 500

	jobs := make([]pool.Job, jobCount)
	for i := 0; i < jobCount; i++ {
		i := i
		jobs[i] = func(context.Context) (int, error) {
			return i + 1, nil
		}
	}

	got := pool.Run(context.Background(), 8, jobs)
	if len(got) != jobCount {
		t.Fatalf("len(Run(...)) = %d, want %d", len(got), jobCount)
	}
	for i, result := range got {
		if result.Index != i || result.Value != i+1 || result.Err != nil {
			t.Fatalf("results[%d] = %+v, want Index %d, Value %d, Err nil", i, result, i, i+1)
		}
	}
}
