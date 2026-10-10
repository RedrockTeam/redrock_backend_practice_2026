# 第六课练习:Web 服务(Gin)

三题,按顺序做。这一课把之前写过的纯函数搬进 HTTP 世界:用 Gin 注册路由、取路径与查询参数、写响应,把磁盘上的静态资源挂到 URL 上,最后用中间件记录日志、拦截请求,并用 JSON 组织响应。

| 题目 | 主题 |
|---|---|
| [test01](test01/) | 注册三条路由:`/ping`、`/user/:id`、`/ping/add` |
| [test02](test02/) | 把 `cat.jpg` 挂到 `/cat`,访问即显示图片 |
| [test03](test03/) | 中间件记录日志、拦截奇数,处理函数用 `c.JSON` 返回 `{"message":"偶数"}` |

```bash
go test ./lesson-06-web/...
```

要真正跑起服务器(用浏览器或 curl 访问),进入某一题的包目录后 `go run .`,服务器监听 `http://localhost:8080`:

```bash
cd lesson-06-web/test01/ping
go run .
```

gin 已经写进仓库根目录的 `go.mod`,`go run` / `go test` 会自动下载依赖,不需要手动安装。
