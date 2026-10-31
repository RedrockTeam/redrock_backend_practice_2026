# 第二课练习：指针、切片、map、对象与 interface

承接第一课 `test08`：你在那里量到了线性查找的代价，也看到两个格式化函数写了同样的循环。这一课先把“值还是句柄”分清楚，再用 map 和 interface 把那两个代价消掉，最后交出一个没有同步的共享 `Cache` 给第三课。

八道题，**按顺序做**，难度逐步上升。题型交替出现：补全、纠错、从零实现。

| 题目 | 包 | 题型 | 主题 | 对应课件 | 预计 |
|---|---|---|---|---|---|
| [test01](test01/) | `pointer` | 补全 + 预测输出 | 数组是值、指针传地址、结构体指针 | 指针和数组 | 30 min |
| [test02](test02/) | `sliceops` | 运行 → 观察 → 修复 | 切片是窗口：`Clone` / 原地改 / `append` 返回值 | 切片 | 30 min |
| [test03](test03/) | `aliasing` | 代码纠错 | 五个切片别名 bug：cap 富余时 append 原地写入、边删边遍历、`&v` | 切片（append、copy、range） | 45 min |
| [test04](test04/) | `scoreboard` | 补全 | `map` 统计总分、排序取 Top | map | 40 min |
| [test05](test05/) | `dns` | **从零实现** | 带反查的 DNS 缓存：两份索引保持一致 | map + 课件小练习 | 60 min |
| [test06](test06/) | `formatter` | 补全 | `Formatter` 接口消掉重复循环 | 对象、interface | 30 min |
| [test07](test07/) | `pitfalls` | 纠错 + 补全 | typed nil、指针接收者、`type switch`、`Stringer`、不可比较 | interface 常见坑 | 60 min |
| [test08](test08/) | `cache` | 补全 + 下一课桥接 | map 是句柄，`Snapshot` 的所有权 | map、对象 | 30 min |

## 八道题是一条线

```text
第一课 test08
  切片线性查找 + 两份重复的格式化循环
        │
test01  数组 b := a 复制全部；切片 b := a 只复制切片头 → 值还是句柄？
        │
test02  切片是窗口：Clone 要复制、range 的 v 是副本、append 要接返回值
        │
test03  调用方继续 append 时，两个窗口互相覆盖 → 收紧 cap、复制、用下标取地址
        │
test04  map：把“按名字查”从扫一遍变成一次查找（第一课代价 1、2）
        │
test05  从零写 DNS 缓存：正查、反查两张表必须同步改
        │
test06  interface：遍历写一次，每行长什么样交给实现（第一课代价 3）
        │
test07  interface 的坑：nil 指针装进 error、值接收者改副本
        │
test08  map 是句柄，Snapshot 要复制；
        没有锁的 Cache → 第三课 goroutine、第四课 test05 账本
```

## 先确认环境

- Go 版本以仓库根目录的 `go.mod` 为准；
- 命令都在**练习仓库根目录**执行；
- 每道题的 Go 包在 `testNN/<包名>/` 下，包名与目录名一致；
- 只改 starter 的 `.go` 文件，**不要改 `_test.go`**（`test01` 的 `demo.go` 也不要改）。

starter 初始是红的，那是题目，不是环境坏了。名字带 `Bridge` 的测试只是把现象摆出来，starter 下也会通过。

## 怎么跑测试

```bash
# 整课
go test ./lesson-02-collections/...

# 一道题，带详细输出（Bridge 测试的 t.Log 只在 -v 下可见）
go test -v ./lesson-02-collections/test03/aliasing

# 只跑一个测试
go test -run TestHeadDoesNotLeakCapacity ./lesson-02-collections/test03/aliasing

# 看竞态（这一课还没有并发，但 test08 的桥接测试会启动两个 goroutine）
go test -race ./lesson-02-collections/...

# 卡死时尽快失败
go test -timeout=10s ./lesson-02-collections/...
```

## 这一课要带走的五句话

1. **数组是值，切片和 map 是句柄。** `b := a` 对数组复制全部元素；对切片只复制 (指针, len, cap)。
2. **`range` 给的 `v` 是副本。** 想改原切片、想取元素地址，都用下标 `s[i]` / `&s[i]`。
3. **`append` 的返回值必须接住；cap 富余时它原地写入。** 要把切片交给别人继续 append，就收紧 cap（`s[a:b:b]`）或复制。
4. **map 的顺序是随机的，要有序就自己排；两份索引要一起改。**
5. **接口值 = (动态类型, 动态值)。** 只有两者都空才 `== nil`；需要改状态的方法用指针接收者。

## 和第三课的连接

`test05` 的 `Resolver` 和 `test08` 的 `Cache` 在**单个执行流**里都正确：每次修改都一口气改完，中间没人能看到“改了一半”的状态。

第三课会从 CPU 讲起：`count++` 在机器上是读、加、写三条指令，两个执行流交错执行时更新会丢。之后再启动真正的 goroutine，第四课用 `Mutex` / `RWMutex` 和 `WaitGroup` 把顺序和所有权补上。

## 评分

每个 `TestXxx` 是一个评分点；一道题的所有测试都通过，这道题才算通过。CI 只跑 `config.json` 里填的那一课，也不跑 `-race`。

## 不要做的事

- 不要为了让测试通过而修改测试或写死返回值；
- 不要用 `append` 的返回值做“原地修改”的假设；
- 不要在 `range` 里改 `v` 或取 `&v` 然后期待原切片变化；
- 不要把 `Snapshot`、`ReverseLookup` 写成直接返回内部数据；
- 返回 `error` 时不要返回一个类型化的 `nil` 指针。
