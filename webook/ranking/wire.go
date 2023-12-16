//go:build wireinject

package main

import (
	"gitee.com/geekbang/basic-go/webook/internal/ranking/grpc"
	"gitee.com/geekbang/basic-go/webook/internal/ranking/ioc"
	"gitee.com/geekbang/basic-go/webook/internal/ranking/repository"
	"gitee.com/geekbang/basic-go/webook/internal/ranking/repository/cache"
	"gitee.com/geekbang/basic-go/webook/internal/ranking/service"
	"github.com/google/wire"
)

var serviceProviderSet = wire.NewSet(
	cache.NewRankingLocalCache,
	cache.NewRedisRankingCache,
	repository.NewCachedRankingRepository,
	service.NewBatchRankingService,
)

var thirdProvider = wire.NewSet(
	ioc.InitRedis,
	ioc.InitInterActiveRpcClient,
	ioc.InitArticleRpcClient,
)

func Init() *App {
	wire.Build(
		thirdProvider,
		serviceProviderSet,
		grpc.NewRankingServiceServer,
		ioc.InitGRPCxServer,
		wire.Struct(new(App), "*"),
	)
	return new(App)
}
