//go:build wireinject

package main

import (
	"gitee.com/geekbang/basic-go/webook/bff/ioc"
	"gitee.com/geekbang/basic-go/webook/bff/web"
	"gitee.com/geekbang/basic-go/webook/bff/web/jwt"
	"gitee.com/geekbang/basic-go/webook/pkg/wego"
	"github.com/google/wire"
)

func InitApp() *wego.App {
	wire.Build(
		ioc.InitLogger,
		ioc.InitRedis,
		ioc.InitEtcdClient,

		web.NewArticleHandler,
		web.NewUserHandler,
		web.NewRewardHandler,
		jwt.NewRedisHandler,

		ioc.InitUserClient,
		ioc.InitIntrClient,
		ioc.InitRewardClient,
		ioc.InitCodeClient,
		ioc.InitArticleClient,
		ioc.InitGinServer,
		wire.Struct(new(wego.App), "WebServer"),
	)
	return new(wego.App)
}
