//go:build wireinject

package main

import (
	"gitee.com/geekbang/basic-go/webook/pkg/wego"
	"gitee.com/geekbang/basic-go/webook/sms/grpc"
	ioc2 "gitee.com/geekbang/basic-go/webook/sms/ioc"
	"github.com/google/wire"
)

func Init() *wego.App {
	wire.Build(
		ioc2.InitSmsTencentService,
		grpc.NewSmsServiceServer,
		ioc2.InitGRPCxServer,
		wire.Struct(new(wego.App), "GRPCServer"),
	)
	return new(wego.App)
}
