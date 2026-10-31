# 题目设计与课堂转场

每个 `testNN/README.md` 是学生看到的任务边界，不含答案；公开测试就是行为契约，CI 只按测试终态给分。这份文档写给出题人，说明每道题为什么在这个位置。

## 第一课：基础（9 题）

前八题的顺序就是课件的顺序，一节对一题，test08 送去第二课；test09 是课后进阶纠错题。

| 题目 | 包 | 对应课件 | 要点 |
|---|---|---|---|
| `test01` | `hellogo` | 第二部分 1〜2.3 | 第一个能跑的程序；强制两行分别用 `Println` 和 `Printf`，把「谁自动换行」钉死 |
| `test02` | `syntax` | 第二部分 2.1〜2.5 | 函数定义、`var`/`:=`、**零值**、导出与未导出常量、`if`/`for`、`Sprintf` |
| `test03` | `translate` | 第二部分 3 | 给一份完整 C 程序，按十步清单整体改写；考分号、大括号、`printf`→`fmt.Printf` |
| `test04` | `hello/calc` + `hello/app` | 第三部分 | 自己拆两个包并互相导入；**导入路径 = module 路径 + 子目录** |
| `test05` | `packagerules` | 第四部分 | 改一棵故意违规的 `testdata` 树：包名与目录不符、导入 `main`、跨模块引用 `internal/`、包循环 |
| `test06` | `visibility` | 第五部分 | 大小写即可见性；代码拆两个文件证明「包内可见，不是文件内可见」；反射挡住「改成大写」的绕过 |
| `test07` | `even` | 第六部分 | 先写测试让它红，再修 `IsEven`；契约测试检查确实覆盖了 `0` 和负数 |
| `test08` | `roster` | 综合 | 见下 |
| `test09` | `habits` | 第二部分 + 第六部分 | 纠错：五个能编译但结果错的 C 习惯——`switch` 不贯穿、整数除法、`len` 数字节、按字节倒序拆坏汉字、`:=` 遮蔽 `err` |

### 为什么 test03 和 test04 必须存在

它们补的是最大的两个洞。第二部分的落脚点是那份十步迁移清单，但如果只出小函数题，学生永远不会真的把一段 C **整体**搬过来——分号、大括号、循环变量作用域这些坑都碰不到。第三部分整节在讲「怎么把代码拆成多个包、导入路径怎么来」，而如果没有一道题让学生亲手写下那行 `import "redrock/teaching-exercises/..."`，这一节就只是读过而已。

这两题都有靠 `go/parser` 读源码的契约测试，因为它们考的东西看返回值分辨不出来：

| 契约测试 | 挡住的捷径 |
|---|---|
| `TestUsesBothPrintlnAndPrintf` | 两行都用 `Printf`，`Println` 的自动换行没体会到 |
| `TestInitialCountReliesOnTheZeroValue` | 直接 `return 0`，绕过零值 |
| `TestTranslationKeepsControlFlow` | 把六行输出写死成六条 `Println` |
| `TestTranslationDropsSemicolons` | C 风格行尾分号留着没删（`for` 头部里的不算） |
| `TestTranslationKeepsSquare` | 把 `square` 展开成字面量，不当函数翻译 |
| `TestAddDoubledGoesThroughHelper` | 自己写 `*2`，没用上同包的小写函数 |

### test08 是整课的出口

前七题都在第一课的知识范围内。`test08` 故意多要一点：学生要用**切片**存名册，用**多返回值** `(Student, bool)` 表达「查到了吗」。这两样第一课没讲，README 就地补——不是打补丁，而是因为后面正好都要用。

题目本身很朴素：往名册里加学生，按 ID 查人，按姓名统计总分，再把名册格式化成两种文本。关键在 `Roster` 里埋了一个 `comparisons` 计数器，于是三道 `TestBridge...` 能**量出**代价而不只是描述它：

