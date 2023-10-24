package repository

import (
	"context"
	"gitee.com/geekbang/basic-go/webook/internal/domain"
)

// HistoryRecordRepository 也就是一个增删改查的事情
type HistoryRecordRepository interface {
	AddRecord(ctx context.Context, r domain.HistoryRecord) error
}
