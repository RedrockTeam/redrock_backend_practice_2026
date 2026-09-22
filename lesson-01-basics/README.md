# 第一课练习：Go 基础、package、可见性

八道题，**按顺序做**。题目顺序就是课件的顺序：第二部分三题、第三部分一题、第四部分一题、第五部分一题、写测试一题，最后一题把你送进第二课。

每道题一个 `testNN/` 目录，里面有独立的 `README.md` 说明任务、对应课件章节和评分点。

| 题目 | 主题 | 对应课件 |
|---|---|---|
| [test01](test01/) | 把 C 的 Hello World 搬到 Go | 第二部分 1〜2.3 |
| [test02](test02/) | 函数、变量、零值、常量、条件与循环 | 第二部分 2.1〜2.5 |
| [test03](test03/) | 把一整个 C 程序按十步清单改写成 Go | 第二部分 3 |
| [test04](test04/) | 拆成多个包，并用完整路径导入 | 第三部分 |
| [test05](test05/) | 修好一个违反 package 规范的项目 | 第四部分 |
| [test06](test06/) | 可见性与不可见性 | 第五部分 |
| [test07](test07/) | 亲手写第一个 Go 测试 | 第六部分 |
| [test08](test08/) | **结课题**：学生名册，桥接第二课 | 综合 |

## 这八道题覆盖了课件的哪些点

| 课件知识点 | 在哪一题练 |
|---|---|
| `package main`、`func main`、`import "fmt"` | test01、test03 |
| `fmt.Println` 自动换行 vs `fmt.Printf` 要自己写 `\n` | test01 |
| 包名小写、导出函数名大写 | test01、test06 |
| 函数定义：`func` 开头，参数「名字 类型」，返回类型在后 | test02、test03、test04 |
| 同类型参数合并 `func add(a, b int) int` | test02 |
| `var` 声明、`:=` 简写 | test02、test03 |
| **零值**：未初始化的 `int` 是 `0` | test02 |
| 常量 `const`，以及常量的可见性 | test02、test06 |
| `if` / `for` 不写小括号，`for` 三段间仍有分号 | test02、test03 |
| 行尾不写分号、`{` 必须同行 | test03 |
| `fmt.Sprintf` | test02、test04、test08 |
| C → Go 十步迁移清单 | test03 |
| 一个目录一个包 | test04、test05 |
| module 路径 + 子目录 = **导入路径** | test04 |
| 跨包调用导出标识符 | test04 |
| 包名规范（小写、无下划线、别叫 `util`/`common`） | test05 |
| `main` 包不能被导入 | test04、test05 |
| `internal/` 的导入边界 | test05 |
| 禁止循环导入 | test05 |
| 大小写决定可见性（函数、变量、常量、类型、字段、方法） | test02、test04、test06 |
| **包内可见，而不是文件内可见**（对比 C 的 `static`） | test04、test06 |
| 同包测试 `package X` vs 外部测试 `package X_test` | test02、test04、test06 |
| 写测试、覆盖边界值 | test07 |
| 切片、多返回值 `(T, bool)` | test08 |

第一部分（什么是后端）是概念课，没有对应的代码题；它的回报在第二课之后的项目里。

## test08 为什么特殊

只用第一课的语法就能做完，但它会让你**量出**线性查找的代价和重复代码的数量——200 条名册、查一个人最坏要比 200 次，两个格式化函数写了两遍同样的循环。第二课的 map 和 interface 正是从这两个数字开始讲。

## 需要改哪些文件

只改各题 `README.md` 里点名的文件。除了这两处例外，**不要修改任何 `*_test.go`**：

- `test07/even/student_test.go`：这道题就是要你写测试；
- `test05/packagerules/testdata/`：这是故意写错的题目素材，需要你去修。

`test03/reference.c` 是题面，不参与编译，也不用改。

## 本地运行

在**仓库根目录**执行：

```bash
go test ./lesson-01-basics/...              # 跑完第一课全部题目
go test ./lesson-01-basics/test04/...       # 只跑一道题
go test -v ./lesson-01-basics/test08/roster # 看详细输出和 t.Log
go test -run TestSumTo ./lesson-01-basics/test02/syntax  # 只跑一个测试函数
go vet ./lesson-01-basics/...               # 顺手查一些静态问题
```

starter 的初始状态**一定是红的**，那是题目，不是环境坏了。`test07` 更是要求你先让它红一次。

## 有几道题会读你的源码

test01、test02、test03、test04 里有一些测试用 `go/parser` 把你的代码解析成语法树来检查，比如「确实用了 `fmt.Printf`」「确实靠零值而不是 `return 0`」「确实没留手写分号」。

这不是为了为难你。这些知识点光看函数返回值是分辨不出来的，而它们恰恰是 C 转 Go 时最容易一直错下去的地方。测试失败时会直接告诉你缺了什么。

## 提交与评分

push 到 main 之后，CI 会把 `config.json` 里 `lesson` 指定的那一课跑一遍，统计每个 `TestXxx` 的终态：

```text
积分 = 通过的测试数 / 全部测试数 × 100，取整
```

结果会生成一份只含结论的 `result.json`，并发送到课程成绩网站。完整流程见仓库根目录的 [README.md](../README.md)。
