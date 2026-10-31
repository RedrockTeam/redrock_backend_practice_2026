# 第三课练习：从 CPU 到 goroutine

承接第二课 `test08`：那个 `Cache` 只在单个执行流里正确。这一课先回答“执行流到底是什么”——从 CPU 一条一条执行指令讲起，到操作系统在线程之间切换，再到 Go 用 `go` 关键字启动成千上万个 goroutine。最后亲眼看到“函数已经返回，goroutine 还没结束”。

七道题，**按顺序做**，难度递增：

| 题目 | 包 | 题型 | 主题 | 对应课件 |
|---|---|---|---|---|
| [test01](test01/) | `riscv` | 补全 | 小端序、位运算、拆一条 RISC-V 指令 | cpu 工作简介 |
| [test02](test02/) | `toycpu` | 从零实现 | 取指 → 译码 → 执行 → 写回 → 更新 PC | cpu 工作简介 |
| [test03](test03/) | `interleave` | 实现 + 观察 | 两个线程共享内存，按调度表交错，数出丢更新 | 协程（进程 / 线程） |
| [test04](test04/) | `spawn` | 补全 | `go` 不阻塞；一万个 goroutine 有多轻 | 协程、协程的 infra、go 关键字 |
| [test05](test05/) | `asyncprint` | 运行观察 → 补全 | 多个 goroutine 写同一个输出 | go 关键字的使用 |
| [test06](test06/) | `parallel` | 改写 | 把求和拆给两个 goroutine，拿回结果 | go 关键字的使用 |
| [test07](test07/) | `missing` | 观察 + 量尺 | 启动 goroutine 却立即返回，结果缺失 | 总结 |

## 七题的关系

```text
test01  一条指令在内存里长什么样，CPU 怎么拆开它
   │
test02  CPU 怎么一条接一条地执行；x++ 原来是三条指令
   │
test03  两个线程共用内存、在任意两条指令之间被切换 → 大多数交错都会丢更新
   │      （第四课 test02 用 Mutex / atomic 解决）
test04  Go 的执行流：go 立即返回，goroutine 便宜到可以开一万个
   │
test05  多个 goroutine 同时往外写
   │
test06  把一份工作拆出去，再把结果拿回来
   │
test07  函数返回得比 goroutine 快 → 结果缺失
        你能量出它，但第三课的工具修不了它
        （第四课 test01 用 WaitGroup 补上“谁等谁”）
```

前三题没有一个 `go` 关键字，但它们解释了后四题里所有“奇怪”的现象：为什么 goroutine 执行顺序不固定、为什么 `x++` 会丢、为什么 `Sleep` 只是在猜。

## 第三课的 `time.Sleep` 怎么用

第三课还没有 `WaitGroup` 和 `channel`，所以 `test05` / `test06` 允许在**函数内部**用 `time.Sleep` 暂时等一会儿，先把现象观察出来。这不是最终答案：

- `test04` 的函数只负责启动，不需要等；
- `test05` 里函数自己等，所以输出能看全；
- `test06` 里函数自己等，所以两段和能加起来；
- `test07` 的 `Collect` 要把结果返回给调用方，Sleep 解决不了“返回之后谁还活着”的问题——这一题不修它，只量它。

第四课会用 `WaitGroup` 和 `channel` 替换掉这些临时等待。

## 怎么跑测试

```bash
# 整课
go test ./lesson-03-goroutines/...

# 一道题，带详细输出
go test -v ./lesson-03-goroutines/test01/riscv
go test -v ./lesson-03-goroutines/test02/toycpu
go test -v ./lesson-03-goroutines/test03/interleave
go test -v ./lesson-03-goroutines/test04/spawn
go test -v ./lesson-03-goroutines/test05/asyncprint
go test -v ./lesson-03-goroutines/test06/parallel
go test -v ./lesson-03-goroutines/test07/missing

# 只看桥接测试打印的数字
go test -v -run TestBridge ./lesson-03-goroutines/...

# test05 通过带锁的 buffer 收输出，-race 应该干净
go test -race ./lesson-03-goroutines/test05/asyncprint

# test06 如果直接用 Sleep 等结果，-race 会报告数据竞争；
# 这不是环境坏了，而是第三课故意留下来的问题，第四课用 WaitGroup 修。
go test -race ./lesson-03-goroutines/test06/parallel
```

starter 初始是红的，那是题目，不是环境坏了。七道题在第三课的知识范围内都**能做到全绿**。

## 这一课要带走的四句话

1. CPU 一次只执行一条指令；`x++` 是读、改、写三条，**任意两条之间都可能被切走**；
2. `go f()` 只是“创建一个 goroutine 并让它变成可运行”，**不会等它**；
3. 主 goroutine 返回，进程就结束，其他 goroutine 会被静默杀死；
4. 循环里的 goroutine 最好把当前值作为参数传进去：`go func(n int){...}(i)`，不要依赖“它应该读到我想要的那个 i”。

## 不要做的事

- 不要修改 `_test.go`；
- 不要改 `test07/missing/collect.go`，它是观察对象，也是第四课 `test01` 的起点；
- 不要在 `answers.go` 里逐个试选项——测试只告诉你对错，依据在课件里。

## 做完之后

`test07` 的 `Collect` 会稳定地少结果，你写的 `Verify` 能把少了多少量出来。想一想：**调用方**怎么知道所有 goroutine 都写完了？

第四课 `test01` 从这个缺口开始：用 `sync.WaitGroup` 建立 `Done → Wait` 的顺序；`test02` 用 `Mutex` / `atomic` 解决 `test03` 里数出来的丢更新；之后再讲 channel 的关闭和所有权。第二课 `test08` 那个没有锁的 `Cache`，会在第四课 `test05` 的账本里被真正管起来。
