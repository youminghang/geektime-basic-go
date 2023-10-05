package service

import (
	"context"
	"gitee.com/geekbang/basic-go/webook/internal/repository"
)

type InteractiveService interface {
	IncrReadCnt(ctx context.Context, biz string, bizId int64) error
	// Like 点赞/取消点赞
	Like(ctx context.Context,
		biz string, bizId int64, uid int64, like bool) error
	// Collect 收藏
	Collect(ctx context.Context, biz string, bizId int64) error
	Get(ctx context.Context, biz string, bizId int64) (int64, error)
}

type interactiveService struct {
	// 直接组合，省事
	repository.InteractiveRepository
}

func NewInteractiveService(repo repository.InteractiveRepository) InteractiveService {
	return &interactiveService{
		InteractiveRepository: repo,
	}
}
