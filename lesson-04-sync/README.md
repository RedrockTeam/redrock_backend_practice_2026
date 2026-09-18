# 第四课练习：channel、sync 与 worker pool

| 题目 | 主题 |
|---|---|
| [test01](test01/) | 合并多个输入 channel，所有输入结束后关闭输出 |
| [test02](test02/) | 用 Mutex 或其他同步方式修复并发计数器 |
| [test03](test03/) | 实现固定 worker 数量的任务池，保留顺序并处理错误与取消 |

结课要求同时通过功能测试和竞态检测。分工记牢：锁只保护共享状态，channel 负责传值和结束信号，WaitGroup 负责等待。

`test02` 修的正是第二课 `test03` 留下的那个 map 竞态。四课到这里连成一条完整的数据流。

```bash
go test ./lesson-04-sync/...
go test -race ./lesson-04-sync/...
```
