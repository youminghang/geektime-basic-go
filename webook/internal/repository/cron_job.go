package repository

import (
	"context"
	"gitee.com/geekbang/basic-go/webook/internal/domain"
	"gitee.com/geekbang/basic-go/webook/internal/repository/dao"
	"time"
)

var ErrNoMoreJob = dao.ErrNoMoreJob

type CronJobRepository interface {
	Preempt(ctx context.Context) (domain.CronJob, error)
	UpdateNextTime(ctx context.Context, id int64, t time.Time) error
}

type PreemptCronJobRepository struct {
	dao dao.JobDAO
}

func (p *PreemptCronJobRepository) Preempt(ctx context.Context) (domain.CronJob, error) {
	j, err := p.dao.Preempt(ctx)
	if err != nil {
		return domain.CronJob{}, err
	}
	return domain.CronJob{
		Id:         j.Id,
		Executor:   j.Executor,
		Cfg:        j.Cfg,
		Expression: j.Expression,
	}, nil
}

func (p *PreemptCronJobRepository) UpdateNextTime(ctx context.Context, id int64, t time.Time) error {
	return p.dao.UpdateNextTime(ctx, id, t)
}
