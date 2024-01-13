//go:build wireinject

package main

import (
	"gitee.com/geekbang/basic-go/webook/pkg/wego"
	"gitee.com/geekbang/basic-go/webook/user/grpc"
	"gitee.com/geekbang/basic-go/webook/user/ioc"
	"gitee.com/geekbang/basic-go/webook/user/repository"
	"gitee.com/geekbang/basic-go/webook/user/repository/cache"
	"gitee.com/geekbang/basic-go/webook/user/repository/dao"
	"gitee.com/geekbang/basic-go/webook/user/service"

	"github.com/google/wire"
)

var thirdProvider = wire.NewSet(
	ioc.InitLogger,
	ioc.InitDB,
	ioc.InitRedis,
	ioc.InitEtcdClient,
)

func Init() *wego.App {
	wire.Build(
		thirdProvider,
		cache.NewRedisUserCache,
		dao.NewGORMUserDAO,
		repository.NewCachedUserRepository,
		service.NewUserService,
		grpc.NewUserServiceServer,
		ioc.InitGRPCxServer,
		wire.Struct(new(wego.App), "GRPCServer"),
	)
	return new(wego.App)
}
