//go:build wireinject

package startup

import (
	"gitee.com/geekbang/basic-go/webook/bff"
	"gitee.com/geekbang/basic-go/webook/bff/jwt"
	repository2 "gitee.com/geekbang/basic-go/webook/interactive/repository"
	cache2 "gitee.com/geekbang/basic-go/webook/interactive/repository/cache"
	dao2 "gitee.com/geekbang/basic-go/webook/interactive/repository/dao"
	service2 "gitee.com/geekbang/basic-go/webook/interactive/service"
	article2 "gitee.com/geekbang/basic-go/webook/internal/events/article"
	"gitee.com/geekbang/basic-go/webook/internal/job"
	"gitee.com/geekbang/basic-go/webook/internal/repository"
	"gitee.com/geekbang/basic-go/webook/internal/repository/cache"
	"gitee.com/geekbang/basic-go/webook/internal/repository/dao"
	"gitee.com/geekbang/basic-go/webook/internal/repository/dao/article"
	"gitee.com/geekbang/basic-go/webook/internal/service"
	"gitee.com/geekbang/basic-go/webook/internal/service/sms"
	"gitee.com/geekbang/basic-go/webook/internal/service/sms/async"
	"gitee.com/geekbang/basic-go/webook/ioc"
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
)

var thirdProvider = wire.NewSet(InitRedis, InitTestDB,
	InitLog,
	NewSyncProducer,
	InitKafka,
)
var userSvcProvider = wire.NewSet(
	dao.NewGORMUserDAO,
	cache.NewRedisUserCache,
	repository.NewCachedUserRepository,
	service.NewUserService)
var articlSvcProvider = wire.NewSet(
	article.NewGORMArticleDAO,
	article2.NewSaramaSyncProducer,
	cache.NewRedisArticleCache,
	repository.NewArticleRepository,
	service.NewArticleService)

var interactiveSvcProvider = wire.NewSet(
	service2.NewInteractiveService,
	repository2.NewCachedInteractiveRepository,
	dao2.NewGORMInteractiveDAO,
	cache2.NewRedisInteractiveCache,
)

var rankServiceProvider = wire.NewSet(
	service.NewBatchRankingService,
	repository.NewCachedRankingRepository,
	cache.NewRedisRankingCache,
	cache.NewRankingLocalCache,
)

var jobProviderSet = wire.NewSet(
	service.NewCronJobService,
	repository.NewPreemptCronJobRepository,
	dao.NewGORMJobDAO)

//go:generate wire
func InitWebServer() *gin.Engine {
	wire.Build(
		thirdProvider,
		userSvcProvider,
		articlSvcProvider,
		interactiveSvcProvider,
		InitInteractiveClient,
		cache.NewRedisCodeCache,
		repository.NewCachedCodeRepository,
		// service 部分
		// 集成测试我们显式指定使用内存实现
		ioc.InitSmsMemoryService,

		// 指定啥也不干的 wechat service
		InitPhantomWechatService,
		service.NewSMSCodeService,
		// handler 部分
		bff.NewUserHandler,
		bff.NewOAuth2WechatHandler,
		bff.NewArticleHandler,
		bff.NewObservabilityHandler,
		jwt.NewRedisHandler,

		// gin 的中间件
		ioc.GinMiddlewares,

		// Web 服务器
		ioc.InitWebServer,
	)
	// 随便返回一个
	return gin.Default()
}

func InitArticleHandler(dao article.ArticleDAO) *bff.ArticleHandler {
	wire.Build(thirdProvider,
		userSvcProvider,
		interactiveSvcProvider,
		article2.NewSaramaSyncProducer,
		cache.NewRedisArticleCache,
		InitInteractiveClient,
		//wire.InterfaceValue(new(article.ArticleDAO), dao),
		repository.NewArticleRepository,
		service.NewArticleService,
		bff.NewArticleHandler)
	return new(bff.ArticleHandler)
}

func InitUserSvc() service.UserService {
	wire.Build(thirdProvider, userSvcProvider)
	return service.NewUserService(nil)
}

func InitAsyncSmsService(svc sms.Service) *async.Service {
	wire.Build(thirdProvider, repository.NewAsyncSMSRepository,
		dao.NewGORMAsyncSmsDAO,
		async.NewService,
	)
	return &async.Service{}
}

func InitRankingService() service.RankingService {
	wire.Build(thirdProvider,
		interactiveSvcProvider,
		articlSvcProvider,
		InitInteractiveClient,
		// 用不上这个 user repo，所以随便搞一个
		wire.InterfaceValue(new(repository.UserRepository),
			&repository.CachedUserRepository{}),
		rankServiceProvider)
	return &service.BatchRankingService{}
}

func InitInteractiveService() service2.InteractiveService {
	wire.Build(thirdProvider, interactiveSvcProvider)
	return service2.NewInteractiveService(nil, nil)
}

func InitJobScheduler() *job.Scheduler {
	wire.Build(jobProviderSet, thirdProvider, job.NewScheduler)
	return &job.Scheduler{}
}

func InitJwtHdl() jwt.Handler {
	wire.Build(thirdProvider, jwt.NewRedisHandler)
	return jwt.NewRedisHandler(nil)
}
