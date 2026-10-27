# 红岩网校 Go 练习仓库

**这是模板仓库，别在这里做题。** 点右上角「Use this template」→「Create a new repository」，复制一份到你自己的账号下，在那一份里做。

每一课一个目录（`lesson-01-basics/` … `lesson-04-sync/`），每道题一个 `testNN/`，题目要求写在各自的 `README.md` 里。你要动的只有 `config.json` 和各题点名的 `.go` 文件。

## 前置准备

### 用模板建出自己的仓库

访问 https://github.com/RedrockTeam/redrock_backend_practice_2026

![image-20261027150254177](assets/image-20261027150254177.png)

![image-20261027150636065](assets/image-20261027150636065.png)

### 把远程仓库clone到本地

这里有两种做法：
1.直接下载源码的zip压缩包，解压后移动到自己专门放代码的地方，但是需要自己绑定远程仓库
![image-20261027150858097](assets/image-20261027150858097.png)

2.使用git命令
![image-20261027151005864](assets/image-20261027151005864.png)

在你想要放代码的目录执行如下命令，会自动下载并绑定远程仓库（把地址里的「你的GitHub用户名」换成你自己的）：

```bash
git clone https://github.com/你的GitHub用户名/redrock_backend_practice_2026.git
```

### 完善个人信息

打开根目录的config.json：

```json
{
  "name": "张三",
  "lesson": "lesson-01-basics",
  "center": "https://homework.blog-manmu.top"
}
```

- `name` 改成你的真实名字
- `lesson` 改成这次要做的课的目录名，必须是根目录名，**括号和中文说明一个字都不能留**，否则直接报错
- `center` 不要动

## 完成题目

### 只改这几类文件

题目要求写在每道题的 `README.md` 里，它点名让你改哪个文件就改哪个。

### 一次只做一课

`config.json` 里的 `lesson` 填哪一课，CI 就只跑那一课。做完一课交一课，下一课再把 `lesson` 改过去。

### 做完一课怎么交

```bash
go test ./lesson-01-basics/...      # 先在本机确认这一课全绿
git add .
git commit -m "完成第一课"
git push
```

`git commit -m` 引号里随便写，是给你自己看的说明。

也可以使用GoLand/VSCode提供的GUI按钮

push 到 **main** 后 GitHub 自动跑测试并上报成绩。**其它分支和 PR 都不会触发**。第一次 `git push` 会弹窗让你登录 GitHub，授权一次即可。

### 怎么算"这题做完了"

一道题下的**所有测试都通过**，这题才算通过。看板上的完成题数就是通过了几道题。所以本机跑 `go test` 时要看到 `ok`，而不是"大部分对"。

### 提交前自查

- `go test ./lesson-0X-.../...` 全绿；
- **例外：`lesson-03` 的 test03 在第三课结束时就是红的**，它是观察题，题面写明"不要改实现，把问题带到第四课"。这一题不算你没做完；
- `config.json` 的 `lesson` 和你这次做的课一致；
- 没改过任何 `*_test.go`（除了 test07 那道题）。

## 作业汇总

**<https://homework.blog-manmu.top>** ——「排行榜」按完成题数排名（通过几题就是几分）。

想看自己每次 CI 的结果：自己仓库顶部的 **Actions** 标签页，点进最近一次运行。

## 本地测试

```bash
go test ./lesson-01-basics/...                # 跑一整课
go test -v ./lesson-01-basics/test08/roster   # 一道题，看详细输出
go test -race ./lesson-04-sync/...            # 竞态自查（CI 不计入成绩）
```

**starter 初始一定是红的**，那是题目还没做完，不是环境坏了。唯一的例外是 `lesson-03` 的 test03：它是观察题，改完前面的题它仍然是红的，这是设计好的。

## 遇到问题

- **CI 红了**：只有两种原因——`lesson` 填错，或测试没过。上传失败不会让 CI 变红。
- **成绩没出来**：上传是「尽力而为」的，失败不影响 CI；报告另存 artifact，可事后补。
- **重复 push**：同一次 commit 重复上传不会重复计分。
- **`lesson-03` 交了还是有题红**：正常，test03 是观察题，第三课结束时本就该红。

