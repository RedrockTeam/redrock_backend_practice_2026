# test02：安全计数器（从零实现）

对应课件：第四课 · 第一部分第 2 节（读改写不是一步）、第五部分（`sync.Mutex`、临界区）、第六部分（`sync/atomic`）、第十部分 LV1。
题型：**从零实现**。

## 背景

第三课 `test03` 里，你在模拟器上一条指令一条指令地看到了丢更新：两个线程都 `LOAD` 到旧值，各自 `ADDI`，再先后 `STORE`，最后只加了一次。第三课 `test06` 的 `Sleep` 版求和也被 `-race` 点了名。

`c.n += delta` 在 Go 里就是那三条指令。这一题把它修好，而且修两遍：一遍用锁，一遍用原子操作。

## 任务

只改 `counter/counter.go`。两个类型对外的方法完全一样：

```go
func (c *Counter) Add(delta int)
func (c *Counter) Value() int
func (c *Counter) Reset()

func (c *AtomicCounter) Add(delta int)
func (c *AtomicCounter) Value() int
func (c *AtomicCounter) Reset()
```

要求：

- **零值可用**：`var c Counter`、`var a AtomicCounter` 不需要构造函数就能直接用；
- 100 个 goroutine 各 `Add(1)` 100 次，`Value()` 必须正好是 10000；
- `Value` 和 `Reset` 也要同步，读和写同时发生时不能出错；
- `Counter` 用 `sync.Mutex`（或 `sync.RWMutex`）；
- `AtomicCounter` 只用 `sync/atomic`，**不能有锁字段**。

结构体现在是空的，字段由你自己加。

## 提示（不是答案）

- `sync.Mutex` 的零值就是一把没锁上的锁，所以只要把它作为字段，零值 `Counter` 就能用；
- 锁只包住读写 `n` 的那一两行。`defer c.mu.Unlock()` 在有 `return` 的方法里最不容易漏；
- Go 1.19 起有类型化的原子类型，比如 `atomic.Int64`，它的零值也能直接用。注意 `int` 和 `int64` 之间要转换；
- 方法接收者必须是指针 `*Counter`。值接收者会把锁和计数一起复制一份，改的是副本；
- 不要复制含锁的结构体（`c2 := *c`），`go vet` 会提示你。

## 评分点

| 测试 | 检查内容 |
|---|---|
| `TestZeroValueAndSequential` | 零值可用；正负 `delta`；`Reset` 后能继续用 |
| `TestConcurrentAddDoesNotLoseUpdates` | 100 × 100 次并发 `Add(1)`，结果正好是 10000，重复 5 轮 |
| `TestConcurrentMixedOperations` | 加、减、读同时进行，最终结果正确 |
| `TestCounterUsesMutexAndAtomicCounterUsesAtomic` | 源码契约：`Counter` 有锁字段；`AtomicCounter` 没有锁、用了 `sync/atomic` |

两个实现由同一段测试代码驱动：测试里定义了一个只含这三个方法的接口，把 `&Counter{}` 和 `&AtomicCounter{}` 都放进去。这是第二课 interface 的用法。

做完后自己跑一次：

```bash
go test -race ./lesson-04-sync/test02/counter
```

CI 不跑 `-race`，但这一题应该是干净的。

## 做完以后

现在你有两把工具，课件第八部分的决策图会告诉你什么时候用哪个：

- **单个变量**（计数、标志）：`atomic` 足够，开销最小；
- **多个变量要一起保持一致**：`atomic` 管不了，要用锁。

第五题 `test05` 的账本就是第二种情况：每个账户的余额和总额要**一起**变。光给 `total` 换成 `atomic.Int64` 修不好它。

下一题 `test03` 先换个方向：不共享内存，而是用 channel 传值。多个 channel 汇成一条时，谁来关闭输出？
