# test01：让 goroutine 真的跑起来

对应课件：第三课 · 协程、go 关键字的使用（匿名函数、参数传入、函数调用对比）。
题型：**运行观察 → 补全代码**。

## 现象

`asyncprint/run.go` 现在用普通 `for` 循环顺序打印：

```text
worker 0
worker 1
...
```

输出是对的，但它一个 goroutine 都没启动。先跑：

```bash
go test -v ./lesson-03-goroutines/test01/asyncprint
```

`TestRunStartsEveryWorker` 会通过，`TestRunUsesGoroutines` 会失败：

```text
run.go 里没有 go 语句
```

这正是这一题要改的地方：用 `go` 启动 worker，而不是顺序调用。

## 任务

只改 `asyncprint/run.go`：

```go
func Run(out io.Writer, workers int)
```

要求：

- 为每个 worker 启动一个 goroutine，向 `out` 写一行 `worker N`；
- 行尾有 `\n`，每个 worker 恰好一行；
- 输出顺序不要求固定；
- 第三课还没有 `WaitGroup`，函数末尾可以用 `time.Sleep(100 * time.Millisecond)` 等它们跑完；
- 空 `workers` 不输出任何内容。

## 提示（不是答案）

```go
for i := 0; i < workers; i++ {
    i := i // 显式把当前值交给闭包
    go func() {
        // 在这里把 worker i 写到 out；具体用 fmt.Fprintf 还是别的，自己决定
    }()
}
// 函数末尾还要让 Run 不要立刻返回：
// 第三课先用 time.Sleep 临时等待，第四课会换成 WaitGroup
```

- `go` 后面必须是**函数调用**：`go f()`、`go func(){}()`；
- 循环变量最好作为参数传进闭包（或者先 `i := i`），不要依赖“等会儿它应该还是这个值”；
- 多个 goroutine 同时写 `out`，测试用了一个带锁的 buffer 接住；真实代码里也要想清楚 `out` 是否并发安全。

## 评分点

| 测试 | 检查内容 |
|---|---|
| `TestRunStartsEveryWorker` | 5 个 worker 每个恰好一行，内容正确 |
| `TestRunHandlesZeroWorkers` | `workers = 0` 不输出 |
| `TestRunUsesGoroutines` | 源码里确实有 `go` 语句 |

## 做完以后

你已经让多个执行流真的跑起来了。注意 `Run` 末尾那个 `Sleep`：它只是“猜它们大概跑完了”。

下一题把一份求和工作拆给两个 goroutine；再下一题 `test03` 会让你看到：当函数必须**返回结果**时，`Sleep` 为什么不够用。
