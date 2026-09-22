# 第四课练习：channel、sync 与 worker pool

| 题目 | 主题 |
|---|---|
| [test01](test01/) | 合并多个输入 channel，所有输入结束后关闭输出 |
| [test02](test02/) | 用 Mutex 或其他同步方式修复并发计数器 |
| [test03](test03/) | 实现固定 worker 数量的任务池，保留顺序并处理错误与取消 |

结课要求功能测试全部通过。另外建议自己跑一次 `go test -race` 自查——**CI 不跑竞态检测，它也不计入成绩**。分工记牢：锁只保护共享状态，channel 负责传值和结束信号，WaitGroup 负责等待。

`test02` 修的正是并发计数器的竞态。四课到这里连成一条完整的数据流。

```bash
go test ./lesson-04-sync/...
go test -race ./lesson-04-sync/...
```
