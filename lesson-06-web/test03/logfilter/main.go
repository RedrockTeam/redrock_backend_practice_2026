package main

import "github.com/gin-gonic/gin"

// NewRouter 注册全部路由并返回引擎。
func NewRouter() *gin.Engine {
	r := gin.New()

	// TODO: 注册 GET /log/:num，并在它前面挂一个中间件。
	// 中间件要做三件事：
	//   1. 用 fmt.Println 打印包含 num 的日志（评分会检查日志）；
	//   2. 用 strconv.Atoi 把 c.Param("num") 转成整数，
	//      转失败就 c.AbortWithStatusJSON(400, gin.H{"error": "num 不是整数"})；
	//   3. num 是奇数就 c.AbortWithStatusJSON(403, gin.H{"error": "奇数被拦截"})，
	//      偶数就 c.Next() 放行。
	// 处理函数收到放行后用 c.JSON 返回 {"message":"偶数"}。
	//
	// 挂载方式二选一：先 r.Use(中间件) 再注册路由，或者
	// r.GET("/log/:num", 中间件, 处理函数) 只对这一条路由生效。

	return r
}

// main 把路由跑在 8080 端口：go run . 之后用浏览器或 curl 访问 http://localhost:8080。
func main() {
	NewRouter().Run(":8080")
}
