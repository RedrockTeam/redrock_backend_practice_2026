# 第三课练习：计算机基础与 goroutine

承接第二课 `test03` 的竞态报告：这一课解释并发到底是怎么跑的。

| 题目 | 主题 |
|---|---|
| [test01](test01/) | 让任务函数返回结果并保持输入顺序 |
| [test02](test02/) | 比较顺序求和与并发拆分的行为 |
| [test03](test03/) | 修复立即返回、丢结果的收集器 |

`test03` 启动了 goroutine 却立即返回，测试会稳定暴露结果缺失。先不要猜答案：用测试和日志确认 goroutine 何时启动、何时结束，以及谁负责等待。第四课用 channel 和 WaitGroup 建立这条边界。

```bash
go test ./lesson-03-goroutines/...
go test -race ./lesson-03-goroutines/...
```
