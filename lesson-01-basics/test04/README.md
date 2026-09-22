# test04：把代码拆成多个包，并正确地导入

对应课件：第一课 · 第三部分（package 组织代码、module 标记边界、导入路径怎么来的）。

前三题都写在同一个包里。真实项目不会这样——第三部分那个 `hello/calc` 的例子才是常态。这一题就是把那个例子在本仓库里真正建出来、跑起来。

## 目录

```text
test04/hello/
  calc/add.go          package calc —— 被使用的一方
  app/app.go           package app  —— 使用别人的一方
```

课件的例子里，使用方是 `main.go`。这里换成普通的 `app` 包，因为第四部分说过：**`main` 包不能被其他包导入**，包括外部测试文件。

## 任务

### 1. `calc/add.go`

| 函数 | 要求 |
|---|---|
| `Add(a, b int) int` | 返回两数之和。大写，别的包能调用 |
| `helper(x int) int` | 返回 `x` 的两倍。小写，只有 `calc` 包内部能用 |
| `AddDoubled(a, b int) int` | 两个数各翻一倍再相加，**必须调用 `helper`**，不要自己写 `*2` |

### 2. `app/app.go`

先补上 `import` 里缺的那一行，再完成两个函数：

| 函数 | 要求 |
|---|---|
| `Sum(a, b int) int` | 调用 `calc.Add` |
| `Describe(a, b int) string` | 返回 `"1 + 2 = 3"` 这样一行，末尾不带换行 |

## 导入路径是怎么算出来的

这是这一题唯一的新知识，也是最容易写错的地方。**导入路径 = module 路径 + 包所在的子目录**。

module 路径写在仓库根目录 `go.mod` 的第一行：

```go
module redrock/teaching-exercises
```

`calc` 包在 `lesson-01-basics/test04/hello/calc/` 目录下，于是：

```go
import "redrock/teaching-exercises/lesson-01-basics/test04/hello/calc"
```

导进来以后，代码里用的名字是**文件开头 `package` 声明的那个名字**，也就是 `calc`：

```go
calc.Add(1, 2)
```

两件常见的错事：

- 写成相对路径 `import "../calc"`——Go 没有这种写法；
- 以为导入路径就是包名 `import "calc"`——那是标准库的形状，自己的包必须写全路径。

包名恰好和目录名一样只是约定（课件第四部分第 3 节），语言并不强制；初学阶段让它们保持一致，能少掉很多坑。

## 为什么要有 `helper`

`AddDoubled` 和 `helper` 在同一个包里，所以 `AddDoubled` 可以直接调用小写的 `helper`。而 `app` 包无论怎么写都调不到它——试着在 `app.go` 里写一句 `calc.helper(1)`，编译器会告诉你 `cannot refer to unexported name`。

这就是课件里说的：C 用 `static` 把可见性关在**文件**里，Go 用小写把可见性关在**包**里。第五部分（test06）会把这条规则完整展开。

## 自己验证

```bash
go test ./lesson-01-basics/test04/...
go test -v ./lesson-01-basics/test04/hello/app
```

## 评分点

| 测试 | 所在包 | 检查内容 |
|---|---|---|
| `TestHelperIsVisibleInsideThePackage` | `calc`（同包测试） | 小写的 `helper` 在包内可以直接调用 |
| `TestAddDoubledGoesThroughHelper` | `calc`（同包测试） | 结果正确，**并且**源码里确实调用了 `helper` |
| `TestSum` | `app_test`（外部测试） | `Sum` 走 `calc.Add` 得到正确结果 |
| `TestDescribe` | `app_test`（外部测试） | 拼出的字符串完全一致 |
| `TestAppImportsCalcByItsFullPath` | `app_test`（外部测试） | `app.go` 用完整路径导入了 `calc` |

这一题有两个包各带测试，所以 CI 上 `test04` 会显示两个测试包的结果。
