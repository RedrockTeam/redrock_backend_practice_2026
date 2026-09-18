# 练习仓库总体架构

## 一次提交会发生什么

```text
学生 fork 本仓库 → 做题 → push 到自己的 fork → 向本仓库发起 Pull Request
                                                        │
                        ┌───────────────────────────────┘
                        ▼
  ① grade.yml（不可信上下文：fork 的 PR 拿不到 secrets，token 只读）
       ├─ 递归发现仓库里所有含 *_test.go 的目录（跳过 testdata/、vendor/）
       ├─ 对每个目录独立执行 go test -json -count=1
       │    ├─ 原始 JSONL / 日志 → test-results/（只进 artifact，不进 result.json）
       │    └─ 逐测试函数的终态 → result.json、score.md、Step Summary
       ├─ go test -race ./...（仅作反馈，不影响分数）
       ├─ PR 事件下记录 pull-request-number.txt
       └─ 上传 artifact `go-test-result`
                        │ workflow_run: completed
                        ▼
  ② report-result.yml（基仓库上下文：可读 secrets，文件始终取自默认分支）
       ├─ 按 run-id 下载 ① 的 artifact
       ├─ stamp_result.py：用事件里可信的值覆盖身份，重算所有聚合
       └─ curl → 成绩网站
                        │
                        ▼
              PR 页面的 check ＋ 成绩网站上的排名
```

两个 workflow 的分工不是为了好看，是被 GitHub 的安全模型强制的：`pull_request` 事件在 fork 上下文里**没有 secrets**，所以上报必须放在 `workflow_run` 里；而 `workflow_run` 永远执行默认分支上的那份 workflow 文件，学生改不动它。

## 为什么不用 LESSON_PATH 裁剪

早期版本用仓库变量 `LESSON_PATH` 只跑当前课。现在改成**每次 push、每次 PR 都把整个仓库跑一遍**：

- 学生做第三课时，第一课的代码不该悄悄坏掉；回归由 CI 兜住；
- 成绩网站需要的是「这个人现在全部课程的状态」，而不是某一次运行恰好覆盖的那部分；
- 单个包编译失败只产生一个 `__build__` 失败项，不会连累其他题。

仓库变量 `TEST_ROOT` 仍然保留，默认 `.`，只在需要临时缩小范围时设置。

## 目录与评分的对应关系

```text
lesson-01-basics/test06/roster/roster_test.go
└── lesson ─┘└─ test ┘└─ package ┘
```

`run_go_tests.py` 的 `attribute()` 把包路径映射成 `(lesson, test)`，`result.json` 按 `lessons[] → tests[] → packages[] → cases[]` 四层嵌套。成绩网站因此可以直接按「第几课第几题」寻址，不需要解析路径字符串。

题目单位叫 `testNN/`，真正的 Go 包放在它下一层且用有意义的名字（`hellogo`、`syntax`、`roster`）。这是刻意的：第一课第四部分刚教完「包名通常与目录名一致，避免 `util`、`common` 这类名字」，如果题目目录本身就是 `package test06`，课件当场自相矛盾。

## 信任边界

`result.json` 是在学生的 fork 运行里生成的，**里面每个字段都是学生可以改的**。`stamp_result.py` 在上报前重建可信部分：

| 字段 | 处理方式 |
|---|---|
| `submission.*` | 整块丢弃，用 `workflow_run` 事件里的 commit / fork 名 / actor 重写 |
| `result_id` | 用可信值重新 sha256 |
| `submission.pull_request` | 调 API 核对该 PR 的 head 是否等于本次 commit，不等则置 `null` |
| `status` / `score` / `summary` | 全部从叶子 `cases[].status` 重算，非法状态值直接剔除 |
| `cases[]` | 保持学生侧结果——测试本来就跑在他们的机器镜像上，这是这套流程固有的 |

换句话说：能改的只有「测试跑出了什么」，改不了「这是谁交的、加起来是多少分」。

## 边界

- 公开练习仓库：题目 README、starter、公开测试、两个 workflow、评分脚本；
- 仓库配置：`RESULT_UPLOAD_ENDPOINT`、`RESULT_UPLOAD_METHOD`、`TEST_ROOT`（变量）与 `RESULT_UPLOAD_TOKEN`（secret）；
- 成绩网站：接收、存储、展示排名，实现不在本仓库；
- 教师课件：只在 `redrock/courseware` 维护，不参与 CI。
