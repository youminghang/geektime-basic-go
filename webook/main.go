package main

import (
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"net/http"
)

func main() {
	// 当前目录的配置
	viper.AddConfigPath("./config")
	server := InitWebServer()
	// 注册路由
	server.GET("/hello", func(ctx *gin.Context) {
		ctx.String(http.StatusOK, "hello, world")
	})
	server.Run(":8080")
}
