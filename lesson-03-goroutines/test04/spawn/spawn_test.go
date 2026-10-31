package spawn_test

import (
	"reflect"
	"runtime"
	"sort"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"redrock/teaching-exercises/lesson-03-goroutines/test04/spawn"
)

// within 在一个新的 goroutine 里调用 f，最多等 d。
// 超时说明 f 被某个 work 卡住了——也就是 go 没有用上。
func within(t *testing.T, d time.Duration, what string, f func()) bool {
	t.Helper()
	done := make(chan struct{})
	go func() {
		defer close(done)
		f()
	}()
	select {
	case <-done:
		return true
	case <-time.After(d):
		t.Errorf("%s 在 %v 内没有返回：work 是不是在当前 goroutine 里被直接调用了？go f() 不会等 f 结束", what, d)
		return false
	}
}

// blocking 返回一个会卡住的 work：先报告自己的 id，再等 cleanup 放行。
// cleanup 放行后会等所有已经进入的 work 退出，免得它们拖到下一个测试里干扰 goroutine 计数。
func blocking(n int) (work func(int), started chan int, cleanup func()) {
	started = make(chan int, n)
	release := make(chan struct{})
	var entered, exited atomic.Int64
	work = func(id int) {
		entered.Add(1)
		defer exited.Add(1)
		started <- id
		<-release
	}
	cleanup = func() {
		close(release)
		deadline := time.Now().Add(5 * time.Second)
		for exited.Load() != entered.Load() && time.Now().Before(deadline) {
			time.Sleep(time.Millisecond)
		}
	}
	return work, started, cleanup
}

func collect(t *testing.T, started chan int, n int) []int {
	t.Helper()
	ids := make([]int, 0, n)
	timeout := time.After(5 * time.Second)
	for len(ids) < n {
		select {
		case id := <-started:
			ids = append(ids, id)
		case <-timeout:
			t.Fatalf("5 秒内只有 %d / %d 个 work 开始执行", len(ids), n)
		}
	}
	sort.Ints(ids)
	return ids
}

func TestSpawnReturnsImmediately(t *testing.T) {
	const n = 5000
	work, started, cleanup := blocking(n)
	defer cleanup()

	begin := time.Now()
	if !within(t, 2*time.Second, "Spawn(5000, 会卡住的 work)", func() { spawn.Spawn(n, work) }) {
		return
	}
	t.Logf("启动 %d 个 goroutine 用了 %v", n, time.Since(begin))

	ids := collect(t, started, n)
	for i, id := range ids {
		if id != i {
			t.Fatalf("排序后第 %d 个 id = %d, want %d；每个编号恰好执行一次", i, id, i)
		}
	}
}

func TestSpawnZero(t *testing.T) {
	called := false
	var mu sync.Mutex
	spawn.Spawn(0, func(int) {
		mu.Lock()
		called = true
		mu.Unlock()
	})
	time.Sleep(10 * time.Millisecond)
	mu.Lock()
	defer mu.Unlock()
	if called {
		t.Fatal("Spawn(0, ...) 调用了 work")
	}
}

func TestMeasureCountsNewGoroutines(t *testing.T) {
	const n = 1000
	work, started, cleanup := blocking(n)
	defer cleanup()

	var delta int
	if !within(t, 2*time.Second, "Measure(1000, 会卡住的 work)", func() { delta = spawn.Measure(n, work) }) {
		return
	}
	// 测试运行时自己也有几个 goroutine 会进出，所以给一点余量。
	if delta < n || delta > n+20 {
		t.Fatalf("Measure(%d, ...) = %d, want 约等于 %d；用 runtime.NumGoroutine() 在 Spawn 前后各读一次", n, delta, n)
	}
	collect(t, started, n)
}

func TestChunks(t *testing.T) {
	cases := []struct {
		n, size int
		want    []spawn.Chunk
	}{
		{10, 3, []spawn.Chunk{{0, 3}, {3, 6}, {6, 9}, {9, 10}}},
		{9, 3, []spawn.Chunk{{0, 3}, {3, 6}, {6, 9}}},
		{5, 10, []spawn.Chunk{{0, 5}}},
		{1, 1, []spawn.Chunk{{0, 1}}},
		{0, 3, nil},
		{-1, 3, nil},
		{10, 0, nil},
		{10, -2, nil},
	}
	for _, c := range cases {
		if got := spawn.Chunks(c.n, c.size); !reflect.DeepEqual(got, c.want) {
			t.Errorf("Chunks(%d, %d) = %v, want %v", c.n, c.size, got, c.want)
		}
	}
}

func TestSpawnChunks(t *testing.T) {
	const n, size = 1000, 64
	release := make(chan struct{})
	defer close(release)

	var mu sync.Mutex
	covered := make([]int, n)
	ran := make(chan struct{}, n)
	work := func(start, end int) {
		mu.Lock()
		for i := start; i < end; i++ {
			covered[i]++
		}
		mu.Unlock()
		ran <- struct{}{}
		<-release
	}

	var count int
	if !within(t, 2*time.Second, "SpawnChunks(1000, 64, 会卡住的 work)", func() { count = spawn.SpawnChunks(n, size, work) }) {
		return
	}
	if count != 16 {
		t.Fatalf("SpawnChunks(1000, 64, ...) = %d, want 16 个 goroutine", count)
	}
	timeout := time.After(5 * time.Second)
	for i := 0; i < count; i++ {
		select {
		case <-ran:
		case <-timeout:
			t.Fatalf("5 秒内只有 %d / %d 段开始执行", i, count)
		}
	}
	mu.Lock()
	defer mu.Unlock()
	for i, times := range covered {
		if times != 1 {
			t.Fatalf("下标 %d 被处理了 %d 次, want 恰好 1 次；各段不能重叠也不能漏", i, times)
		}
	}

	if got := spawn.SpawnChunks(0, 64, work); got != 0 {
		t.Fatalf("SpawnChunks(0, ...) = %d, want 0", got)
	}
}

// TestBridgeGoroutinesAreCheap 不检查你的实现，只读出 runtime 的数字给你看：
// 一万个“卡住”的 goroutine 大约占多少内存。
func TestBridgeGoroutinesAreCheap(t *testing.T) {
	const n = 10000
	release := make(chan struct{})
	ready := make(chan struct{}, n)

	var before, after runtime.MemStats
	runtime.ReadMemStats(&before)
	for i := 0; i < n; i++ {
		go func() {
			ready <- struct{}{}
			<-release
		}()
	}
	for i := 0; i < n; i++ {
		<-ready
	}
	runtime.ReadMemStats(&after)
	alive := runtime.NumGoroutine()
	close(release)

	t.Logf("此刻进程里有 %d 个 goroutine", alive)
	if after.StackInuse > before.StackInuse {
		t.Logf("栈内存增加了约 %d KB，平均每个 goroutine 约 %d 字节",
			(after.StackInuse-before.StackInuse)/1024, (after.StackInuse-before.StackInuse)/n)
	}
	t.Log("对比：一个 OS 线程的默认栈通常是 MB 级（Linux 上常见 8 MB），一万个线程就是几十 GB 的虚拟内存。")
	t.Log("goroutine 的栈从几 KB 起步、按需增长，由 Go runtime 在用户态调度——这就是课件里“协程更轻”的含义。")
}
