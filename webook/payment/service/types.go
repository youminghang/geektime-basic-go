package service

import (
	"context"
	"gitee.com/geekbang/basic-go/webook/payment/domain"
)

type Payment interface {
	Pay(ctx context.Context, payment domain.Payment) error
}
