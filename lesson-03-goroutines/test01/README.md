# test01：亲手拆一条 RISC-V 指令

对应课件：第三课 · 计算机基础知识介绍 · cpu 工作简介（编译 → 汇编 → 机器码 → 小端序 → 取指与译码）。
题型：**补全代码（入门）**。

## 背景

课件里那张大图把 `int a = b + c;` 一路追到了 CPU：

```text
add x8, x5, x7  →  0x00728433  →  内存里是 33 84 72 00  →  CPU 取出来、拆字段、执行
```

图是看懂了，但“拆字段”到底怎么拆？这一题让你用 Go 把它写出来。你写的就是 CPU 译码器在做的事。

## 先补一个新语法：位运算

第一课没讲位运算，这一题必须用。只有四个：

| 运算 | 含义 | 例子（二进制） |
|---|---|---|
| `x >> n` | 右移 n 位，低位丢掉 | `0b1101_0000 >> 4 == 0b1101` |
| `x << n` | 左移 n 位，低位补 0 | `0b1101 << 4 == 0b1101_0000` |
| `x & m` | 按位与：只留下 m 里为 1 的那些位 | `0b1101_0110 & 0b1111 == 0b0110` |
| `x \| m` | 按位或：把两个数的 1 合到一起 | `0b1100_0000 \| 0b0011 == 0b1100_0011` |

两个固定套路：

```go
field := (word >> start) & mask   // 取出从第 start 位开始的一段；mask 是“宽度个 1”
word  |= (field & mask) << start  // 把一段放回去
```

“宽度个 1”：3 位是 `0x7`（`0b111`），5 位是 `0x1f`，7 位是 `0x7f`。

还有一个类型的坑：`byte` 只有 8 位，`byte(0x33) << 8` 会直接溢出成 0。**先 `uint32(b[1])`，再左移。**

## 任务

编辑 `riscv/riscv.go`，完成全部 `TODO`：

| 函数 | 要求 |
|---|---|
| `LittleEndian(b [4]byte) uint32` | 小端序：`b[0]` 是最低 8 位，`b[3]` 是最高 8 位 |
| `Decode(word uint32) Instr` | 按下面的 R 型格式拆出六个字段 |
| `Encode(in Instr) uint32` | `Decode` 的逆运算；字段超出位宽时只保留低位 |
| `Disasm(in Instr) string` | 认出 `add` 和 `sub`，其它一律 `unknown 0x%08x` |

R 型格式（课件原图）：

```text
31      25 24  20 19  15 14  12 11   7 6      0
+---------+------+------+------+------+--------+
| funct7  | rs2  | rs1  |funct3|  rd  | opcode |
+---------+------+------+------+------+--------+
   7 位     5 位   5 位   3 位   5 位    7 位
```

`Disasm` 的规则：

| opcode | funct3 | funct7 | 输出 |
|---|---|---|---|
| `0x33` | `0` | `0x00` | `add x<rd>, x<rs1>, x<rs2>` |
| `0x33` | `0` | `0x20` | `sub x<rd>, x<rs1>, x<rs2>` |
| 其它 | | | `unknown 0x` + 8 位小写十六进制机器码（需要先 `Encode` 回去） |

## 提示（不是答案）

- `Decode` 的 `Opcode` 已经写好了，照着它的形状写剩下五个；每个字段只要知道“从第几位开始、有几位”；
- 写完 `Decode` 先跑 `TestDecodeCoursewareAdd`，它用的就是课件里那个数；
- `%08x` 会补足 8 位并用小写十六进制；
- `Instr` 是可比较的结构体，测试直接用 `==` 比较——所以多出来的字段值一点都不能错。

## 自己验证

```bash
go test -v ./lesson-03-goroutines/test01/riscv
go test -run TestDisasm ./lesson-03-goroutines/test01/riscv
```

## 评分点

| 测试 | 检查内容 |
|---|---|
| `TestLittleEndian` | 课件里三组内存字节拼回机器码；`0xff` 不溢出 |
| `TestDecodeCoursewareAdd` | `0x00728433` 拆成 `add x8, x5, x7` 的六个字段 |
| `TestDecodeEveryField` | 每个字段取最大值时不串到相邻字段 |
| `TestEncodeIsTheInverseOfDecode` | `Encode(Decode(w)) == w`；超宽字段被截断 |
| `TestDisasm` | `add` / `sub` / 三种 `unknown` |
| `TestBridgeFetchDecodeFromMemory` | 从一段字节里每 4 字节取一条、译成汇编 |

## 做完以后

`TestBridgeFetchDecodeFromMemory` 里那个 `for pc := 0; ...; pc += 4` 就是 CPU 的“取指 → 译码 → 更新 PC”。它还少了一步：**执行**。

下一题 `test02` 给你一台玩具 CPU，把“执行”和“写回”也补上——你会写出一个真正跑起来的取指-译码-执行循环。
