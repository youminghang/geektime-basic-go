package migrator

import (
	"gitee.com/geekbang/basic-go/webook/migrator/events"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Fixer[T Entity] struct {
	srcDB   *gorm.DB
	dstDB   *gorm.DB
	columns []string
}

func (f *Fixer[T]) Fix(event events.InconsistentEvent) error {
	var src T
	err := f.srcDB.Where("id = ?", event.ID).First(&src).Error
	if err != nil {
		return err
	}
	return f.dstDB.Clauses(&clause.OnConflict{
		// 我们需要 Entity 告诉我们，修复哪些数据
		DoUpdates: clause.AssignmentColumns(f.columns),
	}).Create(&src).Error
}
