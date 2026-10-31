# test09（进阶）：找出五个从 C 带过来的坏习惯

对应课件：第一课 · 第二部分（语法迁移）、第二部分 3（十步清单之外的坑）。
题型：**代码纠错**。

test03 把一个 C 程序整段翻译成了 Go。翻译完能编译、能跑，不等于没有问题——有些写法在 C 里是对的，搬到 Go 里**能编译、但结果悄悄不对**。这一题给你五个这样的函数，每个只藏了一个坑。

## 任务

只改 `habits/habits.go`，不要改函数签名，不要改 `_test.go`。

先跑，看哪里红：

```bash
go test -v ./lesson-01-basics/test09/habits
```

| 函数 | 应有的行为 |
|---|---|
| `IsWeekend(day)` | `6`、`7` 返回 `true`，其余返回 `false` |
| `Average(scores)` | `[90 85]` 返回 `87.5`；空切片返回 `0` |
| `CountChars(s)` | 数**字符**：`"红岩网校"` 是 `4` |
| `Reverse(s)` | 按**字符**倒序：`"Go语言"` 变成 `"言语oG"` |
| `SumScores(items)` | 全部合法时返回总和和 `nil`；有一项不是整数时返回 `0` 和那个错误 |

## 需要补的一点新知识

`SumScores` 用到了标准库 `strconv`：

```go
n, err := strconv.Atoi("90")   // n = 90, err = nil
n, err = strconv.Atoi("八十")   // n = 0,  err != nil
```

Go 里的字符串是一串 **UTF-8 字节**。英文字母一个字节，汉字三个字节。`len(s)` 返回字节数；想按字符处理，可以先转成 `[]rune`（每个 `rune` 是一个字符），或者用 `unicode/utf8` 包。

## 提示（不是答案）

每个坑对应一个问题，先自己回答：

1. C 的 `switch` 不写 `break` 会怎样？Go 的呢？一个 `case` 怎么同时匹配两个值？
2. `total / len(scores)` 的两边是什么类型？结果在转成 `float64` **之前**还是**之后**算出来的？
3. `len("红岩网校")` 等于多少？为什么？
4. 把 3 个字节组成的汉字按字节倒序，得到的还是那个汉字吗？
5. 循环里的 `n, err := ...` 中，这个 `err` 和函数开头 `var err error` 的 `err` 是同一个变量吗？（这个现象叫 **shadowing**：内层作用域用 `:=` 新建了一个同名变量，把外层的遮住了。）

## 评分点

| 测试 | 检查内容 |
|---|---|
| `TestIsWeekend` | 一周七天的结果都正确 |
| `TestAverageKeepsFraction` | 平均分保留小数，空切片返回 0 |
| `TestCountCharsCountsCharactersNotBytes` | 中英文混合时按字符计数 |
| `TestReverseKeepsCharactersIntact` | 按字符倒序，汉字不被拆坏 |
| `TestSumScoresAddsValidInput` | 合法输入正确求和 |
| `TestSumScoresReportsBadInput` | 非法输入时错误能传出循环 |

## 做完以后

这五个坑有一个共同点：**编译器不会报错**。所以第六部分讲的 `go test` 才重要——test07 里你写测试抓住了 `IsEven` 的漏洞，这一题是测试替你抓住了五个。

第五个坑（`:=` 遮住外层的 `err`）在后面几课会反复出现：第二课的 interface 题、第四课的并发题都会返回 `error`，写错一个 `:=`，错误就静悄悄地消失了。另外，第二课 test01 会讲指针和数组，`[]rune` 和 `[]byte` 这样的切片正是那一课的主角。
