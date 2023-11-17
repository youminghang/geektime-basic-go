package service

import (
	"context"
	"gitee.com/geekbang/basic-go/webook/internal/domain"
	"gitee.com/geekbang/basic-go/webook/internal/repository"
	"gitee.com/geekbang/basic-go/webook/pkg/logger"
	"time"
)

var ErrNoMoreJob = repository.ErrNoMoreJob

//go:generate mockgen -source=./cron_job.go -package=svcmocks -destination=mocks/cron_job.mock.go CronJobService
type CronJobService interface {
	Preempt(ctx context.Context) (domain.CronJob, error)
	ResetNextTime(ctx context.Context, job domain.CronJob) error
	AddJob(ctx context.Context, j domain.CronJob) error
}

type cronJobService struct {
	repo            repository.CronJobRepository
	l               logger.LoggerV1
	refreshInterval time.Duration
}

func NewCronJobService(
	repo repository.CronJobRepository,
	l logger.LoggerV1) CronJobService {
	return &cronJobService{
		repo:            repo,
		l:               l,
		refreshInterval: time.Second * 10,
	}
}

func (s *cronJobService) AddJob(ctx context.Context, j domain.CronJob) error {
	j.NextTime = j.Next(time.Now())
	return s.repo.AddJob(ctx, j)
}

func (s *cronJobService) Preempt(ctx context.Context) (domain.CronJob, error) {
	j, err := s.repo.Preempt(ctx)
	if err != nil {
		return domain.CronJob{}, err
	}
	ch := make(chan struct{})
	go func() {
		// 这边要启动一个 goroutine 开始续约，也就是在持续占有期间
		// 假定说我们这里是十秒钟续约一次
		ticker := time.NewTicker(s.refreshInterval)
		defer ticker.Stop()
		for {
			select {
			case <-ch:
				// 退出续约循环
				return
			case <-ticker.C:
				s.refresh(j.Id)
			}
		}
	}()
	// 只能调用一次，也就是放弃续约。这时候要把状态还原回去
	j.CancelFunc = func() {
		close(ch)
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		err := s.repo.Release(ctx, j.Id)
		if err != nil {
			s.l.Error("释放任务失败",
				logger.Error(err),
				logger.Int64("id", j.Id))
		}
	}
	return j, nil
}

func (s *cronJobService) refresh(id int64) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	err := s.repo.UpdateUtime(ctx, id)
	if err != nil {
		s.l.Error("续约失败",
			logger.Int64("jid", id),
			logger.Error(err))
	}
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
