package lazyfetch_test

import (
	"context"
	"errors"
	"runtime"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"redrock/teaching-exercises/lesson-04-sync/test06/lazyfetch"
)

func slowLoad(calls *atomic.Int32) func() lazyfetch.Config {
	return func() lazyfetch.Config {
		calls.Add(1)
		// 放大「检查」和「赋值」之间的窗口：真实世界里这是读文件、连配置中心。
		time.Sleep(20 * time.Millisecond)
		return lazyfetch.Config{Endpoint: "svc"}
	}
}

func echoFetch(endpoint, key string) string { return endpoint + "/" + key }

func TestFetchReturnsValue(t *testing.T) {
	var calls atomic.Int32
	c := lazyfetch.New(slowLoad(&calls), echoFetch)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	for _, key := range []string{"a", "b"} {
		got, err := c.Fetch(ctx, key)
		if err != nil {
			t.Fatalf("Fetch(%q) error = %v, want nil", key, err)
		}
		if want := "svc/" + key; got != want {
			t.Fatalf("Fetch(%q) = %q, want %q", key, got, want)
		}
	}
	if got := c.Config().Endpoint; got != "svc" {
		t.Fatalf("Config().Endpoint = %q, want %q", got, "svc")
	}
	if n := calls.Load(); n != 1 {
		t.Fatalf("顺序调用下 load 被执行了 %d 次, want 1", n)
	}
}

func TestConfigLoadsOnceUnderConcurrency(t *testing.T) {
	for round := 0; round < 3; round++ {
		var calls atomic.Int32
		c := lazyfetch.New(slowLoad(&calls), echoFetch)

		const goroutines = 100
		var wg sync.WaitGroup
		start := make(chan struct{})
		wg.Add(goroutines)
		for i := 0; i < goroutines; i++ {
			go func() {
				defer wg.Done()
				<-start
				if got := c.Config().Endpoint; got != "svc" {
					t.Errorf("Config().Endpoint = %q, want %q", got, "svc")
				}
			}()
		}
		close(start)
		wg.Wait()

		if n := calls.Load(); n != 1 {
			t.Fatalf("round %d: 100 个 goroutine 同时第一次调用 Config()，load 被执行了 %d 次, want 1；"+
				"「先检查 nil 再赋值」不是一步，几个 goroutine 会同时看到 nil。课件第六部分：sync.Once",
				round, n)
		}
	}
}

func TestFetchTimesOutWithWrappedError(t *testing.T) {
	release := make(chan struct{})
	defer close(release)
	blocking := func(endpoint, key string) string {
		<-release
		return endpoint + "/" + key
	}

	var calls atomic.Int32
	c := lazyfetch.New(slowLoad(&calls), blocking)
	c.Config() // 先把配置加载好，下面只计 Fetch 自己的时间

	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	begin := time.Now()
	got, err := c.Fetch(ctx, "user-42")
	elapsed := time.Since(begin)

	if elapsed > 500*time.Millisecond {
		t.Fatalf("超时 50ms 的 Fetch 用了 %v 才返回；要用 select 同时等结果和 ctx.Done()", elapsed)
	}
	if got != "" {
		t.Fatalf("超时时 Fetch 返回值 = %q, want 空字符串", got)
	}
	if err == nil {
		t.Fatal("超时时 Fetch 的 error = nil")
	}
	if !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("errors.Is(err, context.DeadlineExceeded) = false, err = %v；包装错误要用 %%w", err)
	}
	if !strings.Contains(err.Error(), "user-42") {
		t.Fatalf("err = %q，错误信息里没有 key \"user-42\"；调用方看日志时要知道是哪一次请求超时", err)
	}
}

func TestTimedOutFetchDoesNotLeakGoroutines(t *testing.T) {
	release := make(chan struct{})
	blocking := func(endpoint, key string) string {
		<-release
		return endpoint + "/" + key
	}

	var calls atomic.Int32
	c := lazyfetch.New(slowLoad(&calls), blocking)
	c.Config()

	baseline := runtime.NumGoroutine()

	const requests = 50
	var wg sync.WaitGroup
	wg.Add(requests)
	for i := 0; i < requests; i++ {
		go func() {
			defer wg.Done()
			ctx, cancel := context.WithTimeout(context.Background(), 20*time.Millisecond)
			defer cancel()
			if _, err := c.Fetch(ctx, "k"); err == nil {
				t.Error("慢服务还没放行，Fetch 却没有超时")
			}
		}()
	}
	wg.Wait()

	// 50 次请求都已经超时返回。现在慢服务终于回结果了：
	// 后台 goroutine 们要把结果交出去然后退出。
	close(release)

	deadline := time.Now().Add(2 * time.Second)
	for {
		n := runtime.NumGoroutine()
		if n <= baseline {
			return
		}
		if time.Now().After(deadline) {
			t.Fatalf("慢服务放行 2 秒后仍有 %d 个 goroutine（请求前 %d 个）；"+
				"Fetch 超时返回后再也没人读 result，后台 goroutine 永远卡在发送上。"+
				"想一想：怎样让这次发送不需要有人接？", n, baseline)
		}
		time.Sleep(10 * time.Millisecond)
	}
}
