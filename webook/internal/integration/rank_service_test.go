package integration

import (
	"context"
	"encoding/json"
	"gitee.com/geekbang/basic-go/webook/interactive/repository/dao"
	"gitee.com/geekbang/basic-go/webook/internal/integration/startup"
	"gitee.com/geekbang/basic-go/webook/internal/repository/dao/article"
	"gitee.com/geekbang/basic-go/webook/internal/service"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
	"testing"
	"time"
)

type RankingServiceTestSuite struct {
	suite.Suite
	db  *gorm.DB
	rdb redis.Cmdable
}

func TestRankService(t *testing.T) {
	suite.Run(t, &RankingServiceTestSuite{})
}

func (r *RankingServiceTestSuite) SetupSuite() {
	r.rdb = startup.InitRedis()
	r.db = startup.InitTestDB()
}

func (r *RankingServiceTestSuite) TearDownTest() {
	err := r.db.Exec("TRUNCATE TABLE `interactives`").Error
	require.NoError(r.T(), err)
	err = r.db.Exec("TRUNCATE TABLE `published_articles`").Error
	require.NoError(r.T(), err)
}

func (r *RankingServiceTestSuite) TestRankTopN() {
	// 设置一分钟过期时间
	svc := startup.InitRankingService().(*service.BatchRankingService)
	svc.BatchSize = 10
	svc.N = 10
	rdb := startup.InitRedis()
	db := startup.InitTestDB()
	testCases := []struct {
		name string
		// 这个测试复杂之处就在于，要准备好数据
		before func(t *testing.T)
		after  func(t *testing.T)

		wantErr error
	}{
		{
			name: "计算成功",
			before: func(t *testing.T) {
				ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
				defer cancel()
				// id 小的，点赞数多，并且发表时间比较晚
				now := time.Now()
				db = db.WithContext(ctx)
				// 准备一百条数据
				for i := 0; i < 100; i++ {
					err := db.Create(&dao.Interactive{
						BizId:   int64(i + 1),
						Biz:     "article",
						LikeCnt: int64(1000 - i*10),
					}).Error
					require.NoError(t, err)
					err = db.Create(article.PublishedArticle{
						Id:    int64(i + 1),
						Utime: now.Add(-time.Duration(i) * time.Hour).Unix(),
					}).Error
					require.NoError(t, err)
				}
			},
			after: func(t *testing.T) {
				ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
				defer cancel()

				vals, err := rdb.Get(ctx, "ranking:article").Bytes()
				require.NoError(t, err)
				var data []int64
				err = json.Unmarshal(vals, &data)
				require.NoError(t, err)
				assert.Equal(t, []int64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}, data)
			},
		},
	}

	for _, tc := range testCases {
		r.T().Run(tc.name, func(t *testing.T) {
			tc.before(t)
			ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
			defer cancel()
			err := svc.RankTopN(ctx)
			assert.Equal(t, tc.wantErr, err)
			tc.after(t)
		})
	}
}
