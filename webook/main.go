package main

import (
	"gitee.com/geekbang/basic-go/webook/internal/web"
	"github.com/gin-gonic/gin"
)

func main() {
	server := gin.Default()
	//v1 := server.Group("/v1")
	//users := server.Group("/users/v1")
	u := web.NewUserHandler()
	//u.RegisterRoutesV1(server.Group("/users"))
	u.RegisterRoutes(server)
	server.Run(":8080")
}
