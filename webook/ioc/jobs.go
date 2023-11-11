package ioc

import (
	"gitee.com/geekbang/basic-go/webook/internal/job"
	"gitee.com/geekbang/basic-go/webook/internal/service"
	"gitee.com/geekbang/basic-go/webook/pkg/logger"
	rlock "github.com/gotomicro/redis-lock"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/robfig/cron/v3"
	"time"
)

func InitRankingJob(svc service.RankingService,
	client *rlock.Client,
	l logger.LoggerV1) *job.RankingJob {
	return job.NewRankingJob(svc, client, l, time.Second*30)
}

func InitJobs(l logger.LoggerV1, rankingJob *job.RankingJob) *cron.Cron {
	bd := job.NewCronJobBuilder(l, prometheus.SummaryOpts{
		Namespace: "geekbang_daming",
		Subsystem: "webook",
		Name:      "cron_job",
		Help:      "定时任务",
	})
	expr := cron.New(cron.WithSeconds())
	_, err := expr.AddJob("@every 1m",
		bd.Build(rankingJob))
	if err != nil {
		panic(err)
	}
	return expr
}
