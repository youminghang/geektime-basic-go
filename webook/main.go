package main

import (
	"gitee.com/geekbang/basic-go/webook/config"
	"gitee.com/geekbang/basic-go/webook/internal/repository"
	"gitee.com/geekbang/basic-go/webook/internal/repository/cache"
	"gitee.com/geekbang/basic-go/webook/internal/repository/dao"
	"gitee.com/geekbang/basic-go/webook/internal/service"
	"gitee.com/geekbang/basic-go/webook/internal/web"
	"gitee.com/geekbang/basic-go/webook/internal/web/middleware"
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/memstore"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"net/http"
	"strings"
	"time"
)

func main() {
	db := initDB()
	redisCmd := initRedis()
	server := initWebServer()
	initUser(server, db, redisCmd)
	//server := gin.Default()
	server.GET("/hello", func(ctx *gin.Context) {
		ctx.String(http.StatusOK, "hello, world")
	})
	server.Run(":8080")
}

func initDB() *gorm.DB {
	db, err := gorm.Open(mysql.Open(config.Config.DB.DSN))
	if err != nil {
		panic(err)
	}
	err = dao.InitTables(db)
	if err != nil {
		panic(err)
	}
	return db
}

func initWebServer() *gin.Engine {
	server := gin.Default()
	//rCfg := config.Config.Redis
	//cmd := redis.NewClient(&redis.Options{
	//	Addr:     rCfg.Addr,
	//	Password: rCfg.Password,
	//	DB:       rCfg.DB,
	//})
	// 一分钟 100 次。
	//server.Use(ratelimit.NewBuilder(cmd, time.Minute, 100).Build())

	server.Use(cors.New(cors.Config{
		AllowCredentials: true,
		// 在使用 JWT 的时候，因为我们使用了 Authorization 的头部，所以要加上
		AllowHeaders: []string{"Content-Type", "Authorization"},
		// 为了 JWT
		ExposeHeaders: []string{"X-Jwt-Token"},
		AllowOriginFunc: func(origin string) bool {
			if strings.HasPrefix(origin, "http://localhost") {
				return true
			}
			return strings.Contains(origin, "your_company.com")
		},
		MaxAge: 12 * time.Hour,
	}))

	// 使用 session 机制登录
	//usingSession(server)
	// 使用 JWT
	usingJWT(server)
	return server
}

func usingJWT(server *gin.Engine) {
	mldBd := &middleware.JWTLoginMiddlewareBuilder{}
	server.Use(mldBd.Build())
}

func usingSession(server *gin.Engine) {
	//store := cookie.NewStore([]byte("secret"))

	// 这是基于内存的实现，第一个参数是 authentication key，最好是 32 或者 64 位
	// 第二个参数是 encryption key
	store := memstore.NewStore([]byte("moyn8y9abnd7q4zkq2m73yw8tu9j5ixm"),
		[]byte("o6jdlo2cb9f9pb6h46fjmllw481ldebj"))
	// 第一个参数是最大空闲连接数量
	// 第二个就是 tcp，你不太可能用 udp
	// 第三、四个 就是连接信息和密码
	// 第五第六就是两个 key
	//store, err := redis.NewStore(16, "tcp",
	//	"localhost:6379", "",
	//	// authentication key, encryption key
	//	[]byte("moyn8y9abnd7q4zkq2m73yw8tu9j5ixm"),
	//	[]byte("o6jdlo2cb9f9pb6h46fjmllw481ldebj"))
	//if err != nil {
	//	panic(err)
	//}

	// cookie 的名字叫做ssid
	server.Use(sessions.Sessions("ssid", store))
	// 登录校验
	login := &middleware.LoginMiddlewareBuilder{}
	server.Use(login.CheckLogin())
}

func initUser(server *gin.Engine, db *gorm.DB, cmd redis.Cmdable) {
	ud := dao.NewUserDAO(db)
	uc := cache.NewUserCache(cmd)
	ur := repository.NewUserRepository(ud, uc)
	us := service.NewUserService(ur)
	c := web.NewUserHandler(us)
	c.RegisterRoutes(server)
}

func initRedis() redis.Cmdable {
	rCfg := config.Config.Redis
	cmd := redis.NewClient(&redis.Options{
		Addr:     rCfg.Addr,
		Password: rCfg.Password,
		DB:       rCfg.DB,
	})
	return cmd
}
