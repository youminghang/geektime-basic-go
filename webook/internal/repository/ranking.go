package repository

import (
	"context"
	"gitee.com/geekbang/basic-go/webook/internal/repository/cache"
)

type RankingRepository interface {
	ReplaceTopN(ctx context.Context, ids []int64) error
	GetTopN(ctx context.Context) ([]int64, error)
}

type CachedRankingRepository struct {
	cache cache.RankingCache
}

func NewCachedRankingRepository(cache cache.RankingCache) RankingRepository {
	return &CachedRankingRepository{cache: cache}
}

func (c *CachedRankingRepository) ReplaceTopN(ctx context.Context,
	ids []int64) error {
	return c.cache.Set(ctx, ids)
}

func (c *CachedRankingRepository) GetTopN(ctx context.Context) ([]int64, error) {
	return c.cache.Get(ctx)
}
