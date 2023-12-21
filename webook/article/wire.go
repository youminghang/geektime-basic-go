//go:build wireinject

package main

import (
	"gitee.com/geekbang/basic-go/webook/article/events"
	"gitee.com/geekbang/basic-go/webook/article/grpc"
	"gitee.com/geekbang/basic-go/webook/article/ioc"
	"gitee.com/geekbang/basic-go/webook/article/repository"
	"gitee.com/geekbang/basic-go/webook/article/repository/cache"
	"gitee.com/geekbang/basic-go/webook/article/repository/dao"
	"gitee.com/geekbang/basic-go/webook/article/service"
	"gitee.com/geekbang/basic-go/webook/pkg/wego"
	"github.com/google/wire"
)

var thirdProvider = wire.NewSet(
	ioc.InitRedis,
	ioc.InitLogger,
	ioc.InitUserRpcClient,
	ioc.InitProducer,
	ioc.InitDB,
)

func Init() *wego.App {
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
		wire.Struct(new(wego.App), "GRPCServer"),
	)
	return new(wego.App)
}
