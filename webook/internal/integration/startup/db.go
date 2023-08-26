package startup

import (
	"gitee.com/geekbang/basic-go/webook/ioc"
	"gorm.io/gorm"
)

var db *gorm.DB

// InitTestDB 测试的话，不用控制并发。等遇到了并发问题再说
func InitTestDB() *gorm.DB {
	if db == nil {
		db = ioc.InitDB()
	}
	return db
}
