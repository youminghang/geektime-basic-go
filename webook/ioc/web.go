package ioc

import (
	"gitee.com/geekbang/basic-go/webook/internal/web"
	ijwt "gitee.com/geekbang/basic-go/webook/internal/web/jwt"
	"gitee.com/geekbang/basic-go/webook/internal/web/middleware"
	"gitee.com/geekbang/basic-go/webook/pkg/ginx"
	"gitee.com/geekbang/basic-go/webook/pkg/ginx/middlewares/metric"
	logger2 "gitee.com/geekbang/basic-go/webook/pkg/logger"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/redis/go-redis/v9"
	"strings"
	"time"
)

func InitWebServer(mdls []gin.HandlerFunc, userHdl *web.UserHandler,
	oauth2WechatHdl *web.OAuth2WechatHandler, articleHdl *web.ArticleHandler) *gin.Engine {
	server := gin.Default()
	server.Use(mdls...)
	userHdl.RegisterRoutes(server)
	articleHdl.RegisterRoutes(server)
	oauth2WechatHdl.RegisterRoutes(server)
	(&web.ObservabilityHandler{}).RegisterRoutes(server)
	return server
}

func InitMiddlewares(redisClient redis.Cmdable,
	l logger2.LoggerV1,
	jwtHdl ijwt.Handler) []gin.HandlerFunc {
	//bd := logger.NewBuilder(func(ctx context.Context, al *logger.AccessLog) {
	//	l.Debug("HTTP请求", logger2.Field{Key: "al", Value: al})
	//}).AllowReqBody(true).AllowRespBody()
	//viper.OnConfigChange(func(in fsnotify.Event) {
	//	ok := viper.GetBool("web.logreq")
	//	bd.AllowReqBody(ok)
	//})
	ginx.InitCounter(prometheus.CounterOpts{
		Namespace: "geekbang_daming",
		Subsystem: "webook",
		Name:      "http_biz_code",
		Help:      "HTTP 的业务错误码",
	})
	return []gin.HandlerFunc{
		corsHdl(),
		(&metric.MiddlewareBuilder{
			Namespace:  "geekbang_daming",
			Subsystem:  "webook",
			Name:       "gin_http",
			Help:       "统计 GIN 的 HTTP 接口",
			InstanceID: "my-instance-1",
		}).Build(),
		//bd.Build(),
		middleware.NewLoginJWTMiddlewareBuilder(jwtHdl).
			IgnorePaths("/users/signup").
			IgnorePaths("/users/refresh_token").
			IgnorePaths("/users/login_sms/code/send").
			IgnorePaths("/users/login_sms").
			IgnorePaths("/oauth2/wechat/authurl").
			IgnorePaths("/oauth2/wechat/callback").
			IgnorePaths("/users/login").
			IgnorePaths("/test/metric").
			Build(),
		//ratelimit.NewBuilder(redisClient, time.Second, 100).Build(),
	}
}

func corsHdl() gin.HandlerFunc {
	return cors.New(cors.Config{
		//AllowOrigins: []string{"*"},
		//AllowMethods: []string{"POST", "GET"},
		AllowHeaders: []string{"Content-Type", "Authorization"},
		// 你不加这个，前端是拿不到的
		ExposeHeaders: []string{"x-jwt-token", "x-refresh-token"},
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
	})
}
