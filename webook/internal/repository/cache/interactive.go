package cache

import "context"

type InteractiveCache interface {

	// IncreaseIfPresent 如果在缓存中有对应的数据，就 +1
	IncreaseIfPresent(ctx context.Context,
		biz string, bizId int64) error
	// Get 查询缓存中数据
	Get(ctx context.Context, biz string, bizId int64) (int64, error)
	Set(ctx context.Context, biz string, bizId int64, cnt int64) error
}
