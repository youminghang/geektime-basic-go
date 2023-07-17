package repository

import (
	"context"
	"gitee.com/geekbang/basic-go/webook/internal/domain"
	"gitee.com/geekbang/basic-go/webook/internal/repository/dao"
)

var ErrUserDuplicateEmail = dao.ErrUserDuplicateEmail
var ErrUserNotFound = dao.ErrDataNotFound

type UserRepository struct {
	dao *dao.UserDAO
}

func NewUserRepository(d *dao.UserDAO) *UserRepository {
	return &UserRepository{
		dao: d,
	}
}

func (ur *UserRepository) Create(ctx context.Context, u domain.User) error {
	err := ur.dao.Insert(ctx, dao.User{
		Email:    u.Email,
		Password: u.Password,
	})
	return err
}

func (ur *UserRepository) FindByEmail(ctx context.Context,
	email string) (domain.User, error) {
	u, err := ur.dao.FindByEmail(ctx, email)

	// 因为我们用的是别名机制，所以这里不用这么写
	//if err == gorm.ErrRecordNotFound {
	//	return ErrUserNotFound
	//}
	return domain.User{
		Id:       u.Id,
		Email:    u.Email,
		Password: u.Password,
	}, err
}

func (ur *UserRepository) FindById(ctx context.Context,
	id int64) (domain.User, error) {
	u, err := ur.dao.FindById(ctx, id)
	return domain.User{
		Email:    u.Email,
		Password: u.Password,
	}, err
}
