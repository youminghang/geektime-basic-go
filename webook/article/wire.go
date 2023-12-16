//go:build wireinject

package main

import (
	"gitee.com/geekbang/basic-go/webook/internal/article/events"
	"gitee.com/geekbang/basic-go/webook/internal/article/grpc"
	"gitee.com/geekbang/basic-go/webook/internal/article/ioc"
	"gitee.com/geekbang/basic-go/webook/internal/article/repository"
	"gitee.com/geekbang/basic-go/webook/internal/article/repository/cache"
	"gitee.com/geekbang/basic-go/webook/internal/article/repository/dao"
	"gitee.com/geekbang/basic-go/webook/internal/article/service"
	"github.com/google/wire"
)

var thirdProvider = wire.NewSet(
	ioc.InitRedis,
	ioc.InitLogger,
	ioc.InitUserRpcClient,
	ioc.InitProducer,
	ioc.InitDB,
)

func Init() *App {
	wire.Build(
		thirdProvider,
		events.NewSaramaSyncProducer,
		cache.NewRedisArticleCache,
		dao.NewGORMArticleDAO,
		repository.NewArticleRepository,
		repository.NewGrpcAuthorRepository,
		service.NewArticleService,
		grpc.NewArticleServiceServer,
		ioc.InitGRPCxServer,
		wire.Struct(new(App), "*"),
	)
	return new(App)
}
