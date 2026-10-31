# 第四课练习：把并发管起来

这一课的八道题对应课件第四课，难度从前往后递增。建议课上做 test01〜test04，课后做 test05〜test08：

| 题目 | 题型 | 主题 | 对应课件 |
|---|---|---|---|
| [test01](test01/) | 运行 → 观察 → 修复 | 等一组 goroutine 全部结束（`sync.WaitGroup`） | 第一部分、第四部分 |
| [test02](test02/) | 从零实现 | 安全计数器：`Mutex` 版与 `atomic` 版 | 第一部分、第五部分、第六部分 |
| [test03](test03/) | 补全代码 | 合并多个 channel、所有权与 `close` | 第二部分、第七部分 |
| [test04](test04/) | 补全代码 | 可取消的流水线：`select` + `context`，不泄漏 | 第三部分、第七部分 |
| [test05](test05/) | 运行 → 观察 → 修复 | 共享内存、`Mutex` / `RWMutex`、快照副本 | 第一部分、第五部分 |
| [test06](test06/) | 代码纠错 | 懒加载（`sync.Once`）、超时与 goroutine 泄漏、`%w` | 第二、三、六部分 |
| [test07](test07/) | 从骨架开始实现 | 固定 worker 数量的任务池、`context` 取消、结果归位 | 第三、四、七部分 |
| [test08](test08/) | 桥接题（第五课预告） | 从 channel 走到字节流：消息边界与拆包 | 第四课全部 + 第五课 `net` |

八道题是一条线：

```text
test01  先建立“谁等谁”的顺序                      ← 第三课 test07 的 Collect
   │
test02  读改写不是一步：锁 / 原子操作              ← 第三课 test03 的交错、test06 的竞态
   │
test03  多个发送方，谁负责关闭输出
   │
test04  消费者提前离开，取消怎么一级级传上去
   │
test05  多个变量要一起一致：锁 + 副本              ← 第二课 test08 的 Snapshot
   │
test06  纠错：只初始化一次、超时后没人接的 channel
   │
test07  把 channel、WaitGroup、context 组合成一个能取消的任务池
   │
test08  最后发现：channel 只能在一个进程里用；
        跨进程/跨机器时，手里只剩一根没有消息边界的字节流 → 第五课
```

`test01` 修的正是第三课 `test07` 留下的 `Collect`：函数返回时 goroutine 还没写完。`test02` 把第三课 `test03` 模拟器里一条条指令交错出来的丢更新变成真实代码。`test05` 里 `Snapshot` 的所有权边界，来自第二课 `test08` 的 `Cache`。

## 先确认环境

- Go 版本以仓库根目录的 `go.mod` 为准（当前 1.25+）；
- 命令都在**练习仓库根目录**执行；
- 每道题的 Go 包在 `testNN/<包名>/` 下，包名与目录名一致；
- 你只改 starter 的 `.go` 文件，**不要改 `_test.go`**。测试就是题目契约，也是评分点。

starter 的初始状态一定是红的，这是题目本身，不是环境坏了。

## 怎么跑测试

```bash
# 整课
go test ./lesson-04-sync/...

# 一道题，带详细输出
go test -v ./lesson-04-sync/test01/runall

# 只跑一个测试函数
go test -run TestRunAllWaitsForEveryTask ./lesson-04-sync/test01/runall

# 并发题建议自查竞态；CI 不跑 -race，也不计入成绩
go test -race ./lesson-04-sync/...

# 偶发问题：重复跑，逼出概率事件
go test -count=10 ./lesson-04-sync/test05/ledger

# 卡死时让测试尽快失败，而不是一直挂着
go test -timeout=10s ./lesson-04-sync/...
```

### 读失败信息

- `--- FAIL: TestXxx` 后面那一行通常就是差异；先看 `want` 和实际值；
- `timed out`、`test timed out`：先找“谁在等谁”。每个阻塞点问一句：**谁会唤醒它？**
- `fatal error: concurrent map writes`：这是 runtime 直接终止进程，不是普通 panic，`recover` 也接不住。它通常说明 map 在没有任何同步的情况下被多个 goroutine 写；
- `WARNING: DATA RACE`：看报告里的两处访问（Read / Previous write）分别在哪个文件的哪一行。**某次结果正确不等于没有竞态**；
- 编译失败先修语法；starter 本身一定能编译，如果你看到一屏语法错误，通常是自己刚改坏的那几行。

### 四句提交前自查

1. 每个 channel 的发送方、接收方、关闭方分别是谁？
2. 每个 goroutine 靠什么退出？每个阻塞点会被谁唤醒？
3. 共享状态的临界区在哪里？锁有没有覆盖完整的读和写？
4. 消费者提前返回时，生产者会不会永远阻塞？

## 交作业

```bash
go test ./lesson-04-sync/...
git add .
git commit -m "完成第四课 test01-test08"
git push
```

push 到 `main` 后，仓库的 CI 会按测试函数的终态评分。**它只跑当前 `config.json` 里填的那一课，也不跑 `-race`**。所以：

- `-race` 是你自己用来找 bug 的，不是用来“刷分”的；
- 测试文件不要改；改了测试等于改了评分点；
- 同一道题的所有测试都通过，这道题才算通过。

## 评分点

每个 `TestXxx` 是一个独立评分点。题目 README 里的“评分点”表列出了每个测试在检查什么；测试函数名就是评分接口，不要重命名。

## 做完之后

如果你按顺序做完，会得到这样一条经验：

| 你遇到的问题 | 用到的工具 |
|---|---|
| 函数返回了，goroutine 还没跑完 | `WaitGroup` / 完成信号 channel |
| 计数器丢更新 | `Mutex` / `atomic` |
| 多个输入，消费者不知道什么时候结束 | `close` + 所有权规则 |
| 消费者提前离开，上游卡住泄漏 | `select` + `context` |
| map 崩溃、快照一半新一半旧 | `Mutex` / `RWMutex` + 副本 |
| 初始化跑了好几次、超时后 goroutine 泄漏 | `sync.Once` + 带缓冲 channel |
| 任务很多、要限制并发、要能取消 | `channel` + `WaitGroup` + `context` |
| 两台机器上的程序怎么交换数据 | **下一课：网络、TCP、HTTP** |

`test08` 已经让你用 `net.Pipe` 摸到了一根“像网线一样的字节管道”。下一课会解释：TCP 为什么也只是字节流、为什么需要 IP 和端口、HTTP 又怎样在这根字节流上划出请求和响应。

## 不要做的事

- 不要用 `time.Sleep` 猜“应该跑完了”：test01 会读源码检查；
- 不要用 `runtime.NumGoroutine()` 写业务逻辑：test04、test06 用它只是为了在测试里观察泄漏；
- 不要为了让测试通过而改测试、删测试或写死返回值；
- 不要把 `sync.Mutex` 复制来复制去，也不要在持有锁时做网络/文件 I/O；
- 不要用 `recover` 去接 `fatal error: concurrent map writes`——它接不住，正确做法是加同步。
