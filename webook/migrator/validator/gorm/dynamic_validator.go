package gorm

import (
	"gitee.com/geekbang/basic-go/webook/migrator"
	"gitee.com/geekbang/basic-go/webook/migrator/events"
	"gitee.com/geekbang/basic-go/webook/pkg/logger"
	"gorm.io/gorm"
)

type DynamicValidator[T migrator.Entity] struct {
	srcFirst *Validator[T]
	dstFirst *Validator[T]
}

func NewDynamicValidator[T migrator.Entity](
	src *gorm.DB,
	dst *gorm.DB,
	l logger.LoggerV1,
	producer events.Producer,
) *DynamicValidator[T] {
	return &DynamicValidator[T]{
		srcFirst: NewValidator[T](src, dst, "SRC", l, producer),
		dstFirst: NewValidator[T](dst, src, "DST", l, producer),
	}
}
