# test03:中间件:记录日志与拦截

对应课件:第六课。

## 任务

编辑 `logfilter/main.go`,注册 `GET /log/:num`,并写一个**中间件**挂在它前面。请求到达后的完整流程:

```text
GET /log/2
  → 中间件:打印日志(包含 num)
  → 中间件:num 是偶数,c.Next() 放行
  → 处理函数:c.JSON 返回 {"message":"偶数"}
```

| 请求 | 中间件的动作 | 响应 |
|---|---|---|
| `/log/2` 这类偶数 | 打印日志,然后放行 | 200,响应体 `{"message":"偶数"}`,由处理函数用 `c.JSON` 写出 |
| `/log/3` 这类奇数 | 打印日志,然后拦截 | 403,响应体 `{"error":"奇数被拦截"}`,处理函数不会执行 |
| `/log/abc` 这类非整数 | 拦截 | 400,响应体 `{"error":"num 不是整数"}` |

只改 `main.go` 一个文件。

## 需要的几样工具

- **中间件**:也是一个 `gin.HandlerFunc`。和处理函数唯一的区别是它可以用 `c.Next()` 让请求继续往后走,或者用 `c.Abort...` 让请求停下。日志用 `fmt.Println` 打印(评分会检查它):

```go
func checkNum(c *gin.Context) {
	fmt.Println("收到请求:", c.Param("num"))
	num, err := strconv.Atoi(c.Param("num"))
	if err != nil {
		c.AbortWithStatusJSON(400, gin.H{"error": "num 不是整数"})
		return
	}
	// TODO: num 是奇数就在这里拦截,偶数放行
	c.Next()
}
```

挂载方式二选一:

```go
r.GET("/log/:num", checkNum, func(c *gin.Context) { /* 处理函数 */ }) // 只对这一条路由生效
r.Use(checkNum)                                                       // 对之后注册的所有路由生效
```

- **c.JSON**:和 `c.String` 一样是写响应体,但它把数据编成 JSON。`gin.H` 是 gin 提供的 map 简写:

```go
c.JSON(200, gin.H{"message": "偶数"}) // 响应体 {"message":"偶数"}
```

- **c.AbortWithStatusJSON**:写一个 JSON 响应并且**中止**——它后面的处理函数和中间件都不会再执行。这就是"拦截"。

## 自己验证

```bash
cd lesson-06-web/test03/logfilter
go run .                       # 监听 http://localhost:8080
```

另开一个终端:

```bash
curl http://localhost:8080/log/2     # {"message":"偶数"}
curl http://localhost:8080/log/3     # {"error":"奇数被拦截"}
curl http://localhost:8080/log/abc   # {"error":"num 不是整数"}
```

奇数被拦截时,终端里仍会打出日志,但处理函数没有机会执行。

## 评分点

| 测试 | 检查内容 |
|---|---|
| `TestLogEven` | 偶数返回 200 和 `{"message":"偶数"}`,且请求期间打出了包含 num 的日志 |
| `TestLogOdd` | 奇数返回 403 和 `{"error":"奇数被拦截"}`,且日志同样打出 num |
| `TestLogNotANumber` | 非整数返回 400 和 `{"error":"num 不是整数"}` |
