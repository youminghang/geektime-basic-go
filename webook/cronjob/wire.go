//go:build wireinject

package main

import (
	"gitee.com/geekbang/basic-go/webook/cronjob/grpc"
	"gitee.com/geekbang/basic-go/webook/cronjob/ioc"
	"gitee.com/geekbang/basic-go/webook/cronjob/repository"
	"gitee.com/geekbang/basic-go/webook/cronjob/repository/dao"
	"gitee.com/geekbang/basic-go/webook/cronjob/service"
	"github.com/google/wire"
)

var serviceProviderSet = wire.NewSet(
	dao.NewGORMJobDAO,
	repository.NewPreemptCronJobRepository,
	service.NewCronJobService)

var thirdProvider = wire.NewSet(
	ioc.InitDB,
	ioc.InitLogger,
)

func Init() *App {
	wire.Build(
		thirdProvider,
		serviceProviderSet,
		grpc.NewCronJobServiceServer,
		ioc.InitGRPCxServer,
		wire.Struct(new(App), "*"),
	)
	return new(App)
}
