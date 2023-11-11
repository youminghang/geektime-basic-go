package repository

import (
	"context"
	"gitee.com/geekbang/basic-go/webook/internal/domain"
	"gitee.com/geekbang/basic-go/webook/internal/repository/cache"
	"github.com/ecodeclub/ekit/syncx/atomicx"
)

type RankingRepository interface {
	ReplaceTopN(ctx context.Context, arts []domain.Article) error
	GetTopN(ctx context.Context) ([]domain.Article, error)
}

type CachedRankingRepository struct {
	redisCache *cache.RedisRankingCache
	localCache *cache.RankingLocalCache
	// 你也可以考虑将这个本地缓存塞进去 RankingCache 里面，作为一个实现
	topN atomicx.Value[[]domain.Article]
}

func NewCachedRankingRepository(
	redisCache *cache.RedisRankingCache,
	localCache *cache.RankingLocalCache) RankingRepository {
	return &CachedRankingRepository{
		redisCache: redisCache,
		localCache: localCache,
	}
}

func (c *CachedRankingRepository) ReplaceTopN(ctx context.Context,
	arts []domain.Article) error {
	// 这一步必然不会出错
	_ = c.localCache.Set(ctx, arts)
	return c.redisCache.Set(ctx, arts)
}

func (c *CachedRankingRepository) GetTopN(ctx context.Context) ([]domain.Article, error) {
	arts, err := c.localCache.Get(ctx)
	if err == nil {
		return arts, nil
	}
	// 回写本地缓存
	arts, err = c.redisCache.Get(ctx)
	if err == nil {
		_ = c.localCache.Set(ctx, arts)
	} else {
		// 这里，我们没有进一步区分是什么原因导致的 Redis 错误
		return c.localCache.ForceGet(ctx)
	}
	return arts, err
}
