# test04：go 不等人，goroutine 很轻

对应课件：第三课 · 协程（协程比线程轻）、协程的 infra（runtime、GMP）、go 关键字的使用（函数调用对比：`work()` 阻塞到返回，`go work()` 立即返回）。
题型：**补全代码（中等偏易）**。

## 现象

`spawn/spawn.go` 里的 `Spawn` 现在是顺序调用：

```go
for i := 0; i < n; i++ {
    work(i)
}
```

测试给它的 `work` 会**一直卡住**，直到测试自己放行。所以顺序版本的 `Spawn` 永远回不来。先跑：

```bash
go test -v ./lesson-03-goroutines/test04/spawn
```

你会看到：

```text
Spawn(5000, 会卡住的 work) 在 2s 内没有返回：work 是不是在当前 goroutine 里被直接调用了？
```

这就是课件里那段对比：

```go
work()      // 当前 goroutine 执行，阻塞到返回
go work()   // 新 goroutine 执行，立即返回
```

## 任务

编辑 `spawn/spawn.go`：

| 函数 | 要求 |
|---|---|
| `Spawn(n, work)` | 为 `0..n-1` 每个编号启动一个 goroutine 执行 `work(id)`，**立刻返回** |
| `Measure(n, work)` | 调用 `Spawn`，返回调用前后 `runtime.NumGoroutine()` 的增量 |
| `Chunks(n, size)` | 把 `[0, n)` 切成每段 `size` 个的左闭右开区间；`n <= 0` 或 `size <= 0` 返回 `nil` |
| `SpawnChunks(n, size, work)` | 每一段启动一个 goroutine 执行 `work(start, end)`，立刻返回启动的个数 |

这一题的函数**都不等 work 结束**——它们只负责“启动”。等待是第四课的事。

## 提示（不是答案）

- `go` 后面必须是函数**调用**：`go work(i)` 可以，`go work` 不行；
- 传给 goroutine 的编号要在**启动那一刻**定下来：`go work(i)` 的参数是当场求值的；用闭包时写成 `go func(id int) { work(id) }(i)`；
- `runtime.NumGoroutine()` 返回此刻进程里活着的 goroutine 个数（包括正在调用它的那一个）；
- `Chunks` 是一个步长为 `size` 的 `for` 循环，最后一段的 `End` 用 `min(start+size, n)` 夹住；
- `SpawnChunks` 可以直接复用 `Chunks`。

## 自己验证

```bash
go test -v ./lesson-03-goroutines/test04/spawn
go test -v -run TestBridge ./lesson-03-goroutines/test04/spawn   # 一万个 goroutine 占多少内存
```

## 评分点

| 测试 | 检查内容 |
|---|---|
| `TestSpawnReturnsImmediately` | 5000 个会卡住的 work，`Spawn` 仍能立刻返回；每个编号恰好执行一次 |
| `TestSpawnZero` | `n = 0` 不调用 work |
| `TestMeasureCountsNewGoroutines` | 增量约等于启动的个数 |
| `TestChunks` | 整除、不整除、只有一段、各种非法输入 |
| `TestSpawnChunks` | 每段一个 goroutine、立刻返回、区间不重不漏 |
| `TestBridgeGoroutinesAreCheap` | 读出一万个 goroutine 的栈内存（只打印，不检查你的实现） |

## 做完以后

`TestBridgeGoroutinesAreCheap` 会告诉你：一万个 goroutine 的栈加起来大约只有几十 MB 甚至更少。换成一万个 OS 线程，光默认栈就要几十 GB 的虚拟内存。这就是课件里“协程运行在用户态、由 runtime 调度”的实际回报。

但 `Spawn` 有一个你现在还没法解决的问题：**它根本不知道 work 什么时候结束**。测试能放行和等待，是因为测试用了第四课的工具。

下一题 `test05` 让 goroutine 往同一个输出里写东西；`test06` 把求和拆给 goroutine，第一次需要“拿回结果”；`test07` 会让你看清楚：不等就拿结果，会拿到什么。
