package fixer

import (
	"errors"
	"gitee.com/geekbang/basic-go/webook/migrator"
	"gitee.com/geekbang/basic-go/webook/migrator/events"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// Fixer T 用 any 也可以。
// 但是我强迫症觉得应该和 validator 保持一直
type Fixer[T migrator.Entity] struct {
	base    *gorm.DB
	target  *gorm.DB
	columns []string
}

func NewFixer[T migrator.Entity](base *gorm.DB,
	target *gorm.DB) (*Fixer[T], error) {
	// 在这里需要查询一下数据库中究竟有哪些列
	var t T
	rows, err := target.Model(&t).Rows()
	if err != nil {
		return nil, err
	}
	columns, err := rows.Columns()
	if err != nil {
		return nil, err
	}
	return &Fixer[T]{
		base:    base,
		target:  target,
		columns: columns,
	}, nil
}

func (f *Fixer[T]) Fix(event events.InconsistentEvent) error {
	var src T
	// 找出数据
	err := f.base.Where("id = ?", event.ID).
		First(&src).Error
	if err != nil {
		return err
	}
	switch err {
	// 找到了数据
	case nil:
		return f.target.Clauses(&clause.OnConflict{
			// 我们需要 Entity 告诉我们，修复哪些数据
			DoUpdates: clause.AssignmentColumns(f.columns),
		}).Create(&src).Error
	case gorm.ErrRecordNotFound:
		return f.target.Delete("id = ?", event.ID).Error
	default:
		return err
	}
}

// FixV1 看上去会更加符合直觉，但是有点多余的代码
func (f *Fixer[T]) FixV1(event events.InconsistentEvent) error {
	switch event.Type {
	case events.InconsistentEventTypeTargetMissing,
		events.InconsistentEventTypeNotEqual:
		var src T
		// 找出数据
		err := f.base.Where("id = ?", event.ID).
			First(&src).Error
		if err != nil {
			return err
		}
		switch err {
		// 找到了数据
		case nil:
			return f.target.Clauses(&clause.OnConflict{
				// 我们需要 Entity 告诉我们，修复哪些数据
				DoUpdates: clause.AssignmentColumns(f.columns),
			}).Create(&src).Error
		case gorm.ErrRecordNotFound:
			return f.target.Delete("id = ?", event.ID).Error
		default:
			return err
		}
	case events.InconsistentEventTypeBaseMissing:
		return f.target.Delete("id = ?", event.ID).Error
	default:
		return errors.New("未知数据不一致类型")
	}
}
