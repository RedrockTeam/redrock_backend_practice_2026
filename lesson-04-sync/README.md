# 第四课练习：channel、sync 与 worker pool

1. 合并多个输入 channel，并在所有输入结束后关闭输出；
2. 用 Mutex 或其他合适的同步方式修复并发计数器；
3. 实现固定 worker 数量的任务池，保留任务顺序并处理错误/取消。

结课要求同时通过功能测试和竞态检测。锁只保护共享状态，channel 负责传值和结束信号，WaitGroup 负责等待。

```bash
go test ./lesson-04-sync/...
go test -race ./lesson-04-sync/...
```
