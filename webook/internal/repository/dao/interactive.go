package dao

import (
	"context"
	"gorm.io/gorm"
)

type InteractiveDAO interface {
	IncrReadCnt(ctx context.Context, biz string, bizId int64) error
	Get(ctx context.Context, biz string, bizId int64) (int64, error)
}

type GORMReadCntDAO struct {
	db *gorm.DB
}

func NewGORMReadCntDAO(db *gorm.DB) InteractiveDAO {
	return &GORMReadCntDAO{
		db: db,
	}
}

func (dao *GORMReadCntDAO) IncrReadCnt(ctx context.Context, biz string, bizId int64) error {
	//TODO implement me
	panic("implement me")
}

func (dao *GORMReadCntDAO) Get(ctx context.Context, biz string, bizId int64) (int64, error) {
	//TODO implement me
	panic("implement me")
}
