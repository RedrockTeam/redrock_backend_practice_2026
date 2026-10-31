# test02：写一台玩具 CPU

对应课件：第三课 · cpu 工作简介（PC、寄存器、“取指 → 译码 → 执行 → 写回 → 更新 PC”）。
题型：**从零实现（中等）**。综合用到第一课的函数与循环、第二课的结构体、方法、切片、`error`。

## 背景

上一题你写了“取指 + 译码”。这一题补上剩下的三步，做一台真的能跑程序的 CPU。为了把注意力放在**循环本身**上，指令已经译好码，存成 `[]Instr`，不用再拆二进制。

`toycpu/toycpu.go` 里只给了类型、常量、错误和三个函数签名，函数体全部要你从零写。

## 规格

### 机器

| 部件 | 类型 | 说明 |
|---|---|---|
| 寄存器 | `Regs [8]int` | `R0`..`R7`，全是普通寄存器 |
| 内存 | `Mem []int` | 下标就是地址 |
| PC | `PC int` | **下一条要执行的指令在 `Program` 里的下标** |
| 程序 | `Program []Instr` | 和数据分开放 |
| 停机标志 | `Halted bool` | 执行过 `HALT` 之后为 `true` |

### 指令

| 操作码 | 效果 | 执行后 PC |
|---|---|---|
| `HALT` | `Halted = true` | **不变** |
| `LOADI` | `R[A] = Imm` | `+1` |
| `LOAD` | `R[A] = Mem[Imm]` | `+1` |
| `STORE` | `Mem[Imm] = R[A]` | `+1` |
| `ADD` | `R[A] = R[B] + R[C]` | `+1` |
| `ADDI` | `R[A] = R[B] + Imm` | `+1` |
| `JNZ` | 无 | `R[A] != 0` 时跳到 `Imm`，否则 `+1` |

只检查**这条指令用到的**字段：比如 `LOADI` 不管 `B`、`C` 是多少。

### `New(program, memSize)`

返回一台新 CPU：`Program` 装好，`Mem` 是 `memSize` 个 0，`PC = 0`，没有停机。

### `Step()`：执行恰好一条

1. 已经停机：什么都不做，返回 `nil`；
2. **取指**：`PC` 不在 `[0, len(Program))` 里 → `ErrPCOutOfRange`；
3. **译码**：操作码不认识（包括零值 `Instr{}`）→ `ErrBadOp`；用到的寄存器编号不在 `[0, 8)` → `ErrBadRegister`；用到的内存地址不在 `[0, len(Mem))` → `ErrBadAddress`；
4. **执行 + 写回**：按上表改寄存器或内存；
5. **更新 PC**。

**出错的那一条不能改动任何状态**（寄存器、内存、PC、停机标志都不变）。所以要先检查、再动手。

### `Run(maxSteps)`

反复 `Step`，返回**成功执行的条数**和错误：

- 执行了 `HALT` → 返回 `(条数, nil)`，`HALT` 本身算一条；
- `Step` 出错 → 返回 `(出错前的条数, 那个错误)`；
- 已经执行满 `maxSteps` 条还没停机 → 返回 `(maxSteps, ErrStepLimit)`；恰好第 `maxSteps` 条是 `HALT` 不算超限；
- 调用时已经停机 → `(0, nil)`。

## 提示（不是答案）

- `Step` 的骨架就是一个 `switch in.Op { case HALT: ... case LOADI: ... default: return ErrBadOp }`；
- 寄存器和地址检查会重复很多次，写两个小写的辅助函数，比如 `validReg(r int) bool`、`c.validAddr(a int) bool`；
- “先检查后动手”最简单的写法：每个 `case` 里先 `if !valid... { return Err... }`，全部通过再赋值；
- `Run` 是一个 `for` 循环：每轮先看是否停机，再看是否已满 `maxSteps`，再 `Step`；
- 测试用 `errors.Is(err, toycpu.ErrBadOp)` 判断，直接 `return ErrBadOp` 就行；想加上下文可以用 `fmt.Errorf("pc %d: %w", c.PC, ErrBadOp)`，`%w` 会保留原错误。

## 自己验证

```bash
go test -v ./lesson-03-goroutines/test02/toycpu
go test -run TestStepEachInstruction ./lesson-03-goroutines/test02/toycpu
```

建议顺序：`TestNew` → `TestStepEachInstruction` → `TestHalt` → `TestRunSumProgram` → 其余。

## 评分点

| 测试 | 检查内容 |
|---|---|
| `TestNew` | 程序、内存、寄存器、PC 的初始状态 |
| `TestStepEachInstruction` | 七种指令各自的效果与 PC 变化 |
| `TestHalt` | 停机时 PC 不动；停机后再 `Step` 什么也不做 |
| `TestRunSumProgram` | 跑一个循环求和程序，结果与执行条数都对 |
| `TestRunStopsAtStepLimit` | 死循环被 `maxSteps` 截住；恰好停机不算超限 |
| `TestRunReportsErrors` | 五种错误都能报出来，条数不含出错那条 |
| `TestErrorLeavesStateUnchanged` | 出错的指令不改动任何状态 |
| `TestBridgeIncrementIsThreeInstructions` | 一句 `x++` 在这台 CPU 上是三条指令 |

## 做完以后

`go test -v -run TestBridge` 会打印 `x++` 的三步：读到寄存器、在寄存器里加 1、写回内存。第 1、2 步之后，新值**只在寄存器里**。

真实 CPU 上，操作系统随时可能在两条指令之间把这个执行流换下去，换另一个线程上来。如果那个线程也在做 `x++` 呢？下一题 `test03` 让两个“线程”共用一块内存，你会亲手数出有多少种交错会把一次加法弄丢。
