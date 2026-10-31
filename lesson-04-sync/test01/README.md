# test01：谁在等谁（运行 → 观察 → 修复）

对应课件：第四课 · 第一部分（main 返回、goroutine 被静默杀死）、第四部分（`sync.WaitGroup`）。
题型：**先运行观察，再修复**。

## 现象

第三课 `test07` 的 `Collect` 启动了 goroutine 却立刻返回，你看到的是一堆 0。这道题把同一件事换成一个“返回结果切片”的函数：`RunAll` 为每个任务启动 goroutine，然后**立刻返回**，没有等任何任务结束。

先不改代码，跑：

```bash
go test -v ./lesson-04-sync/test01/runall
```

你会看到类似：

```text
round 0: 第 1 个结果 = 0, want 2；RunAll 是不是没等 goroutine 结束就返回了？
```

再跑几次：

```bash
go test -count=10 ./lesson-04-sync/test01/runall
```

有时候某几个结果碰巧写进去了，但下一轮又不一样。这正是并发问题的典型特征：**结果依赖调度时序**。

再看竞态检测器：

```bash
go test -race ./lesson-04-sync/test01/runall
```

报告里会同时出现读和写同一个 `result[i]` 的位置。想想：读发生在 `RunAll` 返回之后，写发生在 goroutine 里，两者之间有没有建立“写一定先于读”的顺序？

## 任务

只改 `runall/runall.go`，实现：

```go
func RunAll(tasks []func() int) []int
```

要求：

- 每个任务**恰好执行一次**；
- 结果写回它原来的下标：`result[i]` 对应 `tasks[i]()`；
- `RunAll` 必须在**所有任务都结束后**才返回；
- 空输入返回长度为 0 的结果；
- **不能用 `time.Sleep` 猜时间**。测试会用 `go/parser` 读源码，发现 `time.Sleep` 直接判失败。

## 提示（不是答案）

- 想一想：`Wait` 在等一个计数器归零，那么计数器应该在**启动 goroutine 之前**还是之后加？
- `Done` 最好放在 `defer` 里：这样任务提前 return 或 panic 的路径也不会漏掉；
- 如果你更喜欢 channel，也可以让每个任务结束时发一个信号，主 goroutine 收满 `len(tasks)` 个再返回。关键是**建立明确的顺序**，不是猜。

课件里 `wg.Add` 放在 `go` 里面的反例，正是这道题要避免的：`Wait` 可能先看到 0，然后直接返回。

## 评分点

| 测试 | 检查内容 |
|---|---|
| `TestRunAllWaitsForEveryTask` | 任务有延迟时，所有结果都已写回；重复多轮 |
| `TestRunAllPreservesOrder` | 结果按下标归位，不按完成顺序追加 |
| `TestRunAllHandlesEmpty` | `nil` 和空切片都返回长度 0 |
| `TestRunAllDoesNotUseSleep` | 源码里没有 `time.Sleep` 的逃避写法 |

## 做完以后

你不是“让输出看起来对了”，而是建立了一条确定的顺序：**所有 `Done` → `Wait` 返回**。这条顺序让 `result` 的写在读之前完成，`-race` 也就不再报告。

下一题 `test02` 处理第三课留下的另一类问题：`test03` 的交错模拟和 `test06` 的 `-race` 报告都指向“读改写不是一步”。先用 `Mutex` 和 `atomic` 写一个不丢更新的计数器，再到 `test03` 看 channel 的关闭。
