package repository

import (
	"context"
	"gitee.com/geekbang/basic-go/webook/internal/domain"
	"gitee.com/geekbang/basic-go/webook/internal/repository/dao"
)

// ArticleAuthorRepository 演示在 service 层面上分流
type ArticleAuthorRepository interface {
	Create(ctx context.Context, art domain.Article) (int64, error)
	Update(ctx context.Context, art domain.Article) error
}

// CachedArticleAuthorRepository 按照道理，这里也是可以搞缓存的
type CachedArticleAuthorRepository struct {
	dao dao.ArticleDAO
}

func NewArticleAuthorRepository(dao dao.ArticleDAO) ArticleAuthorRepository {
	return &CachedArticleAuthorRepository{
		dao: dao,
	}
}

func (repo *CachedArticleAuthorRepository) Create(ctx context.Context,
	art domain.Article) (int64, error) {
	return repo.dao.Create(ctx, repo.toEntity(art))
}

func (repo *CachedArticleAuthorRepository) Update(ctx context.Context,
	art domain.Article) error {
	return repo.dao.UpdateById(ctx, repo.toEntity(art))
}

func (repo *CachedArticleAuthorRepository) toEntity(art domain.Article) dao.Article {
	return dao.Article{
		Id:       art.Id,
		Title:    art.Title,
		Content:  art.Content,
		AuthorId: art.Author.Id,
	}
}
