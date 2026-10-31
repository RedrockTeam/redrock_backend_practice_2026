# test04：可取消的流水线（补全代码）

对应课件：第四课 · 第三部分（`select`、`context` 取消、"结果或取消"循环骨架）、第七部分第 1 小节（流水线）、第十部分 LV2。
题型：**补全已有代码**。

## 现象

`pipeline/pipeline.go` 是课件里的三级流水线：

```text
Generate ──► Square ──► Filter ──► 消费者
   │            │          │
close(out)   close(out)  close(out)
```

每一级都已经 `defer close(out)`，消费者把结果读完时，它完全正确。问题出在**消费者提前不读了**：

```go
out := Square(ctx, Generate(ctx, 1, 2, ..., 10000))
<-out
<-out
cancel() // 我只要前两个
```

先跑：

```bash
go test -v ./lesson-04-sync/test04/pipeline
```

你会看到 `TestPipelineRunsToCompletion` 通过，其余几个失败，比如：

```text
cancel() 1 秒后仍有 5 个 goroutine（开始前 3 个）；消费者不读了，上游还卡在 out <- v 上
```

每一级都传了 `ctx`，但没有一级去看它。上游 goroutine 卡在发送上，永远不会退出，这就是 **goroutine 泄漏**。服务端每个请求泄漏两个，跑一天就是几百万个。

## 任务

只改 `pipeline/pipeline.go`，三个函数签名不变：

```go
func Generate(ctx context.Context, nums ...int) <-chan int
func Square(ctx context.Context, in <-chan int) <-chan int
func Filter(ctx context.Context, in <-chan int, keep func(int) bool) <-chan int
```

要求：

- 不取消时行为不变：值的顺序、每一级的 `close(out)` 都保持；
- 每一级的**发送**都能被 `ctx` 取消：消费者不读了，发送方不能永远卡住；
- 每一级的**接收**也能被 `ctx` 取消：上游一直不发也不关时，本级也要能退出；
- 取消后 1 秒内所有 goroutine 退出，各级输出都被关闭；
- `ctx` 一开始就已取消时，几乎不产生任何值就关闭。

## 提示（不是答案）

- 课件第三部分第 6 节的"结果或取消"骨架就是这一题的答案形状：一个 `for` 里套 `select`，一个分支是正常工作，另一个分支是 `<-ctx.Done()`，后者直接 `return`；
- `for n := range in` 只在 `in` 关闭时才结束，它听不到取消。要改成 `select` 里的 `case n, ok := <-in:`，`!ok` 表示上游关了；
- 接收到一个值以后，**发送它**又是一个可能卡住的点，所以发送也要放进 `select`；
- `return` 会触发 `defer close(out)`，下游的接收因此结束，取消会这样一级一级传下去；
- `Square` 和 `Filter` 的循环几乎一样。你可以写一个小写的公共函数让它们共用，这是可选的。

## 评分点

| 测试 | 检查内容 |
|---|---|
| `TestPipelineRunsToCompletion` | 不取消时结果正确、输出关闭；空输入立刻关闭 |
| `TestCancelStopsTheWholePipeline` | 读两个值后 `cancel()`、不再读，1 秒内 goroutine 数回落 |
| `TestCancelClosesOutputPromptly` | `cancel()` 后很快关闭输出，不会把剩下的值都算完 |
| `TestCancelUnblocksAStageWaitingOnInput` | 上游既不发也不关时，`cancel()` 也能让 `Square` 退出 |
| `TestAlreadyCanceledContext` | 已取消的 `ctx` 几乎不产生值 |

测试用 `runtime.NumGoroutine()` 观察 goroutine 有没有回落。这个方法只适合在测试里辅助观察，不要用它写业务逻辑。

## 做完以后

到这里，channel 的三件事都齐了：

| 题目 | 解决的问题 |
|---|---|
| `test03` fan-in | 多个发送方，谁来关闭公共输出 |
| `test04` 流水线 | 消费者提前离开，上游怎么停下来 |
| `test07` 任务池 | 把上面两件事和 `WaitGroup` 组合起来 |

下一题 `test05` 回到共享内存：账本的余额 map 和总额要被很多 goroutine 同时读写。那里没有"传过去就归你"的所有权可用，只能靠锁。