- `TestBridgeLookupCostGrowsWithRosterSize`：200 条名册，命中首条 1 次比较，命中末条 200 次，查不到 200 次；
- `TestBridgeAggregationScansWholeRoster`：每问一次「这个名字总分多少」就重扫 200 条，问两次 400 次——切片没有任何记忆；
- `TestBridgeTwoFormattersShareOneSignature`：两个格式化函数被塞进同一个 `[]func(*Roster) string` 里依次调用，说明它们已经能被统一调用了。

这三个测试在做对时**是通过的**，痛感来自被断言出来的具体数字和 `t.Log` 的旁白，而不是人为制造的失败。README 末尾把它们直接兑换成下一课：

| 这一题遇到的 | 第二课的答案 |
|---|---|
| 查一个人要扫 N 条 | `map[string]Student`，`student, ok := m[id]` 和刚写的 `(Student, bool)` 同一个形状 |
| 每次统计都重扫一遍 | `map[string]int`，一次遍历攒好所有名字的总分 |
| 两个格式化函数重复同一个循环 | `interface`，遍历写一次，每行长什么样交给不同实现 |

顺带也把 `FormatText` / `FormatCSV` 写成了自由函数而不是方法——它们签名完全相同，离 interface 只差一步。

## 题量与题型

每课 7〜9 题，单题 30〜60 分钟，一课总量约 5〜7 小时，足够覆盖课上练习加一周课后。每课内部难度递增，并且三种题型交替出现：

| 题型 | 学生拿到什么 | 例子 |
|---|---|---|
| 补全 | 结构和签名已给，`TODO` 处填逻辑 | L2 test04 `scoreboard`、L4 test04 `pipeline` |
| 代码纠错 | 能编译、看起来差不多对，但埋了 bug | L1 test09 `habits`、L2 test03 `aliasing`、L4 test06 `lazyfetch` |
| 从零实现 | 只有类型与签名、返回零值，规格全在 README | L2 test05 `dns`、L3 test02 `toycpu`、L4 test02 `counter` |
| 运行观察 → 修复 / 观察作答 | 先跑出现象，再修或作答 | L4 test01 `runall`、L3 test07 `missing` |
| 写测试 | 学生自己写测试抓 bug | L1 test07 `even` |

所有题都以「push 后 workflow 全绿」为完成标准，**没有设计成必然红的题**。

## 第二课：指针、集合与抽象（8 题）

接第一课 `test08`，顺序是「值还是句柄 → 切片 → map → interface → 交出共享 map」：

| 题目 | 包 | 题型 | 要点 |
|---|---|---|---|
| `test01` | `pointer` | 补全 + 预测输出 | 数组是值、`*[3]int`、结构体指针；预测写进常量，测试对比 `demo.go` 的实际输出 |
| `test02` | `sliceops` | 观察 → 修复 | `Clone` 独立副本、下标写回、接住 `append` 返回值 |
| `test03` | `aliasing` | 纠错 | cap 富余时 append 原地写入、两次 append 写同一格、边删边遍历、`&v`、分段共享 cap |
| `test04` | `scoreboard` | 补全 | map 一次建表，`Top` 分数降序、姓名升序；第一课线性查找的答案 |
| `test05` | `dns` | 从零实现 | 正查 + 反查两张 map 同步维护，大小写与末尾 `.` 归一；桥到第五课 DNS |
| `test06` | `formatter` | 补全 | `Formatter` 接口，`FormatAll` 只写一次遍历；第一课重复循环的答案 |
| `test07` | `pitfalls` | 纠错 + 补全 | typed nil 进 `error`、值接收者改副本、`type switch`、`Stringer`、不可比较值的安全比较 |
| `test08` | `cache` | 补全 + 桥接 | map 是句柄，`Snapshot` 必须复制；两个**各自独立**的 Cache 并发安全，同一个 Cache 留给第三、四课 |

（评分脚本对「跑起来但没走到终态」的崩溃有专门处理：例如 `fatal error: concurrent map writes` 会打死进程、不产生任何 pass/fail 事件，此时脚本把它归因到触发它的那个测试函数，而不是笼统记成编译失败。第四课账本题会用到这条。）

## 第三课：从 CPU 到 goroutine（7 题）

课件前半是计算机基础，后半是 goroutine。前三题不写一个 `go`，但解释后四题的所有现象：

