# 第二课练习：切片、map、对象与 interface

承接第一课 `test08`：你在那里量到了线性查找的代价，也看到两个格式化函数写了同样的循环。这一课用四种工具把这些代价和重复消掉：

| 题目 | 主题 | 对应课件 |
|---|---|---|
| [test01](test01/) | 切片的窗口、副本、`append` 与 `range` | 指针和数组、切片 |
| [test02](test02/) | `map` 统计总分、排序取 Top | map |
| [test03](test03/) | 对象、方法、`Formatter` 接口 | 对象、interface |
| [test04](test04/) | map 是句柄、`Snapshot` 的所有权，接到第三课 | map、对象 |

四道题是一条线：

```text
第一课 test08
  切片线性查找 + 两份重复的格式化循环
        │
test01  先分清切片是窗口还是数据：Clone / 原地改 / append
        │
test02  map 把“按学号找人”从扫一遍变成一次查找
        │
test03  interface 把“遍历一次、每行长什么样”拆开，消掉重复循环
        │
test04  map 变量只是句柄；Snapshot 要复制，
        否则调用方就成了没有拿锁的写者 → 第三课 goroutine
```

## 先确认环境

- Go 版本以仓库根目录的 `go.mod` 为准；
- 命令都在**练习仓库根目录**执行；
- 每道题的 Go 包在 `testNN/<包名>/` 下，包名与目录名一致；
- 只改 starter 的 `.go` 文件，**不要改 `_test.go`**。

starter 初始是红的，那是题目，不是环境坏了。

## 怎么跑测试

```bash
# 整课
go test ./lesson-02-collections/...

# 一道题，带详细输出
go test -v ./lesson-02-collections/test01/sliceops

# 只跑一个测试
go test -run TestCloneIsIndependent ./lesson-02-collections/test01/sliceops

# 看竞态（这一课还没有并发，但 test04 的桥接测试会启动两个 goroutine）
go test -race ./lesson-02-collections/...

# 卡死时尽快失败
go test -timeout=10s ./lesson-02-collections/...
```

## 这一课要带走的四句话

1. **切片赋值复制的是窗口，不是数据。** `b := a` 之后改 `b[0]`，`a[0]` 也会变；要独立副本就 `make + copy`。
2. **`range` 给的 `v` 是副本。** 想改原切片必须用下标 `s[i]`。
3. **`append` 的返回值必须接住。** 它可能扩容，返回新的切片头。
4. **map 也是句柄。** 函数里改 map 的内容外部可见；`Snapshot` 想把状态交出去，就必须复制。

## 和第一课的连接

第一课 `test08` 里的 `FindByID` 是线性查找，`Comparisons()` 量出了 200 条名册最坏要比较 200 次。`test02` 的 `map` 把“查人”变成一次键查找；`test03` 的接口把两个格式化函数收进同一个 `[]Formatter`。

## 做完之后

`test04` 的 `Cache` 在**单个执行流**里是正确的。它的 `Put` / `Get` 没有锁，`Snapshot` 如果不复制，调用方还能绕过缓存直接改 map。

第三课会启动多个 goroutine，让同一个 `Cache` 同时被读写；第四课再用 `Mutex` / `RWMutex` 和 `WaitGroup` 把顺序和所有权补上。

## 评分

每个 `TestXxx` 是一个评分点；一道题的所有测试都通过，这道题才算通过。CI 只跑 `config.json` 里填的那一课，也不跑 `-race`。

## 不要做的事

- 不要为了让测试通过而修改测试或写死返回值；
- 不要用 `append` 的返回值做“原地修改”的假设；
- 不要在 `range` 里改 `v` 然后期待原切片变化；
- 不要把 `Snapshot` 直接写成 `return c.values`。
