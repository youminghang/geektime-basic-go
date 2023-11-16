package job

import (
	"gitee.com/geekbang/basic-go/webook/pkg/logger"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/robfig/cron/v3"
	"strconv"
	"time"
)

// CronJobBuilder 根据需要加各种监控
// 这种 Builder 写法是为了避开 prometheus 重复注册的问题
// 也可以用来组装不同的装饰器，比较灵活
type CronJobBuilder struct {
	vector *prometheus.SummaryVec
	l      logger.LoggerV1
}

func NewCronJobBuilder(l logger.LoggerV1,
	opt prometheus.SummaryOpts) *CronJobBuilder {
	vector := prometheus.NewSummaryVec(opt,
		[]string{"name", "success"})
	prometheus.MustRegister(vector)
	return &CronJobBuilder{vector: vector, l: l}
}

func (m *CronJobBuilder) Build(job Job) cron.Job {
	name := job.Name()
	return cronJobAdapterFunc(func() {
		start := time.Now()
		m.l.Debug("任务开始",
			logger.String("name", name),
			logger.String("time", start.String()),
		)
		err := job.Run()
		duration := time.Since(start)
		if err != nil {
			m.l.Error("任务执行失败",
				logger.String("name", name),
				logger.Error(err))
		}
		m.l.Debug("任务结束",
			logger.String("name", name))
		m.vector.WithLabelValues(name,
			strconv.FormatBool(err == nil)).
			Observe(float64(duration.Milliseconds()))
	})
}

var _ cron.Job = (*cronJobAdapterFunc)(nil)

type cronJobAdapterFunc func()

func (c cronJobAdapterFunc) Run() {
	c()
}
