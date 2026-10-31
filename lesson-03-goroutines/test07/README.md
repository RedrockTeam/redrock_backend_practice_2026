# test07：函数返回了，goroutine 还没写完

对应课件：第三课 · 协程（进程 / 线程 / 协程）、协程的 infra（GMP）、总结（“你会看到一些目前没有办法解决的问题”）。
题型：**观察题 + 写量尺**。这一题是第三课的出口，也是第四课 `test01` 的入口。

> `collect.go` 里的 `Collect` 是**观察对象**，这一题不改它。
> 你要做的是：写一个量尺把问题量出来，然后回答几道观察 / 概念题。

## 现象

`missing/collect.go` 为每个元素启动一个 goroutine，把平方写进 `results[i]`，然后**立刻返回**：

```go
for i, value := range values {
    go func() {
        results[i] = value * value
    }()
}
return results
```

`results` 已经分配好了，长度也对，但很多位置还是零值。**不是 append 丢了，是写还没来得及发生。**

## 任务

### 1. 写量尺：`verify.go`

```go
func Verify(values, got []int) int
```

返回 `got` 里有多少个位置**不等于** `values[i] * values[i]`：

- `got` 比 `values` 短时，缺的那些位置都算不对；
- `got` 比 `values` 长时，多出来的部分不管。

写完后跑桥接测试，看 `Collect` 每一轮少了多少个：

```bash
go test -v -run TestBridgeObserveCollect ./lesson-03-goroutines/test07/missing
go test -v -count=3 -run TestBridgeObserveCollect ./lesson-03-goroutines/test07/missing
```

### 2. 回答观察 / 概念题：`answers.go`

把 `Q1`..`Q6` 填成 `"A"`、`"B"` 或 `"C"`。测试只告诉你哪一题不对，不告诉你正确选项——请回到课件找依据，而不是逐个试。

**Q1（观察）** `Collect` 返回的那一刻，它启动的 200 个 goroutine：

- A. 一定都已经写完了，只是写的值不对
- B. 一定一个都还没开始执行
- C. 不确定：有的可能写完了，有的可能还没开始，取决于调度

**Q2（观察）** 在 `return results` 前加一句 `time.Sleep(100 * time.Millisecond)`：

- A. 彻底解决了，因为 100ms 对 200 个乘法绰绰有余
- B. 测试大概率会变绿，但它只是让写“很可能”先发生，并没有建立“写完才返回”的保证
- C. 会让结果更糟，因为 Sleep 会打断正在运行的 goroutine

**Q3（概念）** `main` 函数返回时，还有别的 goroutine 没跑完：

- A. 进程退出，这些 goroutine 被直接结束，不会有任何提示
- B. runtime 会等所有 goroutine 结束后再退出进程
- C. 这些 goroutine 会脱离进程，在后台继续运行

**Q4（概念）** 同一个 Go 进程里的两个 goroutine：

- A. 各有一份独立的地址空间，要共享数据必须走 IPC
- B. 共享同一个栈，所以局部变量也是共享的
- C. 共享同一个地址空间（堆、全局变量），但各有自己的栈

**Q5（概念）** 谁决定某个 goroutine 此刻在哪个 OS 线程上运行？

- A. 操作系统内核直接调度每一个 goroutine
- B. Go runtime 的调度器：G 通过 P 被放到 M（OS 线程）上执行
- C. 程序员在 `go` 语句里指定

**Q6（概念，回顾 test03）** 两个 goroutine 各对同一个 `int` 做 1000 次 `x++`，结果常常小于 2000：

- A. 因为 `x++` 在 CPU 上是读、改、写三步，两步之间可能被另一个执行流插进来，后写的覆盖先写的
- B. 因为 goroutine 太多，runtime 会随机丢弃一部分
- C. 因为 `int` 在并发时会溢出

## 提示（不是答案）

- `Verify` 就是一个 `for i := range values` 的循环，越界前先判断 `i < len(got)`；
- Q1：看 `go test -count=3` 的输出，每一轮的数字一样吗？
- Q2：Sleep 是在“建立顺序”，还是在“猜时间”？机器变慢、数据变大之后呢？
- Q3 / Q4 / Q5：课件“协程”与“协程的 infra”两节；
- Q6：`test03` 的 `TestOutcomesOfTwoIncrements`。

## 评分点

| 测试 | 检查内容 |
|---|---|
| `TestVerifyCountsWrongPositions` | `Verify` 的九种输入 |
| `TestObservationAnswers` | Q1–Q3 |
| `TestConceptAnswers` | Q4–Q6 |
| `TestBridgeObserveCollect` | 用你的 `Verify` 打印每轮缺失数（不要求 `Collect` 正确） |
| `TestCollectHandlesEmpty` | 空输入返回空切片 |
| `TestCollectDoesNotUseSleep` | `collect.go` 里没有被加上 `time.Sleep` |

## 第四课怎么接

你现在能**量出**问题，但还没法**修**它：第三课的工具里，没有一样能让调用方知道“所有 goroutine 都写完了”。

第四课 `test01` 的第一句话就是这个问题：**启动了一组 goroutine，怎么知道它们什么时候全部结束？** 答案是 `sync.WaitGroup`：

```text
Add：启动前记下“有几个人”
Done：每个人结束时报告“我完成了”
Wait：等到计数归零才返回
```

到那时，你在这里写的 `Verify` 量出来的数字应该稳定是 0。

而 Q6 那种“读-改-写被插队”的问题，第四课 `test02` 会用 `sync.Mutex` / `atomic` 解决。
