package gorm

import (
	"context"
	"gitee.com/geekbang/basic-go/webook/migrator"
	"gitee.com/geekbang/basic-go/webook/migrator/events"
	"gitee.com/geekbang/basic-go/webook/pkg/logger"
	"github.com/ecodeclub/ekit/slice"
	"gorm.io/gorm"
	"time"
)

type Validator[T migrator.Entity] struct {
	base   *gorm.DB
	target *gorm.DB

	// 这边需要告知，是以 SRC 为准，还是以 DST 为准
	// 修复数据需要知道
	direction string

	batchSize int

	l        logger.LoggerV1
	producer events.Producer

	minUtime int64
	// utime 最大值距离当前时间戳多远
	nowDelta int64
}

func NewValidator[T migrator.Entity](
	base *gorm.DB,
	target *gorm.DB,
	direction string,
	l logger.LoggerV1,
	producer events.Producer,
) *Validator[T] {
	return &Validator[T]{
		base:      base,
		target:    target,
		direction: direction,
		l:         l,
		producer:  producer,
		batchSize: 100,
		minUtime:  0,
		nowDelta:  0,
	}
}

// Validate 执行校验。
// 分成两步：
// 1. from => to
func (v *Validator[T]) Validate(ctx context.Context) error {
	err := v.baseToTarget(ctx)
	if err != nil {
		return err
	}
	return v.targetToBase(ctx)
}

// baseToTarget 从 first 到 second 的验证
func (v *Validator[T]) baseToTarget(ctx context.Context) error {
	offset := -1
	base := v.base.WithContext(ctx)
	target := v.target.WithContext(ctx)
	for {
		// 前面的入口是唯一的，所以在这里自增代码会更加简洁
		offset++
		var src T
		// 这里假定主键的规范都是叫做 id，基本上大部分公司都有这种规范
		// 在叠加了 utime 作为查询条件之后，就可以达成增量校验的效果。
		maxUtime := time.Now().UnixMilli() - v.nowDelta

		err := base.WithContext(ctx).
			Where("utime <? AND utime > ?", maxUtime, v.minUtime).
			Order("id").
			Offset(offset).First(&src).Error
		if err == gorm.ErrRecordNotFound {
			// 已经没有数据了
			return nil
		}
		if err != nil {
			v.l.Error("src => dst 查询源表失败", logger.Error(err))
			continue
		}
		var dst T
		err = target.WithContext(ctx).
			Where("id=?", src.ID()).First(&dst).Error
		// 这边要考虑不同的 error
		switch err {
		case gorm.ErrRecordNotFound:
			v.notify(src.ID(), events.InconsistentEventTypeTargetMissing)
		case nil:
			// 查询到了数据
			equal := src.CompareTo(dst)
			if !equal {
				v.notify(src.ID(), events.InconsistentEventTypeNotEqual)
			}
		default:
			v.l.Error("src => dst 查询目标表失败", logger.Error(err))
			continue
		}
	}
}

// targetToBase 反过来，执行 target 到 base 的验证
// 这是为了找出 dst 中多余的数据
func (v *Validator[T]) targetToBase(ctx context.Context) error {
	// 这个我们只需要找出 src 中不存在的 id 就可以了
	offset := -v.batchSize
	base := v.base.WithContext(ctx)
	target := v.target.WithContext(ctx)
	for {
		offset += v.batchSize
		var ts []T
		err := target.Model(new(T)).Select("id").Offset(offset).
			Limit(v.batchSize).Find(&ts).Error
		if err == gorm.ErrRecordNotFound {
			return nil
		}
		if err != nil {
			v.l.Error("dst => src 查询目标表失败", logger.Error(err))
			continue
		}
		ids := slice.Map(ts, func(idx int, src T) int64 {
			return src.ID()
		})
		var srcTs []T
		err = base.Select("id").Where("id IN ?", ids).Find(&srcTs).Error
		switch err {
		case gorm.ErrRecordNotFound:
			// 说明 ids 全部没有
			v.notifySrcMissing(ts)
		case nil:
			// 计算差集
			missing := slice.DiffSetFunc(ts, srcTs, func(src, dst T) bool {
				return src.ID() == dst.ID()
			})
			v.notifySrcMissing(missing)
		default:
			v.l.Error("dst => src 查询源表失败", logger.Error(err))
		}

		if len(ts) < v.batchSize {
			// 数据没了
			return nil
		}
	}
}

// 上报不一致的数据
func (v *Validator[T]) notify(id int64, typ string) {
	// 这里我们要单独控制超时时间
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	evt := events.InconsistentEvent{
		Direction: v.direction,
		ID:        id,
		Type:      typ,
	}

	err := v.producer.ProduceInconsistentEvent(ctx, evt)
	if err != nil {
		v.l.Error("发送消息失败", logger.Error(err),
			logger.Field{Key: "event", Value: evt})
	}
}

func (v *Validator[T]) notifySrcMissing(ts []T) {
	for _, t := range ts {
		v.notify(t.ID(), events.InconsistentEventTypeBaseMissing)
	}
}
