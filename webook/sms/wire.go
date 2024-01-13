//go:build wireinject

package main

import (
	"gitee.com/geekbang/basic-go/webook/pkg/wego"
	"gitee.com/geekbang/basic-go/webook/sms/grpc"
	"gitee.com/geekbang/basic-go/webook/sms/ioc"
	"github.com/google/wire"
)

func Init() *wego.App {
	wire.Build(
		ioc.InitLogger,
		ioc.InitEtcdClient,
		ioc.InitSmsTencentService,
		grpc.NewSmsServiceServer,
		ioc.InitGRPCxServer,
		wire.Struct(new(wego.App), "GRPCServer"),
	)
	return new(wego.App)
}
