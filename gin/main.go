package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func main() {
	server := gin.Default()

	server.Use(func(context *gin.Context) {
		println("第一个 middleware")
	}, func(context *gin.Context) {
		println("第二个 middleware")
	})

	server.GET("/hello", func(ctx *gin.Context) {
		ctx.String(http.StatusOK, "hello, world")
	})

	server.GET("/users/:name", func(context *gin.Context) {
		name := context.Param("name")
		context.String(http.StatusOK, "这是你传过来的名字 %s", name)
	})

	server.GET("/order", func(context *gin.Context) {
		// 查询参数
		id := context.Query("id")
		context.String(http.StatusOK, "你传过来的 ID 是 %s", id)
	})

	server.GET("/views/*.html", func(context *gin.Context) {
		path := context.Param(".html")
		context.String(http.StatusOK, "匹配上的值是 %s", path)
	})

	// 这种路由是不合法的
	//server.GET("/invalid/*", func(context *gin.Context) {
	//})

	// 这种路由也是不合法的
	//server.GET("/invalid/*/b", func(context *gin.Context) {
	//	//context.String(http.StatusOK, "")
	//})

	// 8080 是启动端口
	server.Run(":8080")
}
