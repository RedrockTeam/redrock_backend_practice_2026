# test02：用 map 消掉线性查找

对应课件：第二课 · map（创建、增删改查、`v, ok := m[k]`、range、句柄传递、key 要求）。
题型：**补全代码**。

## 背景

第一课 `test08` 的 `FindByID` 每次都从头扫名册，200 条名册查最后一个要比较 200 次。这一题用 `map` 把“按名字查总分”和“做排行榜”写成一次建表、一次排序。

## 任务

编辑 `scoreboard/scoreboard.go`：

```go
func Totals(records map[string][]int) map[string]int
func Top(totals map[string]int, n int) []Entry
```

### Totals

- 返回一张**新 map**：学生名 -> 所有成绩之和；
- 某个学生成绩切片为空或 `nil` 时，总分为 0，但键仍保留；
- `records` 本身不能被修改；
- `records == nil` 时返回空 map（长度为 0 即可）。

### Top

- 把 `totals` 转成 `[]Entry`；
- 按**分数降序**排序；分数相同按**姓名升序**；
- `n <= 0` 返回空结果；`n` 大于人数时返回全部；
- 不修改 `totals`。

## 提示（不是答案）

- 建表：`make(map[string]int, len(records))`；
- 读 map 区分“没有这个键”和“值是零值”，用 `v, ok := m[k]`；
- range map 的顺序是随机的，排行榜必须显式排序；
- 排序用 `sort.Slice`，比较函数先比分数、再比姓名；
- 截取前 n 个之前先把 n 夹到 `len(entries)` 以内。

## 评分点

| 测试 | 检查内容 |
|---|---|
| `TestTotalsSumsAndDoesNotMutateInput` | 累加正确；输入不被修改 |
| `TestTotalsEmptyAndNil` | 空成绩、`nil` 成绩、`nil` 输入 |
| `TestTopSortsByScoreThenName` | 分数降序、同分姓名升序 |
| `TestTopEdgeCases` | `n <= 0`、`n` 超过人数、`nil` map |
| `TestTopDoesNotMutateTotals` | 排行榜不修改原 map |

## 做完以后

`Totals` 把“每次查人扫一遍”变成“建表时扫一遍，之后按 key 直接查”。但排行榜需要顺序，map 本身不保留顺序——所以最后仍然要排序。

下一题 `test03` 处理第一课留下的另一半：两个格式化函数写了同样的循环，用 `interface` 把循环只写一次。
