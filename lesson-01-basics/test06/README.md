# test06（结课题）：用第一课的工具做一份学生名册

对应课件：第一课 · 第二部分（函数、循环）、第五部分（可见性）。
**这一题是第一课到第二课的桥。** 它只用第一课学过的东西就能做完，但做完之后你会亲手量到两个「不太对劲」的地方——第二课的 map 和 interface 就是来解决它们的。

## 先补两个新语法

第一课的课件没讲这两样，但它们是本题的前提。只有这么多，够用了。

### 1. 切片：一串同类型的值

```go
var students []Student                        // 零值是 nil，长度 0，可以直接 append
students = append(students, Student{ID: "S001"}) // 追加一个元素，返回新的切片
length := len(students)                       // 长度
first := students[0]                          // 下标访问，从 0 开始

for index, student := range students {        // 遍历：index 是下标，student 是这一份的副本
	fmt.Println(index, student.Name)
}

for _, student := range students {            // 不需要下标就用 _ 丢掉它
	fmt.Println(student.Name)
}
```

你可以先把它当成 C 里「长度能自己变的数组」。`append` 必须写成 `s = append(s, x)`，因为它返回的是新切片。

### 2. 多返回值：Go 用返回值表达「有没有」

C 里查不到通常返回 `-1` 或 `NULL`。Go 直接返回两个值：

```go
func FindByID(id string) (Student, bool) {  // 返回类型写成一个括号里的列表
	return Student{}, false                 // Student{} 是结构体的零值
}

student, ok := r.FindByID("S001")           // 一次接收两个返回值
if !ok {
	// 没找到
}
```

`ok` 这个名字是 Go 的惯例，第二课查 map 的时候还会见到一模一样的写法。

## 任务

编辑 `roster/roster.go`，完成全部 `TODO`。

```go
type Student struct {          // 三个字段都导出：别的包要能构造和读取
	ID    string
	Name  string
	Score int
}

type Roster struct {           // 两个字段都不导出：别的包只能通过方法操作
	students    []Student
	comparisons int
}
```

| 函数 | 要求 |
|---|---|
| `New() *Roster` | 返回一份空名册 |
| `(*Roster) Add(Student)` | 追加到末尾；`Name` 两端的空白要用 `strings.TrimSpace` 清掉 |
| `(*Roster) Len() int` | 记录条数 |
| `(*Roster) FindByID(id) (Student, bool)` | 按学号线性查找；找不到返回零值和 `false` |
| `(*Roster) TotalScoreByName(name) int` | 所有同名学生的分数之和；名册里允许同名 |
| `(*Roster) Comparisons() int` | 累计比较次数，见下 |
| `FormatText(*Roster) string` | 每行 `S001 Lin 90`，行尾 `\n`；空名册返回 `""` |
| `FormatCSV(*Roster) string` | 首行表头 `id,name,score`，之后每行 `S001,Lin,90`；空名册只返回表头行 |

### 关于 comparisons

`FindByID` 和 `TotalScoreByName` 每**检查一条记录**，就要把 `r.comparisons` 加 1。

```go
for _, student := range r.students {
	r.comparisons++            // 先记一次，再比较
	if student.ID == id {
		return student, true
	}
}
```

所以：命中第 k 条就是 k 次，查不到就是 `Len()` 次。这不是业务功能，而是这道题的**量尺**——`TestBridge...` 系列测试会用它把线性查找的代价读出来。

`comparisons` 是小写字段，外部测试包读不到，只能走 `Comparisons()` 方法。这正是 test04 的可见性规则第一次用在有意义的地方：内部计数器不对外暴露写入能力。

### 为什么两个 Format 是普通函数，不是方法

留意它们的签名：

```go
func FormatText(r *Roster) string
func FormatCSV(r *Roster) string
```

**完全一样。** 所以可以把它们放进同一个切片，用同一段代码调用：

```go
for _, format := range []func(*Roster) string{FormatText, FormatCSV} {
	fmt.Print(format(r))
}
```

这是故意的，先记住这个形状。

## 自己验证

```bash
go test ./lesson-01-basics/test06/...
go test -v -run TestBridge ./lesson-01-basics/test06/roster   # 看桥接测试打印的数字
```

`-v` 会把 `t.Log` 的内容打出来，这一题的重点信息都在那里。

## 评分点

| 测试 | 检查内容 |
|---|---|
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

第二课还会让你看到 map 的另一面：它**不是并发安全的**。两个 goroutine 同时写同一个 map，程序会直接崩。那道题会故意留着这个问题，把你送到第三课的 goroutine 和第四课的锁。

所以链条是这样的：

```text
第一课  切片 + 结构体 + 可见性     ← 你在这里，量到了线性代价和重复代码
第二课  map 消掉线性查找，interface 消掉重复循环，然后撞上 map 的并发问题
第三课  goroutine：并发到底是怎么跑的，为什么会撞
第四课  channel / Mutex / WaitGroup：把并发管起来
```
