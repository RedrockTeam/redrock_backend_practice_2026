# 红岩网校 Go 练习仓库

本目录是 GitHub Classroom template 的本地镜像，独立 Git 仓库。每个练习目录都有 starter 代码和公开测试；部分 bridge 题刻意失败，用于引出下一节课。不要把讲师隐藏测试或评分服务 token 放进本仓库。

## 课程顺序

| 目录 | 先完成 | bridge 现象 |
| --- | --- | --- |
| `lesson-01-basics` | 变量、函数、struct、slice | 线性查找和重复格式化不易扩展 |
| `lesson-02-collections` | map、方法、interface | 并发访问 map 在 `-race` 中失败 |
| `lesson-03-goroutines` | goroutine、等待、结果收集 | goroutine 没有可靠等待导致结果缺失 |
| `lesson-04-sync` | channel、Mutex、WaitGroup、worker pool | 功能测试和竞态测试共同验收 |

## 本地运行

```bash
go test ./lesson-01-basics/...
go test ./lesson-02-collections/...
go test -race ./lesson-02-collections/03-map-race/...
```

starter 的 bridge 题失败是预期现象。提交前运行当前小题的测试，并阅读测试名称、失败值和 `-race` 报告。

## Classroom Actions

`.github/workflows/grade.yml` 会保存完整 `go test -json` artifact、运行竞态检测，并在设置 `SCORE_API_URL` 后调用评分服务。需要在 Classroom 仓库的 Actions secrets 中配置：

- `SCORE_API_URL`：例如 `https://score.example.edu/api/grade`；
- `SCORE_API_TOKEN`：由部署者生成，脚本只放在请求头中。

工作流中的 `student` 默认使用 `GITHUB_ACTOR`，正式课程也可以改为 Classroom 分配的匿名编号。
