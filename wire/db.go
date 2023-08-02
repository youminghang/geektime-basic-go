package wire

import "gorm.io/gorm"

func InitDB() *gorm.DB {
	// 因为主要是演示 wire，所以这里我们随便写一下
	return &gorm.DB{}
}
