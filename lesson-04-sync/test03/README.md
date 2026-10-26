# test03：共享账本（运行 → 观察 → 修复）

对应课件：第四课 · 第一部分（计数丢失、并发写 map）、第五部分（`Mutex` / `RWMutex`、临界区）。
题型：**运行观察现象，再修复**。

## 现象

`ledger/ledger.go` 里的 `Ledger` 没有任何同步。它有三个问题，而且在不同的运行方式下会露出不同的脸：

```bash
# 1. 顺序测试：先看 Snapshot 的所有权问题
go test -v ./lesson-04-sync/test03/ledger

# 2. 重复跑：让并发问题更容易出现
go test -count=10 ./lesson-04-sync/test03/ledger

# 3. 竞态检测：看两处访问分别在哪
go test -race ./lesson-04-sync/test03/ledger
```

可能看到：

```text
修改 snapshot 后 Balance("alice") = 999, want 10；Snapshot 不能把内部 map 直接交出去
```

也可能直接看到：

```text
fatal error: concurrent map writes
```

或者：

```text
WARNING: DATA RACE
```

这三个现象分别对应：

| 现象 | 本质 |
|---|---|
| snapshot 被改，账本也变了 | `Snapshot` 把内部 map 的别名交了出去 |
| `fatal error: concurrent map writes` | 多个 goroutine 同时写同一个 map，runtime 直接终止进程 |
| 总额少了 / 数据竞争报告 | `total += amount` 是读-改-写三步，并发时可能丢更新 |

## 任务

只改 `ledger/ledger.go`，让 `Ledger` 在并发下仍然正确：

```go
func (l *Ledger) Deposit(account string, amount int)
func (l *Ledger) Balance(account string) int
func (l *Ledger) Total() int
func (l *Ledger) Snapshot() (map[string]int, int)
```

要求：

- 零值 `var l Ledger` 可以直接使用；
- 多个 goroutine 同时 `Deposit` 时，不能丢更新；
- 读也要保护：`Balance`、`Total`、`Snapshot` 不能读到写了一半的状态；
- `Snapshot` 返回的 map 是**副本**：调用方改它不能影响账本；
- 空账本的 `Snapshot` 也要返回一个可写的空 map，而不是 `nil`；
- 临界区尽量只包住共享状态的读写，不要把无关的计算或 I/O 放进去。

## 提示（不是答案）

- `sync.RWMutex` 允许多个读者同时进入；写者独占。读多写少时很合适，但**先写对，再考虑优化**；
- 不要只锁 `Deposit`：`Balance`、`Total`、`Snapshot` 也要拿锁，否则它们读到的可能是“半个新、半个旧”的状态；
- `Snapshot` 里先在锁内把 map 复制出来，再返回；复制完就可以解锁；
- `map` 的零值是 `nil`，不能直接写；可以在第一次 `Deposit` 时再初始化；
- 不要复制含锁的结构体，也不要在持锁时做耗时操作。

## 评分点

| 测试 | 检查内容 |
|---|---|
| `TestZeroValueLedger` | 零值可用；初始余额和总额为 0；能正常存款 |
| `TestSequentialDeposits` | 多账户、多次存款的余额与总额 |
| `TestSnapshotIsACopy` | 修改返回的 map 不影响账本；空账本返回可写 map |
| `TestConcurrentUseKeepsInvariants` | 多写多读下总额准确；每份快照满足“账户之和 == total”；重复多轮 |

## 做完以后

这道题真正要建立的是一个习惯：**看到共享变量，先问“谁会在什么时候读写它”，再决定锁的范围**。

- 只锁写、不锁读，等于没同步；
- 锁住整个函数，可能把并发退化成串行；
- 只保护 map、不保护 total，会得到“map 对、总额错”的怪状态；
- `Snapshot` 不复制，调用方就成了“没有拿锁的写者”。

下一题 `test04` 把这些工具组合起来：用 channel 发任务、用 `WaitGroup` 等 worker、用 `context` 取消、用切片下标保证结果顺序。
