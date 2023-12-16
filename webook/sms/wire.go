//go:build wireinject

package main

import (
	"gitee.com/geekbang/basic-go/webook/internal/sms/grpc"
	ioc2 "gitee.com/geekbang/basic-go/webook/internal/sms/ioc"
	"github.com/google/wire"
)

func Init() *App {
	wire.Build(
		ioc2.InitSmsTencentService,
		grpc.NewSmsServiceServer,
		ioc2.InitGRPCxServer,
		wire.Struct(new(App), "*"),
	)
	return new(App)
}
