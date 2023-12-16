//go:build wireinject

package main

import (
	"gitee.com/geekbang/basic-go/webook/internal/code/grpc"
	"gitee.com/geekbang/basic-go/webook/internal/code/ioc"
	"gitee.com/geekbang/basic-go/webook/internal/code/repository"
	"gitee.com/geekbang/basic-go/webook/internal/code/repository/cache"
	"gitee.com/geekbang/basic-go/webook/internal/code/service"
	"github.com/google/wire"
)

var thirdProvider = wire.NewSet(
	ioc.InitRedis,
	ioc.InitSmsRpcClient,
)

func Init() *App {
	wire.Build(
		thirdProvider,
		cache.NewRedisCodeCache,
		repository.NewCachedCodeRepository,
		service.NewSMSCodeService,
		grpc.NewCodeServiceServer,
		ioc.InitGRPCxServer,
		wire.Struct(new(App), "*"),
	)
	return new(App)
}
