package repository

import (
	"context"
	"gitee.com/geekbang/basic-go/webook/internal/repository/cache"
	"gitee.com/geekbang/basic-go/webook/internal/repository/dao"
	"gitee.com/geekbang/basic-go/webook/pkg/logger"
)

type InteractiveRepository interface {
	IncrReadCnt(ctx context.Context, biz string, bizId int64) error
	Get(ctx context.Context, biz string, bizId int64) (int64, error)
}

type CachedReadCntRepository struct {
	cache cache.InteractiveCache
	dao   dao.InteractiveDAO
	l     logger.LoggerV1
}

func (c *CachedReadCntRepository) IncrReadCnt(ctx context.Context,
	biz string, bizId int64) error {

	err := c.cache.IncreaseIfPresent(ctx, biz, bizId)
	if err != nil {
		return err
	}
	// 这边会有部分失败引起的不一致的问题，但是你其实不需要解决，
	// 因为阅读数不准确完全没有问题
	return c.dao.IncrReadCnt(ctx, biz, bizId)
}

func (c *CachedReadCntRepository) Get(ctx context.Context,
	biz string, bizId int64) (int64, error) {
	cnt, err := c.cache.Get(ctx, biz, bizId)
	if err == nil {
		return cnt, nil
	}
	cnt, err = c.dao.Get(ctx, biz, bizId)
	if err == nil {
		if er := c.cache.Set(ctx, biz, bizId, cnt); er != nil {
			c.l.Error("回写缓存失败",
				logger.Int64("bizId", bizId),
				logger.String("biz", biz),
				logger.Error(er))
		}
	}
	return cnt, err
}

func NewCachedReadCntRepository(dao dao.InteractiveDAO,
	cache cache.InteractiveCache, l logger.LoggerV1) InteractiveRepository {
	return &CachedReadCntRepository{
		dao:   dao,
		cache: cache,
		l:     l,
	}
}
