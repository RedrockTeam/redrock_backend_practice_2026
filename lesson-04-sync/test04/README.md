# test04：可取消的任务池（从骨架开始实现）

对应课件：第四课 · 第三部分（`select` / `context`）、第四部分（`WaitGroup`）、第七部分第 3 小节（固定 worker 数量的任务池）。
题型：**从骨架开始实现**。

## 任务

实现 `pool/pool.go` 里的：

```go
type Job func(context.Context) (int, error)

type Result struct {
    Index int
    Value int
    Err   error
}

func Run(ctx context.Context, workers int, jobs []Job) []Result
```

这是一个固定 worker 数量的任务池。约定如下：

| 约定 | 说明 |
|---|---|
| 结果数量 | `len(results) == len(jobs)`，即使取消也要给没跑的任务留位置 |
| 结果顺序 | `results[i].Index == i`，值属于 `jobs[i]`，与完成先后无关 |
| 并发上限 | 同一时刻正在执行的任务数不超过 `workers` |
| 取消 | `ctx` 取消后不再领取新任务；已开始的任务让它自然结束 |
| 未启动的任务 | `Err` 设为 `ctx.Err()`，不要留下零值假装成功 |
| 错误 | 某个任务返回错误只写进它自己的 `Err`，不自动取消其他任务 |
| 边界 | `workers <= 0` 或 `jobs` 为空时返回空结果，不执行任务 |
| 不泄漏 | 所有启动的 goroutine 都要有退出路径 |

starter 现在能编译，但只返回一排零值；所有测试都会红。

## 先跑一次

```bash
go test -v ./lesson-04-sync/test04/pool
```

失败信息会分别指出：结果下标不对、没有 worker 启动、取消后没有 `ctx.Err()` 等。不要一次改完所有东西，先让 `TestRunReturnsOrderedResults` 变绿，再加并发限制，最后处理取消。

## 设计提示（不是完整代码）

想清楚四件事，再动手：

1. **任务从哪里来？** 一个 `jobs` channel，或者一个“下一个下标”的共享变量（后者要同步）。用 channel 更贴合本课。
2. **worker 怎么退出？** 任务 channel 关闭后，`for index := range jobCh` 自然结束；每个 worker 结束时 `Done`。
3. **谁来关闭任务 channel？** 投递任务的人。投递时也要能响应取消：`select { case jobCh <- i: case <-ctx.Done(): ... }`。
4. **谁负责等待？** 一个 `WaitGroup` 等所有 worker 退出，`Run` 在 `Wait` 返回后再返回结果。不要在 worker 里 `Wait`，也不要在投递者里等自己。

几个容易踩的坑：

- `Add` 在启动 worker **之前**调用；
- 结果切片按下标写，不要用 `append` 收集，否则顺序会随完成顺序变化；
- 取消后，已经发出去但还没被 worker 取走的任务，也要有一个确定的结局；
- 投递者本身也是一个 goroutine，`Run` 要等它退出，否则它可能在 `Run` 返回后还在写结果；
- 不要用 `time.Sleep` 控制并发数。

## 评分点

| 测试 | 检查内容 |
|---|---|
| `TestRunReturnsOrderedResults` | 每个任务都执行；结果按 `Index` 归位；错误为空 |
| `TestRunRespectsWorkerLimit` | 峰值并发正好不超过 `workers`；不会只启动一个 worker |
| `TestRunKeepsPerJobErrors` | 每个任务保留自己的错误；一个失败不取消其他任务 |
| `TestRunCancellationDoesNotHang` | 取消后 2 秒内返回；未完成任务带 `ctx.Err()` |
| `TestRunAlreadyCanceledAndEdgeCases` | 已取消时不再执行任务；`workers = 0`、`jobs = nil` 的边界 |
| `TestRunManyJobs` | 500 个任务时数量、顺序、值都正确 |

## 做完以后

你现在有一个能限制并发、能取消、能保留顺序的任务池。它只在一个进程里跑，任务通过 channel 送进来。

如果任务来自**另一台机器**呢？channel 不能跨进程，`net.Conn` 才行。下一题先把“连接”抽象成一根字节流，你会看到：一旦离开 channel，“消息”这个概念就消失了，只剩字节。
