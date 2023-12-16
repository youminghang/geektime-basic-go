//go:build wireinject

package main

import (
	"gitee.com/geekbang/basic-go/webook/internal/user/grpc"
	"gitee.com/geekbang/basic-go/webook/internal/user/ioc"
	"gitee.com/geekbang/basic-go/webook/internal/user/repository"
	"gitee.com/geekbang/basic-go/webook/internal/user/repository/cache"
	"gitee.com/geekbang/basic-go/webook/internal/user/repository/dao"
	"gitee.com/geekbang/basic-go/webook/internal/user/service"

	"github.com/google/wire"
)

var thirdProvider = wire.NewSet(
	ioc.InitLogger,
	ioc.InitDB,
	ioc.InitRedis,
)

func Init() *App {
	wire.Build(
		thirdProvider,
		cache.NewRedisUserCache,
		dao.NewGORMUserDAO,
		repository.NewCachedUserRepository,
		service.NewUserService,
		grpc.NewUserServiceServer,
		ioc.InitGRPCxServer,
		wire.Struct(new(App), "*"),
	)
	return new(App)
}
