package job

import (
	"context"
	"errors"
	"gitee.com/geekbang/basic-go/webook/internal/domain"
	"gitee.com/geekbang/basic-go/webook/internal/service"
	"gitee.com/geekbang/basic-go/webook/pkg/logger"
	"golang.org/x/sync/semaphore"
	"time"
)

type Executor interface {
	Name() string
	// Exec ctx 是整个任务调度的上下文
	// 当从 ctx.Done 有信号的时候，就需要考虑结束执行
	// 具体实现来控制
	Exec(ctx context.Context, j domain.CronJob) error
}

type LocalFuncExecutor struct {
	funcs map[string]func(ctx context.Context, j domain.CronJob) error
}

func NewLocalFuncExecutor() *LocalFuncExecutor {
	return &LocalFuncExecutor{funcs: make(map[string]func(ctx context.Context, j domain.CronJob) error)}
}

func (l *LocalFuncExecutor) AddLocalFunc(name string,
	fn func(ctx context.Context, j domain.CronJob) error) {
	l.funcs[name] = fn
}

func (l *LocalFuncExecutor) Name() string {
	return "local"
}

func (l *LocalFuncExecutor) Exec(ctx context.Context, j domain.CronJob) error {
	fn, ok := l.funcs[j.Name]
	if !ok {
		return errors.New("是不是忘记注册本地方法了？")
	}
	return fn(ctx, j)
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

func NewScheduler(svc service.CronJobService, l logger.LoggerV1) *Scheduler {
	return &Scheduler{
		execs:     make(map[string]Executor, 8),
		interval:  time.Second,
		svc:       svc,
		dbTimeout: time.Second,
		l:         l,
		// 假如说最多只有 100 个在运行
		limiter: semaphore.NewWeighted(100),
	}
}

func (s *Scheduler) RegisterJob(ctx context.Context, j CronJob) error {
	return s.svc.AddJob(ctx, j)
}

func (s *Scheduler) RegisterExecutor(exec Executor) {
	s.execs[exec.Name()] = exec
}

// Start 开始调度。当被取消，或者超时的时候，就会结束调度
func (s *Scheduler) Start(ctx context.Context) error {
	for {
		if ctx.Err() != nil {
			// 已经超时了，或者被取消运行，大多数时候，都是被取消了，或者说关闭了
			return ctx.Err()
		}
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
			// 这里可以考虑睡眠一段时间
			// 你也可以进一步细分不同的错误，如果是可以容忍的错误，
			// 就继续，不然就直接 return
			time.Sleep(s.interval)
			continue
		}
		exec, ok := s.execs[j.Executor]
		if !ok {
			// 不支持的执行方式。
			// 比如说，这里要求的runner是调用 gRPC，我们就不支持
			s.l.Error("不支持的Executor方式")
			j.CancelFunc()
			continue
		}
		// 要单独开一个 goroutine 来执行，这样我们就可以进入下一个循环了
		go func() {
			defer func() {
				s.limiter.Release(1)
				j.CancelFunc()
			}()

			err1 := exec.Exec(ctx, j)
			if err1 != nil {
				s.l.Error("调度任务执行失败",
					logger.Int64("id", j.Id),
					logger.Error(err1))
				return
			}
			err1 = s.svc.ResetNextTime(ctx, j)
			if err1 != nil {
				s.l.Error("更新下一次的执行失败", logger.Error(err1))
			}
		}()
	}
}

// CronJob 使用别名来做一个解耦
// 后续万一我们要加字段，就很方便扩展
type CronJob = domain.CronJob
