package service

import (
	"context"
	"errors"
	"gitee.com/geekbang/basic-go/webook/internal/domain"
	"gitee.com/geekbang/basic-go/webook/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

var ErrUserDuplicateEmail = repository.ErrUserDuplicate
var ErrInvalidUserOrPassword = errors.New("邮箱或者密码不正确")

type UserService interface {
	Signup(ctx context.Context, u domain.User) error
	FindOrCreate(ctx context.Context, phone string) (domain.User, error)
	Login(ctx context.Context, email, password string) (domain.User, error)
	Profile(ctx context.Context, id int64) (domain.User, error)
}

type userService struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) UserService {
	return &userService{
		repo: repo,
	}
}

func (svc *userService) Signup(ctx context.Context, u domain.User) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hash)
	return svc.repo.Create(ctx, u)
}

// FindOrCreate 如果手机号不存在，那么会初始化一个用户
func (svc *userService) FindOrCreate(ctx context.Context,
	phone string) (domain.User, error) {
	// 这是一种优化写法
	// 大部分人会命中这个分支
	u, err := svc.repo.FindByPhone(ctx, phone)
	if err != repository.ErrUserNotFound {
		return u, err
	}
	// 要执行注册
	err = svc.repo.Create(ctx, domain.User{
		Phone: phone,
	})
	// 注册有问题，但是又不是用户手机号码冲突，说明是系统错误
	if err != nil && err != repository.ErrUserDuplicate {
		return domain.User{}, err
	}
	// 主从模式下，这里要从主库中读取，暂时我们不需要考虑
	return svc.repo.FindByPhone(ctx, phone)
}

func (svc *userService) Login(ctx context.Context,
	email, password string) (domain.User, error) {
	u, err := svc.repo.FindByEmail(ctx, email)
	if err == repository.ErrUserNotFound {
		return domain.User{}, ErrInvalidUserOrPassword
	}
	err = bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	if err != nil {
		return domain.User{}, ErrInvalidUserOrPassword
	}
	return u, nil
}

func (svc *userService) Profile(ctx context.Context,
	id int64) (domain.User, error) {
	return svc.repo.FindById(ctx, id)
}
