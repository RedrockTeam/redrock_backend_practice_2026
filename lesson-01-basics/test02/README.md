# test02：函数、变量、条件与循环

对应课件：第一课 · 第二部分 2.1〜2.5（函数定义、分号、打印、变量声明、条件与循环）。

## 任务

编辑 `syntax/syntax.go`，完成全部 `TODO`：

| 函数 | 要求 |
|---|---|
| `add(a, b int) int` | 返回 `a + b`。注意它是**小写**的包内函数 |
| `ClassifyNumber(n int) string` | `n > 0` 返回 `"positive"`，`n == 0` 返回 `"zero"`，`n < 0` 返回 `"negative"` |
| `SumTo(n int) int` | 返回 `1 + 2 + ... + n`；`n <= 0` 时返回 `0` |
| `Describe(name string, age int) string` | 返回 `"小明 今年 18 岁"` 这样一行，**末尾不带换行** |
| `Clamp(score int) int` | 把分数收进 `[minScore, MaxScore]`；直接用这两个常量，别写死 `0` 和 `100` |
| `InitialCount() int` | 用 `var` 声明一个 `int` 后**不赋值**直接返回，靠 Go 的零值拿到 `0` |

不要修改 `Add`，也不要修改任何 `_test.go` 文件。

## 需要注意的 C → Go 差异

```c
int add(int a, int b) { return a + b; }   // C：返回类型在前，参数是「类型 名字」
```

```go
func add(a, b int) int { return a + b }   // Go：func 开头，返回类型在后，参数是「名字 类型」
```

- 变量：`var sum int = 0`，函数内可以简写成 `sum := 0`。
- **零值**：`var` 声明但不赋值时，Go 会给一个零值——`int` 是 `0`，`string` 是 `""`，`bool` 是 `false`。
  C 里不初始化的局部变量是垃圾值，Go 不是，这是两边最容易踩的差别之一。
  `InitialCount` 考的就是这个，所以直接 `return 0` 不算过（测试会读你的源码）。
- **常量**：`const MaxScore = 100` 导出，`const minScore = 0` 不导出。可见性对常量一样适用。
- 条件和循环的条件**不写小括号**：`if n > 0 {`、`for i := 1; i <= n; i++ {`。
- 左大括号 `{` 必须和 `func` / `if` / `for` 在同一行，否则自动插入的分号会让编译失败。
- 行尾不写分号，但 `for` 的三段之间仍然用分号。
- `fmt.Sprintf` 的格式化规则和 `Printf` 一样，只是把结果当返回值：`fmt.Sprintf("%s 今年 %d 岁", name, age)`。

## 这一题同时在演示可见性

这个目录里有两个测试文件，它们的 `package` 声明不同：

- `syntax_test.go` 声明 `package syntax_test`，是**外部测试包**，只能调用大写的 `Add`；
- `syntax_internal_test.go` 声明 `package syntax`，和被测代码**同包**，因此可以直接调用小写的 `add`。

同一个目录允许出现 `包名` 和 `包名_test` 两个包，这是课件第四部分「一个目录一个包」提到的唯一例外。可见性的完整规则在 test06 练。

## 自己验证

```bash
go test ./lesson-01-basics/test02/...
go test -run TestSumTo ./lesson-01-basics/test02/syntax
```

## 评分点

| 测试 | 检查内容 |
|---|---|
| `TestAdd` | `Add` 经由 `add` 返回正确的和 |
| `TestClassifyNumber` | 三个分支都正确 |
| `TestSumTo` | 正常累加，且 `n <= 0` 时返回 0 |
| `TestDescribe` | 格式化字符串完全一致 |
| `TestPackagePrivateAdd` | 同包测试能直接调用小写 `add` |
| `TestMaxScoreIsExported` | 外部包能读到大写常量 `MaxScore` |
| `TestClamp` | 上下界和区间内都正确 |
| `TestPackagePrivateMinScore` | 同包测试能直接读小写常量 `minScore` |
| `TestInitialCountReliesOnTheZeroValue` | 返回 0，**并且**源码里确实写了不赋值的 `var` 声明 |
