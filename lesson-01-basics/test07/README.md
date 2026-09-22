# test07：亲手写第一个 Go 测试

对应课件：第一课 · 第六部分（练习与提交流程）。

前面四题你都是「被测试检查」，这一题你来写测试。整个课程的评分都建立在 `go test` 上，所以你至少要能读懂、写出一个测试。

## 任务

分两步，顺序很重要。

### 第一步：写测试，看它失败

`even/even.go` 里的 `IsEven` **是错的**，它漏掉了一整类输入。先编辑 `even/student_test.go`，补完 `TestIsEven`：

- 至少检查 `IsEven(2)`、`IsEven(3)`、`IsEven(0)`、`IsEven(-2)` 四种输入；
- 结果和预期不一致时，用 `t.Error` / `t.Errorf` / `t.Fatal` / `t.Fatalf` 报告失败。

期望行为：`2` 是偶数，`3` 不是，`0` 是偶数，`-2` 也是偶数。

写完先跑一次：

```bash
go test -v ./lesson-01-basics/test07/even
```

你应该看到自己写的测试**失败**了。这一步是题目的一部分，不要跳过——先让测试红，才能确认它真的在检查东西。

### 第二步：修实现，让测试变绿

回到 `even/even.go`，修掉 `IsEven`，再跑一次，直到全绿。

## 一个 Go 测试长什么样

```go
package even          // 和被测代码同包，可以直接调用包内的一切

import "testing"      // 测试用的标准库

func TestIsEven(t *testing.T) {   // 必须以 Test 开头，参数是 *testing.T
	if !IsEven(2) {
		t.Errorf("IsEven(2) = false, want true")
	}
}
```

规则：

- 文件名必须以 `_test.go` 结尾，否则 `go test` 不会认它；
- 函数名必须以 `Test` 开头，后面接大写字母，参数是 `*testing.T`；
- 不调用任何报告方法就算通过。**测试是靠「报告失败」工作的，不是靠 `return`。**

| 方法 | 行为 |
|---|---|
| `t.Error` / `t.Errorf` | 标记失败，但继续执行后面的检查 |
| `t.Fatal` / `t.Fatalf` | 标记失败并立刻结束这个测试函数 |
| `t.Log` / `t.Logf` | 只打印，不影响结果；`go test -v` 才看得到 |

写失败信息的惯例是 `函数调用 = 实际值, want 期望值`，这样失败输出本身就说明了问题。

## 为什么还有一个 contract_test.go

`contract_test.go` 用 `go/parser` 读你的 `student_test.go`，确认你真的写了那四个输入的检查，而不是只把 `t.Fatal("TODO")` 那行删掉。它不检查你的代码风格，只检查覆盖到的输入和有没有调用失败报告方法。

## 评分点

| 测试 | 检查内容 |
|---|---|
| `TestIsEven` | 你写的测试本身要通过（也就是 `IsEven` 已被修好） |
| `TestStudentWroteIsEvenCases` | `TestIsEven` 覆盖了 `2 / 3 / 0 / -2`，并且会报告失败 |
