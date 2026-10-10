# test01:用 Gin 注册三条路由

对应课件:第六课。

## 任务

编辑 `ping/main.go`,在 `NewRouter()` 里补全三条路由,`main()` 已经将其跑在 `:8080` 上:

| 路由 | 响应体 |
|---|---|
| `GET /ping` | `pong` |
| `GET /user/:id` | 路径参数 `id` 本身,如 `/user/42` 返回 `42` |
| `GET /ping/add?num1=&num2=` | 两个整数之和,如 `?num1=3&num2=5` 返回 `8` |


## 需要的几样工具

- `c.String(code, ...)`:写纯文本响应体,参数用法和 `fmt.Sprintf` 一样;
- `c.Param("id")`:取 `:id` 这类**路径参数**,例如:

```go
r.GET("/product/:price", func(c *gin.Context) {
	c.String(200, c.Param("price"))
})
```

- `c.Query("num1")`:取 `?num1=` 这类**查询参数**,拿到的是字符串。两个数相加之前先用 `strconv.Atoi` 转成整数,加完之后再用 `fmt.Sprintf` 拼回字符串写进响应体。

## 自己验证

```bash
cd lesson-06-web/test01/ping
go run .                       # 监听 http://localhost:8080
```

另开一个终端:

```bash
curl http://localhost:8080/ping
curl http://localhost:8080/user/42
curl "http://localhost:8080/ping/add?num1=3&num2=5"
```

浏览器里直接访问这三个地址也一样。

## 评分点

| 测试 | 检查内容 |
|---|---|
| `TestPing` | `/ping` 返回 200,响应体恰好是 `pong` |
| `TestUserByID` | `/user/:id` 原样返回路径里的 id(测试会试多个不同的 id,写死一个答不了) |
| `TestAdd` | `/ping/add` 返回两数之和(含 0 和负数) |
