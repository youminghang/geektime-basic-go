package integration

import (
	"context"
	searchv1 "gitee.com/geekbang/basic-go/webook/api/proto/gen/search/v1"
	"gitee.com/geekbang/basic-go/webook/search/grpc"
	"gitee.com/geekbang/basic-go/webook/search/integration/startup"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"testing"
	"time"
)

type SearchTestSuite struct {
	suite.Suite
	searchSvc *grpc.SearchServiceServer
	syncSvc   *grpc.SyncServiceServer
}

func (s *SearchTestSuite) SetupSuite() {
	s.searchSvc = startup.InitSearchServer()
	s.syncSvc = startup.InitSyncServer()
}

func (s *SearchTestSuite) TestSearch() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	_, err := s.syncSvc.InputUser(ctx, &searchv1.InputUserRequest{
		User: &searchv1.User{
			Id:       123,
			Nickname: "Tom White",
		},
	})
	require.NoError(s.T(), err)
	_, err = s.syncSvc.InputArticle(ctx, &searchv1.InputArticleRequest{
		Article: &searchv1.Article{
			Id:     123,
			Title:  "Tom 的小秘密",
			Status: 2,
		},
	})
	require.NoError(s.T(), err)
	_, err = s.syncSvc.InputArticle(ctx, &searchv1.InputArticleRequest{
		Article: &searchv1.Article{
			Id:      124,
			Content: "这是内容，Tom 的小秘密",
			Status:  2,
		},
	})
	require.NoError(s.T(), err)
	resp, err := s.searchSvc.Search(ctx, &searchv1.SearchRequest{
		Expression: "Tom 内容",
	})
	require.NoError(s.T(), err)
	assert.Equal(s.T(), 1, len(resp.User.Users))
	assert.Equal(s.T(), 2, len(resp.Article.Articles))
}

func TestSearchService(t *testing.T) {
	suite.Run(t, new(SearchTestSuite))
}
