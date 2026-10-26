# test02：把求和拆给多个 goroutine

对应课件：第三课 · 协程、go 关键字的使用。
题型：**从顺序版本改写为并发版本**。

## 现象

`parallel/sum.go` 现在是顺序求和：

```go
for _, value := range values {
    total += value
}
```

结果正确，但没有用到 goroutine。跑：

```bash
go test -v ./lesson-03-goroutines/test02/parallel
```

`TestSum` 会通过，`TestSumUsesGoroutines` 会失败：

```text
sum.go 里没有 go 语句
```

## 任务

只改 `parallel/sum.go`：

```go
func Sum(values []int) int
```

要求：

- 把 `values` 拆成至少两段；
- 每段交给一个 goroutine 计算；
- 把各段结果加起来，返回总和；
- 空切片返回 0，不修改输入；
- 第三课允许用 `time.Sleep(100 * time.Millisecond)` 等 goroutine 跑完；第四课会换成 `WaitGroup`。

## 提示（不是答案）

- 先算中点：`middle := len(values) / 2`；
- 用一个长度固定的小切片存各段结果，例如 `parts := make([]int, 2)`；
- 每个 goroutine 只写自己那一格：`parts[0] = sum(values[:middle])`、`parts[1] = sum(values[middle:])`。写不同下标不会竞争；
- 最后 `return parts[0] + parts[1]`；
- 如果两个 goroutine 都去改同一个 `total`，就会产生数据竞争，结果也可能丢更新。

## 评分点

| 测试 | 检查内容 |
|---|---|
| `TestSum` | 正数、负数、单个、空、`nil`、1000 个元素 |
| `TestSumDoesNotModifyInput` | 求和过程不修改输入切片 |
| `TestSumUsesGoroutines` | 源码里确实有 `go` 语句 |

## Sleep 的代价

如果你用 `time.Sleep` 等两个 goroutine，功能测试能过，但：

```bash
go test -race ./lesson-03-goroutines/test02/parallel
```

通常会报告数据竞争：`Sleep` 没有在“goroutine 写 parts”和“主 goroutine 读 parts”之间建立 happens-before 关系。它只是让读通常发生在写之后，并没有消除竞争。

这不是让你现在修——第四课的 `WaitGroup` 才会给出正确写法。

## 做完以后

你已经有“拆分 → 并行 → 汇总”的形状了。但 `Sleep` 仍然是在猜时间：

- 如果某一段算得比 100ms 慢，结果会少；
- 如果机器很慢，测试就会偶发失败；
- 真实业务里不能靠“睡一会儿”保证正确性。

下一题 `test03` 把这个问题放到一个**必须返回结果**的函数里，让你稳定地看到“没等”的后果。
