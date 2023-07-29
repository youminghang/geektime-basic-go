package repository

import (
	"context"
	"gitee.com/geekbang/basic-go/webook/internal/repository/cache"
)

var (
	ErrCodeVerifyTooManyTimes = cache.ErrCodeVerifyTooManyTimes
	ErrCodeSendTooMany        = cache.ErrCodeSendTooMany
)

type CodeRepository struct {
	cache *cache.CodeCache
}

func NewCodeRepository(c *cache.CodeCache) *CodeRepository {
	return &CodeRepository{
		cache: c,
	}
}

func (repo *CodeRepository) Store(ctx context.Context,
	biz string,
	phone string,
	code string) error {
	err := repo.cache.Set(ctx, biz, phone, code)
	return err
}

// Verify 比较验证码。如果验证码相等，那么删除；
func (repo *CodeRepository) Verify(ctx context.Context,
	biz string, phone string, inputCode string) (bool, error) {
	return repo.cache.Verify(ctx, biz, phone, inputCode)
}
