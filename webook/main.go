package main

import (
	"gitee.com/geekbang/basic-go/webook/internal/web"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"strings"
	"time"
)

func main() {
	server := gin.Default()

	server.Use(func(ctx *gin.Context) {
		println("这是第一个 middleware")
	})

	server.Use(func(ctx *gin.Context) {
		println("这是第二个 middleware")
	})

	server.Use(cors.New(cors.Config{
		//AllowOrigins: []string{"*"},
		//AllowMethods: []string{"POST", "GET"},
		AllowHeaders: []string{"Content-Type", "Authorization"},
		//ExposeHeaders: []string{"x-jwt-token"},
		// 是否允许你带 cookie 之类的东西
		AllowCredentials: true,
		AllowOriginFunc: func(origin string) bool {
			if strings.HasPrefix(origin, "http://localhost") {
				// 你的开发环境
				return true
			}
			return strings.Contains(origin, "yourcompany.com")
		},
		MaxAge: 12 * time.Hour,
	}))

	//v1 := server.Group("/v1")
	//users := server.Group("/users/v1")
	u := web.NewUserHandler()
	//u.RegisterRoutesV1(server.Group("/users"))
	u.RegisterRoutes(server)
	server.Run(":8080")
}
