package main

import "github.com/gin-gonic/gin"

// NewRouter 注册全部路由并返回引擎。
//
//	GET /ping       响应体是 pong
//	GET /user/:id   响应体是路径参数 id 本身
//	GET /ping/add   响应体是 num1+num2 的和（两个查询参数都是整数）
func NewRouter() *gin.Engine {
	r := gin.New()

	// TODO: 请注册三条路由并且配置相应的handler

	// TODO: GET /ping

	// TODO: GET /user/:id

	// TODO: GET /ping/add

	return r
}

// main 把路由跑在 8080 端口：go run . 之后用浏览器或 curl 访问 http://localhost:8080。
func main() {
	NewRouter().Run(":8080")
}
