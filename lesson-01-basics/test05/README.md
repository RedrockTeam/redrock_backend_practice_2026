# test05：修好一个违反 package 规范的项目

对应课件：第一课 · 第三部分（package 与 module）、第四部分（包规范）。

## 任务

`packagerules/testdata/` 下有一个**故意写错**的小项目，你要把它改到符合 Go 的包规范。

```text
packagerules/testdata/
  package-project/            # module example.com/hello
    go.mod
    main.go                   # package main
    common/name.go            # 包名过于宽泛
    user_service/login.go     # 与同目录另一个文件包名不一致
    user_service/register.go  # 包名大写 + 目录名带下划线
    order/order.go            # 与 user 包互相导入
    worker/worker.go          # 导入了 package main
    internal/secret/secret.go # 只能被 example.com/hello 目录树导入
  external-module/            # module example.net/outside，另一个 module
    outside.go                # 越过边界导入了别人的 internal 包
```

每个要改的地方都留了 `// TODO` 注释。一共有五类问题：

1. **同一目录只能有一个包名**：`user_service/login.go` 和 `user_service/register.go` 现在声明了不同的包。
2. **包名规范**：小写、简短、有意义，不用下划线和驼峰，不用 `common` / `util` 这类宽泛名字。本题额外要求**包名和目录名保持一致**，所以改包名时要一起改目录名，并更新所有 `import` 路径。
3. **`main` 包不能被导入**：`worker/worker.go` 导入了 `example.com/hello`。
4. **`internal` 边界**：`example.com/hello/internal/secret` 只能被 `example.com/hello` 及其子目录导入，`external-module/` 是另一个 module，不可以。
5. **禁止循环导入**：现在 `order` 和用户包互相导入，需要重新设计成单向依赖。

## 为什么题目放在 `testdata/` 里

`go build` 和 `go test` 会**跳过** `testdata/` 目录，所以这里可以安全地保存一个编译不通过、包名写错的项目，而不影响整个练习仓库。评分测试用 `go/parser` 读取这些文件的 `package` 声明和 `import` 列表来判断，因此：

- 你改完 `testdata/` 里的代码就能被重新检查；
- 但 `go run ./testdata/package-project` 是跑不起来的，不要尝试。

## 改动提示

- 目录改名后，所有引用它的 `import` 路径都要改，否则会被判为「导入了不存在的本地路径」。
- 解决循环导入的常规办法：抽出一个更底层的公共包，或者反转其中一个方向的依赖，让依赖关系变成单向。
- `_ "some/path"` 是空导入，只为了触发这个包的初始化。本题里它的作用只是「制造一条依赖关系」，删掉不需要的那条也是合法解法。

## 自己验证

```bash
go test ./lesson-01-basics/test05/...
go test -v -run TestPackageNames ./lesson-01-basics/test05/packagerules
```

失败信息会直接指出是哪个文件、哪条规则。

## 评分点

| 测试 | 检查内容 |
|---|---|
| `TestPackageNamesAndDirectories` | 一个目录一个包；包名小写规范；包名与目录名一致；根目录是 `package main` |
| `TestPackageImportBoundaries` | 导入路径真实存在；没人导入 `main` 包；`internal` 不被跨 module 导入 |
| `TestPackagesDoNotImportEachOtherInACycle` | 依赖图中没有环 |
