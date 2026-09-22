# 第一课练习：Go 基础、package、可见性

六道题，建议按顺序做。每道题一个 `testNN/` 目录，目录里有独立的 `README.md` 说明任务、对应课件章节和评分点。

| 题目 | 主题 | 对应课件 |
|---|---|---|
| [test01](test01/) | 把 C 的 Hello World 搬到 Go | 第二部分 1〜3 |
| [test02](test02/) | 函数、变量、条件与循环 | 第二部分 2.1〜2.5 |
| [test03](test03/) | 修好一个违反 package 规范的项目 | 第三、四部分 |
| [test04](test04/) | 可见性与不可见性 | 第五部分 |
| [test05](test05/) | 亲手写第一个 Go 测试 | 第六部分 |
| [test06](test06/) | **结课题**：学生名册，桥接第二课 | 第二、五部分 |

`test06` 是刻意设计的收尾：只用第一课的语法就能做完，但它会让你**量出**线性查找的代价和重复代码的数量，第二课的 map 与 interface 正是从这两个数字开始讲。

## 需要改哪些文件

只改各题 `README.md` 里点名的文件。除了这两处例外，**不要修改任何 `*_test.go`**：

- `test05/even/student_test.go`：这道题就是要你写测试；
- `test03/packagerules/testdata/`：这是故意写错的题目素材，需要你去修。

## 本地运行

在**仓库根目录**执行：

```bash
go test ./lesson-01-basics/...              # 跑完第一课全部题目
go test ./lesson-01-basics/test03/...       # 只跑一道题
go test -v ./lesson-01-basics/test06/roster # 看详细输出和 t.Log
go test -run TestSumTo ./lesson-01-basics/test02/syntax  # 只跑一个测试函数
go vet ./lesson-01-basics/...               # 顺手查一些静态问题
```

starter 的初始状态**一定是红的**，那是题目，不是环境坏了。`test05` 更是要求你先让它红一次。

## 提交与评分

push 到 main 之后，CI 会把 `config.json` 里 `lesson` 指定的那一课跑一遍，统计每个 `TestXxx` 的终态：

```text
积分 = 通过的测试数 / 全部测试数 × 100，取整
```

结果会生成一份只含结论的 `result.json`，并发送到课程成绩网站。完整流程见仓库根目录的 [README.md](../README.md)。
