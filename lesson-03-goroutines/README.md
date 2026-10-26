# 第三课练习：goroutine、并发拆分与缺失的等待

承接第二课 `test04`：那个 `Cache` 只在单个执行流里正确。这一课先把多个执行流真正跑起来，并亲眼看到“函数已经返回，goroutine 还没结束”。

| 题目 | 主题 | 对应课件 |
|---|---|---|
| [test01](test01/) | `go` 关键字、闭包参数、先让输出发生 | 协程、`go` 关键字的使用 |
| [test02](test02/) | 把求和拆给两个 goroutine | 协程、`go` 关键字的使用 |
| [test03](test03/) | 启动 goroutine 却立即返回，结果缺失 | 总结里的“你目前没法解决的问题” |

## 三题的关系

```text
test01  让 goroutine 真的启动，看到输出
   │
test02  把一份工作拆给多个 goroutine
   │
test03  函数返回得比 goroutine 快 → 结果缺失
        这一题在第三课结束时仍然是红的；
        第四课 test01 用 WaitGroup 补上“谁等谁”
```

## 第三课的 `time.Sleep` 怎么用

第三课还没有 `WaitGroup` 和 `channel`，所以 `test01` / `test02` 允许在**函数内部**用 `time.Sleep` 暂时等一会儿，先把现象观察出来。这不是最终答案：

- `test01` 里函数自己等，所以输出能看全；
- `test02` 里函数自己等，所以两段和能加起来；
- `test03` 里函数**不能等**——它要把结果返回给调用方，Sleep 解决不了“返回之后谁还活着”的问题。

`test03` 的测试会直接拒绝 `time.Sleep` 的修法。第四课会用 `WaitGroup` 和 `channel` 替换掉这些临时等待。

## 怎么跑测试

```bash
# 整课（test03 预期是红的）
go test ./lesson-03-goroutines/...

# 看每个 worker 的输出 / 失败信息
go test -v ./lesson-03-goroutines/test01/asyncprint
go test -v ./lesson-03-goroutines/test02/parallel
go test -v ./lesson-03-goroutines/test03/missing

# test01 通过带锁的 buffer 收输出，-race 应该干净
go test -race ./lesson-03-goroutines/test01/asyncprint

# test02 如果直接用 Sleep 等结果，-race 会报告数据竞争；
# 这不是环境坏了，而是第三课故意留下来的问题，第四课用 WaitGroup 修。
go test -race ./lesson-03-goroutines/test02/parallel
```

starter 初始是红的；`test03` 的红是**题目的一部分**，不是让你在第三课硬修。

## 这一课要带走的三句话

1. `go f()` 只是“创建一个 goroutine 并让它变成可运行”，**不会等它**；
2. 主 goroutine 返回，进程就结束，其他 goroutine 会被静默杀死；
3. 循环里的 goroutine 最好把当前值作为参数传进去：`go func(n int){...}(i)`，不要依赖“它应该读到我想要的那个 i”。

## 做完之后

`test03` 的 `Collect` 会稳定地返回一堆 0。先不要用 `Sleep` 掩盖：想一想，**调用方**怎么知道所有 goroutine 都写完了？

第四课 `test01` 会从这个缺口开始：用 `sync.WaitGroup` 建立 `Done → Wait` 的顺序；`test02` 再讲 channel 的关闭和所有权。
