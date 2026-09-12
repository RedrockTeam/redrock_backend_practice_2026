# 练习仓库总体架构

```text
学生接受某一课的 Classroom assignment
              │ push / Pull Request
              ▼
GitHub Actions（setup-go）
  ├─ go test -json "./${LESSON_PATH}/..."
  │    └─ 统计 pass / fail / skip，写入 Step Summary 和 score.md
  ├─ go test -race "./${LESSON_PATH}/..."（反馈）
  └─ 上传 test-results.json 与 score.md
              │
              ▼
Classroom check 结果
```

`LESSON_PATH` 是仓库变量，例如 `lesson-03-goroutines`。这样四课可以从同一份练习素材生成四个 Classroom template，学生每次只看到当前课的测试。课件仓库不参与运行，评分不需要自定义 Web 服务；将来若要做排行榜，读取 GitHub check-run 或 artifact 即可。

## 边界

- 公开仓库：题目、starter、公开测试、workflow；
- Classroom 私有配置：隐藏测试、仓库变量和必要的权限；
- CI：执行 Go 测试、计算通过率、上传结果；
- 教师课件：只在 `redrock/courseware` 仓库维护。
