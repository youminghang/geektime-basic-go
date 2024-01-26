package dao

import (
	"context"
	"encoding/json"
	"github.com/olivere/elastic/v7"
	"strconv"
	"strings"
)

const ArticleIndexName = "article_index"

type Article struct {
	Id      int64  `json:"id"`
	Title   string `json:"title"`
	Status  int32  `json:"status"`
	Content string `json:"content"`
}

type ArticleElasticDAO struct {
	client *elastic.Client
}

func NewArticleElasticDAO(client *elastic.Client) ArticleDAO {
	return &ArticleElasticDAO{client: client}
}

func (h *ArticleElasticDAO) Search(ctx context.Context, keywords []string) ([]Article, error) {
	queryString := strings.Join(keywords, " ")
	query := elastic.NewBoolQuery().Must(
		elastic.NewBoolQuery().Should(
			elastic.NewMatchQuery("title", queryString),
			elastic.NewMatchQuery("content", queryString)),
		elastic.NewTermQuery("status", 2))
	resp, err := h.client.Search(ArticleIndexName).Query(query).Do(ctx)
	if err != nil {
		return nil, err
	}
	res := make([]Article, 0, len(resp.Hits.Hits))
	for _, hit := range resp.Hits.Hits {
		var ele Article
		err = json.Unmarshal(hit.Source, &ele)
		res = append(res, ele)
	}
	return res, nil
}

func NewArticleRepository(client *elastic.Client) ArticleDAO {
	return &ArticleElasticDAO{
		client: client,
	}
}
func (h *ArticleElasticDAO) InputArticle(ctx context.Context, art Article) error {
	_, err := h.client.Index().
		Index(ArticleIndexName).
		Id(strconv.FormatInt(art.Id, 10)).
		BodyJson(art).Do(ctx)
	return err
}
