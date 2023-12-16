package ioc

import (
	"gitee.com/geekbang/basic-go/webook/bff"
	ijwt "gitee.com/geekbang/basic-go/webook/bff/jwt"
	"gitee.com/geekbang/basic-go/webook/bff/middleware"
	"gitee.com/geekbang/basic-go/webook/pkg/ginx"
	"gitee.com/geekbang/basic-go/webook/pkg/ginx/middleware/metrics"
	"gitee.com/geekbang/basic-go/webook/pkg/logger"
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/memstore"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/redis/go-redis/v9"
	otelgin "go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"strings"
	"time"
)

func InitWebServer(funcs []gin.HandlerFunc,
	userHdl *bff.UserHandler,
	artHdl *bff.ArticleHandler,
	obHdl *bff.ObservabilityHandler,
	oauth2Hdl *bff.OAuth2WechatHandler, l logger.LoggerV1) *gin.Engine {
	ginx.SetLogger(l)
	server := gin.Default()
	server.Use(funcs...)
	// 注册路由
	userHdl.RegisterRoutes(server)
	artHdl.RegisterRoutes(server)
	oauth2Hdl.RegisterRoutes(server)
	obHdl.RegisterRoutes(server)
	return server
}

func GinMiddlewares(cmd redis.Cmdable,
	hdl ijwt.Handler, l logger.LoggerV1) []gin.HandlerFunc {
	pb := &metrics.PrometheusBuilder{
		Namespace:  "geekbang_daming",
		Subsystem:  "webook",
		Name:       "gin_http",
		InstanceID: "my-instance-1",
		Help:       "GIN 中 HTTP 请求",
	}
	ginx.InitCounter(prometheus.CounterOpts{
		Namespace: "geekbang_daming",
		Subsystem: "webook",
		Name:      "http_biz_code",
		Help:      "GIN 中 HTTP 请求",
		ConstLabels: map[string]string{
			"instance_id": "my-instance-1",
		},
	})
	return []gin.HandlerFunc{
		//ratelimit.NewBuilder(cmd, time.Minute, 100).BuildResponseTime(),
		corsHandler(),
		pb.BuildResponseTime(),
		pb.BuildActiveRequest(),
		otelgin.Middleware("webook"),
		// 使用 JWT
		middleware.NewJWTLoginMiddlewareBuilder(hdl).Build(),
		//accesslog.NewMiddlewareBuilder(func(ctx context.Context, al accesslog.AccessLog) {
		//	// 设置为 DEBUG 级别
		//	l.Debug("GIN 收到请求", logger.Field{
		//		Key:   "req",
		//		Value: al,
		//	})
		//}).AllowReqBody().AllowRespBody().Build(),
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
