# test07：interface 的常见坑

对应课件：第二课 · interface（底层实现、空接口与类型断言、`type switch`、`Stringer`、常见坑三则）。
题型：**代码纠错 + 补全**（综合）。

## 背景

`test06` 用 interface 消掉了重复循环，那是 interface 好用的一面。这一题看它容易出错的一面，四个坑全部来自课件“常见坑”和“接口的常见用法”。

## 任务

题目分在三个文件里，只改这三个文件，不改测试。

### 1. `validate.go`：装了 nil 指针的接口不是 nil

`Validate` 校验通过时“返回了 nil”，可调用方的 `err != nil` 却是 `true`。先写一个小 `main` 或在测试里打印 `err == nil` 和 `%T` 看看，再修。

要求：合法时返回真正的 `nil`；不合法时返回 `*ValidationError`，`Field` 为 `"name"` 或 `"age"`，两者都不合法时报告 `name`。

### 2. `tally.go`：值接收者改的是副本

`Tally` 满足 `Incrementer`，调用 `Inc()` 也不报错，但计数永远是 0。修好它，让 `NewTally()` 返回的计数器通过接口调用也能正确计数。

改完 `Inc` 后，`NewTally` 很可能编译不过——读懂编译器的报错，它说的正是课件“坑 2”。

### 3. `describe.go`：type switch、Stringer、不可比较的接口

- 给 `Point` 实现 `fmt.Stringer`，`fmt.Sprint(Point{1, 2})` 输出 `(1, 2)`；
- 用 `type switch` 实现 `Describe`，匹配顺序见源码注释（**顺序有意义**：`Point` 既是 `Point` 也是 `fmt.Stringer`）；
- 修好 `SafeEqual`：两边是切片、map 这类不可比较类型时返回 `false, false`，不能 panic。

## 提示（不是答案）

- 接口值内部是 `(动态类型, 动态值)` 两个字段；只有两个都为空，接口才 `== nil`；
- 返回 `error` 的函数，成功路径直接写 `return nil`，不要先声明一个具体类型的指针变量再返回它；
- 方法集：`T` 的方法集只包含值接收者方法；`*T` 的方法集包含两种。接口要求的方法有指针接收者时，只有 `*T` 满足它；
- `type switch` 里 `case error:` 这样写接口类型是允许的，会按顺序逐个尝试；
- 判断能不能比较有两条路：`reflect.TypeOf(x).Comparable()`，或者用 `defer` + `recover` 接住比较时的 panic。注意 `reflect.TypeOf(nil)` 返回 `nil`。

## 评分点

| 测试 | 检查内容 |
|---|---|
| `TestValidateReturnsRealNil` | 合法输入返回真正的 `nil` |
| `TestValidateReportsField` | 错误类型与字段；边界 `Age = 0` 合法 |
| `TestTallyCountsThroughInterface` | `NewTally()` 经接口调用能计数 |
| `TestPointerSatisfiesIncrementer` | `*Tally` 满足接口并能计数 |
| `TestPointImplementsStringer` | `Point` 实现 `fmt.Stringer` |
| `TestDescribeTypeSwitch` | 七种分支与匹配顺序 |
| `TestSafeEqualComparable` | 可比较类型的比较结果 |
| `TestSafeEqualUncomparableDoesNotPanic` | 切片、map 不 panic，返回 `false, false` |
| `TestBridgeInterfaceHoldsTypeAndValue` | 观察：接口的动态类型与动态值（不检查你的实现） |

## 做完以后

这四个坑有一个共同点：**接口变量里装的是什么，光看变量类型看不出来**。`err` 的静态类型是 `error`，里面可能是 `nil`，也可能是一个 `nil` 指针；`any` 里可能是 `int`，也可能是一个比较就 panic 的切片。

另外，`Tally` 的教训会在第四课再出现一次：带锁的计数器**必须**用指针接收者，否则每次调用锁的都是一把副本锁。

下一题 `test08` 是这一课的结课题：把 map、方法、所有权放在一起，做一个会被第三、四课继续使用的 `Cache`。
