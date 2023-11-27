package dao

import (
	"gitee.com/geekbang/basic-go/webook/interactive/repository/dao"
	"gitee.com/geekbang/basic-go/webook/internal/repository/dao/article"
	"gorm.io/gorm"
)

func InitTables(db *gorm.DB) error {
	return db.AutoMigrate(&User{}, &article.Article{},
		&article.PublishedArticle{},
		&article.PublishedArticleV1{},
		&AsyncSms{},
		&dao.Interactive{},
		&dao.UserLikeBiz{},
		&dao.Collection{},
		&dao.UserCollectionBiz{},
		&Job{},
	)
}
