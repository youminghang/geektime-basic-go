//go:build wireinject

package main

import (
	"gitee.com/geekbang/basic-go/webook/internal/oauth2/grpc"
	ioc2 "gitee.com/geekbang/basic-go/webook/internal/oauth2/ioc"
	"github.com/google/wire"
)

var thirdProvider = wire.NewSet(
	ioc2.InitLogger,
)

func Init() *App {
	wire.Build(
		thirdProvider,
		ioc2.InitPrometheus,
		grpc.NewOauth2ServiceServer,
		ioc2.InitGRPCxServer,
		wire.Struct(new(App), "*"),
	)
	return new(App)
}
