//go:build wireinject

package wire

import (
	"gitee.com/geekbang/basic-go/wire/repository"
	"gitee.com/geekbang/basic-go/wire/repository/dao"
	"github.com/google/wire"
)

func InitUserRepository() *repository.UserRepository {
	wire.Build(repository.NewUserRepository, dao.NewUserDAO, InitDB)
	// 下面随便返回一个
	return &repository.UserRepository{}
}
