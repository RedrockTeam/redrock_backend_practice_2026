# 红岩网校 Go 练习仓库

学生侧练习仓库，和教师课件仓库 `redrock/courseware` 完全独立。这里放题目说明、starter 代码、公开测试和 CI；不放课件讲稿、答案和评分服务实现。

## 目录

| 路径 | 内容 |
|---|---|
| `lesson-01-basics/` | 第一课：Go 基础、语法迁移、package 与 module、可见性、写测试，共 8 题 |
| `lesson-02-collections/` | 第二课：map 与 interface，最后观察 map 竞态 |
| `lesson-03-goroutines/` | 第三课：计算机基础、goroutine、结果收集 |
| `lesson-04-sync/` | 第四课：channel、Mutex、WaitGroup、worker pool |
| `config.json` | **你要填的文件**：姓名、本次课次；中心站点地址老师已预置 |
| `.github/workflows/grade.yml` | 唯一的 workflow：校验配置 → 按课次评分 → 上报 |
| `.github/scripts/` | 评分与上报的脚本，可在本机直接运行 |

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
① 点仓库右上角的「Use this template」建出属于你自己的仓库
② git clone 到本机
③ 填好仓库根目录的 config.json（见下一节）
④ 按 lesson-0X-*/README.md 的顺序做题，本机 go test 跑通
⑤ commit 并 push 到 main
```

push 到 **main** 之后 CI 自己跑完三件事：校验 `config.json` → **只跑你填写的那一课** → 把成绩上报到中心站点。

其它分支和 Pull Request **都不会触发**——同一份作业只会产生一份成绩，不会互相打架。想在别的分支上试，用本机的 `go test` 就好。

> 用模板建出来的仓库不是 fork，GitHub Actions 默认就是开着的，不需要额外去开启。

## config.json

三个字段，你只需要填前两个：

```json
{
  "name": "张三",
  "lesson": "lesson-01-basics",
  "center": "https://grade.example.com"
}
```

| 字段 | 必填 | 说明 |
|---|---|---|
| `name` | 是 | 姓名，看板上显示用 |
| `lesson` | 是 | 本次完成的课次，取值必须是仓库里的课节目录名，如 `lesson-01-basics` |
| `center` | 否 | 中心站点地址，老师已预置，一般不用改 |

`lesson` 决定这次跑哪一课。它**缺失或填错时 CI 直接报错中止**——既不跑测试，也不上报：没有明确课次的成绩无法归入任何一课。报错信息会列出所有可选课次，改好重新 push 即可。

## 评分规则

CI 递归发现指定课次下**所有**包含 `_test.go` 的目录（跳过 `testdata/`、`vendor/`），逐个执行 `go test -json`，然后按测试函数的终态计分：

```text
积分 = 通过的测试数 / 全部测试数 × 100，取整
```

- 每个 `TestXxx` 是一个独立评分点，子测试归并到父测试；
- 某一道题编译失败只会产生一个 `__build__` 失败项，**不会**影响其他题目继续评分；
- 跑起来但没走到终态（例如并发写 map 把进程打死）会归因到触发它的那个测试函数，同样记为失败。

`starter` 的初始状态一定是红的，那是题目，不是环境坏了。

## 上报的载荷

CI 把结果压成下面这份 JSON 发给中心站点。**只留看板用得到的两件事**：这是哪一课、每道题做完没有。分数、用例名、包名、耗时、测试输出一概不进载荷——它们留在 `result.json` 与 artifact `go-test-<sha>` 里，需要时去那里取。

```json
{
  "repo_url": "https://github.com/zhangsan/redrock_backend_practice_2026",
  "commit": "9f2c1ab7c3d4e5f60718293a4b5c6d7e8f901234",
  "ref": "refs/heads/main",
  "event": "push",
  "config": { "name": "张三", "lesson": "lesson-01-basics" },
  "result": {
    "lesson": "lesson-01-basics",
    "tests": [
      { "test": "test01", "status": "fail" },
      { "test": "test02", "status": "pass" },
      { "test": "test03", "status": "pass" },
      { "test": "test04", "status": "fail" },
      { "test": "test05", "status": "pass" },
      { "test": "test06", "status": "pass" }
    ]
  }
}
```

- `status` 只有 `pass` 与 `fail` 两种：一道题要**全部用例通过**才算 `pass`；
- `config` 是 `config.json` 去掉 `center` 后的内容——`center` 是站点自己的地址，不必回传；
- 站点的「完成题目数」就是 `tests` 里 `status == "pass"` 的条数，题目全集由站点从模板仓库获取。

## 上报的可靠性

- **成绩是「尽力而为」的**：上传最多退避重试 3 次，失败**不影响流水线结论**，报告另随 artifact 留存，可事后补取；
- **同一个仓库 + 同一次 commit 重复上传不会产生重复记录**（中心站点按它去重）；
- 站点**只收集和展示**，不执行你的代码、也不做复算。因此报告上的数字是你自己仓库 CI 产出的**自报值**，看板按这个口径如实展示。

## CI 什么时候会变红

只有两种原因：

1. `config.json` 的课次缺失或填错 —— 校验步骤直接失败，后面的步骤全部跳过；
2. 测试没过（含超时）。

上传失败、产物上传失败都**不会**让 CI 变红。

## 本地运行

```bash
go test ./...                                   # 全部课次
go test ./lesson-01-basics/...                  # 一整课
go test -v ./lesson-01-basics/test06/roster     # 一道题，带详细输出
go test -run TestSumTo ./lesson-01-basics/test02/syntax   # 一个测试函数
go test -race ./lesson-04-sync/...              # 竞态自查（CI 不跑，不计入成绩）
go vet ./...

# 在本地复现 CI 的评分过程（先填好 config.json）
python3 .github/scripts/check_config.py
python3 .github/scripts/run_go_tests.py --root lesson-01-basics
```

Go 版本写在 `go.mod` 里（当前 1.25），CI 的 `setup-go` 直接跟随它，不会出现"本地能编译、CI 编译失败"。
