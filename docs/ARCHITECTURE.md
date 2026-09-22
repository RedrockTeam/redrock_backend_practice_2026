# 练习仓库总体架构

## 一次提交会发生什么

```text
学生用「Use this template」建出自己的仓库 → 做题 → push 到 main
        │（只有默认分支的 push 会触发；其它分支与 PR 都不跑）
        ▼
grade.yml（唯一的 workflow）
  ① check_config.py
       ├─ 读仓库根目录的 config.json：lesson（必填）、name、center
       ├─ 课次缺失或非法 → 打 ::error:: 并以退出码 1 中止整个 job
       └─ 把 LESSON_PATH / TEST_ROOT / CENTER_API 写进 GITHUB_ENV
  ② run_go_tests.py --root <课次>
       ├─ 递归发现该课次下所有含 *_test.go 的目录（跳过 testdata/、vendor/）
       ├─ 对每个目录独立执行 go test -json -count=1 -timeout=120s
       │    ├─ 原始 JSONL / 日志 → test-results/（只进 artifact，不进载荷）
       │    └─ 逐测试函数的终态 → result.json、score.md、Step Summary
  ③ notify_center.py（仅 push 事件，尽力而为）
       ├─ 把 result.json 压成上报载荷：课次 + 逐题编号与结论（pass / fail）
       └─ POST /api/v1/reports（最多 3 次退避重试，失败不影响流水线）
  ④ 上传 artifact go-test-<sha>
```

## 为什么按课次裁剪

`config.json` 的 `lesson` 决定这次跑哪一课，其余课次不跑：

- 反馈聚焦：学生做第三课时，看到的只有第三课的逐题结果，不会被其它课次的红叉淹没；
- 课次是计分的归属。一份「跑了全部课次」的报告无法归入任何一课，站点也没法用它更新某一课的成绩；
- 因此 `lesson` 缺失或填错时**直接报错中止**，而不是退回「全跑一遍」。

代价要说清楚：**没有跨课次的回归保护**——学生把前面课次改坏了，本次运行不会发现。站点侧仍能看到「该课次最近一次」的成绩，坏掉的那一次只有在再次提交该课次时才会体现。

## 目录与评分的对应关系

```text
lesson-01-basics/test06/roster/roster_test.go
└── lesson ─┘└─ test ┘└─ package ┘
```

`run_go_tests.py` 的 `attribute()` 把包路径映射成 `(lesson, test)`，`result.json` 按 `lessons[] → tests[] → packages[] → cases[]` 四层嵌套。成绩网站因此可以直接按「第几课第几题」寻址，不需要解析路径字符串。

题目单位叫 `testNN/`，真正的 Go 包放在它下一层且用有意义的名字（`hellogo`、`syntax`、`roster`）。这是刻意的：第一课第四部分刚教完「包名通常与目录名一致，避免 `util`、`common` 这类名字」，如果题目目录本身就是 `package test06`，课件当场自相矛盾。

## 成绩的性质

**报告是学生仓库自己的 CI 产出的自报值。** 学生对自己的仓库有管理权限，理论上可以改测试、改脚本、改分数；本项目**不做防作弊设计**，也不引入额外的校验机制。

划分是清楚的：

| 环节 | 谁负责 |
|---|---|
| 跑测试、算分 | 学生仓库的 CI（本仓库的脚本） |
| 接收、存储、聚合、展示 | 中心站点，原样接收，**不执行学生代码、也不复算** |
| 认定成绩 | 教师（看板只如实呈现上报值） |

看板按「成绩未经复核」的口径展示即可。

## 边界

- 公开练习仓库：题目 README、starter、公开测试、一个 workflow、三个脚本、`config.json` 模板；
- 学生唯一的输入：`config.json` 的 `name` 与 `lesson`（`center` 由教师预置）；
- 中心站点：接收、存储、聚合、展示，是独立项目，接口契约由它自己维护；
- 教师课件：只在 `redrock/courseware` 维护，不参与 CI。
