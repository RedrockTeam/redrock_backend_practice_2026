# test06：用 interface 消掉重复循环

对应课件：第二课 · 对象、interface（隐式实现、方法集、多态、错误返回）。
题型：**补全代码**。

## 背景

第一课 `test08` 的 `FormatText` 和 `FormatCSV` 签名完全一样，却各自写了一遍“遍历名册、拼每行”的循环。这一题把“怎么遍历”留在 `FormatAll` 里，把“每行长什么样”交给不同的 `Formatter`。

## 任务

编辑 `formatter/formatter.go`：

```go
type Formatter interface {
    Format(Student) (string, error)
}

func FormatAll(formatters []Formatter, s Student) ([]string, error)
```

要求：

- `TextFormatter.Format` 返回 `"Lin: 95"`；
- `JSONFormatter.Format` 保持现有的 `encoding/json` 实现；
- `FormatAll` 按顺序调用每个 formatter，把结果收集到 `[]string`；
- 遇到第一个错误立刻返回**已经收集到的结果**和这个错误；
- 空 `formatters` 返回空结果和 `nil`。

不要修改 `_test.go`。

## 提示（不是答案）

- `Formatter` 是接口，`TextFormatter` 和 `JSONFormatter` 只要实现了 `Format` 就自动满足它，不需要写 `implements`；
- 传 `[]Formatter{TextFormatter{}, JSONFormatter{}}` 时，方法接收者是值类型也满足接口；
- `FormatAll` 里用 `for _, f := range formatters`，调用 `f.Format(s)`；
- 错误处理：先判断 `err != nil`，再 append 结果（顺序不同会导致“失败前的结果”不一致）；
- 文本格式化用 `fmt.Sprintf("%s: %d", ...)`。

## 评分点

| 测试 | 检查内容 |
|---|---|
| `TestTextAndJSONFormatters` | 两种具体格式化的输出 |
| `TestFormatAllUsesTheInterface` | 同一段代码通过 `[]Formatter` 调用两种实现 |
| `TestFormatAllStopsOnError` | 遇到错误停止，并保留之前的结果 |
| `TestFormatAllHandlesEmpty` | 空 formatter 列表 |

## 做完以后

你已经把“变化的部分”（每行长什么样）和“不变的部分”（按顺序遍历、收集、处理错误）分开了。这是接口最实用的形态，也是后面 Web 框架里 middleware、数据库驱动、日志输出的常见结构。

下一题 `test07` 看 interface 容易出错的一面：装了 nil 指针的 `error`、值接收者改副本、不可比较的 `any`。之后 `test08` 回到 map 做结课题。
