package ioc

import (
	"context"
	"gitee.com/geekbang/basic-go/webook/internal/web"
	ijwt "gitee.com/geekbang/basic-go/webook/internal/web/jwt"
	"gitee.com/geekbang/basic-go/webook/internal/web/middleware"
	"gitee.com/geekbang/basic-go/webook/pkg/ginx"
	"gitee.com/geekbang/basic-go/webook/pkg/ginx/middleware/accesslog"
	"gitee.com/geekbang/basic-go/webook/pkg/ginx/middleware/ratelimit"
	"gitee.com/geekbang/basic-go/webook/pkg/logger"
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/memstore"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"strings"
	"time"
)

func InitWebServer(funcs []gin.HandlerFunc,
	userHdl *web.UserHandler,
	artHdl *web.ArticleHandler,
	oauth2Hdl *web.OAuth2WechatHandler, l logger.LoggerV1) *gin.Engine {
	ginx.SetLogger(l)
	server := gin.Default()
	server.Use(funcs...)
	// 注册路由
	userHdl.RegisterRoutes(server)
	artHdl.RegisterRoutes(server)
	oauth2Hdl.RegisterRoutes(server)
	return server
}

func GinMiddlewares(cmd redis.Cmdable,
	hdl ijwt.Handler, l logger.LoggerV1) []gin.HandlerFunc {
	return []gin.HandlerFunc{
		ratelimit.NewBuilder(cmd, time.Minute, 100).Build(),
		corsHandler(),

		// 使用 JWT
		middleware.NewJWTLoginMiddlewareBuilder(hdl).Build(),
		accesslog.NewMiddlewareBuilder(func(ctx context.Context, al accesslog.AccessLog) {
			// 设置为 DEBUG 级别
			l.Debug("GIN 收到请求", logger.Field{
				Key:   "req",
				Value: al,
			})
		}).AllowReqBody().AllowRespBody().Build(),
		// 使用session 登录校验
		//sessionHandlerFunc(),
		//middleware.NewLoginMiddlewareBuilder().CheckLogin(),
	}
}

func corsHandler() gin.HandlerFunc {
	return cors.New(cors.Config{
		AllowCredentials: true,
		// 在使用 JWT 的时候，因为我们使用了 Authorization 的头部，所以要加上
		AllowHeaders: []string{"Content-Type", "Authorization"},
		// 为了 JWT，长短 token 的设置
		ExposeHeaders: []string{"X-Jwt-Token", "X-Refresh-Token"},
		AllowOriginFunc: func(origin string) bool {
			if strings.HasPrefix(origin, "http://localhost") {
				return true
			}
			return strings.Contains(origin, "your_company.com")
		},
		MaxAge: 12 * time.Hour,
	})
}

func sessionHandlerFunc() gin.HandlerFunc {
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
	return sessions.Sessions("ssid", store)
}
