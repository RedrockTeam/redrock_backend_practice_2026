package pipeline_test

import (
	"context"
	"reflect"
	"runtime"
	"testing"
	"time"

	"redrock/teaching-exercises/lesson-04-sync/test04/pipeline"
)

// drain 一直读 ch，直到它被关闭或超时。返回读到的值和是否已关闭。
func drain(ch <-chan int, timeout time.Duration) ([]int, bool) {
	var got []int
	deadline := time.After(timeout)
	for {
		select {
		case v, ok := <-ch:
			if !ok {
				return got, true
			}
			got = append(got, v)
		case <-deadline:
			return got, false
		}
	}
}

// waitGoroutines 轮询，直到 goroutine 数回落到 baseline 以内或超时。
func waitGoroutines(baseline int, timeout time.Duration) (int, bool) {
	deadline := time.Now().Add(timeout)
	for {
		n := runtime.NumGoroutine()
		if n <= baseline {
			return n, true
		}
		if time.Now().After(deadline) {
			return n, false
		}
		time.Sleep(10 * time.Millisecond)
	}
}

func numbers(n int) []int {
	nums := make([]int, n)
	for i := range nums {
		nums[i] = i + 1
	}
	return nums
}

func TestPipelineRunsToCompletion(t *testing.T) {
	ctx := context.Background()
	even := func(v int) bool { return v%2 == 0 }

	out := pipeline.Filter(ctx, pipeline.Square(ctx, pipeline.Generate(ctx, numbers(10)...)), even)
	got, closed := drain(out, 2*time.Second)
	if !closed {
		t.Fatalf("2 秒内输出没有关闭；每一级都要 defer close(out)。已读到 %v", got)
	}
	want := []int{4, 16, 36, 64, 100}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("Filter(Square(Generate(1..10)), 偶数) = %v, want %v", got, want)
	}

	empty, closed := drain(pipeline.Square(ctx, pipeline.Generate(ctx)), time.Second)
	if !closed || len(empty) != 0 {
		t.Fatalf("Generate() 没有参数时 = %v, closed=%v；want 空且立刻关闭", empty, closed)
	}
}

func TestCancelStopsTheWholePipeline(t *testing.T) {
	baseline := runtime.NumGoroutine()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	out := pipeline.Square(ctx, pipeline.Generate(ctx, numbers(10000)...))
	first, second := <-out, <-out
	if first != 1 || second != 4 {
		t.Fatalf("前两个结果 = %d, %d, want 1, 4", first, second)
	}

	// 消费者拿够了，取消后**不再读**。上游每一级都必须自己退出。
	cancel()

	if n, ok := waitGoroutines(baseline, time.Second); !ok {
		t.Fatalf("cancel() 1 秒后仍有 %d 个 goroutine（开始前 %d 个）；"+
			"消费者不读了，上游还卡在 out <- v 上。发送要写成 select { case out <- v: case <-ctx.Done(): return }",
			n, baseline)
	}
}

func TestCancelClosesOutputPromptly(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	out := pipeline.Square(ctx, pipeline.Generate(ctx, numbers(10000)...))
	<-out
	<-out
	cancel()

	rest, closed := drain(out, time.Second)
	if !closed {
		t.Fatal("cancel() 后 1 秒内输出没有关闭")
	}
	if len(rest) > 10 {
		t.Fatalf("cancel() 之后又收到了 %d 个值；取消后各级应该很快停止，而不是把剩下的 9998 个都算完", len(rest))
	}
}

func TestCancelUnblocksAStageWaitingOnInput(t *testing.T) {
	baseline := runtime.NumGoroutine()
	ctx, cancel := context.WithCancel(context.Background())

	// 上游既不发送也不关闭：Square 只能卡在接收上。
	stuck := make(chan int)
	out := pipeline.Square(ctx, stuck)

	cancel()
	if _, closed := drain(out, time.Second); !closed {
		t.Fatal("上游一直不发也不关时，cancel() 1 秒后 Square 的输出仍没关闭；" +
			"接收也要能被取消：select { case v, ok := <-in: ... case <-ctx.Done(): return }")
	}
	if n, ok := waitGoroutines(baseline, time.Second); !ok {
		t.Fatalf("cancel() 后仍有 %d 个 goroutine（开始前 %d 个）", n, baseline)
	}
}

func TestAlreadyCanceledContext(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	out := pipeline.Filter(ctx,
		pipeline.Square(ctx, pipeline.Generate(ctx, numbers(1000)...)),
		func(int) bool { return true })
	got, closed := drain(out, time.Second)
	if !closed {
		t.Fatal("ctx 一开始就已取消时，1 秒内输出没有关闭")
	}
	if len(got) > 10 {
		t.Fatalf("ctx 一开始就已取消，却收到了 %d 个值；各级应该几乎立刻停止", len(got))
	}
}
