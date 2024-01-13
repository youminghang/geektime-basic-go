//go:build wireinject

package main

import (
	"gitee.com/geekbang/basic-go/webook/code/grpc"
	"gitee.com/geekbang/basic-go/webook/code/ioc"
	"gitee.com/geekbang/basic-go/webook/code/repository"
	"gitee.com/geekbang/basic-go/webook/code/repository/cache"
	"gitee.com/geekbang/basic-go/webook/code/service"
	"github.com/google/wire"
)

var thirdProvider = wire.NewSet(
	ioc.InitRedis,
	ioc.InitEtcdClient,
	ioc.InitLogger,
)

func Init() *App {
	wire.Build(
		thirdProvider,
		ioc.InitSmsRpcClient,
		cache.NewRedisCodeCache,
		repository.NewCachedCodeRepository,
		service.NewSMSCodeService,
		grpc.NewCodeServiceServer,
		ioc.InitGRPCxServer,
		wire.Struct(new(App), "*"),
	)
	return new(App)
}
