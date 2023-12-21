package service

import (
	"context"
	articlev1 "gitee.com/geekbang/basic-go/webook/api/proto/gen/article/v1"
	intrv1 "gitee.com/geekbang/basic-go/webook/api/proto/gen/intr/v1"
	"gitee.com/geekbang/basic-go/webook/ranking/domain"
	"gitee.com/geekbang/basic-go/webook/ranking/repository"
	"github.com/ecodeclub/ekit/queue"
	"github.com/ecodeclub/ekit/slice"
	"google.golang.org/protobuf/types/known/timestamppb"
	"math"
	"time"
)

//go:generate mockgen -source=./ranking.go -package=svcmocks -destination=./mocks/ranking.mock.go RankingService
type RankingService interface {
	// RankTopN 计算 TopN
	RankTopN(ctx context.Context) error
	// TopN 返回业务的 ID
	TopN(ctx context.Context) ([]domain.Article, error)
}

// BatchRankingService 分批计算
type BatchRankingService struct {
	intrSvc intrv1.InteractiveServiceClient
	artSvc  articlev1.ArticleServiceClient
	repo    repository.RankingRepository
	// 为了测试，不得已暴露出去
	BatchSize int
	N         int
	// 将来扩展，以及支持测试
	scoreFunc func(likeCnt int64, utime time.Time) float64
}

func NewBatchRankingService(
	intrSvc intrv1.InteractiveServiceClient,
	artSvc articlev1.ArticleServiceClient,
	repo repository.RankingRepository) RankingService {
	res := &BatchRankingService{
		intrSvc:   intrSvc,
		artSvc:    artSvc,
		repo:      repo,
		BatchSize: 100,
		N:         100,
	}
	res.scoreFunc = res.score
	return res
}

func (a *BatchRankingService) RankTopN(ctx context.Context) error {
	arts, err := a.rankTopN(ctx)
	if err != nil {
		return err
	}
	// 准备放到缓存里面
	return a.repo.ReplaceTopN(ctx, arts)
}

func (a *BatchRankingService) rankTopN(ctx context.Context) ([]domain.Article, error) {
	now := time.Now()
	// 只计算七天内的，因为超过七天的我们可以认为绝对不可能成为热榜了
	// 如果一个批次里面 utime 最小已经是七天之前的，我们就中断当前计算
	ddl := now.Add(-time.Hour * 24 * 7)
	offset := 0
	type Score struct {
		art   domain.Article
		score float64
	}
	// 这是一个优先级队列，维持住了 topN 的 id。
	que := queue.NewPriorityQueue[Score](a.N,
		func(src Score, dst Score) int {
			if src.score > dst.score {
				return 1
			} else if src.score == dst.score {
				return 0
			} else {
				return -1
			}
		})

	for {
		arts, err := a.artSvc.ListPub(ctx, &articlev1.ListPubRequest{
			StartTime: timestamppb.New(now),
			Offset:    int32(offset),
			Limit:     int32(a.BatchSize),
		})
		if err != nil {
			return nil, err
		}
		// 转化成 domain Article
		domainArts := make([]domain.Article, 0, len(arts.Articles))
		for _, art := range arts.Articles {
			domainArts = append(domainArts, articleToDomain(art))
		}

		artIds := slice.Map[domain.Article, int64](domainArts, func(idx int, src domain.Article) int64 {
			return src.Id
		})
		intrResp, err := a.intrSvc.GetByIds(ctx, &intrv1.GetByIdsRequest{
			Biz: "article", Ids: artIds,
		})
		if err != nil {
			return nil, err
		}
		minScore := float64(0)
		for _, art := range domainArts {
			intr, ok := intrResp.GetIntrs()[art.Id]
			if !ok {
				continue
			}
			score := a.scoreFunc(intr.LikeCnt, art.Utime)
			if score > minScore {
				ele := Score{art: art, score: score}
				err = que.Enqueue(ele)
				if err == queue.ErrOutOfCapacity {
					_, _ = que.Dequeue()
					err = que.Enqueue(ele)
				}
			} else {
				minScore = score
			}
		}
		if len(domainArts) == 0 || len(domainArts) < a.BatchSize ||
			domainArts[len(domainArts)-1].Utime.Before(ddl) {
			break
		}
		offset = offset + len(domainArts)
	}
	ql := que.Len()
	res := make([]domain.Article, ql)
	for i := ql - 1; i >= 0; i-- {
		val, _ := que.Dequeue()
		res[i] = val.art
	}
	return res, nil
}

// 这里不需要提前抽象算法，因为正常一家公司的算法都是固定的，不会今天切换到这里，明天切换到那里
func (a *BatchRankingService) score(likeCnt int64, utime time.Time) float64 {
	// 这个 factor 也可以做成一个参数
	const factor = 1.5
	return float64(likeCnt-1) /
		math.Pow(time.Since(utime).Hours()+2, factor)
}

func (a *BatchRankingService) TopN(ctx context.Context) ([]domain.Article, error) {
	return a.repo.GetTopN(ctx)
}

func articleToDomain(article *articlev1.Article) domain.Article {
	domainArticle := domain.Article{}
	if article != nil {
		domainArticle.Id = article.GetId()
		domainArticle.Title = article.GetTitle()
		domainArticle.Status = domain.ArticleStatus(article.Status)
		domainArticle.Content = article.Content
		domainArticle.Author = domain.Author{
			Id:   article.GetAuthor().GetId(),
			Name: article.GetAuthor().GetName(),
		}
		domainArticle.Ctime = article.Ctime.AsTime()
		domainArticle.Utime = article.Utime.AsTime()
	}
	return domainArticle
}
