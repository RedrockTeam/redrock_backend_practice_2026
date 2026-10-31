# test08（结课题）：map 是句柄，Snapshot 要有所有权边界

对应课件：第二课 · map（句柄传递、`delete`、range、零值）、对象与方法。
题型：**补全代码 + 下一课桥接**。

## 背景

`Cache` 用 map 存“键 -> 值”。map 变量本身只是一个句柄：把它传给函数，函数里改内容，外面看得到；但 `Snapshot` 如果直接 `return c.values`，调用方拿到的就是内部 map，可以绕过 `Cache` 直接写。

这一题先把**单执行流**下的行为做对；并发留到第三、四课。

## 任务

编辑 `cache/cache.go`：

| 方法 | 要求 |
|---|---|
| `New() *Cache` | 返回空缓存 |
| `(*Cache) Put(key, value)` | 写入/覆盖；零值 `Cache` 也能用 |
| `(*Cache) Get(key) (int, bool)` | 不存在返回 `0, false` |
| `(*Cache) Delete(key)` | 删除；键不存在不 panic |
| `(*Cache) Len() int` | 键的数量 |
| `(*Cache) Snapshot() map[string]int` | 返回**独立副本**；空缓存也要返回可写的空 map |

## 提示（不是答案）

- map 的零值是 `nil`，写之前要 `make`；`Put` 里可以 lazy 初始化；
- 读不存在的键返回零值，必须用 `v, ok :=` 区分；
- `delete(m, key)` 删除不存在的键不会 panic；
- `Snapshot` 用 `make(map[string]int, len(c.values))` 再逐个复制；
- 不要写 `return c.values`。

## 评分点

| 测试 | 检查内容 |
|---|---|
| `TestPutGetDelete` | 增删改查、覆盖、缺失键、删除不存在键 |
| `TestZeroValueCache` | `var c Cache` 可以直接 `Put` / `Get` / `Len` |
| `TestSnapshotIsACopy` | 修改快照不影响缓存；空缓存返回可写空 map |
| `TestBridgeIndependentCachesAreSafe` | 两个 goroutine 各写各的 Cache 互不影响 |

## 做完以后：下一课的问题

`TestBridgeIndependentCachesAreSafe` 里两个 goroutine 各用各的 `Cache`，所以安全。如果它们写的是**同一个** `Cache` 呢？

```text
goroutine A: c.Put("Lin", 95)
goroutine B: c.Put("Lin", 100)
```

map 不是并发安全的：轻则丢更新，重则 `fatal error: concurrent map writes`。而且 `Snapshot` 即使复制了 map，复制过程中也可能看到“写了一半”的状态。

第三课先在 `test03` 用一台玩具 CPU 模拟两个执行流交错执行 `count++`，亲眼看到丢更新；再在 `test05`〜`test07` 把 goroutine 跑起来，看清“函数返回了，协程还没结束”的问题；第四课再用 `WaitGroup`、`Mutex` / `RWMutex` 把这些问题管起来。你在这一题写的 `Snapshot` 所有权边界，会在第四课 `test05` 的账本里原样再用一次。
