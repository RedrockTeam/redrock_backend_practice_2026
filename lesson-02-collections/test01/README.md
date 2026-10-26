# test01：切片是窗口，不是数据

对应课件：第二课 · 指针和数组、切片（`make` / `copy` / `append` / `range`）。
题型：**运行 → 观察 → 修复**。

## 现象

`sliceops/sliceops.go` 里有四个函数，其中三个藏着典型的切片坑。先跑：

```bash
go test -v ./lesson-02-collections/test01/sliceops
```

你会看到 `Clone` 改副本把原件也改了、`DoubleInPlace` 原地乘 2 没生效、`AppendOne` 返回的长度没变。分别对应：

| 现象 | 原因 |
|---|---|
| `clone[0] = 100` 之后 `base[0]` 也变了 | `return values` 只复制了切片头，底层数组共享 |
| `range` 里 `value *= 2` 原切片没变 | `v` 是元素副本 |
| `append(values, x)` 之后长度没变 | 丢了 `append` 的返回值 |

## 任务

只改 `sliceops/sliceops.go`：

| 函数 | 要求 |
|---|---|
| `Clone(values []int) []int` | 返回独立副本；修改副本不影响原切片；`nil` 仍返回 `nil` |
| `DoubleInPlace(values []int)` | 把每个元素乘 2，**原地修改** |
| `AppendOne(values []int, x int) []int` | 返回追加 `x` 后的新切片，**接住 `append` 的返回值** |
| `Sum(values []int) int` | 求和；空切片返回 0 |

## 提示（不是答案）

- 复制切片：`make([]int, len(values))` + `copy`，或者 `append([]int(nil), values...)`；
- 原地修改必须拿到下标：`for i := range values { values[i] *= 2 }`；
- `append` 返回的是**新的切片头**，要 `return append(values, x)`；
- `nil` 切片可以 `range`、可以 `len`，直接返回 `nil` 也没问题。

## 评分点

| 测试 | 检查内容 |
|---|---|
| `TestCloneIsIndependent` | 改副本不影响原件；`Clone(nil)` 返回 `nil` |
| `TestDoubleInPlace` | 原地修改生效；`nil` 调用不 panic |
| `TestAppendOne` | 返回长度正确；原切片可见部分不变；`nil` 可追加 |
| `TestSum` | 正数、负数、空、`nil` |
| `TestBridgeSliceAssignmentIsAWindow` | 观察 `b := a` 共享底层数组（不检查你的实现） |

## 做完以后

你已经亲手区分了“切片头”和“底层数组”。下一题 `test02` 的 `map` 也是句柄：函数里改 map 内容，外部能看到；但 map 不能像切片一样用 `append`，也不能直接比较。

> 用 `-race` 跑这一题不会报错，因为这里还没有并发；第三课才会把共享状态放到多个 goroutine 下。
