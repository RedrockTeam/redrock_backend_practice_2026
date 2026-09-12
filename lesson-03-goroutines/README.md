# 第三课练习：goroutine

1. 让任务函数返回结果并保持输入顺序；
2. 比较顺序求和与并发拆分的行为；
3. 修复 `03-missing-results` 中立即返回的结果收集器。

第三题先不要猜答案：用测试和日志确认 goroutine 何时启动、何时结束，以及谁负责等待。第四课会用 channel 和 WaitGroup 建立这个边界。

```bash
go test ./lesson-03-goroutines/...
```
