# test03：两个线程，一个 x++，丢了一次

对应课件：第三课 · 协程（进程 / 线程：线程共享地址空间，但各有自己的 PC、寄存器和栈）、cpu 工作简介。
题型：**实现模拟器 + 观察（较难）**。

## 背景

上一题 `test02` 的桥接测试说明了一件事：`x++` 在 CPU 上是三条指令。

```text
LOAD  Reg, [x]     读
ADDI  Reg, 1       改
STORE [x], Reg     写
```

课件说线程是“进程内部的执行流”：**同一个进程里的线程共享内存，但每个线程有自己私有的 PC 和寄存器**。操作系统随时可以在任意两条指令之间换线程。

这一题不用真的线程，而是用一个**确定的调度表**模拟“操作系统决定下一步让谁跑”。这样同一种交错可以重放任意多次，你能数清楚到底有多少种交错会出错。

## 任务

编辑 `interleave/interleave.go`，完成全部 `TODO`。

### 机器

- 每个线程有自己的 `PC`（从 0 开始）和**一个**寄存器 `Reg`（初始为 0）；
- 所有线程共用同一块 `memory []int`。

| 操作码 | 效果 |
|---|---|
| `LOAD` | `Reg = Mem[Addr]` |
| `ADDI` | `Reg = Reg + Imm` |
| `STORE` | `Mem[Addr] = Reg` |
| `INCR` | `Mem[Addr] = Mem[Addr] + Imm`，**一步完成** |

### 函数

| 函数 | 要求 |
|---|---|
| `Increment(addr)` | 返回三条指令：`LOAD addr`、`ADDI 1`、`STORE addr` |
| `AtomicIncrement(addr)` | 返回一条指令：`INCR addr, 1` |
| `Run(programs, memory, schedule)` | 按调度表交错执行，返回最终内存的**副本**；不修改 `memory` |
| `Schedules(lengths)` | 按字典序列出所有合法调度表 |
| `Outcomes(programs, memory, addr)` | 每种合法调度各跑一次，统计 `Mem[addr]` 的最终值：值 → 次数 |

`schedule[k] = i` 的意思是“第 k 步，让线程 i 执行它的下一条指令”。例如两个线程各做一次 `x++`：

```text
schedule = [0 0 0 1 1 1]      线程 0 做完，线程 1 再做     → x = 2
schedule = [0 1 0 1 0 1]      两个线程都先读到 0          → x = 1
```

`Run` 遇到下面的情况返回对应错误（`errors.Is` 判断）：

| 情况 | 错误 |
|---|---|
| 调度表里的线程编号不存在 | `ErrBadThread` |
| 调度了一个已经执行完的线程 | `ErrThreadDone` |
| 调度表走完了，还有线程没执行完 | `ErrUnfinished` |
| `Addr` 不在内存范围内 | `ErrBadAddress` |
| 操作码不认识 | `ErrBadOp` |

“合法调度表”：线程 `i` 恰好出现 `lengths[i]` 次。两个线程各三条指令时，就是从 6 个位置里挑 3 个给线程 0，一共 20 种。

## 提示（不是答案）

- `Run` 里先 `make` + `copy` 出一份内存，后面只改这一份；
- 每个线程的状态用一个小结构体存：`type thread struct { pc, reg int }`，再开一个 `[]thread`；
- `Schedules` 用递归最自然：手里有一个“每个线程还剩几条”的切片，每一步按 `0, 1, 2…` 的顺序试着挑一个还有剩余的线程，挑完递归下去，回来再恢复——按这个顺序挑，结果天然就是字典序；
- 递归时把当前前缀存进结果之前**要复制一份**，否则后面的修改会改到已经存进去的切片（第二课 test02 的窗口问题）；
- `Outcomes` = `Schedules` + 对每个调度调用 `Run` + 一个 `map[int]int` 计数。

## 自己验证

```bash
go test -v ./lesson-03-goroutines/test03/interleave
go test -v -run TestOutcomes ./lesson-03-goroutines/test03/interleave   # 看分布
```

## 评分点

| 测试 | 检查内容 |
|---|---|
| `TestIncrementPrograms` | 两种加法的指令序列 |
| `TestRunSerialSchedules` | 串行调度得到 2 |
| `TestRunInterleavedScheduleLosesAnUpdate` | 三种交错都只得到 1 |
| `TestRunKeepsThreadsPrivateAndInputIntact` | 寄存器各自私有；不修改传入的内存 |
| `TestRunErrors` | 五种错误 |
| `TestSchedules` | 字典序、数量（20、1680）、每种调度合法且不重复 |
| `TestOutcomesOfTwoIncrements` | 两个线程：20 种调度的结果分布 |
| `TestOutcomesOfThreeIncrements` | 三个线程：1680 种调度的结果分布 |
| `TestBridgeAtomicIncrementNeverLoses` | 换成“一步完成”的加法，所有调度都正确 |

## 做完以后

跑 `go test -v -run TestOutcomes`，看看那几个数字：两个线程各加一次，20 种交错里只有 2 种是对的；三个线程时，1680 种里只有 6 种是对的。

这不是“运气差才会碰到”，而是**大多数交错都是错的**。真实程序里你没法控制调度表，所以：

- 读-改-写必须在别人眼里像“一步”——这就是 `TestBridgeAtomicIncrementNeverLoses` 里的 `INCR`；
- 第四课 `test02` 的 `sync.Mutex` 和 `atomic` 就是在真实 Go 程序里做到这一点的工具。

接下来 `test04` 回到真正的 Go：用 `go` 关键字启动成千上万个 goroutine，看看它们到底有多轻。
