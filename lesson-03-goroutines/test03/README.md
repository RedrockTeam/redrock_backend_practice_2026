# test03：函数返回了，goroutine 还没写完

对应课件：第三课 · 总结（“你会看到一些目前没有办法解决的问题”）。
题型：**观察题 / 第四课桥接题**。

> **这一题在第三课结束时仍然是红的，这是设计好的。**
> 不要用 `time.Sleep` 硬修；第四课 `test01` 会用 `WaitGroup` 补上。

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

跑：

```bash
go test -v ./lesson-03-goroutines/test03/missing
```

你会看到类似：

```text
第 0 个结果 = 0, want 1；Collect 在 goroutine 写完之前就返回了
```

`results` 已经分配好了，长度也对，但大多数位置还是零值。**不是 append 丢了，是写还没来得及发生。**

## 为什么 Sleep 不算修好

你可以在 `Collect` 里加 `time.Sleep`，让测试碰巧变绿。但：

- 调用方拿到结果的时间仍然取决于“睡多久”；
- 数据量变大、机器变慢，Sleep 就不够；
- 这不是建立顺序，是猜顺序。

所以测试里有一道源码契约 `TestCollectDoesNotUseSleep`：实现文件里出现 `time.Sleep` 直接判失败。

## 任务（第三课只要求观察）

1. 先用 `go test -count=5` 看它是不是稳定地缺结果；
2. 想一想：**调用方**怎么知道所有 goroutine 都写完了？
3. 不要改实现；把这个问题带到第四课。

## 评分点

| 测试 | 检查内容 |
|---|---|
| `TestCollectWaitsForEveryResult` | 200 个元素、重复 5 轮，每个平方都要正确 |
| `TestCollectHandlesEmpty` | 空输入返回空切片 |
| `TestCollectDoesNotUseSleep` | 源码里没有 `time.Sleep` 的逃避写法 |

## 第四课怎么接

第四课 `test01` 的第一句话就是这个问题：**启动了一组 goroutine，怎么知道它们什么时候全部结束？**

答案不是 Sleep，而是 `sync.WaitGroup`：

```text
Add：启动前记下“有几个人”
Done：每个人结束时报告“我完成了”
Wait：等到计数归零才返回
```

`Collect` 在第四课会有一个用 `WaitGroup` 的版本；那时 `results[i]` 的写在 `return` 之前就完成了，测试才会稳定变绿。
