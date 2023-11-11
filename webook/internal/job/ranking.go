package job

import (
	"context"
	"gitee.com/geekbang/basic-go/webook/internal/service"
	"gitee.com/geekbang/basic-go/webook/pkg/logger"
	rlock "github.com/gotomicro/redis-lock"
	"time"
)

type RankingJob struct {
	svc service.RankingService
	// 一次运行的超时时间
	timeout    time.Duration
	lockClient *rlock.Client
	l          logger.LoggerV1
	key        string
}

func NewRankingJob(svc service.RankingService,
	lockClient *rlock.Client,
	l logger.LoggerV1,
	timeout time.Duration) *RankingJob {
	return &RankingJob{svc: svc,
		lockClient: lockClient,
		timeout:    timeout,
		key:        "job:ranking",
		l:          l,
	}
}

func (r *RankingJob) Name() string {
	return "ranking"
}

func (r *RankingJob) Run() error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()
	// 加锁本身，我们使用一个ctx
	// 本身我们这里设计的就是要在 r.timeout 内计算完成
	// 刚好也做成分布式锁的超时时间
	lock, err := r.lockClient.Lock(ctx, r.key, r.timeout,
		// 每隔 100ms 重试一次，每次重试的超时时间是 1s
		&rlock.FixIntervalRetry{
			Interval: time.Millisecond * 100,
			Max:      3,
		}, time.Second)
	// 我们这里不需要处理 error，因为大部分情况下，可以相信别的节点会继续拿锁
	if err != nil {
		return err
	}

	defer func() {
		// 释放锁，再来一个独立的 ctx
		ctx, cancel = context.WithTimeout(context.Background(), time.Second*3)
		defer cancel()
		// 释放锁失败，但是也不需要作什么，因为这个分布式锁会在过期时间之后自动释放
		err = lock.Unlock(ctx)
		if err != nil {
			r.l.Error("释放分布式锁失败",
				logger.Error(err),
				logger.String("name", r.Name()))
		}
	}()
	return r.run()
}

func (r *RankingJob) run() error {
	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()
	return r.svc.RankTopN(ctx)
}

var _ Job = (*RankingJob)(nil)
