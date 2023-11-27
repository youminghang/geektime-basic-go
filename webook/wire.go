//go:build wireinject

package main

import (
	"gitee.com/geekbang/basic-go/webook/interactive/events"
	repository2 "gitee.com/geekbang/basic-go/webook/interactive/repository"
	cache2 "gitee.com/geekbang/basic-go/webook/interactive/repository/cache"
	dao2 "gitee.com/geekbang/basic-go/webook/interactive/repository/dao"
	service2 "gitee.com/geekbang/basic-go/webook/interactive/service"
	article2 "gitee.com/geekbang/basic-go/webook/internal/events/article"
	"gitee.com/geekbang/basic-go/webook/internal/repository"
	"gitee.com/geekbang/basic-go/webook/internal/repository/cache"
	"gitee.com/geekbang/basic-go/webook/internal/repository/dao"
	"gitee.com/geekbang/basic-go/webook/internal/repository/dao/article"
	"gitee.com/geekbang/basic-go/webook/internal/service"
	"gitee.com/geekbang/basic-go/webook/internal/web"
	ijwt "gitee.com/geekbang/basic-go/webook/internal/web/jwt"
	"gitee.com/geekbang/basic-go/webook/ioc"
	"github.com/google/wire"
)

var rankServiceProvider = wire.NewSet(
	service.NewBatchRankingService,
	repository.NewCachedRankingRepository,
	cache.NewRedisRankingCache,
	cache.NewRankingLocalCache,
)

func InitApp() *App {
	wire.Build(
		ioc.InitRedis, ioc.InitDB,
		ioc.InitLogger,
		ioc.InitKafka,
		ioc.InitRLockClient,
		ioc.NewSyncProducer,

		rankServiceProvider,
		ioc.InitJobs,
		ioc.InitRankingJob,

		// DAO 部分
		dao.NewGORMUserDAO,
		dao2.NewGORMInteractiveDAO,
		article.NewGORMArticleDAO,

		// Cache 部分
		cache.NewRedisUserCache,
		cache.NewRedisCodeCache,
		cache.NewRedisArticleCache,
		cache2.NewRedisInteractiveCache,

		// repository 部分
		repository.NewCachedUserRepository,
		repository.NewCachedCodeRepository,
		repository.NewArticleRepository,
		repository2.NewCachedInteractiveRepository,

		// events 部分
		article2.NewSaramaSyncProducer,
		events.NewInteractiveReadEventConsumer,
		ioc.NewConsumers,

		// service 部分
		ioc.InitSmsService,
		ioc.InitWechatService,
		service.NewSMSCodeService,
		service.NewUserService,
		service.NewArticleService,
		service2.NewInteractiveService,

		// handler 部分
		ijwt.NewRedisHandler,
		web.NewUserHandler,
		web.NewArticleHandler,
		web.NewOAuth2WechatHandler,
		web.NewObservabilityHandler,

		// gin 的中间件
		ioc.GinMiddlewares,

		// Web 服务器
		ioc.InitWebServer,

		wire.Struct(new(App), "*"),
	)
	// 随便返回一个
	return new(App)
}
