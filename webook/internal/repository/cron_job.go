package repository

import (
	"context"
	"gitee.com/geekbang/basic-go/webook/internal/domain"
	"time"
)

type CronJobRepository interface {
	Preempt(ctx context.Context) (domain.CronJob, error)
	UpdateNextTime(ctx context.Context, id int64, t time.Time) error
}
