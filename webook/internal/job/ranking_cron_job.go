package job

import (
	"context"
	"gitee.com/geekbang/basic-go/webook/internal/service"
	"gitee.com/geekbang/basic-go/webook/pkg/logger"
	"golang.org/x/sync/semaphore"
	"time"
)

type Executor interface {
	// Exec ctx 是整个任务调度的上下文
	// 当从 ctx.Done 有信号的时候，就需要考虑结束执行
	// 具体实现来控制
	Exec(ctx context.Context, cfg string) error
}

type LocalFuncExecutor struct {
	funcs map[string]func(ctx context.Context) error
}

// Scheduler 调度器
type Scheduler struct {
	execs     map[string]Executor
	interval  time.Duration
	svc       service.CronJobService
	dbTimeout time.Duration
	l         logger.LoggerV1
	limiter   *semaphore.Weighted
}

func (s *Scheduler) Add(j Job) {

}

// Start 开始调度。当被取消，或者超时的时候，就会结束调度
func (s *Scheduler) Start(ctx context.Context) error {
	for {
		err := s.limiter.Acquire(ctx, 1)
		if err != nil {
			// 正常来说，只有 ctx 超时或者取消才会进来这里
			return err
		}
		// 抢占，获得可以运行的资格
		// 数据库查询的时候，超时时间是要短的
		dbCtx, cancel := context.WithTimeout(ctx, s.dbTimeout)
		j, err := s.svc.Preempt(dbCtx)
		cancel()
		if err != nil {
			// 没有抢占到，进入下一个循环
			continue
		}
		exec, ok := s.execs[j.Executor]
		if !ok {
			// 不支持的执行方式。
			// 比如说，这里要求的runner是调用 gRPC，我们就不支持
			s.l.Error("不支持的Executor方式")
		}
		// 要单独开一个 goroutine 来执行，这样我们就可以进入下一个循环了
		go func() {
			// 不支持的 runner
			err1 := exec.Exec(ctx, j.Cfg)
			if err1 == nil {

			}
			s.limiter.Release(1)
		}()
	}
}
