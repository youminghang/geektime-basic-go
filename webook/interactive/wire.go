//go:build wireinject

package main

import (
	"gitee.com/geekbang/basic-go/webook/interactive/events"
	"gitee.com/geekbang/basic-go/webook/interactive/grpc"
	"gitee.com/geekbang/basic-go/webook/interactive/ioc"
	"gitee.com/geekbang/basic-go/webook/interactive/repository"
	"gitee.com/geekbang/basic-go/webook/interactive/repository/cache"
	"gitee.com/geekbang/basic-go/webook/interactive/repository/dao"
	"gitee.com/geekbang/basic-go/webook/interactive/service"
	"github.com/google/wire"
)

var serviceProviderSet = wire.NewSet(
	dao.NewGORMInteractiveDAO,
	cache.NewRedisInteractiveCache,
	repository.NewCachedInteractiveRepository,
	service.NewInteractiveService)

var thirdProvider = wire.NewSet(
	ioc.InitRedis, ioc.InitDB,
	ioc.InitLogger,
	ioc.InitKafka,
)

func Init() *App {
	wire.Build(
		thirdProvider,
		serviceProviderSet,
		grpc.NewInteractiveServiceServer,
		events.NewInteractiveReadEventConsumer,
		ioc.InitGRPCxServer,
		ioc.NewConsumers,
		wire.Struct(new(App), "*"),
	)
	return new(App)
}
