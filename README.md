# 红岩网校 Go 练习仓库

这是 GitHub Classroom 使用的学生侧仓库，和教师课件仓库完全独立。这里放题目说明、starter 代码、公开测试和 CI；不放课件讲稿、答案、隐藏测试或自定义评分服务。

## 目录

- `lesson-01-basics/`：`var`、`func`、`struct`、slice 和 Git 提交练习；
- `lesson-02-collections/`：map、方法、interface，最后观察 map 竞态；
- `lesson-03-goroutines/`：计算机基础、goroutine、等待和结果收集；
- `lesson-04-sync/`：channel、Mutex、WaitGroup 和 worker pool；
- `.github/workflows/grade.yml`：唯一评分入口，复杂性留在 CI；
- `ai-suggestions/`：本地 AI 讨论稿，已被 `.gitignore` 排除。

## 简单评分

CI 运行 `go test -json ./...`，按终态为 `pass` 的测试数计算：

```
积分 = 通过测试数 /（通过 + 失败 + 跳过）× 100，取整
```

编译失败或没有可统计测试时为 0 分。GitHub Classroom 直接把 workflow check 作为提交结果；若需要排行榜，后续只需读取各仓库的 check-run 结论或导出的 JSON，不需要在这里维护评分服务器。

## 本地运行

```bash
go test ./lesson-01-basics/...
go test ./lesson-02-collections/...
go test -race ./lesson-02-collections/03-map-race/...
```

bridge 题的失败是课程设计的一部分：先记录现象，再在下一课用新概念修复。
