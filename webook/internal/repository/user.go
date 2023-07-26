package repository

import (
	"context"
	"gitee.com/geekbang/basic-go/webook/internal/domain"
	"gitee.com/geekbang/basic-go/webook/internal/repository/cache"
	"gitee.com/geekbang/basic-go/webook/internal/repository/dao"
)

var ErrUserDuplicateEmail = dao.ErrUserDuplicateEmail
var ErrUserNotFound = dao.ErrDataNotFound

type UserRepository struct {
	dao   *dao.UserDAO
	cache *cache.UserCache
}

func NewUserRepository(d *dao.UserDAO, c *cache.UserCache) *UserRepository {
	return &UserRepository{
		dao:   d,
		cache: c,
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
	u, err := ur.cache.Get(ctx, id)
	// 注意这里的处理方式
	if err == nil {
		return u, err
	}
	ue, err := ur.dao.FindById(ctx, id)
	if err != nil {
		return domain.User{}, err
	}
	u = domain.User{
		Id:       ue.Id,
		Email:    ue.Email,
		Password: ue.Password,
	}
	// 忽略掉这里的错误
	_ = ur.cache.Set(ctx, u)
	return u, nil
}

func (ur *UserRepository) FindByIdV1(ctx context.Context,
	id int64) (domain.User, error) {
	u, err := ur.cache.Get(ctx, id)
	switch err {
	case nil:
		return u, err
	case cache.ErrKeyNotExist:
		ue, err := ur.dao.FindById(ctx, id)
		if err != nil {
			return domain.User{}, err
		}
		u = domain.User{
			Id:       ue.Id,
			Email:    ue.Email,
			Password: ue.Password,
		}
		// 忽略掉这里的错误
		_ = ur.cache.Set(ctx, u)
		return u, nil
	default:
		return domain.User{}, err
	}
}
