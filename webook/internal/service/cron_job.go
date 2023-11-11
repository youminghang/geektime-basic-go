package service

import (
	"context"
	"gitee.com/geekbang/basic-go/webook/internal/domain"
	"gitee.com/geekbang/basic-go/webook/internal/repository"
	"time"
)

type CronJobService interface {
	Preempt(ctx context.Context) (domain.CronJob, error)
	ResetNextTime(ctx context.Context, job domain.CronJob) error
}

type cronJobService struct {
	repo repository.CronJobRepository
}

func (s *cronJobService) Preempt(ctx context.Context) (domain.CronJob, error) {
	return s.repo.Preempt(ctx)
}

func (s *cronJobService) ResetNextTime(ctx context.Context,
	jd domain.CronJob) error {
	// 计算下一次的时间
	t := jd.Next(time.Now())
	// 我们认为这是不需要继续执行了
	if !t.IsZero() {
		return s.repo.UpdateNextTime(ctx, jd.Id, t)
	}
	return nil
}
