//go:build wireinject

package startup

import (
	"gitee.com/geekbang/basic-go/webook/interactive/grpc"
	"gitee.com/geekbang/basic-go/webook/interactive/repository"
	"gitee.com/geekbang/basic-go/webook/interactive/repository/cache"
	"gitee.com/geekbang/basic-go/webook/interactive/repository/dao"
	"gitee.com/geekbang/basic-go/webook/interactive/service"
	"github.com/google/wire"
)

var thirdProvider = wire.NewSet(
	InitRedis, InitTestDB,
	InitLog,
	InitKafka,
)

func InitGRPCServer() *grpc.InteractiveServiceServer {
	wire.Build(
		grpc.NewInteractiveServiceServer,
		thirdProvider,
		dao.NewGORMInteractiveDAO,
		cache.NewRedisInteractiveCache,
		repository.NewCachedInteractiveRepository,
		service.NewInteractiveService,
	)
	return new(grpc.InteractiveServiceServer)
}
