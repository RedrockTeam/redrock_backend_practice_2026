package main

import "github.com/gin-gonic/gin"

// NewRouter 注册全部路由并返回引擎。
func NewRouter() *gin.Engine {
	r := gin.New()

	// TODO: GET /cat —— 把同目录下的 cat.jpg 作为响应内容返回。

	return r
}

// main 把路由跑在 8080 端口：go run . 之后访问 http://localhost:8080/cat。
func main() {
	NewRouter().Run(":8080")
}
