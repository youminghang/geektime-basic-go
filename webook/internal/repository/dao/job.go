package dao

import (
	"context"
	"gorm.io/gorm"
	"time"
)

var ErrNoMoreJob = gorm.ErrRecordNotFound

type JobDAO interface {
	Preempt(ctx context.Context) (Job, error)
	UpdateNextTime(ctx context.Context, id int64, t time.Time) error
}

type GORMJobDAO struct {
	db *gorm.DB
}

func (dao *GORMJobDAO) Preempt(ctx context.Context) (Job, error) {
	db := dao.db.WithContext(ctx)
	for {
		db.First("ctime")
	}
}

func (dao *GORMJobDAO) UpdateNextTime(ctx context.Context, id int64, t time.Time) error {
	//TODO implement me
	panic("implement me")
}

type Job struct {
	Id         int64 `gorm:"primaryKey,autoIncrement"`
	Name       string
	Executor   string
	Cfg        string
	Expression string
	NextTime   int64 `gorm:"index"`
	Ctime      int64
	Utime      int64
}
