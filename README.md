# 红岩网校 Go 练习仓库

学生侧练习仓库，和教师课件仓库 `redrock/courseware` 完全独立。这里放题目说明、starter 代码、公开测试和 CI；不放课件讲稿、答案和评分服务实现。

## 目录

| 路径 | 内容 |
|---|---|
| `lesson-01-basics/` | 第一课：Go 基础、package、可见性、写测试，共 6 题 |
| `lesson-02-collections/` | 第二课：map 与 interface，最后观察 map 竞态 |
| `lesson-03-goroutines/` | 第三课：计算机基础、goroutine、结果收集 |
| `lesson-04-sync/` | 第四课：channel、Mutex、WaitGroup、worker pool |
| `.github/workflows/grade.yml` | 跑测试、产出结果 |
| `.github/workflows/report-result.yml` | 把结果上报到成绩网站 |

每一课的目录结构都一样：

```text
lesson-01-basics/
  README.md          课次总览
  test01/
    README.md        这一题要做什么、对应课件哪一节、评分点是什么
    hellogo/         真正的 Go 包，包名和目录名一致
  test02/
  ...
```

题目单位是 `testNN/`，Go 包放在它下面一层。这样每道题的说明、素材和包边界互不干扰，评分结果也能精确定位到「第几课第几题」。

## 怎么做题、怎么交

```text
① Fork 本仓库到你自己的账号
② git clone 你的 fork，新建一个分支
③ 按 lesson-01-basics/README.md 的顺序做题，本地 go test 跑通
④ commit 并 push 到你的 fork
⑤ 向本仓库发起 Pull Request
⑥ CI 自动全量跑一遍所有测试，结果发送到成绩网站
```

第 ⑥ 步不需要你做任何事。PR 页面下方的 check 会显示通过与否，点进去能看到每道题的得分表。

## 评分规则

CI 会递归发现仓库里**所有**包含 `_test.go` 的目录（跳过 `testdata/`、`vendor/`），逐个执行 `go test -json`，然后按测试函数的终态计分：

```text
积分 = 通过的测试数 / 全部测试数 × 100，取整
```

- 每个 `TestXxx` 是一个独立评分点，子测试归并到父测试；
- 某一道题编译失败只会产生一个 `__build__` 失败项，**不会**影响其他题目继续评分；
- 每次 push 和每次 Pull Request 都把所有课次完整跑一遍，不做按课次的裁剪。

`starter` 的初始状态一定是红的，那是题目，不是环境坏了。

## result.json：只包含结果

CI 产出一份 `result.json`，里面**只有结论**——状态、分数、每个测试函数的 pass/fail/skip，没有任何测试输出、日志或本地路径。完整的 `go test -json` 原始流保存在同一次运行的 `go-test-result` artifact 里，只给调试用。

```json
{
  "schema_version": "2.0",
  "status": "fail",
  "score": 92,
  "submission": {
    "result_id": "4ee19990932cb71a",
    "repository": "redrock/teaching-exercises",
    "head_repository": "student42/teaching-exercises",
    "commit": "9f2c...",
    "pull_request": 128,
    "actor": "student42",
    "verified": true
  },
  "summary": { "total": 26, "passed": 24, "failed": 2, "skipped": 0 },
  "lessons": [
    {
      "lesson": "lesson-01-basics",
      "status": "fail",
      "score": 92,
      "tests": [
        {
          "test": "test06",
          "status": "fail",
          "score": 80,
          "packages": [
            {
              "package": "lesson-01-basics/test06/roster",
              "status": "fail",
              "cases": [{ "name": "TestFindByID", "status": "fail", "ms": 0 }]
            }
          ]
        }
      ]
    }
  ]
}
```

## 上报是怎么做到安全的

fork 发来的 Pull Request 在 `pull_request` 事件里**拿不到本仓库的 secrets**，这是 GitHub 的设计。所以评分和上报拆成了两个 workflow：

| workflow | 触发 | 上下文 | 职责 |
|---|---|---|---|
| `grade.yml` | `push` / `pull_request` | fork 的不可信上下文，无 secrets | 跑测试，产出 `result.json`，存成 artifact |
| `report-result.yml` | `workflow_run` | 基仓库上下文，可读 secrets | 取 artifact，核验身份，发送到成绩网站 |

`workflow_run` 始终执行默认分支上的那份 workflow 文件，学生改不了它。上报前 `stamp_result.py` 会：

- 用事件里可信的 commit、fork 名、发起人覆盖 `submission`，并据此重算 `result_id`；
- 调 API 核验 PR 号确实指向同一个 commit，对不上就丢弃这个号；
- 从每个测试的 pass/fail 重新计算所有 `status`、`score` 和 `summary`，改总分没有用。

配置项（未配 endpoint 时整个上报 workflow 直接跳过，不影响评分）：

| 类型 | 名称 | 说明 |
|---|---|---|
| 仓库变量 | `RESULT_UPLOAD_ENDPOINT` | 成绩网站收件路由，如 `https://grade.example.edu/api/v1/results` |
| 仓库变量 | `RESULT_UPLOAD_METHOD` | 可选，默认 `POST` |
| 仓库变量 | `TEST_ROOT` | 可选，默认 `.`（全仓库）。只在需要临时缩小范围时设置 |
| 仓库 secret | `RESULT_UPLOAD_TOKEN` | 可选，作为 `Authorization: Bearer` 发送 |

请求带 `Idempotency-Key: <result_id>`，重跑同一次运行不会产生重复记录。

## 本地运行

```bash
go test ./...                                   # 全部课次
go test ./lesson-01-basics/...                  # 一整课
go test -v ./lesson-01-basics/test06/roster     # 一道题，带详细输出
go test -run TestSumTo ./lesson-01-basics/test02/syntax   # 一个测试函数
go test -race ./...                             # 竞态检测
go vet ./...

# 在本地复现 CI 的评分过程
python3 .github/scripts/run_go_tests.py --root .
```

`-race` 的失败是课程设计的一部分：第二课会让你看到 map 的竞态，第四课再修它。
