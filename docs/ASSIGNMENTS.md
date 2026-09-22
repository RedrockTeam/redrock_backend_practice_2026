# 题目设计与课堂转场

每个 `testNN/README.md` 是学生看到的任务边界，不含答案；公开测试就是行为契约，CI 只按测试终态给分。这份文档写给出题人，说明每道题为什么在这个位置。

## 第一课：基础（8 题）

习题顺序就是课件的顺序，一节对一题，最后一题送去第二课。

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

## 第二课：集合

`test01` 用 map 统计多次成绩，`test02` 让两个 formatter 落到同一个 interface 上（正好是 test06 留下的那个形状）。

（评分脚本对「跑起来但没走到终态」的崩溃有专门处理：例如 `fatal error: concurrent map writes` 会打死进程、不产生任何 pass/fail 事件，此时脚本把它归因到触发它的那个测试函数，而不是笼统记成编译失败。第三、四课遇到并发 bug 时会用到这条。）

## 第三课：goroutine

`test01` 先保持顺序，`test02` 比较顺序与并发的耗时，`test03` 启动 goroutine 却立即返回，测试稳定地暴露结果缺失。这一课把第二课的 map 放到并发视角下重新看：共享内存是什么、为什么会乱、谁负责等待。

## 第四课：同步

`test01` 合并多个 channel，`test02` 修一个并发计数器，`test03` 实现固定 worker 数量的任务池。结课题要求功能测试全部通过，`-race` 由学生自己跑、不计入成绩，四课的概念在这里连成一条完整的数据流。

## 出题时的几条约束

- 题目单位固定为 `testNN/`，Go 包放在它下一层，包名与目录名一致；
- 每个 `testNN/` 必须有 README，说明做什么、对应课件哪一节、评分点是哪几个测试函数；
- starter 必须能编译——学生该看到的是测试失败，不是一屏语法错误。用 `return fmt.Sprintf("")` 这类「能过编译但显然不对」的占位，别用 `_ = fmt.Sprintf`；
- 公开测试的函数名就是评分点，改名等于改评分接口；
- 需要故意留 bug 的题（如 `test07`），配一个契约测试检查学生的测试确实覆盖了边界，否则「写一个恒真的测试」也能拿分；
- 知识点如果「看返回值分辨不出来」（零值、用了哪个打印函数、有没有留分号），就用 `go/parser` 读源码判，并在失败信息里直接说明缺了什么。