| 题目 | 包 | 题型 | 要点 |
|---|---|---|---|
| `test01` | `riscv` | 补全 | 小端序拼字、位运算拆 R 型字段、反汇编 `add x8, x5, x7`（课件里的 `0x00728433`） |
| `test02` | `toycpu` | 从零实现 | 取指 → 译码 → 执行 → 写回 → 更新 PC；最大步数、非法指令与越界错误 |
| `test03` | `interleave` | 实现 + 观察 | 两三个线程共享内存按调度表交错，枚举全部调度：两线程 20 种里 18 种丢更新；原子 `INCR` 全部正确 → 第四课 test02 |
| `test04` | `spawn` | 补全 | `go` 立即返回、`runtime.NumGoroutine` 量一万个 goroutine 的代价 |
| `test05` | `asyncprint` | 观察 → 补全 | 源码契约要求出现 `go` 语句；允许函数内 `Sleep` |
| `test06` | `parallel` | 改写 | 拆成两段并行求和；`Sleep` 版在 `-race` 下报竞争，是第四课的入口 |
| `test07` | `missing` | 观察 + 作答 | `Collect` 的 bug **保留不修**；学生写 `Verify` 量出缺失数，并在 `answers.go` 回答观察题（测试只存答案摘要，不泄露选项）。第四课 `test01` 修它 |

`test07` 以前是「第三课结束时仍然红」的设计，和「以 workflow 跑通为准」冲突，因此改成了能在第三课知识范围内全绿的观察题；`-race` 下桥接测试自动跳过，避免学生自查时被 `Collect` 自身的竞争误报。

## 第四课：同步（8 题）

建议课上做 test01〜04，课后做 test05〜08：

| 题目 | 包 | 题型 | 要点 |
|---|---|---|---|
| `test01` | `runall` | 观察 → 修复 | 修第三课 `test07` 的“缺失的等待”，建立 `Done → Wait`；源码契约禁止 `Sleep` |
| `test02` | `counter` | 从零实现 | 课件 LV1：`Mutex` 版与 `atomic` 版计数器，100×100 并发；接第三课 `test03` 的丢更新 |
| `test03` | `fanin` | 补全 | 所有输入关闭、所有转发 goroutine 退出后由 `Merge` 单点 `close` |
| `test04` | `pipeline` | 补全 | 课件 LV2：Generate → Square → Filter，每级收发都可取消，`NumGoroutine` 回落验证不泄漏 |
| `test05` | `ledger` | 观察 → 修复 | 丢更新、并发写 map 进程退出、`Snapshot` 交出内部 map；接第二课 `test08` |
| `test06` | `lazyfetch` | 纠错 | check-then-act 懒加载改 `sync.Once`；超时错误用 `%w` 包装；超时后无人接收的 channel 导致泄漏 |
| `test07` | `pool` | 从骨架实现 | 固定 worker、可取消、结果按下标归位 |
| `test08` | `wire` | 桥接 | `net.Pipe` 字节流没有消息边界，接第五课 |

`-race` 由学生自己跑、不计入成绩。四课的概念在这里连成一条完整的数据流，最后通过 `test08` 接到第五课的 `net`、TCP 与 HTTP。

## 出题时的几条约束

- 题目单位固定为 `testNN/`，Go 包放在它下一层，包名与目录名一致；
- 每个 `testNN/` 必须有 README，说明做什么、对应课件哪一节、评分点是哪几个测试函数；
- starter 必须能编译——学生该看到的是测试失败，不是一屏语法错误。用 `return fmt.Sprintf("")` 这类「能过编译但显然不对」的占位，别用 `_ = fmt.Sprintf`；
- 公开测试的函数名就是评分点，改名等于改评分接口；
- 需要故意留 bug 的题（如 `test07`），配一个契约测试检查学生的测试确实覆盖了边界，否则「写一个恒真的测试」也能拿分；
- 知识点如果「看返回值分辨不出来」（零值、用了哪个打印函数、有没有留分号），就用 `go/parser` 读源码判，并在失败信息里直接说明缺了什么。
