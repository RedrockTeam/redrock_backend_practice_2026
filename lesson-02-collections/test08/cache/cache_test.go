package cache_test

import (
	"reflect"
	"sync"
	"testing"

	"redrock/teaching-exercises/lesson-02-collections/test08/cache"
)

func TestPutGetDelete(t *testing.T) {
	c := cache.New()

	if value, ok := c.Get("missing"); ok || value != 0 {
		t.Fatalf(`Get("missing") = %d, %v, want 0, false`, value, ok)
	}

	c.Put("Lin", 95)
	if value, ok := c.Get("Lin"); !ok || value != 95 {
		t.Fatalf(`Get("Lin") = %d, %v, want 95, true`, value, ok)
	}

	c.Put("Lin", 100)
	if value, _ := c.Get("Lin"); value != 100 {
		t.Fatalf(`覆盖后 Get("Lin") = %d, want 100`, value)
	}

	c.Delete("Lin")
	if _, ok := c.Get("Lin"); ok {
		t.Fatal(`Delete("Lin") 之后仍然能查到`)
	}
	c.Delete("never-existed") // 不应该 panic
}

func TestZeroValueCache(t *testing.T) {
	var c cache.Cache
	c.Put("Lin", 95)

	if value, ok := c.Get("Lin"); !ok || value != 95 {
		t.Fatalf("零值 Cache.Get(\"Lin\") = %d, %v, want 95, true", value, ok)
	}
	if got := c.Len(); got != 1 {
		t.Fatalf("零值 Cache.Len() = %d, want 1", got)
	}
}

func TestSnapshotIsACopy(t *testing.T) {
	c := cache.New()
	c.Put("Lin", 95)
	c.Put("Wang", 80)

	snapshot := c.Snapshot()
	if snapshot == nil {
		t.Fatal("Snapshot() 返回 nil；空缓存也应该返回可写的空 map")
	}
	if !reflect.DeepEqual(snapshot, map[string]int{"Lin": 95, "Wang": 80}) {
		t.Fatalf("Snapshot() = %#v", snapshot)
	}

	snapshot["Lin"] = 0
	snapshot["Zhao"] = 70

	if value, _ := c.Get("Lin"); value != 95 {
		t.Fatalf(`修改 snapshot 后 Get("Lin") = %d, want 95；Snapshot 必须复制 map`, value)
	}
	if _, ok := c.Get("Zhao"); ok {
		t.Fatal("给 snapshot 新增键之后 Cache 也出现了这个键；Snapshot 不能是内部 map 的别名")
	}
	if got := c.Len(); got != 2 {
		t.Fatalf("修改 snapshot 后 Cache.Len() = %d, want 2", got)
	}

	empty := cache.New().Snapshot()
	if empty == nil {
		t.Fatal("空 Cache 的 Snapshot() 返回 nil，want 空 map")
	}
}

// TestBridgeIndependentCachesAreSafe 只是把「两个 goroutine 各用各的 Cache」
// 摆出来：它们互不影响。下一课会问：如果两个 goroutine 写的是同一个 Cache 呢？
func TestBridgeIndependentCachesAreSafe(t *testing.T) {
	a := cache.New()
	b := cache.New()

	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		a.Put("shared-key", 1)
	}()
	go func() {
		defer wg.Done()
		b.Put("shared-key", 2)
	}()
	wg.Wait()

	if value, _ := a.Get("shared-key"); value != 1 {
		t.Fatalf("缓存 a 的值 = %d, want 1", value)
	}
	if value, _ := b.Get("shared-key"); value != 2 {
		t.Fatalf("缓存 b 的值 = %d, want 2", value)
	}

	t.Log("两个 goroutine 各用各的 Cache 是安全的，因为它们没有共享内存。")
	t.Log("如果它们同时写同一个 Cache，map 就会被并发读写；下一课先看 goroutine，第四课再解决同步。")
}
