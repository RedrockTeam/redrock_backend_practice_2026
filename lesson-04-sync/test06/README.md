# test06：懒加载与超时（代码纠错）

对应课件：第四课 · 第三部分第 4〜5 节（超时、`context`）、第二部分（无缓冲与有缓冲 channel）、第六部分（`sync.Once`）。
题型：**代码纠错**：给你一段"看起来能用"的代码，找出并修好里面的三个 bug。

## 背景

`lazyfetch/lazyfetch.go` 是后端里常见的客户端写法：

- 配置很贵（读文件、连配置中心），所以**第一次用到时才加载**；
- 下游服务很慢，所以每次请求都放进后台 goroutine，调用方用 `ctx` 设超时。

单个 goroutine、服务又很快的时候，它完全正常。这一题要你在并发和超时的情况下找出它的问题。

## 先跑

```bash
go test -v ./lesson-04-sync/test06/lazyfetch
```

`TestFetchReturnsValue` 会通过。另外三个各指向一个 bug，先读失败信息，再回到代码里找对应的那几行：

| 测试 | 它在暗示什么 |
|---|---|
| `TestConfigLoadsOnceUnderConcurrency` | 100 个 goroutine 同时第一次调用 `Config()`，`load` 被执行了很多次 |
| `TestFetchTimesOutWithWrappedError` | 超时的时候，错误信息里看不出是哪一次请求 |
| `TestTimedOutFetchDoesNotLeakGoroutines` | 请求都超时返回了，慢服务也终于回了结果，后台 goroutine 却一个都没退出 |

还可以用竞态检测器看第一个 bug：

```bash
go test -race -run TestConfigLoadsOnce ./lesson-04-sync/test06/lazyfetch
```

## 任务

只改 `lazyfetch/lazyfetch.go`，`New`、`Config`、`Fetch` 的签名不变。

1. **懒加载只做一次**：多个 goroutine 同时第一次调用 `Config()` 时，`load` 只执行一次，所有人拿到同一份配置；
2. **超时错误要能追查**：`ctx` 先结束时返回 `""` 和一个错误，这个错误
   - 信息里包含这次请求的 `key`；
   - 满足 `errors.Is(err, ctx.Err())`，也就是超时的时候 `errors.Is(err, context.DeadlineExceeded)` 为真；
3. **不泄漏**：`Fetch` 超时返回以后，后台 goroutine 在慢服务回结果时必须能退出。

## 提示（不是答案）

**bug 1。** `if c.cfg == nil { ... c.cfg = ... }` 是"先检查、再行动"。检查和赋值之间隔着一个耗时 20ms 的 `load`，这段时间里所有 goroutine 看到的都是 `nil`。这和第三课 `test03` 的交错是同一类问题，只是被放大了。课件第六部分有一个工具专门做"并发安全地只执行一次"。用它的话，`cfg` 也就不需要再是指针了。

**bug 2。** 第二课 `test06` 里你已经写过返回 `error` 的接口。要给错误加上下文、又不丢掉原来的错误，用 `fmt.Errorf` 和哪个格式动词？

**bug 3。** 想清楚这条 channel 的一生：

```text
后台 goroutine:  result <- value    ← 无缓冲：要等有人接才能完成
Fetch:           select 选中 ctx.Done() 后 return，再也没人读 result
```

发送方等一个永远不会出现的接收方，于是永远卡住。课件第二部分讲过两种 channel，其中一种的发送**不需要**接收方同时在场。容量要多大，才能保证这一次发送一定不会阻塞？

修 bug 3 时不要关闭 `result`，也不要在 `Fetch` 里再开一个 goroutine 去读它。

## 评分点

| 测试 | 检查内容 |
|---|---|
| `TestFetchReturnsValue` | 正常路径：返回 `endpoint/key`，`load` 只调用一次 |
| `TestConfigLoadsOnceUnderConcurrency` | 100 个 goroutine 并发首次调用，`load` 恰好执行 1 次，重复 3 轮 |
| `TestFetchTimesOutWithWrappedError` | 50ms 超时及时返回；错误带 `key`，并且 `errors.Is(err, context.DeadlineExceeded)` |
| `TestTimedOutFetchDoesNotLeakGoroutines` | 50 次超时请求之后放行慢服务，2 秒内 goroutine 数回落 |

## 做完以后

这一题里的三个 bug 都不会在"单个用户点一下"的时候出现，只会在并发、慢下游、超时同时出现时冒出来，也就是线上真正的样子。

- `sync.Once` 只保证**创建一次**。创建出来的东西之后要是还会被改，仍然要加锁；
- 只要 goroutine 的发送可能没人接收，就要问一句：谁来接？接不到的时候它怎么退出？
- 第五课的 HTTP 客户端每个请求都带 `ctx`，超时错误也是这样一层层用 `%w` 包上来的。

下一题 `test07` 把本课的工具全部组合起来：固定数量的 worker、可取消、结果按下标归位的任务池。
