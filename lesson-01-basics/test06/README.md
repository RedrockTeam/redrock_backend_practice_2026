# test06：可见性与不可见性

对应课件：第一课 · 第五部分（可见性的全部 8 小节）。

## 任务

`visibility/` 目录下有**两个源文件**，它们属于同一个包 `visibility`：

- `user.go`：`User` 类型和它的方法；
- `stage.go`：包内的小写常量 `adultAge` 和小写函数 `stage`。

完成全部 `TODO`：

| 标识符 | 所在文件 | 要求 |
|---|---|---|
| `NewUser(name, age)` | `user.go` | 同时保存 `name` 和 `age` |
| `(*User) Age()` | `user.go` | 返回小写字段 `age` |
| `(*User) Grow()` | `user.go` | 调用包内小写方法 `grow`，不要直接写 `u.age++` |
| `(*User) grow()` | `user.go` | 把 `age` 加 1 |
| `(*User) Stage()` | `user.go` | 调用 `stage(u.age)` 得到阶段字符串 |
| `stage(age)` | `stage.go` | `age >= adultAge` 返回 `"adult"`，否则返回 `"minor"` |

不要修改 `_test.go` 文件，也不要把小写标识符改成大写。

## 这一题真正要你体会的

Go 没有 `public` / `private` 关键字，只有一条规则：**首字母是不是大写**。

```go
type User struct {
	Name string // 大写：所有包可见
	age  int    // 小写：只有 visibility 包内可见
}
```

在其他包里：

```go
u := visibility.NewUser("小登", 18)
u.Name    // 可以
u.Age()   // 可以
u.age     // 编译错误：cannot refer to unexported name
u.grow()  // 编译错误
```

### 小写是「包内可见」，不是 C 的「文件内可见」

这是最容易搞错的一点。C 的 `static` 把可见范围锁在当前 `.c` 文件里；Go 的小写标识符在**整个包**内都能用，跨文件也可以。

所以 `user.go` 里的 `Stage()` 能直接调用 `stage.go` 里的 `stage()`，不需要任何声明或 `#include`。这就是这一题特意拆成两个文件的原因。

### 两个测试文件，两种视角

| 文件 | `package` 声明 | 能看到什么 |
|---|---|---|
| `visibility_internal_test.go` | `visibility` | 小写字段 `age`、小写方法 `grow`、小写函数 `stage`、小写常量 `adultAge` |
| `visibility_external_test.go` | `visibility_test` | 只有 `User`、`NewUser`、`Name`、`Age`、`Grow`、`Stage` |

`TestAgeFieldStaysUnexported` 会用反射确认 `age` 还是小写的。反射能「看到」未导出字段的存在，但普通代码依然不能读写它——这正好说明可见性是编译期规则。

## 自己验证

```bash
go test ./lesson-01-basics/test06/...
go test -v ./lesson-01-basics/test06/visibility
```

## 评分点

| 测试 | 检查内容 |
|---|---|
<<<<<<< HEAD
| `TestAddAndLen` | `New` / `Add` / `Len` |
| `TestAddTrimsName` | `Add` 清理了名字两端空白 |
| `TestFindByID` | 命中与未命中，未命中时返回零值 |
| `TestTotalScoreByName` | 同名累加，不存在的名字返回 0 |
| `TestFormatText` | 文本输出逐字节一致 |
| `TestFormatCSV` | CSV 输出逐字节一致 |
| `TestFormatEmptyRoster` | 空名册的两种输出 |
| `TestBridgeLookupCostGrowsWithRosterSize` | 比较次数 = 命中位置，未命中 = 全长 |
| `TestBridgeAggregationScansWholeRoster` | 每次统计都重扫全表，两次的代价一样 |
| `TestBridgeTwoFormattersShareOneSignature` | 两个格式化函数可以被同一段代码调用 |

## 做完以后：你刚刚量到了什么

跑一遍 `go test -v -run TestBridge ./lesson-01-basics/test06/roster`，把输出留着。三件事：

**1. 查一个人，最坏要看完所有人。**
200 条名册，查最后一个用 200 次比较，查不存在的学号也是 200 次。名册涨到 20000 条，代价就涨 100 倍。你已经把这个数字测出来了，它不是猜的。

**2. 每问一次统计，就要重扫一遍。**
`TotalScoreByName("student-7")` 问两次，就扫了 400 条。切片没有任何「记住上次结果」的能力，因为它只知道顺序，不知道「谁是谁」。

**3. 两种输出，两个几乎一样的循环。**
`FormatText` 和 `FormatCSV` 都在做「遍历名册 → 拼每一行 → 接起来」，只有中间那一行长什么样不同。现在是两种格式；如果产品又要 JSON、又要 Markdown 表格，就是四份同样的循环。

## 第二课接着这三件事讲

| 你在这一题遇到的 | 第二课的答案 |
|---|---|
| 查一个人要扫 N 条 | `map[string]Student`：按 key 直接定位，`student, ok := m[id]`——和你刚写的 `(Student, bool)` 是同一个形状 |
| 每次统计都重扫一遍 | `map[string]int`：一次遍历把所有名字的总分都攒好，之后每次查都是一步 |
| 两个格式化函数重复同一个循环 | `interface`：把「怎么遍历」写一次，把「每行长什么样」交给不同实现；你已经见过它的雏形——两个函数共享一个签名 |

所以链条是这样的：

```text
第一课  切片 + 结构体 + 可见性     ← 你在这里，量到了线性代价和重复代码
第二课  map 消掉线性查找，interface 消掉重复循环
第三课  goroutine：并发到底是怎么跑的，谁负责等待
第四课  channel / Mutex / WaitGroup：把并发管起来
```
=======
| `TestPackagePrivateGrow` | 同包可以调用 `grow()`，`age` 被正确递增 |
| `TestPackagePrivateStage` | 同包可以调用另一个文件里的 `stage()` |
| `TestExportedUserAPI` | 外部包通过 `NewUser` / `Name` / `Age` / `Grow` 正常工作 |
| `TestExportedStage` | `Stage()` 正确转发到包内的 `stage()` |
| `TestAgeFieldStaysUnexported` | `age` 保持未导出，没有新增导出字段 `Age` |
>>>>>>> 1c50096 (feat(test): add lesson 01 test)
